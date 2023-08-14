package request

import "github.com/automated-pen-testing/api/pkg/models/project"

type (
	SetRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	ProjectRequest struct {
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Host        string       `json:"host"`
		Port        int          `json:"port"`
		HTTPSecure  bool         `json:"http_secure"`
		Endpoints   []string     `json:"endpoints,omitempty"`
		Labels      []SetRequest `json:"labels,omitempty"`
		Params      []SetRequest `json:"params,omitempty"`
	}
)

func (p ProjectRequest) ToModel(namespaceID uint, creator string) *project.Project {
	params := make([]*project.ParamSet, 0)
	labels := make([]*project.LabelSet, 0)
	endpoints := make([]*project.EndpointSet, 0)

	for _, item := range p.Params {
		params = append(params, &project.ParamSet{
			Key:   item.Key,
			Value: item.Value,
		})
	}

	for _, item := range p.Labels {
		labels = append(labels, &project.LabelSet{
			Key:   item.Key,
			Value: item.Value,
		})
	}

	for _, item := range p.Endpoints {
		endpoints = append(endpoints, &project.EndpointSet{
			Endpoint: item,
		})
	}

	return &project.Project{
		Name:        p.Name,
		Host:        p.Host,
		Port:        p.Port,
		Description: p.Description,
		HTTPSecure:  p.HTTPSecure,
		Creator:     creator,
		NamespaceID: namespaceID,
		Params:      params,
		Labels:      labels,
		Endpoints:   endpoints,
	}
}
