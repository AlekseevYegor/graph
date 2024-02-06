#!/bin/bash

#Local
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=graph_db_user
export DB_PASSWORD=graph_db_user
export DB_NAME=graph
export DB_SCHEMA=graph
export SSL_MODE=false

go build -o graphs .

if [ $? -ne 0 ]; then
        echo "Build failed"
        exit $?
fi

./graphs


