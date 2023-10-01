package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ptaas-tool/base-api/internal/config/ftp"
	"github.com/ptaas-tool/base-api/internal/core/ai"
	"github.com/ptaas-tool/base-api/internal/core/scanner"
	"github.com/ptaas-tool/base-api/internal/utils/crypto"
	"github.com/ptaas-tool/base-api/pkg/client"
	"github.com/ptaas-tool/base-api/pkg/enum"
	"github.com/ptaas-tool/base-api/pkg/models"
	"github.com/ptaas-tool/base-api/pkg/models/document"
)

// worker is the smallest unit of our core
type worker struct {
	channel chan int
	rerun   chan int
	done    chan int
	cfg     ftp.Config
	client  client.HTTPClient
	models  *models.Interface
	ai      *ai.AI
}

type (
	// executeRequest is used to call ftp system
	executeRequest struct {
		Param      string `json:"param"`
		Path       string `json:"path"`
		DocumentID uint   `json:"document_id"`
	}
)

// work method will do the logic of penetration testing
func (w worker) work() error {
	for {
		// get project id from channel
		id := <-w.channel

		projectID := uint(id)

		// manifests
		manifests := make([]string, 0)

		// make http request to ftp client in order to get attacks
		rsp, err := w.client.Get(w.cfg.Host)
		if err != nil {
			log.Println(fmt.Errorf("[worker.work] failed to get attacks error=%w", err))

			w.exit(id)

			continue
		}

		if er := json.NewDecoder(rsp.Body).Decode(&manifests); er != nil {
			log.Println(fmt.Errorf("[worker.work] failed to parse attacks error=%w", er))

			w.exit(id)

			continue
		}

		// get project from db
		project, er := w.models.Projects.GetByID(projectID)
		if er != nil {
			log.Println(fmt.Errorf("[worker.work] failed to get project error=%w", er))

			w.exit(id)
		}

		// set http or https
		prefix := "http"
		if project.HTTPSecure {
			prefix = "http"
		}

		// start scanner
		vulnerabilities, err := scanner.Scan(fmt.Sprintf("%s://%s:%d", prefix, project.Host, project.Port))
		if err != nil {
			log.Println(fmt.Errorf("[worker.work] failed to scan host error=%w", err))
		}

		// get attacks from ai module
		attacks := w.ai.GetAttacks(manifests, vulnerabilities)

		docs := make([]*document.Document, 0)

		// create document
		for _, attack := range attacks {
			// create document
			doc := &document.Document{
				ProjectID:   projectID,
				Instruction: attack,
				ExecutedBy:  project.Creator,
				Result:      enum.ResultNotSet,
				Status:      enum.StatusInit,
			}

			if e := w.models.Documents.Create(doc); e != nil {
				log.Println(fmt.Errorf("[worker.work] failed to create document error=%w", e))

				continue
			}

			docs = append(docs, doc)
		}

		// perform each attack
		for _, doc := range docs {
			start := time.Now()

			// update doc status
			doc.Status = enum.StatusPending
			doc.Result = enum.ResultUnknown
			_ = w.models.Documents.Update(doc)

			// create ftp request
			tmp := executeRequest{
				Param:      project.Host,
				Path:       doc.Instruction,
				DocumentID: doc.ID,
			}

			// send ftp request
			var buffer bytes.Buffer
			if e := json.NewEncoder(&buffer).Encode(tmp); e != nil {
				log.Println(fmt.Errorf("[worker.work] failed to create request error=%w", e))

				continue
			}

			address := fmt.Sprintf("%s/execute", w.cfg.Host)
			headers := []string{
				"Content-Type:application/json",
				fmt.Sprintf("x-token:%s", crypto.GetMD5Hash(w.cfg.Secret)),
			}

			// update document based of response
			if response, httpError := w.client.Post(address, &buffer, headers...); httpError != nil {
				log.Println(fmt.Errorf("[worker.work] failed to execute script error=%w", httpError))

				doc.Result = enum.ResultFailed
				doc.Status = enum.StatusFailed
			} else {
				if response.StatusCode == 200 {
					type rsp struct {
						Code int `json:"code"`
					}

					r := rsp{}

					if err := json.NewDecoder(response.Body).Decode(&r); err == nil {
						if r.Code != 0 {
							doc.Result = enum.ResultFailed
						} else {
							doc.Result = enum.ResultSuccessful
						}
					}
				}

				doc.Status = enum.StatusDone
			}

			doc.ExecutionTime = time.Now().Sub(start)

			_ = w.models.Documents.Update(doc)
		}

		w.exit(id)
	}
}

func (w worker) exit(id int) {
	w.done <- id
}
