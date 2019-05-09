package template

import (
	"io/ioutil"
	"log"
)

//Templates is
var (
	Templates          map[string]string
	defaultTemplateDir = "/etc/fist/templates"
)

func loadDefault() {
	Templates["Deployment"] = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      name: {{ .Name }}
  template:
    metadata:
      labels:
        name: {{ .Name }}
    spec:
      containers:
      - name: {{ .Name }}
        image: {{ .Image }}
        command: {{ .Command }}
        imagePullPolicy: {{ .ImagePolicy }}
        ports:
        - containerPort: {{.Port}}
	`
}

//LoadTemplates is
func LoadTemplates(dir string) error {
	if Templates == nil {
		Templates = make(map[string]string)
	}
	loadDefault()
	if dir == "" {
		dir = defaultTemplateDir
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		t := readFile(file.Name)
		if t == "" {
			continue
		}
		Templates[file.Name()] = readFile(file.Name)
	}
}

func readFile(name string) string {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)
}
