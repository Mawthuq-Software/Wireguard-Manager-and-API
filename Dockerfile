# Filename: Dockerfile 
FROM golang:buster
WORKDIR /usr/src/wireguard-manager-and-api
COPY . .

RUN echo 'deb http://ftp.debian.org/debian buster-backports main' | tee /etc/apt/sources.list.d/buster-backports.list
RUN apt-get update
RUN apt-get install sudo -y \
    wireguard \
    dkms \
	git \
	gnupg \ 
	ifupdown \
	iproute2 \
	iptables \
	iputils-ping \
	jq \
	libc6 \
	libelf-dev \
	net-tools \
	openresolv \
	systemctl
RUN go build main.go
RUN adduser --disabled-password --gecos '' docker
RUN adduser docker sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN \
mkdir /app && \
echo "**** install CoreDNS ****" && \
COREDNS_VERSION=$(curl -sX GET "https://api.github.com/repos/coredns/coredns/releases/latest" \
	| awk '/tag_name/{print $4;exit}' FS='[""]' | awk '{print substr($1,2); }') && \
 curl -o \
	/tmp/coredns.tar.gz -L \
	"https://github.com/coredns/coredns/releases/download/v${COREDNS_VERSION}/coredns_${COREDNS_VERSION}_linux_amd64.tgz" && \
 tar xf \
	/tmp/coredns.tar.gz -C \
	/app && \
 echo "**** clean up ****" && \
 rm -rf \
	/tmp/* \
	/var/lib/apt/lists/* \
	/var/tmp/*
RUN mv Corefile /app/Corefile 
USER docker
RUN sudo mv services/coredns.service /etc/systemd/system/
RUN sudo systemctl daemon-reload
RUN sudo chmod +x services/start.sh  

EXPOSE 8443 51820/udp
CMD ["sudo", "./services/start.sh"]

USER docker
