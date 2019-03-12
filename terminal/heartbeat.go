package terminal

import (
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Heartbeater is
type Heartbeater interface {
	//terminal deployment is in sealyun-tty namespace
	UpdateTimestamp(clientSet *kubernetes.Clientset) error
	//need delete deployment and service in sealyun-tty if it timeout
	CleanTerminalJob(clientSet *kubernetes.Clientset) error, chan bool
}

type terHeartbeater struct {
	terminalID string
	namespace  string
}

//NewHeartbeater is
func NewHeartbeater(tid string, namespace string) Heartbeater {
	var hbInterface Heartbeater
	hbInterface = &terHeartbeater{namespace: namespace, terminalID: tid}
	return hbInterface
}

func (t *terHeartbeater) CleanTerminalJob(clientSet *kubernetes.Clientset, stopped chan bool) error {
	//get deploy of terminalId
	deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
	deploymentResult, err := deploymentsClient.Get(t.terminalID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	//get svc of terminalId
	serviceClient := clientSet.CoreV1().Services(t.namespace)
	_, err = serviceClient.Get(t.terminalID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	lastHeartbeat, _ := strconv.Atoi(deploymentResult.Annotations["fist.seayun.com/lastHeartbeat"])
	heartbeatTime := time.Now().Unix() - int64(lastHeartbeat)
	if heartbeatTime > 600 { // time.Minute() * 10 = 600		
		deletePolicy := metav1.DeletePropagationForeground
		//delete deploy
		if err := deploymentsClient.Delete(t.terminalID, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			return err
		}
		//delete svc
		if err := serviceClient.Delete(t.terminalID, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			return err
		}
		stopped <- true 
	}


	return nil
}

func (t *terHeartbeater) UpdateTimestamp(clientSet *kubernetes.Clientset) error {
	deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
	result, err := deploymentsClient.Get(t.terminalID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	result.Annotations["fist.seayun.com/lastHeartbeat"] = timestamp
	_, err = deploymentsClient.Update(result)
	if err != nil {
		return err
	}
	return nil
}
