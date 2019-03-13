package tools

import (
	"flag"
	"fmt"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func GetK8sClient() (*kubernetes.Clientset, error) {
	client := instanceSingleK8sClient()
	if client != nil {
		return client, nil
	} else {
		return nil, ErrK8sClientInitFailed
	}
}

//CreateNamespace create namespace
func CreateNamespace(namespace string) error {
	client := instanceSingleK8sClient()
	_, err := client.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		_, err = client.CoreV1().Namespaces().Create(&apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//GetSecrets is get sercrets
func GetSecrets(namespace string, name string) (*v1.Secret, error) {
	client := instanceSingleK8sClient()
	secret, err := client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return secret, err
}
