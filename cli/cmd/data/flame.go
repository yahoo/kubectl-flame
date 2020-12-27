package data

import "k8s.io/cli-runtime/pkg/genericclioptions"

type FlameConfig struct {
	TargetConfig *TargetDetails
	JobConfig    *JobDetails
	ConfigFlags  *genericclioptions.ConfigFlags
}
