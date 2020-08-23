# kubectl flame :fire:

A kubectl plugin that allows you to profile production applications with low-overhead by generating
[FlameGraphs](http://www.brendangregg.com/flamegraphs.html)

Running `kubectlf-flame` does **not** require any modification to existing pods.
## Table of Contents

- [Requirements](#requirements)
- [Usage](#usage)
- [Installing](#installing)
- [How it works](#how-it-works)
- [Contribute](#contribute)
- [Maintainers](#maintainers)
- [License](#license)

## Requirements
* Currently, only Java applications are supported. (Golang support coming soon!)
* Kubernetes cluster that use Docker as the container runtime (tested on GKE, EKS and AKS)

## Usage
### Profiling Kubernetes Pod
In order to profile pod `mypod` for 1 minute and save the flamegraph as `/tmp/flamegraph.svg` run:
```shell
kubectl flame mypod -t 1m -f /tmp/flamegraph.svg
```
### Profiling Alpine based container
Profiling alpine based containers require using `--alpine` flag:
```shell
kubectl flame mypod -t 1m -f /tmp/flamegraph.svg --alpine
```
### Profiling sidecar container
Pods that contains more than one container require specifying the target container as an argument:
```shell
kubectl flame mypod -t 1m -f /tmp/flamegraph.svg mycontainer
```

## Installing

### Krew

You can install `kubectl flame` using the [Krew](https://github.com/kubernetes-sigs/krew), the package manager for kubectl plugins.

Once you have [Krew installed](https://krew.sigs.k8s.io/docs/user-guide/setup/install/) just run:

```bash
kubectl krew install flame
```

### Pre-built binaries
See the release page for the full list of pre-built assets.

## How it works
`kubectl-flame` launch a Kubernetes Job on the same node as the target pod.
Under the hood `kubectl-flame` use [async-profiler](https://github.com/jvm-profiling-tools/async-profiler) in order to generate flame graphs for Java applications. 
Interaction with the target JVM is done via a shared `/tmp` folder.
Other languages support (such as the upcoming Golang support) will be based on [ebpf profiling](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter).

## Contribute
Please refer to [the contributing.md file](Contributing.md) for information about how to get involved. We welcome issues, questions, and pull requests.

## Maintainers
- Eden Federman: efederman@verizonmedia.com

## License
This project is licensed under the terms of the [Apache 2.0](LICENSE-Apache-2.0) open source license. Please refer to [LICENSE](LICENSE) for the full terms.
