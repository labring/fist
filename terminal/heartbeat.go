package main

import (
	"k8s.io/client-go/kubernetes"
)

//Heartbeater is
type Heartbeater interface {
	UpdateTimestamp(clientSet *kubernetes.Clientset, terminalID string)
	CleanTerminalJob(clientSet *kubernetes.Clientset)
}
