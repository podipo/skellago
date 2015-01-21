# Skella Back End by [Podipo](http://podipo.com/)

<div style="text-align: center;">
	<img width="150" style="float: left; margin: 0 20px 2px 0;"  src="http://podipo.github.io/skella/images/Skella-logo-300.png" /> 
</div>

The [Skella back end](https://github.com/podipo/skellago/) is a skeleton for web API back ends.  It is designed to work with projects based on the [Skella front end](https://github.com/podipo/skella/).

# Technologies

The Skella back end compiles [go](http://golang.org/) code to produce a web server running in a [Docker](https://www.docker.com) container.  The default production environment is [CoreOS](https://coreos.com/) running in an [CloudFormation](https://aws.amazon.com/cloudformation/) stack on [AWS](http://aws.amazon.com/), but the docker images are portable to other environments.

The Skella back end is designed to produce fantastic web APIs at great speed. It does not handle the user experience side of things other than to serve up static files.  You will want to use one of the many [excellent front end toolkits](https://github.com/podipo/skella/) to produce your HTML, Javascript, and CSS, which the skella back end will then happily serve.

The web API is written on top of the [Negroni](http://negroni.codegangsta.io/) framework.  Negroni has the nice property that it plays nicely with [go's net/http](http://golang.org/pkg/net/http/), so mixing the Skella back end with other go web packages (e.g. a websocket event server) is not hard.

Skella's web API provides a JSON description of itself which is used to automatically create a [Backbone.js](http://backbonejs.org/) wrapper.  It is through Backbone Models and Collections that the Skella front end connects to the Skella back end.  The JSON description is language and framework agnostic, so users of [other client frameworks](http://vanilla-js.com/) should have no problem using Skella.

# Features

- API resource library
- Persistence layer using [QBS](https://github.com/coocood/qbs) and PostgreSQL
- User records and authentication
- API description resource
- Backbone.js wrapper
- Go API client
- Integration with the [Skella front end](https://github.com/podipo/skella/)

# Installation

The Skella back end development environment is CoreOS running in [Vagrant](https://www.vagrantup.com/) hosted on OS X or Linux.  You need a working Vagrant installation, but you don't need to directly install CoreOS as it will be running in a Vagrant managed VM.

You will need a docker client on your host OS.  Follow the instructions in the [Docker installation guide](https://docs.docker.com/installation/#installation) for your operating system.

If you don't already have git then follow the [git download instructions](http://www.git-scm.com/downloads).

If you don't already have make then you'll need to install it.  On OS X, that usually means installing X Code.

Now open up a terminal to check out the code:

	git clone https://github.com/podipo/skellago.git
	cd skellago

Open [https://discovery.etcd.io/new](https://discovery.etcd.io/new) in your browser and cut and paste the resulting URL into the `discovery` field of config/vagrantfile-user-data so that your wee CoreOS cluster can find itself using etcd.

Now fire up CoreOS (and thus Docker) using Vagrant:

	vagrant up

Note: At each of the next few steps, the build system will need to download a lot of images for CoreOS, Postgres, Debian, etc.  This takes a while the first time but have no fear, those downloads happen only on the first run and after that development and management is quite snappy.

To talk to docker you'll need to export a variable naming your docker host:

	export DOCKER_HOST=tcp://:2375 # Or wherever docker is running

Now run this to set up your go third party libraries and then build your project:

	make

The first time you kick off a build process, docker will go fetch the containers which build go and host postgres.  So, the first time you run these commands they will be slow, but after that they will be fast.

To test that everything is built and ready, run the following:

	make start_postgres # fire up the database container
	make start_api      # fire up the API container
	make install_demo   # load up the DB with some example users

Now point your browser at [127.0.0.1:9000/api/0.1.0/schema](http://127.0.0.1:9000/api/0.1.0/schema) and you should see JSON describing the API endpoints.

# Development

If this is your first time using the Skella back end and you just want to see it in action, the easiest thing to do is to set up the Skella front end in a directory next to the skellago directory.  Go follow the instructions on the [Skella front end readme](https://github.com/podipo/skella/) to build the Skella front end and then use `make cycle_api` to rebuild the Skella back end.  The Skella back end assumes that the `skella` and `skellago` directories are next to each other and it serves `skella/dist/index.html` when you hit [127.0.0.1:9000](http://127.0.0.1:9000/).

So, assuming that the Skella front end is now being served by the back end and that you ran the `make install_demo` target listed above, you should be able to authenticate with the back end from [127.0.0.1:9000/login/](http://127.0.0.1:9000/login/) using the email `alice@example.com` and the password `1234`.

To connect directly to the database with psql:

	make psql

To watch the logging on the API service's stdout:

	make watch_api

# Adding API resources

This skeleton project assumes that you're going to add your own API endpoints and fire up your own special API.  The easiest way to get started is to modify [api.go](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/api/api.go) with a few example resources, using the [user](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/user_api.go) and [schema](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/schema.go) resources as examples.

Your normal dev cycle is tweak some go code, rebuild the image, replace any running api container with the new image, and then watch the logs.  Do that using this command:

	make cycle_api
	# This assumes that you have already started a postgres container with `make start_postgres`

# Testing

The Skella back end uses the normal go testing system and includes several handy features for setting up a test DB, a test web API, and a client to exercise the API.  To see how that's done, check out *_test.go files like [user_test.go](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/user_test.go).

To run the tests:

	make test

# Cleaning up

To stop the containers:

	make stop_api
	make stop_postgres

To stop containers and delete the compiled binaries:

	make clean
	
# Todo

- figure out a file persistence story on AWS
- figure out QBS migrations

# Possible future features

- unified logging (syslogd, papertrail, [riemann](http://riemann.io/), [heka](https://blog.mozilla.org/services/2013/04/30/introducing-heka/), [sensu](http://sensuapp.org/), [nagios](http://www.nagios.org/))
- cross machine file storage (S3, nfs, etc)
- backup and restoration
- example project
- websocket resource pubsub and rpc, perhaps [wamp](http://wamp.ws/spec/)
- rate limiting
- brute force auth delays
- monitoring
- CDN

# License

This project is an effort of [Podipo](http://podipo.com/) but depends on a HUGE ecosystem of open source code.  So, what kind of people would we be if we kept Skella all to ourselves?

This project is licensed under the [MIT open source license](http://opensource.org/licenses/MIT).

See the included [LICENSE](https://github.com/podipo/skellago/blob/master/LICENSE) for details.
