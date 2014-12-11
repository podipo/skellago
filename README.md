# Skellago

[Skellago](https://github.com/podipo/skellago/) is a skeleton for web API back ends.  It is designed to work with front ends created by [Skella](https://github.com/podipo/skella/).

# Technologies

Back end logic is written in [Go](http://golang.org/) using the [Negroni](http://negroni.codegangsta.io/) framework and is hosted in [Docker](https://www.docker.com) containers.  Skellago also provides a JSON description of the API as well as a [Backbone.js](http://backbonejs.org/) wrapper.

# Installation

First, install the docker client.

	git clone https://github.com/podipo/skellago.git
	cd skellago
	vagrant up # If you are on OS X and need boot2docker 
	export DOCKER_HOST=tcp://:2375 # Or wherever docker is running
	make go_get_dependencies
	make

# Development

	make image_postgres start_postgres
	make start_api

	make stop_api stop_postgres

	make cycle_api
	make psql
	make watch_api

# Testing

	make test

# To-do

## Ops

- deployment to CI, AWS, aor GAE
- backup and restoration

## Go

- consider [design foundations](https://github.com/interagent/http-api-design/blob/master/README.md)

## Features

- API resource library
- Persistence layer
- User records and authentication
- API description resource
- Backbone.js wrapper
- Service of Skella front end

## Possible future features

- Binary file (esp image) handling
- Example project
- [golint](https://github.com/golang/lint)'ing
- better cleanup after building a container
- option to include skella dist files in api container
- improve container starting and stopping
- integration with docker registries
- Command line API tool
- Websocket resource events

# License

This project is an effort of [Podipo](http://podipo.com/) but depends on a HUGE ecosystem of open source code.  So, what kind of people would we be if we kept Skella all to ourselves?

This project is licensed under the [MIT open source license](http://opensource.org/licenses/MIT).

See the included [LICENSE](https://github.com/podipo/skellago/blob/master/LICENSE) for details.
