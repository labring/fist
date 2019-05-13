# template for [palm](https://github.com/fanux/palm)


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

# Read template from file
Template read templates from file when it start, so user can set it own templates. Also we can build some default template in bin file.

# Use cases
This will be useful for kubeconfig render, namespace quota render, yaml file render ...
