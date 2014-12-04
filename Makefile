.PHONY: clean compile_api collect_api image_api start_api stop_api go_get_dependencies

# Generally, this compiles go using a build container and then builds docker images with the results 

# Remote container tags
# TODO: publish skellago specific containers
BUILD_TAG := igneoussystems/build:2
DOCKER_CLIENT_TAG := igneoussystems/docker-client:1.3.1

# Local container tags
API_TAG := api:dev

# The list of paths to build with Go
API_PKGS := podipo.com/skellago/...

# The prefix for running one-off commands in transient Docker containers
DKR_COMMAND := scripts/docker_command.sh

# The prefix for running commands in the build container
DKR_BUILD := $(DKR_COMMAND) $(BUILD_TAG)

# The prefix for running commands in the docker client container
DKR_CLIENT  := $(DKR_COMMAND) $(DOCKER_CLIENT_TAG)

export GOPATH=go

all: go_get_dependencies image_api

go_get_dependencies:
	go get github.com/codegangsta/negroni
	#go get github.com/goincremental/negroni-sessions
	go get github.com/gorilla/mux
	go get github.com/golang/glog

clean: stop_api
	-rm -rf go/bin go/pkg deploy collect
	-rm -rf go/src/github.com go/src/labix.org
	-docker rmi -f $(API_TAG)

compile_api: 
	$(DKR_BUILD) go install -v $(API_PKGS)

collect_api: compile_api
	$(DKR_CLIENT) /skellago/scripts/container/create_artifact.sh api

image_api: collect_api
	$(DKR_CLIENT) /skellago/scripts/container/prepare_image.sh /skellago/containers/api /skellago/collect/api-artifact.tar.gz /skellago/deploy/containers/api
	$(DKR_CLIENT) docker build -q --rm -t $(API_TAG) /skellago/deploy/containers/api

start_api: stop_api
	docker run -d -p 9000:9000 $(API_TAG)

stop_api:
	scripts/container_by_image.sh stop $(API_TAG)
	scripts/container_by_image.sh rm $(API_TAG)
