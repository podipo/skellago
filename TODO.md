# Ops

- rename create_api_artifact.sh to create_artifact.sh
- rename the dist dir to collect
- figure out third party go libs
- integration with docker registries
- create a skellago container for building and stop using Igneous's
- DB container
- file storage
- backup and restoration
- deployment to CI, AWS, aor GAE
- better cleanup after building a container
- improve container starting and stopping

# Go

- search again for existing API toolkits, considering [design foundations](https://github.com/interagent/http-api-design/blob/master/README.md).
- start the common API+User code in its own package with an eye toward a separate repo

# V1 Features

- API resource lib
- Persistence layer
- Glue for resources based on persistence records
- User records and authentication
- API description resource
- Backbone.js wrapper
- Example project

# Possible features

- Command line API tool
- Websocket resource events