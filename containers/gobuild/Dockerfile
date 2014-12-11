FROM golang:1.3.1-cross

MAINTAINER trevor@podipo.com

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update
RUN apt-get install -y aufs-tools ca-certificates curl git iptables xz-utils

ENV DOCKER_VERSION 1.3.1

RUN curl -SL https://get.docker.io/builds/Linux/x86_64/docker-$DOCKER_VERSION -o /usr/bin/docker && chmod +x /usr/bin/docker

# This hack from https://registry.hub.docker.com/u/igneoussystems/docker-client/dockerfile/
ENV DOCKER_HOST tcp://172.17.42.1:2375
