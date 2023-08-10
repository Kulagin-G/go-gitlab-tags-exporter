package gitlab

import (
	"context"
	"github.com/xanzy/go-gitlab"
	. "go-gitlab-tags-exporter/internal/config"
	"go-gitlab-tags-exporter/internal/lib/logger/sl"
	"golang.org/x/exp/slog"
	"net/http"
	"sync"
)

const (
	tagsPerPage = 100
)

type Gitlab struct {
	cfg    *Config
	log    *slog.Logger
	client *gitlab.Client
}

func NewClient(cfg *Config, log *slog.Logger) (*Gitlab, error) {
	client, err := gitlab.NewClient(cfg.Exporter.GitlabApiToken,
		gitlab.WithBaseURL(cfg.Exporter.GitlabUrl),
		gitlab.WithCustomRetryMax(cfg.Exporter.GitlabApiRetryMax),
		gitlab.WithCustomRetryWaitMinMax(1, 5),
		gitlab.WithCustomRetry(func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			if resp != nil && (resp.StatusCode == 500 || resp.StatusCode == 503) {
				log.Warn("Gitlab API error, retrying",
					slog.String("message", resp.Status),
					slog.Int("status_code", resp.StatusCode),
				)
				return true, nil
			}
			return false, nil
		}),
	)

	if err != nil {
		log.Error("Error during Gitlab client initialization", sl.Err(err))

		return nil, err
	}

	gl := &Gitlab{
		cfg:    cfg,
		log:    log,
		client: client,
	}

	return gl, nil
}

func (g *Gitlab) tags(path string, ctx context.Context) ([]*gitlab.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, g.cfg.Exporter.GoroutinesTimeout)
	defer cancel()

	opt := &gitlab.ListTagsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: tagsPerPage,
		},
	}

	var allTags []*gitlab.Tag

	for {
		tags, resp, err := g.client.Tags.ListTags(path, opt, gitlab.WithContext(ctxTimeout))

		if err != nil {
			return nil, err
		}

		g.log.Debug("Getting tags", slog.String("path", path), slog.Int("page", resp.NextPage))

		allTags = append(allTags, tags...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allTags, nil
}

func (g *Gitlab) AllProjectsTags() []ProjectData {
	projectsCount := len(g.cfg.Projects)
	data := make([]ProjectData, 0, projectsCount)
	mu := sync.Mutex{}
	ctx := context.Background()

	limiter := make(chan struct{}, g.cfg.Exporter.GoroutinesMax)

	wg := sync.WaitGroup{}
	wg.Add(projectsCount)

	for _, project := range g.cfg.Projects {
		limiter <- struct{}{}
		go func(project *ProjectConfigs) {
			g.log.Info("Getting tags from project", slog.String("name", project.Name))

			t, err := g.tags(project.Path, ctx)

			if err != nil {
				g.log.Error("Finish goroutine for project", slog.String("name", project.Name), sl.Err(err))
				wg.Done()
				<-limiter

				return
			}

			if len(t) == 0 {
				g.log.Warn("Project has no tags", slog.String("name", project.Name))
				wg.Done()
				<-limiter

				return
			}

			mu.Lock()
			data = append(data, ProjectData{
				Name:       project.Name,
				Path:       project.Path,
				Repository: project.Repository,
				Tags:       t,
			})
			mu.Unlock()

			wg.Done()
			<-limiter
		}(project)
	}

	wg.Wait()

	return data
}
