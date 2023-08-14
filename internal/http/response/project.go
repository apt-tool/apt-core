package response

import (
	"fmt"
	"time"

	"github.com/automated-pen-testing/api/pkg/models/project"
)

type (
	SetResponse struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	ProjectResponse struct {
		ID          uint                `json:"id"`
		Name        string              `json:"name"`
		Description string              `json:"description"`
		Host        string              `json:"host"`
		Creator     string              `json:"creator"`
		CreatedAt   time.Time           `json:"created_at"`
		Labels      []*SetResponse      `json:"labels"`
		Documents   []*DocumentResponse `json:"documents"`
	}
)

func (p ProjectResponse) DTO(project *project.Project) *ProjectResponse {
	p.ID = project.ID
	p.Name = project.Name
	p.Description = project.Description
	p.Creator = project.Creator
	p.CreatedAt = project.CreatedAt

	p.Host = p.createHost(project.Host, project.Port, project.HTTPSecure)

	list1 := make([]*SetResponse, 0)

	for _, item := range project.Labels {
		list1 = append(list1, &SetResponse{
			Key:   item.Key,
			Value: item.Value,
		})
	}

	list := make([]*DocumentResponse, 0)

	for _, item := range project.Documents {
		list = append(list, DocumentResponse{}.DTO(item))
	}

	p.Labels = list1
	p.Documents = list

	return &p
}

func (p ProjectResponse) createHost(host string, port int, secure bool) string {
	tmp := "http"
	if secure {
		tmp = "https"
	}

	return fmt.Sprintf("%s://%s:%d", tmp, host, port)
}
