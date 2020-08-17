//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/VerizonMedia/kubectl-flame/api"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/kubernetes"
	batchv1 "k8s.io/api/batch/v1"
)

type ApiEventsHandler struct {
	Job    *batchv1.Job
	Target *data.TargetDetails
}

func (h *ApiEventsHandler) Handle(events chan string, done chan bool, ctx context.Context) {
	for eventString := range events {
		event, err := api.ParseEvent(eventString)
		if err != nil {
			fmt.Printf("Got invalid event: %s\n", err)
		} else {
			switch data := event.(type) {
			case *api.ErrorData:
				fmt.Printf("Error: %s\n", data.Reason)
			case *api.FlameGraphData:
				h.createFlameGraph(data)
			case *api.ProgressData:
				h.reportProgress(data, done, ctx)
			default:
				fmt.Printf("Unrecognized event type: %T!\n", data)
			}
		}
	}
}

func (h *ApiEventsHandler) createFlameGraph(data *api.FlameGraphData) {
	decodedData, err := base64.StdEncoding.DecodeString(data.EncodedFile)
	if err != nil {
		fmt.Printf("Failed to decode flamegraph: %v\n", err)
		return
	}

	err = ioutil.WriteFile(h.Target.FileName, decodedData, 0777)
	if err != nil {
		fmt.Printf("Failed to write flamegraph file: %v\n", err)
	}
}

func (h *ApiEventsHandler) reportProgress(data *api.ProgressData, done chan bool, ctx context.Context) {
	if data.Stage == api.Started {
		fmt.Printf("Profiling ... ")
	} else if data.Stage == api.Ended {
		_ = kubernetes.DeleteProfilingJob(h.Job, h.Target, ctx)
		fmt.Printf("✔\nFlameGraph saved to: %s 🔥\n", h.Target.FileName)
		done <- true
	}
}
