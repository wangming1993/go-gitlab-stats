package lib

import (
	"log"
	"regexp"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type MR struct {
	Id              int
	PID             int
	Title           string
	Description     string
	CommentCount    int
	Comment         []*Comment
	Url             string
	CommentsDisplay string `json:"-"`
}

type Comment struct {
	Note   string
	Author string
}

func (m *MR) AllComments() []*Comment {
	opt := &gitlab.GetMergeRequestCommentsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: MAX_SIZE,
		},
	}
	resp, _, err := Augmentum.MergeRequests.GetMergeRequestComments(m.PID, m.Id, opt)
	//log.Println("[Comments]", resp)
	if err != nil {
		log.Println(err)
		return nil
	}
	var cms []*Comment
	for _, c := range resp {
		if isIgnoredComments(c.Note) {
			continue
		}
		cms = append(cms, &Comment{c.Note, c.Author.Name})
	}
	return cms
}

func (m *MR) Stats() {
	if SettinsService.StatsComments() {
		cms := m.AllComments()
		m.Comment = cms
		m.CommentCount = len(cms)
	}
	if m.CommentCount > 0 {
		m.CommentsDisplay = "block"
	} else {
		m.CommentsDisplay = "none"
	}
}

func isIgnoredComments(title string) bool {
	matched, _ := regexp.MatchString("^Added [0-9]+ commit", title)
	return matched ||
		strings.HasPrefix(title, "Status changed to") ||
		strings.HasPrefix(title, "mentioned in") ||
		strings.HasPrefix(title, "Reassigned to")
}
