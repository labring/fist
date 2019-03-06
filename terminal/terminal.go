package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//consts
const (
	TTYnameapace     = "sealyun-tty"
	DefaultApiserver = "https://kubernetes.default.svc.cluster.local:443" //or https://10.96.0.1:443
	kubeTTYimage     = "fanux/fist-tty-tools:v1.0.0"
)

//Terminal is
type Terminal struct {
	User         string
	Apiserver    string // just using default apiserver
	UserToken    string
	Namespace    string // the kubeconfig default context namespace
	TerminalID   string
	EndPoint     string
	WithoutToken bool // if true, mount the kubeconfig file, using ttyd instead the start-terminal.sh
}

func newUUID() string {
	u := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, u); err != nil {
		panic(err)
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u)
}

//Create a terminal
func (t *Terminal) Create() error {
	t.Apiserver = DefaultApiserver
	t.TerminalID = newUUID()

	//create tty deployment and service
	return CreateTTYcontainer(t)
}

//CreateTTYcontainer is
func CreateTTYcontainer(t *Terminal) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	var re int32
	re = 1
	_, err = clientset.CoreV1().Namespaces().Get(TTYnameapace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(&apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: TTYnameapace,
			},
		})
		if err != nil {
			return err
		}
	}

	client := clientset.AppsV1().Deployments(TTYnameapace)
	_, err = client.Create(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deploy-" + t.TerminalID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &re,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"TerminalID": t.TerminalID,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"TerminalID": t.TerminalID,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Env: []apiv1.EnvVar{
								{Name: "APISERVER", Value: DefaultApiserver},
								{Name: "USER_TOKEN", Value: t.UserToken},
								{Name: "NAMESPACE", Value: "default"},
								{Name: "USER_NAME", Value: t.User},
								{Name: "TERMINAL_ID", Value: t.TerminalID},
							},
							Name:  "tty",
							Image: kubeTTYimage,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	service, err := clientset.CoreV1().Services(TTYnameapace).Create(&apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc-" + t.TerminalID,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"TerminalID": t.TerminalID,
			},
			Type: "NodePort",
			Ports: []apiv1.ServicePort{
				{Name: "tty", Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: apiv1.Protocol("TCP")},
			},
		},
	})
	if err != nil {
		return err
	}
	t.EndPoint = fmt.Sprintf("%d", service.Spec.Ports[0].NodePort)
	return nil
}

//LoadTerminalID is
func LoadTerminalID() error {
	//TODO
	return nil
}
