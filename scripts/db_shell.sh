#/bin/bash

docker run -it --link pg:postgres --rm -e POSTGRES_USER=$1 -e PGPASSWORD=$2 postgres sh -c 'exec psql -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U "$POSTGRES_USER"'
