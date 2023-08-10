package metrics

import (
	"github.com/xanzy/go-gitlab"
	. "go-gitlab-tags-exporter/internal/config"
	. "go-gitlab-tags-exporter/internal/gitlab/v4"
	"go-gitlab-tags-exporter/internal/lib/logger"
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	t.Log("[TEST]: Check ParseTags returns expected slice of ProjectDataSorted.")

	cfg := &Config{
		Exporter: &Exporter{
			GitlabApiToken:    "test-token",
			GitlabUrl:         "test-url",
			GitlabApiRetryMax: 5,
			GoroutinesTimeout: 5,
			GoroutinesMax:     5,
			LogLevel:          "debug",
		},
		MetricsOptions: &MetricsOptions{
			ReleaseTagRegex:          "^[0-9]+\\.[0-9]+\\.[0-9]+$",
			ReleaseCandidateTagRegex: "^[0-9]+\\.[0-9]+\\.[0-9]+-rc$",
		},
		Projects: []*ProjectConfigs{
			{
				Name: "test-project1",
				Path: "test-group1/test-project1",
			},
			{
				Name: "test-project2",
				Path: "test-group2/test-project2",
			},
			{
				Name: "test-project3",
				Path: "test-group3/test-project3",
			},
		},
	}

	log := logger.SetupLogger(cfg)

	data := []ProjectData{
		{
			Name:       "test-project1",
			Path:       "test-group1/test-project1",
			Repository: "",
			Tags: []*gitlab.Tag{
				{
					Name: "vvv1.0.0",
				},
				{
					Name: "1.10.0",
				},
				{
					Name: "1.19.0-rc",
				},
				{
					Name: "1.20.0-rc",
				},
				{
					Name: "999.999.999",
				},
			},
		},
		{
			Name:       "test-project2",
			Path:       "test-group2/test-project2",
			Repository: "",
			Tags: []*gitlab.Tag{
				{
					Name: "v2.0.0",
				},
				{
					Name: "2.10.0",
				},
				{
					Name: "1.1001.0-rc",
				},
				{
					Name: "1.100.0-rc",
				},
				{
					Name: "1.999.999",
				},
			},
		},
	}

	want := []ProjectDataSorted{
		{
			Name:       "test-project1",
			Path:       "test-group1/test-project1",
			Repository: "",
			LatestReleaseTag: &gitlab.Tag{
				Name: "999.999.999",
			},
			LatestReleaseCandidateTag: &gitlab.Tag{
				Name: "1.20.0-rc",
			},
		},
		{
			Name:       "test-project2",
			Path:       "test-group2/test-project2",
			Repository: "",
			LatestReleaseTag: &gitlab.Tag{
				Name: "2.10.0",
			},
			LatestReleaseCandidateTag: &gitlab.Tag{
				Name: "1.1001.0-rc",
			},
		},
	}

	gen := NewGenerator(cfg, log)

	result := gen.ParseTags(data)

	if !reflect.DeepEqual(want, result) {
		t.Errorf("Returned %#+v, want %#+v", result, want)
	}
}

func TestParseTagsNegative(t *testing.T) {
	t.Log("[TEST]: Check ParseTags can handle empty slice of ProjectData.")

	cfg := &Config{
		Exporter: &Exporter{
			GitlabApiToken:    "test-token",
			GitlabUrl:         "test-url",
			GitlabApiRetryMax: 5,
			GoroutinesTimeout: 5,
			GoroutinesMax:     5,
			LogLevel:          "debug",
		},
		MetricsOptions: &MetricsOptions{
			ReleaseTagRegex:          "^[0-9]+\\.[0-9]+\\.[0-9]+$",
			ReleaseCandidateTagRegex: "^[0-9]+\\.[0-9]+\\.[0-9]+-rc$",
		},
		Projects: []*ProjectConfigs{
			{
				Name: "test-project1",
				Path: "test-group1/test-project1",
			},
			{
				Name: "test-project2",
				Path: "test-group2/test-project2",
			},
			{
				Name: "test-project3",
				Path: "test-group3/test-project3",
			},
		},
	}

	log := logger.SetupLogger(cfg)

	data := []ProjectData{}

	want := []ProjectDataSorted{}

	gen := NewGenerator(cfg, log)

	result := gen.ParseTags(data)

	if !reflect.DeepEqual(want, result) {
		t.Errorf("Returned %#+v, want %#+v", result, want)
	}
}
