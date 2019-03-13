rm ssl/*
sh gencert.sh
sleep 3
sh secret.sh
kubectl create -f auth.yaml
kubectl create -f secret.yaml
mkdir /etc/kubernetes/pki/fist/ || true
cp ssl/ca.pem /etc/kubernetes/pki/fist/

echo "wait for auth service sleep 15s... "
sleep 15 

echo '  [WARN] edit kube-apiserver.yaml and add oidc config, if auth service not ready, apiserver may start failed!'
sed '/- kube-apiserver/a\    - --oidc-groups-claim=groups' -i /etc/kubernetes/manifests/kube-apiserver.yaml
sed '/- kube-apiserver/a\    - --oidc-username-claim=name' -i /etc/kubernetes/manifests/kube-apiserver.yaml
sed '/- kube-apiserver/a\    - --oidc-ca-file=/etc/kubernetes/pki/fist/ca.pem' -i /etc/kubernetes/manifests/kube-apiserver.yaml
sed '/- kube-apiserver/a\    - --oidc-client-id=sealyun-fist' -i /etc/kubernetes/manifests/kube-apiserver.yaml
sed '/- kube-apiserver/a\    - --oidc-issuer-url=https://fist.sealyun.svc.cluster.local:8080' -i /etc/kubernetes/manifests/kube-apiserver.yaml
