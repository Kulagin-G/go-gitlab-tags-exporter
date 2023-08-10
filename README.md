# go-gitlab-tags-exporter
Simple Go application that converts the latest Gitlab's release and release-candidate tags to Prometheus metrics.

# Metrics example
```bash
# HELP gitlab_tag_latest_info Returns the latest tag for a repo based on tag type
# TYPE gitlab_tag_latest_info gauge
gitlab_tag_latest_info{project_name="infrastructure",repository="",tag_name="1.75.0",tag_type="latest_release"} 1
gitlab_tag_latest_info{project_name="infrastructure",repository="",tag_name="1.76.0-rc",tag_type="latest_release_candidate"} 1
gitlab_tag_latest_info{project_name="experiments",repository="",tag_name="1.28.0-rc",tag_type="latest_release_candidate"} 1
gitlab_tag_latest_info{project_name="experiments",repository="",tag_name="1.29.0",tag_type="latest_release"} 1
# HELP gitlab_tag_parsing_duration_seconds Returns the time it took to parse all tags
# TYPE gitlab_tag_parsing_duration_seconds gauge
gitlab_tag_parsing_duration_seconds 1.181482208
````

# Local run without compiling
```bash
cd ./go-gitlab-tags-exporter/cmd/
go mod tidy
export GITLAB_API_TOKEN=<your_token>
export EXPORTER_CONFIG_PATH=<path_to_config>
go run main.go
```
For compiling options see Dockerfile.

# Tests
There are several unit tests were implemented:
```bash
cd ./go-gitlab-tags-exporter
go test -v ./...
```
# Config description
```yaml
---
#
# Main exporter configuration.
#
exporter:
  address: 0.0.0.0  # Address to bind.
  port: 8090  # Port to bind.
  logLevel: info  # info, debug
  gitlabUrl: https://gitlab.com  # Gitlab URL.
  gitlabApiRetryMax: 5  # Max number of retries to reach Gitlab API.
  goroutinesMax: 1000  # Max number of goroutines.
  goroutinesTimeout: 60s  # Timeout for goroutines in sec.
  metricsEndpoint: /metrics  # Metrics endpoint.
  health:  # Health endpoint.
    port: 8091
    name: healthz

metricsOptions:
  releaseTagPattern: '^\d+\.\d+\.\d+$'  # Regex pattern for release tags.
  releaseCandidateTagPattern: '^\d+\.\d+\.\d+-rc$'  # Regex pattern for release candidate tags.

#
# Gitlab projects configuration.
#
projects:
  - name: infrastructure
    path: sre/common/infrastructure
```

# Sealed secrets

Helm chart supports `SealedSecret` for using encrypted data in `values.yaml`:
```yaml
secrets:
  sealedSecret:
    enabled: true
    annotations:
      sealedsecrets.bitnami.com/cluster-wide: "true"
  data:
    GITLAB_API_TOKEN: "AABxxxx"
```
If you don't need to hide sensitivity data:
```yaml
secrets:
  defaultSecret:
    enabled: true
    annotations: {}
  data:
    GITLAB_API_TOKEN: "TOKENxxx"
```