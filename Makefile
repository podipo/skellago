.PHONY: clean compile_api dist_api image_api

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


all: image_api

clean: stop_api
	-rm -rf go/bin go/pkg deploy dist
	-docker rmi -f $(API_TAG)

compile_api: 
	$(DKR_BUILD) go install -v $(API_PKGS)

dist_api: compile_api
	$(DKR_CLIENT) /skellago/scripts/container/create_api_artifact.sh api

image_api: dist_api
	$(DKR_CLIENT) /skellago/scripts/container/prepare_image.sh /skellago/containers/api /skellago/dist/api-artifact.tar.gz /skellago/deploy/containers/api
	$(DKR_CLIENT) docker build --rm -t $(API_TAG) /skellago/deploy/containers/api

start_api: stop_api
	docker run -d -p 8000:8000 $(API_TAG)

stop_api:
	scripts/container_by_image.sh stop $(API_TAG)
	scripts/container_by_image.sh rm $(API_TAG)
