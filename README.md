# Policy Reporter UI

## Motivation

Policy Reporter supports different kinds of targets and a metrics Endpoint to provide as many information about your (Cluster)PolicyReports as possible, but this targets and also the dashboards have external dependencies. You need tools like Grafana, Kibana or Discord to make information visible.

In the most production clusters are monitoring solutions available, so it is not big deal. If you don't have any (supported) monitoring solution, it is additional work and it needs not insignificant resources.

To make Policy Reporter more accessible, this additional tool adds a standalone, minimal UI. It offers the same information as the provided Grafana Dashboards and is installable as optional SubChart in the Policy Reporter Helm Chart.

## Requirements

Policy Reporter runs with the enabled REST API endpoints.

## Configuration

### Flags

* `-backend` is the URL to the Policy Reporter REST API
* `-port` is the default port for the UI (`8080` by default)
* `-dev` enables development mode and adds CORS Headers to the REST API Endpoints

## Running with Docker

```bash
docker run -p 8080:8080 --rm fjogeleit/policy-reporter-ui -backend http://host.docker.internal:8081
```

## Screenshots

![Dashboard](https://github.com/kyverno/policy-reporter-ui/blob/main/docs/images/dashboard.png?raw=true)

![Policy Reports](https://github.com/kyverno/policy-reporter-ui/blob/main/docs/images/policy-report.png?raw=true)

![ClusterPolicyReports](https://github.com/kyverno/policy-reporter-ui/blob/main/docs/images/cluster-policy-report.png?raw=true)
