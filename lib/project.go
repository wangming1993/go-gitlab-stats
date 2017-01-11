package lib

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Project struct {
	Id           int
	Url          string
	Name         string
	Description  string
	CommitCount  int
	Contributes  Contributes `json:"top10"`
	AuthorCount  int
	CreatedAt    time.Time
	Commits      []*Commit      `json:"-"`
	Authors      map[string]int `json:"-"`
	MR           []*MR
	MrCount      int
	CommentCount int
}

type Projects []*Project

func (p Projects) Len() int {
	return len(p)
}

func (p Projects) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Projects) Less(i, j int) bool {
	return p[i].CreatedAt.Unix() > p[j].CreatedAt.Unix()
}

type Contribute struct {
	Name  string
	Total int
}

type Contributes []Contribute

func (c Contributes) Len() int {
	return len(c)
}

func (c Contributes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Contributes) Less(i, j int) bool {
	return c[i].Total > c[j].Total
}

func ToProject(p *gitlab.Project) *Project {
	return &Project{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   *p.CreatedAt,
		Url:         p.WebURL,
	}
}

func (p *Project) WithCommits(c []*Commit) {
	p.Commits = c
	p.CommitCount = len(c)
}

func (p *Project) WithAuthors(authors map[string]int) {
	p.Authors = authors
	var cs Contributes
	for name, total := range authors {
		cs = append(cs, Contribute{Name: name, Total: total})
	}
	sort.Sort(cs)

	p.AuthorCount = cs.Len()
	p.Contributes = cs
	top := SettinsService.TopContributor()
	if cs.Len() > top {
		p.Contributes = cs[0:top]
	}
}

func (p *Project) Stats() error {
	cms, err := Commits(p)
	if err != nil {
		log.Println(err)
		return err
	}
	p.WithCommits(cms)

	authors := GetContributes(cms)
	p.WithAuthors(authors)

	mrs := p.AllMR()
	p.MR = mrs
	p.MrCount = len(mrs)

	p.CommentCount = p.GetCommentCount()

	return nil
}

func (p *Project) AllMR() []*MR {
	resp, err := p.getAllMr()
	if err != nil {
		log.Println(err)
		return nil
	}
	var mrs []*MR = make([]*MR, len(resp))

	for index, m := range resp {
		mr := &MR{
			Id:          m.ID,
			PID:         p.Id,
			Title:       m.Title,
			Description: m.Description,
			Url:         fmt.Sprintf("%s/merge_requests/%d", p.Url, m.IID),
		}
		func(index int, mr *MR) {
			mr.Stats()
			mrs[index] = mr
		}(index, mr)

	}

	return mrs
}

func (p *Project) getMrByPage(page int) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	opt := &gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: MAX_SIZE,
			Page:    page,
		},
		Sort:    gitlab.String("desc"),
		OrderBy: gitlab.String("created_at"),
	}
	resp, header, err := Augmentum.MergeRequests.ListMergeRequests(p.Id, opt)
	if err != nil {
		log.Println(err)
		return nil, header, err
	}
	return resp, header, nil
}

func (p *Project) getAllMr() ([]*gitlab.MergeRequest, error) {
	page := 1
	var mrs []*gitlab.MergeRequest

	for {
		resp, _, err := p.getMrByPage(page)

		for _, c := range resp {
			if c.CreatedAt.Unix() < StatsStartTime().Unix() {
				return mrs, err
			}
			if c.CreatedAt.Unix() > StatsEndTime().Unix() {
				continue
			}
			mrs = append(mrs, c)
		}

		page++
		if len(resp) < MAX_SIZE || err != nil {
			return mrs, err
		}
	}
}

func (p *Project) GetCommitCount() int {
	resp, header, err := getCommitsByPage(p, 1)
	if err != nil {
		log.Println(err)
		return 0
	}
	last := header.LastPage
	if last > 1 {
		resp, _, err = getCommitsByPage(p, last)
		return (last-1)*MAX_SIZE + len(resp)
	}
	return len(resp)
}

func (p *Project) GetMrCount() int {
	resp, header, err := p.getMrByPage(1)
	if err != nil {
		log.Println(err)
		return 0
	}
	last := header.LastPage
	log.Println("last:", header)
	if last > 1 {
		resp, _, err = p.getMrByPage(last)
		return (last-1)*MAX_SIZE + len(resp)
	}
	return len(resp)
}

func (p *Project) CountStats() {
	p.CommitCount = p.GetCommitCount()
	p.MrCount = p.GetMrCount()
}

func (p *Project) GetCommentCount() int {
	var commentCount int
	for _, mr := range p.MR {
		commentCount += mr.CommentCount
	}
	return commentCount
}
