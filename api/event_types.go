//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package api

import (
	"encoding/json"
	"time"
)

type EventType string
type ProgressStage string

const (
	Error      EventType = "error"
	FlameGraph EventType = "flamegraph"
	Progress   EventType = "progress"

	Started ProgressStage = "started"
	Ended   ProgressStage = "ended"
)

type Event struct {
	Type EventType        `json:"type"`
	Data *json.RawMessage `json:"data"`
}

type ErrorData struct {
	Reason string `json:"reason"`
}

type FlameGraphData struct {
	EncodedFile string `json:"encoded_file"`
}

type ProgressData struct {
	Time  time.Time     `json:"time"`
	Stage ProgressStage `json:"stage"`
}

var typeToData = map[EventType]interface{}{
	Error:      &ErrorData{},
	FlameGraph: &FlameGraphData{},
	Progress:   &ProgressData{}}

func GetDataStructByType(t EventType) interface{} {
	return typeToData[t]
}
