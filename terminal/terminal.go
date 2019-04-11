package terminal

import (
	"fmt"
	"os"

	"github.com/wonderivan/logger"

	"github.com/fanux/fist/tools"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	labelsMap                map[string]string
)

func (t *Terminal) buildLabelsMap() {
	//label
	labelsMap = map[string]string{
		"TerminalID": t.TerminalID,
		"TerminalNS": t.Namespace,
	}
	if t.CookieUserName != "" {
		labelsMap["TerminalUN"] = t.CookieUserName
	}
}

func (q *ListQuery) fetchLabelsMap() map[string]string {
	//append map
	selectMap := map[string]string{}
	if q.CookieUserName != "" {
		selectMap["TerminalUN"] = q.CookieUserName
	}
	if q.TerminalID != "" {
		selectMap["TerminalID"] = q.TerminalID
	}
	if q.Namespace != "" {
		selectMap["TerminalNS"] = q.Namespace
	}
	return selectMap
}

//Terminal is return obj
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
	//
	CookieUserName string `json:"cookieUserName,omitempty"`
}

//ListQuery is query param
type ListQuery struct {
	CookieUserName string `json:"cookieUserName,omitempty"`
	TerminalID     string `json:"terminalID,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
}

func newTerminal() *Terminal {
	return &Terminal{
		Namespace:    "default",
		WithoutToken: false,
		Apiserver:    DefaultApiserver,
		TTYKubeImage: DefaultKubeTTYimage,
	}
}
func newListQuery() *ListQuery {
	return &ListQuery{}
}

//Query a terminal
func (q *ListQuery) Query() ([]*Terminal, error) {
	//append map
	selectMap := q.fetchLabelsMap()
	clientset := tools.GetK8sClient()
	//get deploy deployClient
	deployClient := clientset.AppsV1().Deployments(DefaultTTYnameapace)
	logger.Info(tools.MapToString(selectMap))
	deployList, err := deployClient.List(metav1.ListOptions{
		LabelSelector: tools.MapToString(selectMap),
	})
	if err != nil {
		return nil, err
	}
	terminalArr := make([]*Terminal, len(deployList.Items))
	for index, deploy := range deployList.Items {
		container := deploy.Spec.Template.Spec.Containers[0]
		terminal := &Terminal{
			Apiserver:      container.Env[0].Value,
			UserToken:      container.Env[1].Value,
			Namespace:      container.Env[2].Value,
			User:           container.Env[3].Value,
			TerminalID:     container.Env[4].Value,
			TTYKubeImage:   container.Image,
			CookieUserName: deploy.Spec.Template.Labels["TerminalUN"],
		}
		service, err := clientset.CoreV1().Services(DefaultTTYnameapace).Get(terminal.TerminalID, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		terminal.EndPoint = fmt.Sprintf("%d", service.Spec.Ports[0].NodePort)
		terminalArr[index] = terminal
	}
	return terminalArr, nil
}

//Create a terminal
func (t *Terminal) Create() error {
	t.TerminalID = DefaultPrefix + tools.NewUUID()
	t.buildLabelsMap()
	//create tty deployment and service
	return createTTYcontainer(t)
}

//createTTYdeploy create deployment
func createTTYdeploy(t *Terminal) error {
	clientset := tools.GetK8sClient()
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
		Name:   t.TerminalID,
		Labels: labelsMap,
	}
	selector = &metav1.LabelSelector{
		MatchLabels: labelsMap,
	}
	templateObjMeta = metav1.ObjectMeta{
		Labels: labelsMap,
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

//createTTYservice tty service
func createTTYservice(t *Terminal) error {
	clientset := tools.GetK8sClient()
	service, err := clientset.CoreV1().Services(DefaultTTYnameapace).Create(&apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   t.TerminalID,
			Labels: labelsMap,
		},
		Spec: apiv1.ServiceSpec{
			Selector: labelsMap,
			Type:     "NodePort",
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
func withoutToken(t *Terminal) error {
	if t.WithoutToken {
		clientset := tools.GetK8sClient()
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
			if token == "" {
				return tools.ErrServiceAccountEmpty
			}
			t.UserToken = token
		} else {
			return tools.ErrServiceAccountNotExists
		}
	}
	return nil
}

//createTTYcontainer is
func createTTYcontainer(t *Terminal) error {
	//process without token
	err := withoutToken(t)
	if err != nil {
		return err
	}
	//create namespace
	err = tools.CreateNamespace(DefaultTTYnameapace)
	if err != nil {
		return err
	}
	//create deploy
	err = createTTYdeploy(t)
	if err != nil {
		return err
	}
	//create service
	err = createTTYservice(t)
	if err != nil {
		return err
	}
	return nil
}
