package template

import (
	"io/ioutil"
	"log"
)

//Templates is
var Templates map[string]string

//LoadTemplates is
func LoadTemplates(dir string) error {
	if Templates == nil {
		Templates = make(map[string]string)
	}
	if dir == "" {
		dir = "/etc/fist/templates"
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
