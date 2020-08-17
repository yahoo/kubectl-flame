//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package main

import (
	"os"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-flame", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	root := cmd.NewFlameCommand(streams)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
