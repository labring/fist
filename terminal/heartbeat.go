package terminal

import (
	"strconv"
	"time"

	"github.com/fanux/fist/tools"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Heartbeater is
type Heartbeater interface {
	//terminal deployment is in sealyun-tty namespace
	UpdateTimestamp() error
	//need delete deployment and service in sealyun-tty if it timeout
	CleanTerminalJob() error
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

func (t *terHeartbeater) CleanTerminalJob() error {
	clientSet := tools.GetK8sClient()
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
	}
	return nil
}

func (t *terHeartbeater) UpdateTimestamp() error {
	clientSet := tools.GetK8sClient()
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
