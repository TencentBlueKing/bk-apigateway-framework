// Package serializer ...
package serializer

import (
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common/probe"
)

// HealthzResponse ...
type HealthzResponse struct {
	Time    string         `json:"time"`
	Healthy bool           `json:"healthy"`
	Fatal   bool           `json:"fatal"`
	Results []probe.Result `json:"results"`
}

// VersionResponse ...
type VersionResponse struct {
	Version     string `json:"version"`
	GitCommit   string `json:"gitCommit"`
	BuildTime   string `json:"buildTime"`
	TmplVersion string `json:"tmplVersion"`
	GoVersion   string `json:"goVersion"`
}
