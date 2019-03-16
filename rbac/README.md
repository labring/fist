# rbac for [palm](https://github.com/fanux/palm)


# Quick start
## Using the http server
### Install 

terminal is sub module of fist, if you already install fist, needn't install this again.

```
cd rbac/deploy
kubectl create ns sealyun
kubectl create -f secret.yaml
kubectl create -f deploy.yaml
```

default admin account

> username : admin
> password : 1f2d1e2e67df
