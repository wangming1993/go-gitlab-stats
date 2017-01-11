package lib

import (
	"os"

	"github.com/cbroglie/mustache"
)

var template string = "templates/project.mustache"

type File struct {
	Name string
}

func (p *Project) WriteHtml() error {
	out, _ := mustache.RenderFile(template,
		map[string]interface{}{
			"Name":         p.Name,
			"Url":          p.Url,
			"Contributes":  p.Contributes,
			"MergeRequest": p.MR,
			"MrCount":      p.MrCount,
			"CommitCount":  p.CommitCount,
			"CommentCount": p.CommentCount,
		},
	)

	file, err := os.Create("htmls/" + p.Name + ".html")
	if err != nil {
		return err
	}

	_, err = file.WriteString(out)
	return err
}

func (g *Group) WriteHtml() error {
	out, _ := mustache.RenderFile("templates/group.mustache",
		map[string]interface{}{
			"Contributes":  g.Contributes,
			"MrCount":      g.MrCount,
			"CommitCount":  g.CommitCount,
			"Url":          g.Url,
			"Name":         g.Name,
			"Description":  g.Description,
			"CommentCount": g.CommentCount,
		},
	)

	file, err := os.Create("htmls/scrm_group.html")
	if err != nil {
		return err
	}

	_, err = file.WriteString(out)
	return err
}
