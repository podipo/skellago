# Skella Back End

The [Skella back end](https://github.com/podipo/skellago/) is a skeleton for web API back ends.  It is designed to work with projects based on the [Skella front end](https://github.com/podipo/skella/).

# Technologies

After all is said and done, the Skella back end creates a [Docker](https://www.docker.com) container running a web API which is backed by another container running PostgreSQL.  The web API is written in [Go](http://golang.org/) using the [Negroni](http://negroni.codegangsta.io/) framework.  The API provides a JSON description of itself which is used to automatically create a [Backbone.js](http://backbonejs.org/) wrapper.  It is through Backbone Models and Collections that the Skella front end connects to the Skella back end.

# Features

- API resource library
- Persistence layer using [QBS](https://github.com/coocood/qbs) and PostgreSQL
- User records and authentication
- API description resource
- Backbone.js wrapper
- Integration with the Skella front end

# Installation

Your development environment will need to be running a linux or OS X.  If you're on Windows, fire up a VM with linux.

You will need a working docker container and the docker client.  Follow the instructions in the [Docker installation guide](https://docs.docker.com/installation/#installation) for your operating system.

If you don't already have git, follow the [git download instructions](http://www.git-scm.com/downloads).

If you don't already have `make`, you'll need to install it.  On OS X, that usually means installing X Code.

Now open up a terminal to check out the code:

	git clone https://github.com/podipo/skellago.git
	cd skellago

Run this if you're on OS X and need to run boot2docker using [Vagrant](https://www.vagrantup.com/):

	vagrant up

On linux, make sure that the docker daemon is running.

No matter your OS, you'll need to export a variable naming your docker host:

	export DOCKER_HOST=tcp://:2375 # Or wherever docker is running

Now run this to set up your go third party libraries and then build your project:

	make go_get_dependencies
	make

The first time you kick off a build process, docker will go fetch the containers which build go and host postgres.  So, the first time you run these commands they will be slow, but after that they will be fast.

To test that everything is built and ready, run the following:

	make start_postgres # fire up the database container
	make start_api      # fire up the API container
	make install_demo   # load up the DB with some example users

Now point your browser at [127.0.0.1:9000/api/schema](http://127.0.0.1:9000/api/schema) and you should see JSON describing the API endpoints.

# Development

If this is your first time using the Skella back end and you just want to see it in action, the easiest thing to do is to set up the Skella front end in a directory next to the skellago directory.  Go follow the instructions on the [Skella front end readme](https://github.com/podipo/skella/) to build skella.  The Skella back end assumes that the `skella` and `skellago` directories are next to each other and it serves `skella/dist/index.html` when you hit [127.0.0.1:9000](http://127.0.0.1:9000/).

So, assuming that the Skella front end is now being served by the back end and that you ran the `install_demo` target listed above, you should be able to authenticate with the back end from [127.0.0.1:9000/login/](http://127.0.0.1:9000/login/) using the email `alice@example.com` and the password `1234`.

To connect directly to the database with psql:

	make psql

To watch the logging on the API service's stdout:

	make watch_api

# Adding API resources

This skeleton project assumes that you're going to add your own API endpoints and fire up your own special API.  The easiest way to get started is to modify [api.go](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/api/api.go) with a few example resources, using the [user](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/user_api.go) and [schema](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/schema.go) resources as examples.

Your normal dev cycle is tweak some go code, rebuild the image, replace any running api container with the new image, and then watch the logs.  Do that using this command:

	make cycle_api
	# This assumes that you already have a running postgres container

# Testing

The Skella back end uses the normal go testing system and includes several handy features for setting up a test DB, a test web API, and a client to exercise the API.  To see how that's done, check out *_test.go files like [user_test.go](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/user_test.go).

To run the tests:

	make test

# Cleaning up

To stop the containers:

	make stop_api
	make stop_postgres
	# or to stop both
	make stop_all

	make clean # stops containers, removes the api image, then deletes the compiled binaries

# To-do

## Ops

- deployment to CI, AWS, aor GAE
- backup and restoration

## Go

- consider [API design foundations](https://github.com/interagent/http-api-design/blob/master/README.md)

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
