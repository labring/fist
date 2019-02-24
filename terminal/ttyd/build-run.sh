docker build -t fanux/kube-ttyd:latest . -f Dockerfile-terminal
docker run --net=host -e APISERVER="https://172.31.12.61:6443" -e USER_TOKEN="XX" -e NAMESPACE="default" -e USER_NAME=fanux -e TERMINAL_ID="uuid" fanux/kube-ttyd:latest
