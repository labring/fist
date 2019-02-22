sh gencert.sh
sh secret.sh
kubectl create -f auth.yaml
mkdir /etc/kubernetes/pki/fist/
cp ssl/ca.pem /etc/kubernetes/pki/fist/
