package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ptaas-tool/base-api/internal/config/ftp"
	scannerCfg "github.com/ptaas-tool/base-api/internal/config/scanner"
	"github.com/ptaas-tool/base-api/internal/core/ai"
	"github.com/ptaas-tool/base-api/internal/core/scanner"
	"github.com/ptaas-tool/base-api/internal/utils/crypto"
	"github.com/ptaas-tool/base-api/pkg/client"
	"github.com/ptaas-tool/base-api/pkg/enum"
	"github.com/ptaas-tool/base-api/pkg/models"
	"github.com/ptaas-tool/base-api/pkg/models/document"
	"github.com/ptaas-tool/base-api/pkg/models/project"
	"github.com/ptaas-tool/base-api/pkg/models/track"
)

// worker is the smallest unit of our core
type worker struct {
	channel chan int
	reruns  chan int
	done    chan int

	cfg     ftp.Config
	scanner scannerCfg.Config
	client  client.HTTPClient
	models  *models.Interface
	ai      *ai.AI
}

type (
	// executeRequest is used to call ftp system
	executeRequest struct {
		Params     []string `json:"params"`
		Path       string   `json:"path"`
		DocumentID uint     `json:"document_id"`
	}
)

// work method will do the logic of penetration testing
func (w worker) work() error {
	for {
		select {
		case id := <-w.channel:
			w.execute(id)
		case id := <-w.reruns:
			w.rerun(id)
		}
	}
}

// rerun method will rerun a specific document
func (w worker) rerun(id int) {
	documentID := uint(id)

	// get document
	oldDoc, err := w.models.Documents.GetByID(documentID)
	if err != nil {
		log.Println(fmt.Errorf("[worker.rerun] failed to get document error=%w", err))

		w.exit(0)

		return
	}

	// get project from db
	projectInstance, er := w.models.Projects.GetByID(oldDoc.ProjectID)
	if er != nil {
		log.Println(fmt.Errorf("[worker.rerun] failed to get project error=%w", er))

		w.exit(0)

		return
	}

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectInstance.ID,
		DocumentID:  documentID,
		Service:     "base-api/worker/rerun",
		Description: "Got rerun request",
		Type:        enum.TrackInProgress,
	})

	// create new document
	doc := &document.Document{
		ProjectID:   oldDoc.ProjectID,
		Instruction: oldDoc.Instruction,
		ExecutedBy:  oldDoc.ExecutedBy,
		Result:      enum.ResultNotSet,
		Status:      enum.StatusInit,
	}

	// create new document
	if e := w.models.Documents.Create(doc); e != nil {
		log.Println(fmt.Errorf("[worker.rerun] failed to create document error=%w", e))

		w.exit(int(projectInstance.ID))

		return
	}

	// update doc status
	doc.Status = enum.StatusPending
	doc.Result = enum.ResultUnknown
	_ = w.models.Documents.Update(doc)

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectInstance.ID,
		DocumentID:  doc.ID,
		Service:     "base-api/worker/rerun",
		Description: fmt.Sprintf("Running the document on `%s` attack.", doc.Instruction),
		Type:        enum.TrackInProgress,
	})

	// execute the doc
	if err := w.executeDoc(projectInstance, doc); err != nil {
		log.Println(fmt.Errorf("[worker.rerun] failed to create request error=%w", err))
	}

	w.exit(int(projectInstance.ID))
}

// execute a project
func (w worker) execute(id int) {
	projectID := uint(id)

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectID,
		Service:     "base-api/worker",
		Description: "Got execute request",
		Type:        enum.TrackInProgress,
	})

	// manifests
	manifests := make([]string, 0)

	// make http request to ftp client in order to get attacks
	rsp, err := w.client.Get(w.cfg.Host)
	if err != nil {
		log.Println(fmt.Errorf("[worker.execute] failed to get attacks error=%w", err))

		w.exit(id)

		return
	}

	if er := json.NewDecoder(rsp.Body).Decode(&manifests); er != nil {
		log.Println(fmt.Errorf("[worker.execute] failed to parse attacks error=%w", er))

		w.exit(id)

		return
	}

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectID,
		Service:     "base-api/worker/execute",
		Description: fmt.Sprintf("Got attacks list from ftp-server: (%s)", strings.Join(manifests, ",")),
		Type:        enum.TrackSuccess,
	})

	// get project from db
	projectInstance, er := w.models.Projects.GetByID(projectID)
	if er != nil {
		log.Println(fmt.Errorf("[worker.execute] failed to get project error=%w", er))

		w.exit(id)

		return
	}

	// set http or https
	prefix := "http"
	if projectInstance.HTTPSecure {
		prefix = "http"
	}

	host := fmt.Sprintf("%s://%s:%d", prefix, projectInstance.Host, projectInstance.Port)

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectID,
		Service:     "base-api/worker/scanner",
		Description: "Running scanner",
		Type:        enum.TrackInProgress,
	})

	// start scanner
	vulnerabilities, err := scanner.Scanner{
		Enable:   w.scanner.Enable,
		Defaults: w.scanner.Defaults,
		Command:  w.scanner.Command,
	}.Scan(map[string]string{
		"host": host,
	})
	if err != nil {
		log.Println(fmt.Errorf("[worker.execute] failed to scan host error=%w", err))
	}

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectID,
		Service:     "base-api/worker/scanner",
		Description: fmt.Sprintf("Vulnerabilities: (%s)", strings.Join(vulnerabilities, ",")),
		Type:        enum.TrackSuccess,
	})

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectID,
		Service:     "base-api/worker/ai",
		Description: "Running AI",
		Type:        enum.TrackInProgress,
	})

	// get attacks from ai module
	attacks := w.ai.GetAttacks(manifests, vulnerabilities)

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   projectID,
		Service:     "base-api/worker/ai",
		Description: fmt.Sprintf("Attacks: (%s)", strings.Join(attacks, ",")),
		Type:        enum.TrackSuccess,
	})

	docs := make([]*document.Document, 0)

	// create document
	for _, attack := range attacks {
		// create document
		doc := &document.Document{
			ProjectID:   projectID,
			Instruction: attack,
			ExecutedBy:  projectInstance.Creator,
			Result:      enum.ResultNotSet,
			Status:      enum.StatusInit,
		}

		if e := w.models.Documents.Create(doc); e != nil {
			log.Println(fmt.Errorf("[worker.execute] failed to create document error=%w", e))

			continue
		}

		docs = append(docs, doc)
	}

	// perform each attack
	for _, doc := range docs {
		_ = w.models.Tracks.Create(&track.Track{
			ProjectID:   projectInstance.ID,
			DocumentID:  doc.ID,
			Service:     "base-api/worker/ftp-server",
			Description: fmt.Sprintf("Running the document on `%s` attack.", doc.Instruction),
			Type:        enum.TrackInProgress,
		})

		if err := w.executeDoc(projectInstance, doc); err != nil {
			log.Println(fmt.Errorf("[worker.execute] failed to create request error=%w", err))
		}
	}

	w.exit(id)
}

// executeDoc will make a call to ftp server to run a script
func (w worker) executeDoc(project *project.Project, doc *document.Document) error {
	start := time.Now()

	// update doc status
	doc.Status = enum.StatusPending
	doc.Result = enum.ResultUnknown
	_ = w.models.Documents.Update(doc)

	// create params for request
	params := map[string]string{
		"host": project.Host,
	}

	// create ftp request
	tmp := executeRequest{
		Params:     []string{},
		Path:       doc.Instruction,
		DocumentID: doc.ID,
	}

	// append params
	for key := range params {
		tmp.Params = append(tmp.Params, fmt.Sprintf("--%s", key), params[key])
	}

	// send ftp request
	var buffer bytes.Buffer
	if e := json.NewEncoder(&buffer).Encode(tmp); e != nil {
		return e
	}

	address := fmt.Sprintf("%s/execute", w.cfg.Host)
	headers := []string{
		"Content-Type:application/json",
		fmt.Sprintf("x-token:%s", crypto.GetMD5Hash(w.cfg.Secret)),
	}

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   project.ID,
		DocumentID:  doc.ID,
		Service:     "base-api/worker/ftp-server",
		Description: fmt.Sprintf("Sending post request for `%s` attack.", doc.Instruction),
		Type:        enum.TrackInProgress,
	})

	// update document based of response
	if response, httpError := w.client.Post(address, &buffer, headers...); httpError != nil {
		log.Println(fmt.Errorf("[worker.executeDoc] failed to execute script error=%w", httpError))

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

	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   project.ID,
		DocumentID:  doc.ID,
		Service:     "base-api/worker/ftp-server",
		Description: fmt.Sprintf("Got response for `%s` attack.", doc.Instruction),
		Type:        enum.TrackSuccess,
	})

	_ = w.models.Documents.Update(doc)

	return nil
}

// exit the current task
func (w worker) exit(id int) {
	_ = w.models.Tracks.Create(&track.Track{
		ProjectID:   uint(id),
		Service:     "base-api/worker",
		Description: fmt.Sprintf("Worker exit for %d.", id),
		Type:        enum.TrackWarning,
	})

	w.done <- id
}
