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

