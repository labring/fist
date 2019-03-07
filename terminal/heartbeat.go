package main

import (
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Heartbeater is
type Heartbeater interface {
	//terminal deployment is in sealyun-tty namespace
	UpdateTimestamp(clientSet *kubernetes.Clientset)
	//need delete deployment and service in sealyun-tty if it timeout
	CleanTerminalJob(clientSet *kubernetes.Clientset)
}

type terHeartbeater struct {
	terminalID string
	namespace  string
}

func (t *terHeartbeater) CleanTerminalJob(clientSet *kubernetes.Clientset) {
	deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
	result, err := deploymentsClient.Get("deploy-"+t.terminalID, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	lastHeartbeat, _ := strconv.Atoi(result.Annotations["fist.seayun.com/lastHeartbeat"])
	heartbeatTime := time.Now().Unix() - int64(lastHeartbeat)
	if heartbeatTime > 600 { // time.Minute() * 10 = 600
		deletePolicy := metav1.DeletePropagationForeground
		if err := deploymentsClient.Delete("deploy-"+t.terminalID, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
	}
}

func (t *terHeartbeater) UpdateTimestamp(clientSet *kubernetes.Clientset) error {
	deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
	result, err := deploymentsClient.Get("deploy-"+t.terminalID, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	result.Annotations["fist.seayun.com/lastHeartbeat"] = timestamp
	_, updateErr := deploymentsClient.Update(result)
	return updateErr

}
