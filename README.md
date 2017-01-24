#Drae
A RESTful API for el Diccionario de la Real Academia Española

[![Build Status](https://travis-ci.org/squat/drae.svg?branch=master)](https://travis-ci.org/squat/drae) [![Go Report Card](https://goreportcard.com/badge/github.com/squat/drae)] (https://goreportcard.com/report/github.com/squat/drae) [![](https://images.microbadger.com/badges/image/squat/drae.svg)](https://microbadger.com/images/squat/drae)
##Installation
````sh
$ go get github.com/squat/drae
````

##Usage
###Define
The `define` command does exactly what you would think: it accepts a string argument and returns a JSON object with the definition for that word.

````sh
$ drae define gato
````

###Api
You can run a drae API server from your machine using the `api` command. By default, this server runs on port *4000*, though this can be configured with the `--port` flag.

````sh
$ drae --port=6969 api
````

You can now open a browser and query the API at `http://localhost:6969/api/`.

###Endpoints
####`/api/<word>`
 The API has only one endpoint: `/api/<word>`. For example, I could search for the definition of `gato` at `http://localhost:4000/api/gato`.
 
 ####`/healthz`
 The `/healthz` endpoint returns a 200 if the API is running.
