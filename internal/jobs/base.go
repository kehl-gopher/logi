package jobs

import (
	"context"
)

type JobHandler struct {
	Type    string
	Payload interface{}
}

func NewJob(typ string, payload interface{}) *JobHandler {
	return &JobHandler{Type: typ, Payload: payload}
}

func StartWorker(ctx context.Context) error {
	return nil
}

func StopWorker() {

}

func RegisterHandler(jobType string, handler JobHandler) {

}
