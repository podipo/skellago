.PHONY: clean cycle_api compile_api collect_api image_api start_api stop_api go_get_deps start_postgres stop_postgres stop_all psql

# Generally, this compiles go using a build container and then builds docker images with the results 

# The remote container which will build the code
BUILD_TAG := podipo/gobuild

# Local container tags
API_TAG := api:dev
API_NAME := api
POSTGRES_TAG := postgres
POSTGRES_NAME := pg
TEST_NAME := test

# TODO: Load these from a config file which is .gitignore'd
FRONT_END_DIR = $(PWD)/../skella/dist
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

all: go_get_deps image_api

clean: stop_all
	-rm -rf go/bin go/pkg deploy collect
	$(DKR_BUILD) docker rmi -f $(API_TAG)

go_get_deps:
	$(DKR_BUILD) /skellago/scripts/container/go_get_deps.sh

clean_deps:
	-rm -rf go/src/github.com go/src/labix.org go/src/code.google.com go/src/golang.org

cycle_api: stop_api image_api start_api watch_api

compile_api: 
	$(DKR_BUILD) go install -v $(API_PKGS)

collect_api: compile_api
	$(DKR_BUILD) /skellago/scripts/container/create_artifact.sh api

image_api: collect_api
	$(DKR_BUILD) /skellago/scripts/container/prepare_image.sh /skellago/containers/api /skellago/collect/api-artifact.tar.gz /skellago/deploy/containers/api
	$(DKR_BUILD) docker build -q --rm -t $(API_TAG) /skellago/deploy/containers/api

start_api: stop_api
	docker run -d $(POSTGRES_ARGS) -v "$(FRONT_END_DIR)":"/opt/root/front_end" -e FRONT_END_DIR=/opt/root/front_end -e SESSION_SECRET="$(SESSION_SECRET)" -p 9000:9000 --link $(POSTGRES_NAME):postgres --name $(API_NAME) $(API_TAG)

stop_api:
	scripts/container_by_image.sh stop $(API_TAG)
	scripts/container_by_image.sh rm $(API_TAG)

watch_api:
	scripts/watch_api.sh

install_demo:
	scripts/install_demo.sh $(POSTGRES_ARGS) --link $(POSTGRES_NAME):postgres $(BUILD_TAG) 

test:
	@DOCKER_FLAGS="$(DOCKER_TEST_ARGS)" $(DKR_BUILD) go test -v $(API_PKGS)

start_postgres:
	docker run -d $(POSTGRES_ARGS) --name $(POSTGRES_NAME) $(POSTGRES_TAG)

stop_postgres:
	scripts/container_by_image.sh stop $(POSTGRES_TAG)
	scripts/container_by_image.sh rm $(POSTGRES_TAG)

psql:
	scripts/db_shell.sh $(POSTGRES_USER) $(POSTGRES_PASSWORD)

stop_all: stop_api stop_postgres

image_gobuild:
	docker build -q --rm -t $(BUILD_TAG) containers/gobuild
