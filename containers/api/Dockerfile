FROM debian:wheezy
MAINTAINER trevor@podipo.com
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -y

EXPOSE 8000

ADD ./start.sh /skellago/start.sh
ADD ./api-artifact.tar.gz /opt/root

WORKDIR /opt/root

CMD ["/skellago/start.sh"]
