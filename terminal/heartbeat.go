package main

import (
	"k8s.io/client-go/kubernetes"
         metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "fmt"
        "time"
        "strconv"
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
  namespace string
}


func (t *terHeartbeater) CleanTerminalJob(clientSet *kubernetes.Clientset) {
    deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
    deletePolicy := metav1.DeletePropagationForeground
    if err := deploymentsClient.Delete(t.terminalID, &metav1.DeleteOptions{
        PropagationPolicy: &deletePolicy,
    }); err != nil {
        panic(err)
    }
}

func (t *terHeartbeater) UpdateTimestamp(clientSet *kubernetes.Clientset) {
    deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
    result, err := deploymentsClient.Get(t.terminalID, metav1.GetOptions{})
    if err != nil {
       panic(err)
    }
    timestamp := strconv.FormatInt( time.Now().Unix(), 10)
    result.Annotations["fist.seayun.com/lastHeartbeat"] = timestamp
    fmt.Println(result.Annotations["fist.seayun.com/lastHeartbeat"])
}
