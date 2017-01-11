package lib

import (
	"log"
	"sort"

	"github.com/xanzy/go-gitlab"
)

type Group struct {
	Name         string
	Url          string
	Description  string
	Ps           Projects       `json:"-"`
	Contributes  Contributes    `json:"top10"`
	Authors      map[string]int `json:"-"`
	CommitCount  int
	AuthorCount  int
	MrCount      int
	CommentCount int
}

var GroupService *Group

func init() {
	GroupService = &Group{}
}

func (*Group) Group(name string) *gitlab.Group {
	groups, _, err := Augmentum.Groups.SearchGroup(name)
	PanicIfError(err)

	if len(groups) != 1 {
		log.Fatalf("[ERROR]: Get group failed, search by name scrm get %d group \n", len(groups))
	}

	return groups[0]
}

func (g *Group) SetWebUrl() {
	g.Url = SettinsService.GetGitlabIndex() + "/groups/" + g.Name
}

func (*Group) Projects(groupId int) ([]*gitlab.Project, error) {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: MAX_SIZE,
		},
		Search: gitlab.String(SettinsService.GroupName()),
	}
	projects, _, err := Augmentum.Projects.ListProjects(opt)
	if err != nil {
		return nil, err
	}

	var SCRMProjects []*gitlab.Project

	for _, p := range projects {
		//&& p.ID == 817
		if p.Namespace.ID == groupId {
			SCRMProjects = append(SCRMProjects, p)
		}
	}

	return SCRMProjects, nil
}

func (g *Group) Stats(proj []*Project) {
	//g.Ps = proj
	authors := make(map[string]int)

	var commitCount int
	var mrCount int
	var commentCount int
	for _, p := range proj {
		commitCount += p.CommitCount
		mrCount += p.MrCount
		commentCount += p.CommentCount
		aus := p.Authors
		for name, t := range aus {
			if _, ok := authors[name]; ok {
				authors[name] += t
			} else {
				authors[name] = 1
			}
		}
	}
	g.Authors = authors
	log.Println(len(authors))

	var cts Contributes
	for name, total := range authors {
		cts = append(cts, Contribute{name, total})
	}
	sort.Sort(cts)
	g.AuthorCount = cts.Len()
	g.CommitCount = commitCount
	g.MrCount = mrCount
	g.CommentCount = commentCount

	g.Contributes = cts
	top := SettinsService.TopContributor()
	if cts.Len() > top {
		g.Contributes = cts[0:top]
	}
}
