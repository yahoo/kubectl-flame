//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package cmd

import (
	"fmt"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"time"
)

const (
	defaultDuration = 1 * time.Minute
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
`
)

var targetDetails data.TargetDetails

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
			if len(args) == 0 {
				cmd.Help()
				return
			}

			targetDetails.PodName = args[0]
			if len(args) > 1 {
				targetDetails.ContainerName = args[1]
			}

			Flame(&targetDetails, options.configFlags)
		},
	}

	cmd.AddCommand(newVersionCommand(streams))
	cmd.Flags().DurationVarP(&targetDetails.Duration, "time", "t", defaultDuration, "Enter max scan Duration")
	cmd.Flags().StringVarP(&targetDetails.FileName, "file", "f", "flamegraph.svg", "Optional file location")
	cmd.Flags().BoolVar(&targetDetails.Alpine, "alpine", false, "Target image is based on Alpine")
	cmd.Flags().BoolVar(&targetDetails.DryRun, "dry-run", false, "simulate profiling")
	options.configFlags.AddFlags(cmd.Flags())

	return cmd
}
func newVersionCommand(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information for kubectl flame",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Fprintln(streams.Out, version.String())
			return nil
		},
	}
	return cmd
}
