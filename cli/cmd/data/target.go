//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package data

import (
	"time"

	"github.com/VerizonMedia/kubectl-flame/api"
)

type TargetDetails struct {
	Namespace          string
	PodName            string
	ContainerName      string
	ContainerId        string
	Event              api.ProfilingEvent
	Duration           time.Duration
	Id                 string
	FileName           string
	Alpine             bool
	DryRun             bool
	Image              string
	DockerPath         string
	Language           api.ProgrammingLanguage
	Pgrep              string
	ImagePullSecret    string
	ServiceAccountName string
}
