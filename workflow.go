package gbdx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Workflow holds the ID and all the associated Acquisitions of requested process.
type Workflow struct {
	ID          string          `json:"identifier"`
	Type        []string        `json:"type"`
}

// WorkflowStatus returns the status of the imagery given the string identifying it.
func (a *Api) WorkflowStatus(wfID string) (*Workflow, error) {

	url := fmt.Sprintf("%s%s", endpoints.workflow, wfID)
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Workflow status returned a bad status code: %s", resp.Status)
	}

	var workflow *Workflow
	err = json.NewDecoder(resp.Body).Decode(&workflow)
	if err != nil {
		return nil, fmt.Errorf("Workflow status failed to decode properly: %v", err)
	}
	return workflow, err
}

// WorkflowHeartbeat checks if the workflow endpoint is alive and well.
func WorkflowHeartbeat() error {

	url := endpoints.workflowHeartbeat

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("workflow heartbeat failed %s: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("workflow heartbeat returned a bad status code: %s", resp.Status)
	}
	defer resp.Body.Close()

	return err
}
