package terminal

import (
	"errors"
	"fmt"
	"os"

	"github.com/fanux/fist/tools"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

//consts
const (
	DefaultTTYnameapace = "sealyun-tty"
	DefaultApiserver    = "https://kubernetes.default.svc.cluster.local:443" //or https://10.96.0.1:443
	DefaultKubeTTYimage = "fanux/fist-tty-tools:v1.0.0"
	DefaultPrefix       = "tid-"
	ClassPathNamespace  = "MY_NAMESPACE"
	ServiceAccountName  = "MY_SA_NAME"

	DefaultNamespace      = "sealyun"
	DefaultServiceAccount = "admin"
)

//vars
var (
	DefaultTTYDeployReplicas = int32(1)
)

//Terminal is
type Terminal struct {
	//input field
	User         string `json:"user,omitempty"`
	UserToken    string `json:"userToken,omitempty"`
	Apiserver    string `json:"apiServer,omitempty"`    // just using default apiserver
	Namespace    string `json:"namespace,omitempty"`    // the kubeconfig default context namespace
	WithoutToken bool   `json:"withoutToken,omitempty"` // if true, mount the kubeconfig file, using ttyd instead the start-terminal.sh
	TTYKubeImage string `json:"ttyKubeImage,omitempty"` //default is  "fanux/fist-tty-tools:v1.0.0"

	//output append field
	TerminalID string `json:"terminalID,omitempty"`
	EndPoint   string `json:"endPoint,omitempty"`
}

func newTerminal() *Terminal {
	return &Terminal{
		Namespace:    "default",
		WithoutToken: false,
		Apiserver:    DefaultApiserver,
		TTYKubeImage: DefaultKubeTTYimage,
	}
}

//Create a terminal
func (t *Terminal) Create() error {
	t.TerminalID = DefaultPrefix + tools.NewUUID()

	//create tty deployment and service
	return CreateTTYcontainer(t)
}

//CreateTTYdeploy create deployment
func CreateTTYdeploy(t *Terminal, clientset *kubernetes.Clientset) error {
	//get deploy deployClient
	deployClient := clientset.AppsV1().Deployments(DefaultTTYnameapace)
	//vars
	var (
		objMeta         metav1.ObjectMeta
		selector        *metav1.LabelSelector
		templateObjMeta metav1.ObjectMeta
		ports           []apiv1.ContainerPort
		env             []apiv1.EnvVar
	)
	//init
	objMeta = metav1.ObjectMeta{
		Name: t.TerminalID,
	}
	selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"TerminalID": t.TerminalID,
		},
	}
	templateObjMeta = metav1.ObjectMeta{
		Labels: map[string]string{
			"TerminalID": t.TerminalID,
		},
	}
	ports = []apiv1.ContainerPort{
		{
			Name:          "http",
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: 8080,
		},
	}
	env = []apiv1.EnvVar{
		{Name: "APISERVER", Value: t.Apiserver},
		{Name: "USER_TOKEN", Value: t.UserToken},
		{Name: "NAMESPACE", Value: t.Namespace},
		{Name: "USER_NAME", Value: t.User},
		{Name: "TERMINAL_ID", Value: t.TerminalID},
	}
	_, err := deployClient.Create(&appsv1.Deployment{
		ObjectMeta: objMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: &DefaultTTYDeployReplicas,
			Selector: selector,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: templateObjMeta,
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Env:   env,
							Name:  "tty",
							Image: t.TTYKubeImage,
							Ports: ports,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

//CreateTTYservice tty service
func CreateTTYservice(t *Terminal, clientset *kubernetes.Clientset) error {
	service, err := clientset.CoreV1().Services(DefaultTTYnameapace).Create(&apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: t.TerminalID,
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

//WithoutToken create terminal without token
func withoutToken(t *Terminal, clientset *kubernetes.Clientset) error {
	if t.WithoutToken {
		//get namespace
		myNamespace := os.Getenv(ClassPathNamespace)
		mySaName := os.Getenv(ServiceAccountName)
		if myNamespace == "" {
			myNamespace = DefaultNamespace
		}
		if mySaName == "" {
			mySaName = DefaultServiceAccount
		}

		t.User = mySaName
		sa, err := clientset.CoreV1().ServiceAccounts(myNamespace).Get(mySaName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		saSercerts := sa.Secrets
		if saSercerts != nil && len(saSercerts) > 0 {
			saTokenName := saSercerts[0].Name
			saTokenSecrets, err := clientset.CoreV1().Secrets(myNamespace).Get(saTokenName, metav1.GetOptions{})
			if err != nil {
				return err
			}
			token := string(saTokenSecrets.Data["token"])
			if err != nil {
				return errors.New("the serviceAccount token is empty")
			}
			t.UserToken = token
		} else {
			return errors.New("the serviceAccount token is not exists")
		}
	}
	return nil
}

//CreateTTYcontainer is
func CreateTTYcontainer(t *Terminal) error {
	//get client of k8s
	clientset, err := tools.GetK8sClient()
	if err != nil {
		return err
	}
	//process without token
	err = withoutToken(t, clientset)
	if err != nil {
		return err
	}
	//create namespace
	err = tools.CreateNamespace(clientset, DefaultTTYnameapace)
	if err != nil {
		return err
	}
	//create deploy
	err = CreateTTYdeploy(t, clientset)
	if err != nil {
		return err
	}
	//create service
	err = CreateTTYservice(t, clientset)
	if err != nil {
		return err
	}
	return nil
}

//LoadTerminalID is
func LoadTerminalID() error {
	//TODO
	return nil
}
