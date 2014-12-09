.PHONY: clean compile_api collect_api image_api start_api stop_api go_get_dependencies image_postgres start_postgres stop_postgres stop_all psql

# Generally, this compiles go using a build container and then builds docker images with the results 

# Remote container tags
# TODO: publish skellago specific containers
BUILD_TAG := igneoussystems/build:2
DOCKER_CLIENT_TAG := igneoussystems/docker-client:1.3.1

# Local container tags
API_TAG := api:dev
API_NAME := api
POSTGRES_TAG := postgres:dev
POSTGRES_NAME := pg
TEST_NAME := test

# TODO: Load these from a config file which is .gitignore'd
POSTGRES_DB_NAME := skella
POSTGRES_TEST_DB_NAME := test
POSTGRES_USER := skella
POSTGRES_PASSWORD := seekret
SESSION_SECRET := "fr0styth3sn0wm@n"

POSTGRES_AUTH_ARGS := -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) 
POSTGRES_ARGS := $(POSTGRES_AUTH_ARGS) -e POSTGRES_DB_NAME=$(POSTGRES_DB_NAME)

DOCKER_TEST_ARGS := $(POSTGRES_AUTH_ARGS) -e POSTGRES_DB_NAME=$(POSTGRES_TEST_DB_NAME) --link $(POSTGRES_NAME):postgres

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
	go get github.com/chai2010/assert
	go get github.com/codegangsta/negroni
	go get github.com/gorilla/mux
	go get github.com/coocood/qbs
	go get github.com/lib/pq
	go get code.google.com/p/go.crypto/bcrypt
	go get github.com/goincremental/negroni-sessions
	go get github.com/nu7hatch/gouuid

clean: stop_all
	-rm -rf go/bin go/pkg deploy collect
	-rm -rf go/src/github.com go/src/labix.org go/src/code.google.com go/src/golang.org
	-docker rmi -f $(API_TAG)

compile_api: 
	$(DKR_BUILD) go install -v $(API_PKGS)

collect_api: compile_api
	$(DKR_CLIENT) /skellago/scripts/container/create_artifact.sh api

image_api: collect_api
	$(DKR_CLIENT) /skellago/scripts/container/prepare_image.sh /skellago/containers/api /skellago/collect/api-artifact.tar.gz /skellago/deploy/containers/api
	$(DKR_CLIENT) docker build -q --rm -t $(API_TAG) /skellago/deploy/containers/api

start_api: stop_api
	docker run -d $(POSTGRES_ARGS) -e SESSION_SECRET="$(SESSION_SECRET)" -p 9000:9000 --link $(POSTGRES_NAME):postgres --name $(API_NAME) $(API_TAG)

stop_api:
	scripts/container_by_image.sh stop $(API_TAG)
	scripts/container_by_image.sh rm $(API_TAG)

install_demo:
	scripts/install_demo.sh $(POSTGRES_ARGS) --link $(POSTGRES_NAME):postgres $(BUILD_TAG) 

test:
	@DOCKER_FLAGS="$(DOCKER_TEST_ARGS)" $(DKR_BUILD) go test -v $(API_PKGS)

image_postgres:
	$(DKR_CLIENT) docker build --rm -t $(POSTGRES_TAG) /skellago/containers/postgres

start_postgres:
	docker run -d $(POSTGRES_ARGS) --name $(POSTGRES_NAME) $(POSTGRES_TAG)

stop_postgres:
	scripts/container_by_image.sh stop $(POSTGRES_TAG)
	scripts/container_by_image.sh rm $(POSTGRES_TAG)

psql:
	scripts/db_shell.sh $(POSTGRES_USER) $(POSTGRES_PASSWORD)

stop_all: stop_api stop_postgres
