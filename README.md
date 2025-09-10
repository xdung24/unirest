---
runme:
  id: 01HQ2WV4N9YCG2C7Q9Y2HPHZNF
  version: v3
---

A very basic REST service for JSON data - enough for prototyping and MVPs!

**Features**:

- no need to set up a database, all data is managed automagically*
- REST paradigm CRUD for multiple entities/namespaces
- JWT authentication
- realtime notifications (HTTP/SSE)
- schema validation
- autogenerates Swagger/OpenAPI specs
- search using jq like syntax (see https://stedolan.github.io/jq/manual/)
- CORS enabled
- easy to deploy as container

**Currently supports**:

- in memory database (map)
- filesystem storage
- sqlite
- postgres
- mysql
- mongodb

## How to

optional params are:

```yaml {"id":"01HQ2WV4N8PPP7JG0DXVZHRZQJ"}
Usage:
  -AUTH_ENABLED=false: enable JWT auth
  -DB_TYPE="memory": db type to use, options: memory | fs | sqlite | postgres | mysql | mongodb
  -DB_PATH="./data": path of the file storage root or sqlite database
  -IP_PORT=":8000": ip:port to expose
  -PG_HOST="0.0.0.0": postgres host (port is 5432)
  -PG_PASS="": postgres password
  -PG_USER="": postgres user
```

Store a new "user" with an ID and some json data:

```sh {"id":"01HQ2WV4N8PPP7JG0DXZR88Y1Q"}
> curl -X POST -d '{"name":"jack","age":25}'  http://localhost:8000/ns/users/1
{"name":"jack","age":25}
```

the value will be validated, but it could be anything (in JSON!)

retrieve later with:

```sh {"id":"01HQ2WV4N9YCG2C7Q9WJP91MQR"}
> curl http://localhost:8000/ns/users/1
{"name":"jack","age":25}
```

## Sample startup

```sh {"id":"01HQ2WV4N9YCG2C7Q9WK7PFXGC"}
# memory
./unirest --DB_DRIVER=memory --AUTH_ENABLED=true --BROKER_ENABLED=true
```

```sh {"id":"01HQ2WV4N9YCG2C7Q9WNDVFPY8"}
# file system
./unirest --DB_DRIVER=fs --DB_PATH=./data/ --AUTH_ENABLED=true --BROKER_ENABLED=true
```

```sh {"id":"01HQ2WV4N9YCG2C7Q9WPDGK0MB"}
# sqlite
./unirest --DB_DRIVER=sqlite --DB_PATH=./data/db.sqlite --AUTH_ENABLED=true --BROKER_ENABLED=true
```

```sh {"id":"01HQ2WV4N9YCG2C7Q9WQ139BSK"}
# redis
./unirest --DB_DRIVER=redis --DB_HOST=localhost:6379 --AUTH_ENABLED=true --BROKER_ENABLED=true
```

```sh {"id":"01HQ2WV4N9YCG2C7Q9WT3DQW4C"}
# postgres
./unirest --DB_DRIVER=postgres --DB_HOST=localhost:5432 --DB_NAME=nettruyen --DB_USER=postgres --DB_PASS=postgres --AUTH_ENABLED=true --BROKER_ENABLED=true
```

```sh {"id":"01HQ2WV4N9YCG2C7Q9WWWMHP98"}
# mysql/mariadb
./unirest --DB_DRIVER=mysql --DB_HOST=localhost:3306 --DB_NAME=nettruyen --DB_USER=divawallet --DB_PASS=divawallet --AUTH_ENABLED=true --BROKER_ENABLED=true
```

```sh {"id":"01HQ2WV4N9YCG2C7Q9WYHJ04WY"}
# mongodb
./unirest --DB_DRIVER=mongo --DB_HOST=localhost:27017 --DB_NAME=nettruyen --AUTH_ENABLED=true --BROKER_ENABLED=true
```

## All operations

Insert/update

```sh {"id":"01HQ2WV4N9YCG2C7Q9X1DR9J3X"}
> curl -X POST -d '{"name":"jack","age":25}'  http://localhost:8000/ns/users/1
{"name":"jack","age":25}
```

Delete

```sh {"id":"01HQ2WV4N9YCG2C7Q9X3D5KMKD"}
> curl -X DELETE http://localhost:8000/ns/users/1
```

Get by ID

```sh {"id":"01HQ2WV4N9YCG2C7Q9X46XK025"}
> curl http://localhost:8000/ns/users/1
{"name":"jack","age":25}
```

Get all values for a namespace

```sh {"id":"01HQ2WV4N9YCG2C7Q9X5SCA6YZ"}
> curl http://localhost:8000/ns/users | jq 
[
  {
    "key": "2",
    "value": {
      "age": 25,
      "name": "john"
    }
  },
  {
    "key": "1",
    "value": {
      "age": 25,
      "name": "jack"
    }
  }
]
```

Get all namespaces

```sh {"id":"01HQ2WV4N9YCG2C7Q9X8J4132T"}
> curl http://localhost:8000/ns
["users"]
```

Delete a namespace

```sh {"id":"01HQ2WV4N9YCG2C7Q9XA3V8CBV"}
> curl -X DELETE http://localhost:8000/ns/users
{}
```

Search by property (jq syntax)

```sh {"id":"01HQ2WV4N9YCG2C7Q9XDSWJZR0"}
> curl http://localhost:8000/search/users?filter="select(.name==\"jack\")"  | jq
{
  "results": [
    {
      "key": "1",
      "value": {
        "age": 25,
        "name": "jack"
      }
    }
  ]
}
```

## Sample load tests

```sh {"id":"01HQ2WV4N9YCG2C7Q9XEFFTCWW"}
# ren load test
k6 run ./tests/get-user-1.js
```

## JWT Authentication

There's a first implementation of JWT authentication. See [documentation about JWT](JWT.md)

## Realtime Notifications

Using HTTP Server Sent Events (SSE) you can get notified when data changes, just need to listen from the /broker endpoint:

```sh {"id":"01HQ2WV4N9YCG2C7Q9XF6DSYTY"}
curl http://localhost:8000/broker
```

and for every insert or delete an event will be triggered:

```sh {"id":"01HQ2WV4N9YCG2C7Q9XHQD9983"}
{"event":"ITEM_ADDED","namespace":"test","key":"1","value":{"name":"john"}}
...
{"event":"ITEM_DELETED","namespace":"test","key":"1"}
...
```

## Swagger/OpenAPI specs

After you add some data, you can generate the specs with:

```sh {"id":"01HQ2WV4N9YCG2C7Q9XNG0B0K1"}
curl -X GET http://localhost:8000/openapi.json
```

or you can just go to http://localhost:8000/swaggerui/ and use it interactively!

## Schema Validation

You can add a schema for a specific namespace, and only correct JSON data will be accepted

To add a schema for the namespace "user", use the one available in schema_sample/:

```sh {"id":"01HQ2WV4N9YCG2C7Q9XNSQW23S"}
curl --data-binary @./schema_sample/user_schema.json http://localhost:8000/schema/users
```

Now only validated "users" will be accepted (see user.json and invalid_user.json under schema_sample/)

## Run as container

First run an instance of Postgres (for example with docker):

```sh {"id":"01HQ2WV4N9YCG2C7Q9XVVGSYYF"}
docker run -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:latest
```

Then run with the right params to connect to the db:

```sh {"id":"01HQ2WV4N9YCG2C7Q9XZ6YK09M"}
DB_TYPE=postgres PG_HOST=0.0.0.0 PG_USER=postgres PG_PASS=mysecretpassword docker run --publish 8000:8000 xdung24/unirest:latest
```

(params can be passed as ENV variables or as command-line ones)

A very quick to run both on docker with docker-compose:

```sh {"id":"01HQ2WV4N9YCG2C7Q9Y0YX0KPR"}
docker-compose up -d
```
