# Skellago

Skellago is a skeleton for web API back ends.  It is designed to work with front ends created by [Skella](https://github.com/podipo/skella/).

# Technologies

Back end logic is written in [Go](http://golang.org/) using the [Negroni](http://negroni.codegangsta.io/) framework and is hosted in [Docker](https://www.docker.com) containers.  Skellago also provides a JSON description of the API as well as a [Backbone.js](http://backbonejs.org/) wrapper.

# Installation

	git clone ...
	cd skellago
	export DOCKER_HOST=tcp://:2375
	vagrant up
	go get -u github.com/codegangsta/negroni

# Development

	make
	make start_api

	make stop_api

# Testing

TBD

# To-do

## Ops

- figure out third party go libs
- integration with docker registries
- create a skellago container for building and stop using Igneous's
- DB container
- file storage
- backup and restoration
- deployment to CI, AWS, aor GAE
- better cleanup after building a container
- improve container starting and stopping

## Go

- search again for existing API toolkits, considering [design foundations](https://github.com/interagent/http-api-design/blob/master/README.md).
- start the common API+User code in its own package with an eye toward a separate repo

## V1 Features

- API resource lib
- Persistence layer
- Glue for resources based on persistence records
- User records and authentication
- API description resource
- Backbone.js wrapper
- Example project

## Possible features

- Command line API tool
- Websocket resource events

# License

This project is an effort of [Podipo](http://podipo.com/) but depends on a HUGE ecosystem of open source code.  So, what kind of people would we be if we kept Skella all to ourselves?

This project is licensed under the [MIT open source license](http://opensource.org/licenses/MIT).

See the included [LICENSE](https://github.com/podipo/skellago/blob/master/LICENSE) for details.
