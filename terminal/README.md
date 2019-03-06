# Kuberntes web terminal

![](./terminal.jpeg)

![](./show-terminal.png)

# Create a terminal
```
clent                     terminal                        terminal-pod
  | {"User", "Token"}        |                                |
  |------------------------->|create a service and deployment |
  |                          |                                |
  |                          |------------------------------->| pod ENVS: kubernetes APIserver IP
  |                          |                                |           User token
  |                          |                                |           Namespace
  |                          |                                |           User name
  | { terminalID,            |                                |           Terminelid
  | EndPoint:"ip:port"}      |      terminal IP:port          | pod Label: LastHeartbeat, time.Now - LastKeepalived > 10, kill it          
  |<-------------------------|<-------------------------------|           
```

# How to stop a terminal container

```
clent                     terminal                        terminal-pod
  | heartbeat (5min)         |                                |
  |------------------------->|update deployment LastHeartbeat |
  |                          |                                |
  |                          |------------------------------->| 
  |                          |                                | 
  |                          |                                | 
  |                          | loop check LastHeartbeat       | 
  |                          |------------------------------->| 
  |                          | if timeout delete deployment   |          
  |                          | and service                    |           
```

# Quick start
## Using the http server
### Install
```
cd terminal/deploy
kubectl create -f rbac.yaml
kubectl create -f deploy.yaml
```

get fist terminal cluster ip

```
[root@iZj6cegflzze2l7fpcqoerZ deploy]# kubectl get svc -n sealyun
NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
fist-terminal   NodePort    10.99.162.142   <none>        8080:31857/TCP   17m
```

Create a terminal from POST, it will return the webterminal NodePort
```
curl -H "Content-Type:application/json" -d '{"User":"fanux", "UserToken":"eyJhbGciOiJSUzI1NiIsImtpZCI6IkNnYzRPVEV5TlRVM0VnWm5hWFJvZFdJIn0.eyJpc3MiOiJodHRwczovL2Zpc3Quc2VhbHl1bi5zdmMuY2x1c3Rlci5sb2NhbDo4MDgwIiwic3ViIjoiQ2djNE9URXlOVFUzRWdabmFYUm9kV0kiLCJhdWQiOiJzZWFseXVuLWZpc3QiLCJleHAiOjE1NTE3MTgzMTgsImlhdCI6MTU1MTM1ODMxOCwiZW1haWwiOiJmaHRqb2JAaG90bWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZ3JvdXBzIjpbImRldiIsInRlc3QiXSwibmFtZSI6ImZhbnV4In0.N8YjE5CqhpPHbtpmNxWIuyIii-VD4vYf8i_JfWpwvJ68SB5JbUM1PN-nXH8xU_biyAB5P1ehNKer2ciMfVA45hyDedzkb6E4A-6B0ezMhU-OL9pCdleq2cJma2R3g64i0YybtQuoOS56kS0isbbmx2NDcL_fluMRdZRPS0Zbx5Pu8CBb60ChXj8r7-zkYuzPFkO1UBi99a3QHR9jBpXz2vBQMjqTPcGddaeUpKP6r80Ib4K9whF3U4rH5gNEOGwgMTe_V4PnEdQ3JbIIiqNJPqIisgxT_HBUqIzHtY3Vsxkjzr2Sj5v0ZWV-gYnZr1ib1bM1uXGK4VOebA5DLbCYaw", "Namespace":"default"}' http://10.99.162.142:8080/terminal
```

OR you can get terminal endpoint in sealyun-tty namespace:
```
[root@iZj6cegflzze2l7fpcqoerZ deploy]# kubectl get svc -n sealyun-tty
NAME                               TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
bc5aaf70741d41c0b4011ebfe7a70f52   NodePort   10.98.60.90   <none>        8080:32219/TCP   16m
```

Using http://IP:32219

## Using bare docker image
imageList:

* fanux/kube-ttyd:latest : only has kubectl and tty
* fanux/fist-tty-tools:v1.0.0 : include golang dev environment, full vim plugins, and vim kubernetes plugins

```
docker run -d --net=host -e APISERVER="https://172.31.12.61:6443" -e USER_TOKEN="XX" -e NAMESPACE="default" -e USER_NAME=fanux -e TERMINAL_ID="uuid" fanux/kube-ttyd:latest
```
Access to :

```
http://yourip:8080
```

OR if you want mount the kubeconfig file instead to use user token:

```
docker run -d --net=host -v /root/.kube/config:/root/.kube/config fanux/kube-ttyd:latest ttyd -p 8080 bash
```


