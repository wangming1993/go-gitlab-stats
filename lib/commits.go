package lib

import (
	"log"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Commit struct {
	Id        string
	Title     string
	Message   string
	Author    string
	CreatedAt string
}

func Commits(p *Project) ([]*Commit, error) {
	cs, err := getAllCommits(p)
	var cm []*Commit
	for _, c := range cs {
		createdAt := c.CreatedAt.Unix()
		if createdAt < StatsStartTime().Unix() || createdAt > StatsEndTime().Unix() {
			continue
		}
		if isMR(c.Title) {
			// Ignore merge request commit message
			continue
		}
		cm = append(cm, ToCommit(c))
	}

	log.Printf("Get %d commits form project %s, pid = %d", len(cm), p.Name, p.Id)
	return cm, err
}

func getCommitsByPage(p *Project, page int) ([]*gitlab.Commit, *gitlab.Response, error) {
	opt := &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: MAX_SIZE,
			Page:    page,
		},
		Since: StatsStartTime(),
		Until: StatsEndTime(),
	}
	cs, header, err := Augmentum.Commits.ListCommits(p.Id, opt)
	return cs, header, err
}

func getAllCommits(p *Project) ([]*gitlab.Commit, error) {
	page := 1
	var cms []*gitlab.Commit

	for {
		resp, _, err := getCommitsByPage(p, page)
		cms = append(cms, resp...)
		page++
		if len(resp) < MAX_SIZE || err != nil {
			return cms, err
		}
	}
}

func ToCommit(c *gitlab.Commit) *Commit {
	return &Commit{
		Id:        c.ID,
		Title:     c.Title,
		Message:   c.Message,
		Author:    c.AuthorName,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}

func isMR(title string) bool {
	return strings.HasPrefix(title, "Merge branch")
}

func GetContributes(cms []*Commit) map[string]int {
	contributes := make(map[string]int)
	for _, c := range cms {
		author := c.Author
		if _, ok := contributes[author]; ok {
			contributes[author]++
		} else {
			contributes[author] = 1
		}
	}
	return contributes
}
