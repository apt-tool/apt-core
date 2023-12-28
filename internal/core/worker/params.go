package worker

import (
	"fmt"
	"strings"

	"github.com/ptaas-tool/base-api/pkg/models/project"
)

// generate request params from Project model
func (w worker) generateParamsFromProject(project *project.Project) []string {
	params := make([]string, 0)

	// create host address
	prefix := "http"
	if project.HTTPSecure {
		prefix = "http"
	}

	host := fmt.Sprintf("%s://%s:%d", prefix, project.Host, project.Port)
	params = append(params, "--host", host)

	// endpoints
	if len(project.Endpoints) > 0 {
		endpoints := make([]string, 0)

		for _, item := range project.Endpoints {
			endpoints = append(endpoints, item.Endpoint)
		}

		params = append(params, "--endpoints", strings.Join(endpoints, ","))
	}

	// other parameters
	if len(project.Params) > 0 {
		p := make([]string, 0)

		for _, item := range project.Params {
			p = append(p, fmt.Sprintf("%s=%s", item.Key, item.Value))
		}

		params = append(params, "--params", strings.Join(p, "&"))
	}

	return params
}
