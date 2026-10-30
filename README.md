# Policy Reporter UI

[![Checks](https://github.com/kyverno/policy-reporter-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/kyverno/policy-reporter-ui/actions/workflows/ci.yml)

Policy Reporter UI is an optional, standalone UI to visualize the processed information of Policy Reporter. It groups related information and adds additional functionality like different graphs, authentication and customizable dashboards on top of it.

<img src="https://kyverno.github.io/policy-reporter-docs/assets/policy-reporter-ui-dark.w6i6-MIf.png" style="border: 1px solid #555" alt="Policy Reporter UI" />

## Feature Overview

* Resource Dashboards
* Policy Dashboards
* Authentication with OAUth2 or OpenIDConnect
* Custom Dashboards
* Multi Cluster Support
* Colorblind Themes

## Documentation

You can find more information about the different features in the [Official Documentation](https://kyverno.github.io/policy-reporter-docs/policy-reporter-ui/introduction.html)

Native HTTPS and optional client-certificate authentication can be configured independently for the UI and metrics listeners. See [TLS configuration and certificate renewal](./docs/tls.md).

## Development

Detailed instructions about the development setup can be found in the [DEVELOPMENT.md](./DEVELOPMENT.md) docs.