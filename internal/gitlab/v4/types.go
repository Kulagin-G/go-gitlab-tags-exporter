package gitlab

import "github.com/xanzy/go-gitlab"

type ProjectData struct {
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
	Repository string `yaml:",omitempty"`
	Tags       []*gitlab.Tag
}

type ProjectDataSorted struct {
	Name                      string `yaml:"name"`
	Path                      string `yaml:"path"`
	Repository                string `yaml:",omitempty"`
	LatestReleaseTag          *gitlab.Tag
	LatestReleaseCandidateTag *gitlab.Tag
}
