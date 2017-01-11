package main

import (
	"log"
	"sync"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/wangming1993/go-gitlab-stats/lib"
)

func init() {
	conf := "conf/dev.toml"

	config, err := toml.LoadFile(conf)
	if err != nil {
		log.Fatalln(err)
	}
	settings := &lib.Settings{}

	user := config.Get("private_token")
	prompt(user, "rivate_token")
	settings.PrivateToken = user.(string)

	gitlabApi := config.Get("gitlab_api")
	prompt(gitlabApi, "gitlab_api")
	settings.GitlabApi = gitlabApi.(string)

	gitlabIndex := config.Get("gitlab_index")
	prompt(gitlabIndex, "gitlab_index")
	settings.GitlabIndex = gitlabIndex.(string)

	topContributor := config.GetDefault("top_contributor", 10)
	settings.TopContributors = int(topContributor.(int64))

	group := config.Get("group")
	prompt(group, "group")
	settings.Group = group.(string)

	statsComments := config.GetDefault("stats_comments", false)
	settings.StatsComment = statsComments.(bool)

	year := config.GetDefault("year", time.Now().Year())
	settings.Year = int(year.(int64))

	lib.Initialize(settings)
}

func main() {
	start := time.Now()

	statsGroup := lib.GroupService.Group(lib.SettinsService.Group)
	if nil == statsGroup {
		return
	}

	groupProjects, err := lib.GroupService.Projects(statsGroup.ID)
	lib.PanicIfError(err)

	var wg sync.WaitGroup
	var projects []*lib.Project = make([]*lib.Project, len(groupProjects))
	for i, p := range groupProjects {
		wg.Add(1)
		project := lib.ToProject(p)
		go func(project *lib.Project, index int) {
			defer wg.Done()
			project.Stats()
			//project.CountStats()
			if project.CommitCount > 0 {
				//lib.WriteJSON(project, project.Name)
				project.WriteHtml()
			}

			projects[index] = project

		}(project, i)
	}
	wg.Wait()

	group := &lib.Group{
		Name:        statsGroup.Name,
		Description: statsGroup.Description,
	}
	group.SetWebUrl()
	group.Stats(projects)
	group.WriteHtml()
	end := time.Now()
	log.Printf("Cost: %d", end.Unix()-start.Unix())

}

func prompt(data interface{}, key string) {
	if data == nil {
		log.Fatalf("[ERROR]: required %s", key)
	}
}
