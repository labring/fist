kubectl delete secret fist -n sealyun
kubectl create secret generic fist --from-file=ssl/cert.pem --from-file=ssl/key.pem -n sealyun

