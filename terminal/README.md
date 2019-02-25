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
