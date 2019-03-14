package tools

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sync"
)

var singleK8sClientInstance *kubernetes.Clientset
var once sync.Once

// install function
func newK8sClient() (*kubernetes.Clientset, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = rest.InClusterConfig()
	if err != nil {
		var kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	// creates the clientSet
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func instanceSingleK8sClient() *kubernetes.Clientset {
	once.Do(func() {
		var err error
		singleK8sClientInstance, err = newK8sClient()
		if err != nil {
			fmt.Printf("kubernetes client init faild. the error info %v", err)
		}
	})
	return singleK8sClientInstance
}

//GetK8sClient get a kubernetes in cluster clientset
func GetK8sClient() *kubernetes.Clientset {
	client := instanceSingleK8sClient()
	if client != nil {
		return client
	} else {
		fmt.Println(ErrK8sClientInitFailed.Error())
		return nil
	}
}
