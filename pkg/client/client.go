package client

// import kubernetes and add kubernetes client

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	// Create a Kubernetes client configuration
	var config *rest.Config
	var err error
	if home := homedir.HomeDir(); home != "" {
		kubeconfig := fmt.Sprintf("%s/.kube/config", home)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create config from kubeconfig: %w", err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
		}
	}

	if config == nil {
		return nil, fmt.Errorf("failed to create Kubernetes config")
	}

	// Create a new Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	return clientset, nil
}
