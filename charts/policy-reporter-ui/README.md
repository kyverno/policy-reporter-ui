# policy-reporter-ui

Trivy Plugin for Policy Reporter UI

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.2](https://img.shields.io/badge/AppVersion-0.0.2-informational?style=flat-square)

## Documentation

You can find detailed Information and Screens about Features and Configurations in the [Documentation](https://kyverno.github.io/policy-reporter).

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| image.registry | string | `"ghcr.io"` | Image registry |
| image.repository | string | `"kyverno/policy-reporter/policy-reporter-ui"` | Image repository |
| image.pullPolicy | string | `"IfNotPresent"` | Image PullPolicy |
| image.tag | string | `""` | Image tag Defaults to `Chart.AppVersion` if omitted |
| replicaCount | int | `1` | Deployment replica count |
| logging.encoding | string | `"console"` | log encoding possible encodings are console and json |
| logging.logLevel | int | `0` | log level default info |
| server.port | int | `8080` | Application port |
| server.logging | bool | `false` | Enables Access logging |
| server.basicAuth.username | string | `""` | HTTP BasicAuth username |
| server.basicAuth.password | string | `""` | HTTP BasicAuth password |
| server.basicAuth.secretRef | string | `""` | Read HTTP BasicAuth credentials from secret |
| policyReporter.host | string | `""` | Host of the Policy Reporter Core App |
| policyReporter.skipTLS | bool | `false` | Skip TLS Verification |
| policyReporter.certificate | string | `""` | TLS Certificate |
| policyReporter.basicAuth.username | string | `""` | HTTP BasicAuth Username |
| policyReporter.basicAuth.password | string | `""` | HTTP BasicAuth Password |
| policyReporter.secretRef | string | `""` | Secret to read the API configuration from supports `host`, `certificate`, `skipTLS`, `username`, `password` key |
| imagePullSecrets | list | `[]` | Image pull secrets for image verification policies, this will define the `--imagePullSecrets` argument |
| nameOverride | string | `""` | Override the name of the chart |
| fullnameOverride | string | `""` | Override the expanded name of the chart |
| serviceAccount.create | bool | `true` | Create ServiceAccount |
| serviceAccount.automount | bool | `true` | Enable ServiceAccount automaount |
| serviceAccount.annotations | object | `{}` | Annotations for the ServiceAccount |
| serviceAccount.name | string | `""` | The ServiceAccount name |
| podAnnotations | object | `{}` | Additional annotations to add to each pod |
| podLabels | object | `{}` | Additional labels to add to each pod |
| updateStrategy | object | `{}` | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| revisionHistoryLimit | int | `10` | The number of revisions to keep |
| podSecurityContext | object | `{"runAsGroup":1234,"runAsUser":1234}` | Security context for the pod |
| envVars | list | `[]` | Allow additional env variables to be added |
| rbac.enabled | bool | `true` | Create RBAC resources |
| securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"runAsUser":1234,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| service.type | string | `"ClusterIP"` | Service type. |
| service.port | int | `8080` | Service port. |
| service.annotations | object | `{}` | Service annotations. |
| service.labels | object | `{}` | Service labels. |
| ingress.enabled | bool | `false` | Create ingress resource. |
| ingress.className | string | `""` | Ingress class name. |
| ingress.labels | object | `{}` | Ingress labels. |
| ingress.annotations | object | `{}` | Ingress annotations. |
| ingress.hosts | list | `[]` | List of ingress host configurations. |
| ingress.tls | list | `[]` | List of ingress TLS configurations. |
| networkPolicy.enabled | bool | `false` | When true, use a NetworkPolicy to allow ingress to the webhook This is useful on clusters using Calico and/or native k8s network policies in a default-deny setup. |
| networkPolicy.egress | list | `[{"ports":[{"port":6443,"protocol":"TCP"}]}]` | A list of valid from selectors according to https://kubernetes.io/docs/concepts/services-networking/network-policies. Enables Kubernetes API Server by default |
| networkPolicy.ingress | list | `[]` | A list of valid from selectors according to https://kubernetes.io/docs/concepts/services-networking/network-policies. |
| resources | object | `{}` |  |
| podDisruptionBudget.minAvailable | int | `1` | Configures the minimum available pods for kyvernoPlugin disruptions. Cannot be used if `maxUnavailable` is set. |
| podDisruptionBudget.maxUnavailable | string | `nil` | Configures the maximum unavailable pods for kyvernoPlugin disruptions. Cannot be used if `minAvailable` is set. |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | List of node taints to tolerate |
| affinity | object | `{}` | Affinity constraints. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)