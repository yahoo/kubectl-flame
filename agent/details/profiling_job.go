//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package details

import (
	"github.com/VerizonMedia/kubectl-flame/api"
	"time"
)

type ProfilingJob struct {
	Duration          time.Duration
	ID                string
	ContainerID       string
	ContainerName     string
	PodUID            string
	Language          api.ProgrammingLanguage
	TargetProcessName string
}
