package tools

import (
	"flag"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

//GetK8sClient get a kubernetes in cluster clientset
func GetK8sClient() (*kubernetes.Clientset, error) {
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

//CreateNamespace create namespace
func CreateNamespace(client *kubernetes.Clientset, namespace string) error {
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

//GetSecrets
func GetSecrets(namespace string, name string, clientset *kubernetes.Clientset) (*v1.Secret, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return secret, err
}
