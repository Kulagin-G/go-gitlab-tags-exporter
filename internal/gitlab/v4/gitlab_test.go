package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	. "go-gitlab-tags-exporter/internal/config"
	"go-gitlab-tags-exporter/internal/lib/logger"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func setup(path string, handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *gitlab.Client) {
	mux := http.NewServeMux()
	mux.HandleFunc(path, handler)

	server := httptest.NewServer(mux)

	client, _ := gitlab.NewClient("",
		gitlab.WithBaseURL(server.URL),
	)

	return server, client
}

func TestAllProjectsTags(t *testing.T) {
	t.Log("[TEST]: Check AllProjectsTags returns expected slice of ProjectData.")

	mock, client := setup("/api/v4/projects/test-group1/test-project1/repository/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[
			{"name": "vvv1.0.0", "message": "", "protected": false, "target": ""},
			{"name": "1.10.0", "message": "", "protected": false, "target": ""},
			{"name": "1.20.0-rc", "message": "", "protected": false, "target": ""},
			{"name": "999.999.999", "message": "", "protected": false, "target": ""}
		]`)
	})

	defer mock.Close()

	cfg := &Config{
		Exporter: &Exporter{
			GitlabApiToken:    "test-token",
			GitlabUrl:         "test-url",
			GitlabApiRetryMax: 5,
			GoroutinesTimeout: 5 * time.Second,
			GoroutinesMax:     5,
			LogLevel:          "debug",
		},
		Projects: []*ProjectConfigs{
			{
				Name: "test-project1",
				Path: "test-group1/test-project1",
			},
		},
	}

	log := logger.SetupLogger(cfg)

	gl := &Gitlab{
		cfg:    cfg,
		log:    log,
		client: client,
	}

	want := []ProjectData{
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
					Name: "1.20.0-rc",
				},
				{
					Name: "999.999.999",
				},
			},
		},
	}
	result := gl.AllProjectsTags()

	if !reflect.DeepEqual(want, result) {
		t.Errorf("Returned %#+v, want %#+v", result, want)
	}
}

func TestAllProjectsTagsNegative(t *testing.T) {
	t.Log("[TEST]: Check AllProjectsTags can handle incorrect data")

	mock, client := setup("/api/v4/projects/test-group99/test-project99/repository/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[
			{"name": "vvv1.0.0", "message": "", "protected": false, "target": ""},
			{"name": "1.10.0", "message": "", "protected": false, "target": ""},
			{"name": "1.20.0-rc", "message": "", "protected": false, "target": ""},
			{"name": "999.999.999", "message": "", "protected": false, "target": ""}
		]`)
	})

	defer mock.Close()

	cfg := &Config{
		Exporter: &Exporter{
			GitlabApiToken:    "test-token",
			GitlabUrl:         "test-url",
			GitlabApiRetryMax: 5,
			GoroutinesTimeout: 5 * time.Second,
			GoroutinesMax:     5,
			LogLevel:          "debug",
		},
		Projects: []*ProjectConfigs{
			{
				Name: "test-project1",
				Path: "test-group1/test-project1",
			},
		},
	}

	log := logger.SetupLogger(cfg)

	gl := &Gitlab{
		cfg:    cfg,
		log:    log,
		client: client,
	}

	want := make([]ProjectData, 0, 1)
	result := gl.AllProjectsTags()

	if !reflect.DeepEqual(want, result) {
		t.Errorf("Returned %#+v, want %#+v", result, want)
	}
}
