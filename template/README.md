# Render kubernetes yaml 
```
  client          server         template file
    |  value json   |                |
    |==============>|  read templates| from disk
    |               |<===============|
    |               |===+            |
    |               |   | render     |
    |  real yaml    |<==+            |
    |<==============|                |
    | apply         |                |
    |=============================================> kubernetes
    |               |                |
```

# Install
> In kubernetes

```
kubectl create -f deploy
```

> Or ship on docker

```
docker run -d -p 8080:8080 -v /etc/fist/templates:/etc/fist/templates lameleg/fist:latest ./fist template
```

# Render a default Deployment
```shell
curl http://localhost:8080/templates?type=text -H "Content-Type:application/json" -d '[
{
	"name":"Deployment",  
	"value": {                       
		"Name":"fist",
		"Image":"sealyun/fist",
		"Replicas":3,
		"Namespace":"sealyun",
		"Command": "['./fist', 'serve']",
		"ImagePolicy":"IfnotPresent",
		"Port":9090}
}
]'
```

If not specified type=text, it will response a string array.

Result:

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fist
  namespace: sealyun
spec:
  replicas: 3
  selector:
    matchLabels:
      name: fist
  template:
    metadata:
      labels:
        name: fist
    spec:
      containers:
      - name: fist
        image: sealyun/fist
        command: [./fist, serve]
        imagePullPolicy: IfnotPresent
        ports:
        - containerPort: 9090
```

# Custom your templates
> Create templates file, default dir is `/etc/fist/templates` (in your container)

```
# cat /etc/fist/templates/Hello-world 
This is my hello world template
{{ .Name }}
{{ .Value }}
```

```shell
curl http://localhost:8080/templates?type=text -H "Content-Type:application/json" -d '[
{
	"name":"Hello-world",  
	"value": {                       
		"Name": "Sealyun",
        "Value": "Hello everyone!"}
}
]'
```

Result:

```
---
This is my hello world template
SealYun
Hello everyone
```

# API
```
[
    {
        "name":"template-name",
        "value": {} # its a interface
    },
]
```

Then we can select a template in fist and render it.

Like if we want to render a deployment:

We named a template called "deploy"
```
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
```

```
{
   "name":"deploy",
   "value": {
        "Name": "sealyun",
        "Replicas":3,
        ...
    }
}
```
It will return the rendered base64 encoded yaml

# Use cases
This will be useful for kubeconfig render, namespace quota render, yaml file render ...

Render any file your want :)
