rm ssl/*
sh gencert.sh
sleep 3
sh secret.sh
kubectl create -f auth.yaml
kubectl create -f secret.yaml
mkdir /etc/kubernetes/pki/fist/ || true
cp ssl/ca.pem /etc/kubernetes/pki/fist/
