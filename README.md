# Skella Back End by [Podipo](http://podipo.com/)

<div style="text-align: center;">
	<img width="150" style="float: left; margin: 0 20px 2px 0;"  src="http://podipo.github.io/skella/images/Skella-logo-300.png" /> 
</div>

The [Skella back end](https://github.com/podipo/skellago/) is a skeleton for web API back ends.  It is designed to work with projects based on the [Skella front end](https://github.com/podipo/skella/).

# Technologies

The Skella back end compiles [go](http://golang.org/) code to produce a web server binary.

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

If you don't already have git then follow the [git download instructions](http://www.git-scm.com/downloads).

If you don't already have make then you'll need to install it.  On OS X, that usually means installing X Code.

You'll also need a PostgreSQL server running somethere the code can reach it.

Now open up a terminal to check out the code:

	git clone https://github.com/podipo/skellago.git
	cd skellago

Edit the Makefile and fix up the POSTGRES_ variables.

Now build, test, and run the service:

	make
	make test
	make install_demo
	make run_api

Now point your browser at [127.0.0.1:9000/api/0.1.0/schema](http://127.0.0.1:9000/api/0.1.0/schema) and you should see JSON describing the API endpoints.

# Development

If this is your first time using the Skella back end and you just want to see it in action, the easiest thing to do is to set up the Skella front end in a directory next to the skellago directory.  Go follow the instructions on the [Skella front end readme](https://github.com/podipo/skella/) to build the Skella front end and then use `make cycle_api` to rebuild the Skella back end.  The Skella back end assumes that the `skella` and `skellago` directories are next to each other and it serves `skella/dist/index.html` when you hit [127.0.0.1:9000](http://127.0.0.1:9000/).

So, assuming that the Skella front end is now being served by the back end and that you ran the `make install_demo` target listed above, you should be able to authenticate with the back end from [127.0.0.1:9000/login/](http://127.0.0.1:9000/login/) using the email `alice@example.com` and the password `1234`.

To connect directly to the database with psql:

	make psql

# Adding API resources

This skeleton project assumes that you're going to add your own API endpoints and fire up your own special API.  The easiest way to get started is to modify [api.go](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/api/api.go) with a few example resources, using the [user](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/user_api.go) and [schema](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/schema.go) resources as examples.

# Testing

The Skella back end uses the normal go testing system and includes several handy features for setting up a test DB, a test web API, and a client to exercise the API.  To see how that's done, check out *_test.go files like [user_test.go](https://github.com/podipo/skellago/blob/master/go/src/podipo.com/skellago/be/user_test.go).

To run the tests:

	make test

# Cleaning up

	make clean
	
# Todo
- vendor our go depenencies, perhaps with godep
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
