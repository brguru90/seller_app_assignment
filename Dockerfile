FROM debian:buster-slim
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update

ARG expose_port

RUN apt update && apt install  openssh-server sudo -y

RUN mkdir -p  ~/.ssh/
RUN echo "RemoteCommand cd /host_mount && bash -l" >> ~/.ssh/config
RUN service ssh start
RUN update-rc.d ssh enable
RUN apt-get install -y --no-install-recommends apt-utils
# RUN apt-get install build-essential  -y
RUN apt-get install wget curl iputils-ping -y

WORKDIR /tmp

RUN curl -OL https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
RUN ls
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin
ENV PATH $PATH:/usr/local/go/bin
RUN go version

LABEL stage=ubuntu_module

EXPOSE ${expose_port} 22 80 443

WORKDIR /app
COPY app /app/

RUN go build -o ./app main.go && chmod +x ./app

CMD ["./app"]
