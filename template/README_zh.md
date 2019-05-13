# 模板使用教程
模板模块可以定义动态的API，且新增一些模板时无需修改代码，仅需要协商好模板格式与value格式即可.

## 应用场景
最重要的应用场景就是帮助用户渲染各种yaml配置，Deployment Configmap Service等等。。。

所以代码设计时肯定不希望增加一个模板时就需要修改代码和API，所以动态API很重要。

还可用于其它一切模板渲染的场景，如渲染kubeconfig文件，namespace配额等等

## 工作原理
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

# 安装
> 使用kubernetes

```
kubectl create -f deploy
```

> 使用docker

```
docker run -d -p 8080:8080 lameleg/fist:latest ./fist template
```

## 渲染模板
> 渲染请求

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
这里的`Deployment` 是内建模板, 用户可以自定义模板

type=text指定渲染结果格式，不指定会把渲染结果字符串放在数组中。

渲染结果：
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

## 自定义模板
> 创建template文件

默认目录：`/etc/fist/templates` (进入容器内部，或者挂载此文件夹)

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

所以可以创建很多模板，value里面的值也是随意调整的，但是一定要与模板对应上。

如此就能解决大家写yaml不方便的问题了

# 获取templates列表
GET /fist/templates
```
[
    {
	    "name": {
            "Key":"fist-deploy.yaml.tmpl",  
            "FormName": "Deployment",                # 用于前端表单渲染
            "Describe": "用户渲染deployment"
        }
        "value": [                       # 这里可用于动态渲染前端模板, value有层级嵌套动态表单暂不太好处理
            {
                "Key": "Name",
                "DefaultValue": "fist",
                "FormName": "应用名称",
                "Describe": "deployment名称"
            },
        ]
    },
]
```
此内容可直接存在文件(/etc/fist/templates/templates.metadata.json)中作为templates的元数据
