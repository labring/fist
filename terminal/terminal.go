package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

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
	kubeTTYimage     = "fanux/kube-ttyd:latest"
)

//Terminal is
type Terminal struct {
	User         string
	Apiserver    string // just using default apiserver
	UserToken    string
	Namespace    string
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

	client := clientset.AppsV1().Deployments(TTYnameapace)
	result, err := client.Create(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: t.TerminalID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
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
			Name: t.TerminalID,
		},
		Spec: apiv1.ServiceSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"TerminalID": t.TerminalID,
				},
			},
			Type: "NodePort",
			Ports: []apiv1.ServicePort{
				{Name: "tty", Port: 8080, TargetPort: 8080, Protocol: apiv1.Protocol{"TCP"}},
			},
		},
	})
	if err != nil {
		return err
	}
	t.EndPoint = fmt.Sprintf("%d", service.Spec.Ports[0].NodePort)
}

//LoadTerminalID is
func LoadTerminalID() error {
	//TODO
}
