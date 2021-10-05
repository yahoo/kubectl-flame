//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/VerizonMedia/kubectl-flame/api"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	defaultDuration = 1 * time.Minute
	defaultEvent    = string(api.Wall)
	flameLong       = `Profile existing applications with low-overhead by generating flame graphs.

These commands help you identify application performance issues. 
`
	flameExamples = `
	# Profile a pod for 5 minutes and save the output as flame.svg file
	%[1]s flame mypod -f flame.svg -t 5m

	# Profile an alpine based container
	%[1]s flame mypod -f flame.svg --alpine

	# Profile specific container container1 from pod mypod in namespace test
	%[1]s flame mypod -f /tmp/flame.svg -n test container1

	# Set custom resource requests and limits for the kubectl-flame pod (default: neither requests nor limits are set)
	%[1]s flame mypod -f flame.svg -cpu.requests 100m -cpu.limits 200m -mem.requests 100Mi -mem.limits 200Mi
`
)

type FlameOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

func NewFlameOptions(streams genericclioptions.IOStreams) *FlameOptions {
	return &FlameOptions{
		configFlags: genericclioptions.NewConfigFlags(false),
		IOStreams:   streams,
	}
}

func NewFlameCommand(streams genericclioptions.IOStreams) *cobra.Command {
	var (
		targetDetails data.TargetDetails
		jobDetails    data.JobDetails
		showVersion   bool
		chosenLang    string
		chosenEvent   string
	)

	options := NewFlameOptions(streams)
	cmd := &cobra.Command{
		Use:                   "flame [pod-name]",
		DisableFlagsInUseLine: true,
		Short:                 "Profile running applications by generating flame graphs.",
		Long:                  flameLong,
		Example:               fmt.Sprintf(flameExamples, "kubectl"),
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.SetOutput(streams.ErrOut)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Fprintln(streams.Out, version.String())
				return
			}

			if len(args) == 0 {
				cmd.Help()
				return
			}

			if err := validateFlags(chosenLang, chosenEvent, &targetDetails, &jobDetails); err != nil {
				fmt.Fprintln(streams.Out, err)
				os.Exit(1)
			}

			targetDetails.PodName = args[0]
			if len(args) > 1 {
				targetDetails.ContainerName = args[1]
			}

			cfg := &data.FlameConfig{
				TargetConfig: &targetDetails,
				JobConfig:    &jobDetails,
				ConfigFlags:  options.configFlags,
			}

			Flame(cfg)
		},
	}

	cmd.Flags().BoolVar(&showVersion, "version", false, "Print version info")

	cmd.Flags().DurationVarP(&targetDetails.Duration, "time", "t", defaultDuration, "Max scan Duration")
	cmd.Flags().StringVarP(&targetDetails.FileName, "file", "f", "flamegraph.svg", "Optional file location")
	cmd.Flags().BoolVar(&targetDetails.Alpine, "alpine", false, "Target image is based on Alpine")
	cmd.Flags().BoolVar(&targetDetails.DryRun, "dry-run", false, "Simulate profiling")
	cmd.Flags().StringVar(&targetDetails.Image, "image", "", "Manually choose agent docker image")
	cmd.Flags().StringVar(&targetDetails.DockerPath, "docker-path", "/var/lib/docker/", "Use a different Docker install path")
	cmd.Flags().StringVarP(&targetDetails.Pgrep, "pgrep", "p", "", "name of the target process")

	cmd.Flags().StringVarP(&chosenLang, "lang", "l", "", fmt.Sprintf("Programming language of "+
		"the target application, choose one of %v", api.AvailableLanguages()))
	cmd.Flags().StringVarP(&chosenEvent, "event", "e", defaultEvent, fmt.Sprintf("Profiling event, choose one of %v",
		api.AvailableEvents()))

	cmd.Flags().StringVar(&jobDetails.RequestConfig.CPU, "cpu.requests", "", "CPU requests of the started profiling container")
	cmd.Flags().StringVar(&jobDetails.RequestConfig.Memory, "mem.requests", "", "Memory requests of the started profiling container")
	cmd.Flags().StringVar(&jobDetails.LimitConfig.CPU, "cpu.limits", "", "CPU limits of the started profiling container")
	cmd.Flags().StringVar(&jobDetails.LimitConfig.Memory, "mem.limits", "", "Memory limits of the started profiling container")
	cmd.Flags().StringVar(&targetDetails.ImagePullSecret, "imagePullSecret", "", "imagePullSecret for agent docker image")
	cmd.Flags().StringVar(&targetDetails.ServiceAccountName, "serviceAccountName", "", "serviceAccountName to be used for profiling container")

	options.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func validateFlags(langString string, eventString string, targetDetails *data.TargetDetails, jobDetails *data.JobDetails) error {
	if langString == "" {
		return fmt.Errorf("use -l flag to select one of the supported languages %s", api.AvailableLanguages())
	}

	if !api.IsSupportedLanguage(langString) {
		return fmt.Errorf("unsupported language, choose one of %s", api.AvailableLanguages())
	}

	if eventString != "" && !api.IsSupportedEvent(eventString) {
		return fmt.Errorf("unsupported event, choose one of %s", api.AvailableEvents())
	}

	targetDetails.Language = api.ProgrammingLanguage(langString)
	targetDetails.Event = api.ProfilingEvent(eventString)

	if _, err := jobDetails.RequestConfig.ParseResources(); err != nil {
		return fmt.Errorf("unable to parse resource requests: %w", err)
	}

	if _, err := jobDetails.LimitConfig.ParseResources(); err != nil {
		return fmt.Errorf("unable to parse resourse limits: %w", err)
	}

	return nil
}
