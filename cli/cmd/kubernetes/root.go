//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package kubernetes

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var clientSet *kubernetes.Clientset

func Connect(clientGetter genericclioptions.RESTClientGetter) (string, error) {
	restConfig, err := clientGetter.ToRESTConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", err
	}

	clientSet = clientset
	ns, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", err
	}

	return ns, nil

	// homeDirectory, err := homedir.Dir()
	// if err != nil {
	// 	return err
	// }
	// kubeconfig := filepath.Join(homeDirectory, ".kube", "config")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	return err
	// }

	// // create the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	return err
	// }

	// clientSet = clientset
	// return nil
}
