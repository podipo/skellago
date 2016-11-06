.PHONY: clean clean_deps go_get_deps lint compile_api install_demo test psql

PORT := 9000
FRONT_END_DIR = $(PWD)/../skella/dist

POSTGRES_USER := trevor
POSTGRES_PASSWORD := seekret
POSTGRES_HOST := localhost
POSTGRES_PORT := 5432

POSTGRES_DB_NAME := skella
POSTGRES_TEST_DB_NAME := skella_test

SESSION_SECRET := "fr0styth3sn0wm@n"

STATIC_DIR := $(PWD)/go/src/podipo.com/skellago/be/static/
FILE_STORAGE_DIR := $(PWD)/file_storage

API_PKGS := podipo.com/skellago/... example.com/api/...

COMMON_POSTGRES_ENVS := POSTGRES_USER=$(POSTGRES_USER) \
						POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
						POSTGRES_HOST=$(POSTGRES_HOST) \
						POSTGRES_PORT=$(POSTGRES_PORT)

API_POSTGRES_ENVS :=	$(COMMON_POSTGRES_ENVS) \
						POSTGRES_DB_NAME=$(POSTGRES_DB_NAME)


TEST_POSTGRES_ENVS := 	$(COMMON_POSTGRES_ENVS) \
						POSTGRES_DB_NAME=$(POSTGRES_TEST_DB_NAME) 

API_RUNTIME_ENVS := 	PORT=$(PORT) \
						STATIC_DIR=$(STATIC_DIR) \
						FILE_STORAGE_DIR=$(FILE_STORAGE_DIR) \
						FRONT_END_DIR=$(FRONT_END_DIR) \
						SESSION_SECRET=$(SESSION_SECRET) \
						$(API_POSTGRES_ENVS)

all: go_get_deps compile_api

clean:
	rm -rf go/bin go/pkg deploy collect

clean_deps:
	rm -rf go/src/github.com go/src/labix.org go/src/code.google.com go/src/golang.org

go_get_deps:
	go get github.com/chai2010/assert
	go get github.com/codegangsta/negroni
	go get github.com/gorilla/mux
	go get github.com/coocood/qbs
	go get github.com/lib/pq
	go get code.google.com/p/go.crypto/bcrypt
	go get github.com/nu7hatch/gouuid
	go get github.com/rs/cors
	go get github.com/goincremental/negroni-sessions
	go get github.com/golang/lint
	go get github.com/nfnt/resize

lint:
	go install github.com/golang/lint/...
	golint podipo.com/...

compile_api: 
	go install -v $(API_PKGS)

run_api: compile_api
	-mkdir $(FILE_STORAGE_DIR)
	$(API_RUNTIME_ENVS) go/bin/api

install_demo:
	-echo "drop database $(POSTGRES_DB_NAME); create database $(POSTGRES_DB_NAME);" | psql
	go install -v podipo.com/skellago/demo
	go install -v example.com/api/example_demo
	$(API_POSTGRES_ENVS) $(GOBIN)/demo
	$(API_POSTGRES_ENVS) $(GOBIN)/example_demo

test:
	-echo "drop database $(POSTGRES_TEST_DB_NAME)" | psql
	$(TEST_POSTGRES_ENVS) go test -v $(API_PKGS)

psql:
	scripts/db_shell.sh $(POSTGRES_USER) $(POSTGRES_PASSWORD)

