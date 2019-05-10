package template

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"
)

//Value is
type Value struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

//Render out put rendered template
/*
value json string

{
	"name":"test-template",
	"value": {
		"Name":"fist",
		"Image":"sealyun/fist",
		"Replicas":3,
		"Namespace":"sealyun",
		"Command": "["./fist", "serve"]",
		"ImagePolicy":"IfnotPresent",
		"Port":9090
	}
}

template:

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
*/
func Render(value []byte) (result string) {
	v := &Value{}
	err := json.Unmarshal(value, v)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	t, ok := Templates[v.Name]
	if !ok {
		log.Fatal("Template not found")
	}

	return Template(v.Value, t)
}

//RenderValue is
func RenderValue(v Value) (result string) {
	t, ok := Templates[v.Name]
	if !ok {
		log.Fatal("Template not found")
	}
	return Template(v.Value, t)
}

//Template is
func Template(value interface{}, temp string) string {
	t := template.Must(template.New("template").Parse(temp))

	res := new([]byte)
	buf := bytes.NewBuffer(*res)

	err := t.Execute(buf, value)
	if err != nil {
		log.Fatal("Render template failed", temp)
		return ""
	}
	return buf.String()
}

/*

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	// Define a template.
	const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.
{{- else}}
It is a shame you couldn't make it to the wedding.
{{- end}}
{{with .Gift -}}
Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`

	j := []byte(`{"Name":"Aunt Mildred", "Gift":"bone china tea set", "Attended":"true"}`)

	r := new(interface{})

	err := json.Unmarshal(j, r)
	if err != nil {
		fmt.Println(err)
	}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Parse(letter))

	// Execute the template for each recipient.
	err = t.Execute(os.Stdout, r)
	if err != nil {
		log.Println("executing template:", err)
	}
}
*/
