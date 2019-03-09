cd ../auth/deploy && sh install.sh
echo "10.106.233.67 fist.sealyun.svc.cluster.local" >> /etc/hosts
cd - && cd ../terminal/deploy && kubectl create -f rbac.yaml && kubectl create -f deploy.yaml
