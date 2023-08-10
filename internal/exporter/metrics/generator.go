package metrics

import (
	"github.com/Masterminds/semver/v3"
	"github.com/xanzy/go-gitlab"
	. "go-gitlab-tags-exporter/internal/config"
	. "go-gitlab-tags-exporter/internal/gitlab/v4"
	"go-gitlab-tags-exporter/internal/lib/logger/sl"
	"golang.org/x/exp/slog"
	. "regexp"
	"sort"
)

const (
	latestReleaseTagType          = "latest_release"
	latestReleaseCandidateTagType = "latest_release_candidate"
)

type Generator struct {
	cfg      *Config
	log      *slog.Logger
	relRegex *Regexp
	rcRegex  *Regexp
}

func NewGenerator(cfg *Config, log *slog.Logger) *Generator {
	return &Generator{
		cfg:      cfg,
		log:      log,
		relRegex: MustCompile(cfg.MetricsOptions.ReleaseTagRegex),
		rcRegex:  MustCompile(cfg.MetricsOptions.ReleaseCandidateTagRegex),
	}
}

func (g *Generator) GenerateData(gl *Gitlab) []ProjectDataSorted {
	allTags := gl.AllProjectsTags()

	if len(allTags) == 0 {
		g.log.Warn("No tags found or during parsing tags error occurred")
		return nil
	}

	return g.ParseTags(allTags)
}

func (g *Generator) ParseTags(data []ProjectData) []ProjectDataSorted {
	dataSorted := make([]ProjectDataSorted, 0, len(data))

	for _, project := range data {
		dataSorted = append(dataSorted, ProjectDataSorted{
			Name:                      project.Name,
			Path:                      project.Path,
			Repository:                project.Repository,
			LatestReleaseTag:          g.sortTags(project.Tags, latestReleaseTagType),
			LatestReleaseCandidateTag: g.sortTags(project.Tags, latestReleaseCandidateTagType),
		})
	}

	return dataSorted
}

func (g *Generator) sortTags(tags []*gitlab.Tag, tagType string) *gitlab.Tag {
	rcTags := make([]*semver.Version, 0, len(tags))
	relTags := make([]*semver.Version, 0, len(tags))

	if len(tags) == 0 {
		return &gitlab.Tag{
			Name: "notFound",
		}
	}

	for _, tag := range tags {
		v, err := semver.NewVersion(tag.Name)
		if err != nil {
			g.log.Warn("Error parsing version", slog.String("tag_name", tag.Name), sl.Err(err))
			continue
		}

		switch tagType {
		case latestReleaseCandidateTagType:
			if g.rcRegex.MatchString(tag.Name) {
				rcTags = append(rcTags, v)
			}
		case latestReleaseTagType:
			if g.relRegex.MatchString(tag.Name) {
				relTags = append(relTags, v)
			}
		}
	}

	switch tagType {
	case latestReleaseCandidateTagType:
		sort.Sort(semver.Collection(rcTags))

		return searchTag(tags, rcTags[len(rcTags)-1])
	case latestReleaseTagType:
		sort.Sort(semver.Collection(relTags))

		return searchTag(tags, relTags[len(relTags)-1])
	}

	g.log.Warn("No tag found for type: %s", tagType)

	return &gitlab.Tag{
		Name: "notFound",
	}
}

func searchTag(tags []*gitlab.Tag, pattern *semver.Version) *gitlab.Tag {
	for _, t := range tags {
		if t.Name == pattern.Original() {
			return t
		}
	}

	return &gitlab.Tag{
		Name: "None",
	}
}
