[![Build Status](https://cloud.drone.io/api/badges/fanux/fist/status.svg)](https://cloud.drone.io/fanux/fist)
[![Go Report Card](https://goreportcard.com/badge/github.com/fanux/fist)](https://goreportcard.com/report/github.com/fanux/fist)

```
                    __                        _____      __ 
   ________  ____ _/ /_  ____  ______        / __(_)____/ /_
  / ___/ _ \/ __ `/ / / / / / / / __ \______/ /_/ / ___/ __/
 (__  )  __/ /_/ / / /_/ / /_/ / / / /_____/ __/ (__  ) /_  
/____/\___/\__,_/_/\__, /\__,_/_/ /_/     /_/ /_/____/\__/  
                  /____/                                    
```

# Fist = (One punch to solve everything)
![](./fist.png)

- [x] A lightweight JWT User token creater. RBAC and PSP manager.
- [x] A powerful webterminal
- [ ] Muti tencent namespace manager
- [ ] Web yaml render
- [ ] kubernets job based pipeline

# Install
```
cd deploy
sh install.sh
```

# Uninstall
```
kubectl delete ns sealyun
kubectl delete ns sealyun-tty
rm -rf /etc/kubernetes/pki/fist
```
and delete oidc config in kube-apiserver.yaml (/etc/kuberentes/manifests/kube-apiserver.yaml)

```
    - --oidc-issuer-url=https://fist.sealyun.svc.cluster.local:8080
    - --oidc-client-id=sealyun-fist
    - --oidc-ca-file=/etc/kubernetes/pki/fist/ca.pem
    - --oidc-username-claim=name
    - --oidc-groups-claim=groups
```

# Auth
Create a kubernetes User token
[README](./auth/README.md)

# Webterminal
![](https://sealyun.com/fist/config-highlight.png)

[terminal show](https://sealyun.com/post/fist-terminal/)

[README](./terminal/README.md)

# Contributing
[Contributing guide](./CONTRIBUTING.md)
