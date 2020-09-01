//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package data

import (
	"github.com/VerizonMedia/kubectl-flame/api"
	"time"
)

type TargetDetails struct {
	Namespace     string
	PodName       string
	ContainerName string
	ContainerId   string
	Duration      time.Duration
	Id            string
	FileName      string
	Alpine        bool
	DryRun        bool
	Image         string
	Language      api.ProgrammingLanguage
	Pgrep         string
}
