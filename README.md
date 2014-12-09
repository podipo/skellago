# Skellago

[Skellago](https://github.com/podipo/skellago/) is a skeleton for web API back ends.  It is designed to work with front ends created by [Skella](https://github.com/podipo/skella/).

# Technologies

Back end logic is written in [Go](http://golang.org/) using the [Negroni](http://negroni.codegangsta.io/) framework and is hosted in [Docker](https://www.docker.com) containers.  Skellago also provides a JSON description of the API as well as a [Backbone.js](http://backbonejs.org/) wrapper.

# Installation

	git clone https://github.com/podipo/skellago.git
	cd skellago
	vagrant up # If you are on OS X and need boot2docker 
	export DOCKER_HOST=tcp://:2375 # Or wherever docker is running
	go get -u github.com/codegangsta/negroni

# Development

	make
	make start_api

	make stop_api

# Testing

TBD

# To-do

## Ops

- configure a DB container
- configure a file storage container
- backup and restoration
- deployment to CI, AWS, aor GAE

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

- [golint](https://github.com/golang/lint)'ing
- create a skellago specific container for building
- better cleanup after building a container
- improve container starting and stopping
- integration with docker registries
- Command line API tool
- Websocket resource events

# License

This project is an effort of [Podipo](http://podipo.com/) but depends on a HUGE ecosystem of open source code.  So, what kind of people would we be if we kept Skella all to ourselves?

This project is licensed under the [MIT open source license](http://opensource.org/licenses/MIT).

See the included [LICENSE](https://github.com/podipo/skellago/blob/master/LICENSE) for details.
