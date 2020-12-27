//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/handler"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/kubernetes"
	v1 "k8s.io/api/core/v1"
)

func Flame(cfg *data.FlameConfig) {
	ns, err := kubernetes.Connect(cfg.ConfigFlags)
	if err != nil {
		log.Fatalf("Failed connecting to kubernetes cluster: %v\n", err)
	}

	p := NewPrinter(cfg.TargetConfig.DryRun)

	cfg.TargetConfig.Namespace = ns
	ctx := context.Background()

	p.Print("Verifying target pod ... ")
	pod, err := kubernetes.GetPodDetails(cfg.TargetConfig.PodName, cfg.TargetConfig.Namespace, ctx)
	if err != nil {
		p.PrintError()
		log.Fatalf(err.Error())
	}

	containerName, err := validatePod(pod, cfg.TargetConfig)
	if err != nil {
		p.PrintError()
		log.Fatalf(err.Error())
	}

	containerId, err := kubernetes.GetContainerId(containerName, pod)
	if err != nil {
		p.PrintError()
		log.Fatalf(err.Error())
	}

	p.PrintSuccess()

	cfg.TargetConfig.ContainerName = containerName
	cfg.TargetConfig.ContainerId = containerId

	p.Print("Launching profiler ... ")
	profileId, job, err := kubernetes.LaunchFlameJob(pod, cfg, ctx)
	if err != nil {
		p.PrintError()
		log.Fatalf(err.Error())
	}

	if cfg.TargetConfig.DryRun {
		return
	}

	cfg.TargetConfig.Id = profileId
	profilerPod, err := kubernetes.WaitForPodStart(cfg.TargetConfig, ctx)
	if err != nil {
		p.PrintError()
		log.Fatalf(err.Error())
	}

	p.PrintSuccess()
	apiHandler := &handler.ApiEventsHandler{
		Job:    job,
		Target: cfg.TargetConfig,
	}
	done, err := kubernetes.GetLogsFromPod(profilerPod, apiHandler, ctx)
	if err != nil {
		p.PrintError()
		fmt.Println(err.Error())
	}

	<-done
}

func validatePod(pod *v1.Pod, targetDetails *data.TargetDetails) (string, error) {
	if pod == nil {
		return "", errors.New(fmt.Sprintf("Could not find pod %s in Namespace %s",
			targetDetails.PodName, targetDetails.Namespace))
	}

	if len(pod.Spec.Containers) != 1 {
		var containerNames []string
		for _, container := range pod.Spec.Containers {
			if container.Name == targetDetails.ContainerName {
				return container.Name, nil // Found given container
			}

			containerNames = append(containerNames, container.Name)
		}

		return "", errors.New(fmt.Sprintf("Could not determine container. please specify one of %v", containerNames))
	}

	return pod.Spec.Containers[0].Name, nil
}
