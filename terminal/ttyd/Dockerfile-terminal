FROM centos:7
RUN yum install -y wget &&  \
    wget https://github.com/tsl0922/ttyd/releases/download/1.4.2/ttyd_linux.x86_64 && \
    chmod +x ttyd_linux.x86_64 && \
    mv ttyd_linux.x86_64 /usr/bin/ttyd

# ENV APISERVER="https://127.0.0.1:6443"
# ENV USER_TOKEN="xxx"
# ENV NAMESPACE="default"
# ENV USER_NAME="fanux"
# ENV TERMINAL_ID="UUID"

RUN wget https://dl.k8s.io/v1.13.3/kubernetes-client-linux-amd64.tar.gz && \
    tar zxvf kubernetes-client-linux-amd64.tar.gz && \
    cp kubernetes/client/bin/kubectl /usr/bin && \
    rm -rf kubernetes-client-linux-amd64.tar.gz kubernetes

COPY start-terminal.sh .
CMD ["sh","./start-terminal.sh"]
