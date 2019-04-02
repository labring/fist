FROM centos:7.6.1810
COPY ./auth/auth .
COPY auth/deploy/start.sh .
ADD /auth/deploy/ssl/ /etc/fist-genKey/
CMD ["sh", "start.sh"]
