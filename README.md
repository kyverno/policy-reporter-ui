# Policy Reporter UI

## Motivation

Policy Reporter supports different kinds of targets and a metrics Endpoint to provide as many informations about your (Cluster)PolicyReports as possible but this targets and also the Dashboards have external dependencies. You need tools like Grafana, Kibana or Discord to make informations visible.

In the most production Clusters are Monitoring solutions available so it is not big deal. If you don't have any monitoring solution running because it is only a test environment or you only want to try things out, its additional work and it needs not insignificant resources.

To make policy reporters more accessible, this additional tool adds a standalone, minimal UI on top of Policy Reporter. It offers the same information as Grafana Dashboards and can be enabled as optional Subchart for PolicyReporter.

## Requirements

Policy Reporter runs with the enabled REST API endpoints.

## Configuration

### Flags

* `-backend` is the URL to the Policy Reporter REST API
* `-port` is the default port for the UI (`8080` by default)

## Running with Docker

```bash
docker run -p 8080:8080 --rm fjogeleit/policy-reporter-ui -backend http://host.docker.internal:8081
```

## Screenshots

![Dashboard](https://github.com/fjogeleit/policy-reporter-ui/blob/main/docs/images/dashboard.png?raw=true)

![Policy Reports](https://github.com/fjogeleit/policy-reporter-ui/blob/main/docs/images/policy-report.png?raw=true)

![ClusterPolicyReports](https://github.com/fjogeleit/policy-reporter-ui/blob/main/docs/images/cluster-policy-report.png?raw=true)