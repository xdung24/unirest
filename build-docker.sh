#!/bin/bash

VERSION=$(cat ./VERSION)
docker buildx create --use --name mybuilder
docker buildx build --progress=plain --platform linux/amd64 -t xdung24/unirest:$VERSION -t xdung24/unirest:latest --push .
