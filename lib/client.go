package lib

import "github.com/xanzy/go-gitlab"

var Augmentum *gitlab.Client

const MAX_SIZE = 100

var SettinsService *Settings

type Settings struct {
	PrivateToken    string
	GitlabIndex     string
	GitlabApi       string
	StatsComment    bool
	TopContributors int
	Group           string
	Year            int
}

func Initialize(settings *Settings) {
	SettinsService = settings
	SettinsService.initialize()
}

func (settings *Settings) initialize() {
	Augmentum = gitlab.NewClient(nil, settings.PrivateToken)
	Augmentum.SetBaseURL(settings.GitlabApi)
}

func (settings *Settings) GetGitlabIndex() string {
	return settings.GitlabIndex
}

func (settings *Settings) StatsComments() bool {
	return settings.StatsComment
}

func (settings *Settings) TopContributor() int {
	return settings.TopContributors
}

func (settings *Settings) StatsYear() int {
	return settings.Year
}

func (settings *Settings) GroupName() string {
	return settings.Group
}
