# Policy Reporter UI

## Motivation

Policy Reporter supports different kinds of targets and a metrics Endpoint to provide as many information about your (Cluster)PolicyReports as possible, but this targets and also the dashboards have external dependencies. You need tools like Grafana, Kibana or Discord to make information visible.

In the most production clusters are monitoring solutions available, so it is not big deal. If you don't have any (supported) monitoring solution, it is additional work and it needs not insignificant resources.

To make Policy Reporter more accessible, this additional tool adds a standalone, minimal UI. It offers the same information as the provided Grafana Dashboards and is installable as optional SubChart in the Policy Reporter Helm Chart.

## Requirements

Running and accessable Policy Reporter Application

## Configuration

### Backend Flags

* `-policy-reporter` is the URL to the Policy Reporter REST API
* `-kyverno-plugin` is the URL to the Policy Reporter Kyverno Plugin REST API
* `-port` is the default port for the UI (`8080` by default)
* `-dev` enables development mode and adds CORS Headers to the REST API Endpoints
* `-no-ui` runs the Go Backend Application without the Frontend - mainly for development purpose

## Running with Docker

```bash
docker run -p 8082:8080 --rm ghcr.io/kyverno/policy-reporter-ui -backend http://host.docker.internal:8080
```

## Build Setup

```bash
# install dependencies
$ npm install

# serve with hot reload at localhost:3000
$ npm run dev

# build for production and launch server
$ npm run build
$ npm run start

# generate static project
$ npm run generate
```

For detailed explanation on how things work, check out the [nuxtjs documentation](https://nuxtjs.org).
