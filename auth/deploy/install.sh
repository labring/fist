rm ssl/*
sh gencert.sh
sleep 3
sh secret.sh
kubectl create -f auth.yaml
mkdir /etc/kubernetes/pki/fist/
cp ssl/ca.pem /etc/kubernetes/pki/fist/
