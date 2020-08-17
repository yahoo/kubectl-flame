//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package main

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/profiler"
	"github.com/VerizonMedia/kubectl-flame/api"
)

func main() {
	args, err := validateArgs()
	if err != nil {
		api.PublishError(err)
		os.Exit(1)
	}

	err = api.PublishEvent(api.Progress, &api.ProgressData{Time: time.Now(), Stage: api.Started})
	if err != nil {
		api.PublishError(err)
		os.Exit(1)
	}

	err = profiler.SetUp(args)
	if err != nil {
		api.PublishError(err)
		os.Exit(1)
	}

	done := handleSignals()
	err = profiler.Invoke(args)
	if err != nil {
		api.PublishError(err)
		os.Exit(1)
	}

	err = api.PublishEvent(api.Progress, &api.ProgressData{Time: time.Now(), Stage: api.Ended})
	if err != nil {
		api.PublishError(err)
		os.Exit(1)
	}

	<-done
}

func validateArgs() (*details.ProfilingJob, error) {
	if len(os.Args) != 6 {
		return nil, errors.New("expected 6 arguments")
	}

	duration, err := time.ParseDuration((os.Args[5]))
	if err != nil {
		return nil, err
	}

	currentJob := &details.ProfilingJob{}
	currentJob.ID = os.Args[1]
	currentJob.PodUID = os.Args[2]
	currentJob.ContainerName = os.Args[3]
	currentJob.ContainerID = strings.Replace(os.Args[4], "docker://", "", 1)
	currentJob.Duration = duration

	return currentJob, nil
}

func handleSignals() chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		os.RemoveAll("/tmp/async-profiler")
		os.Remove("/tmp")
		done <- true
	}()

	return done
}
