# Native TLS and client certificates

The UI/API listener and the metrics listener have independent TLS settings. Both
use HTTP by default; metrics also remains disabled unless `server.metrics.enabled`
is set. TLS and client-certificate authentication each require their own `enabled`
flag. Enabling client authentication requires TLS. Invalid settings stop startup;
there is no fallback to HTTP or to unauthenticated HTTPS.

For example, serve the UI over HTTPS and require a certificate from the ingress,
while allowing the metrics scraper to use ordinary HTTPS with its own certificate:

```yaml
server:
  port: 8080
  tls:
    enabled: true
    certFile: /etc/ui-tls/tls.crt
    keyFile: /etc/ui-tls/tls.key
    clientAuth:
      enabled: true
      caFile: /etc/ui-clients/ca.crt
  metrics:
    enabled: true
    port: 9090
    tls:
      enabled: true
      certFile: /etc/metrics-tls/tls.crt
      keyFile: /etc/metrics-tls/tls.key
```

Omit `clientAuth` for HTTPS without client authentication. Either listener can
instead use HTTP independently by omitting its `tls` block. Set
`server.metrics.tls.clientAuth.enabled` and `caFile` to protect metrics with its
own client trust bundle. Supplying TLS certificate paths with TLS disabled, or a
client CA path with client authentication disabled, is rejected to catch
contradictory configuration.

The existing environment-variable loader supports the same configuration, for
example `PR_SERVER_TLS_ENABLED=true`, `PR_SERVER_TLS_CERTFILE`,
`PR_SERVER_TLS_KEYFILE`, `PR_SERVER_TLS_CLIENTAUTH_ENABLED=true`, and
`PR_SERVER_TLS_CLIENTAUTH_CAFILE`. Use the `PR_SERVER_METRICS_TLS_` prefix for the
corresponding metrics settings.

TLS requires PEM certificate and matching private-key files and supports TLS 1.2
or newer. `certFile` may contain a certificate chain, leaf first. Client
authentication requires a PEM CA bundle and accepts only client certificates
validated against that bundle; system roots are not added. In the default `required` client-authentication mode, missing, expired,
untrusted, or unsuitable client certificates fail the TLS handshake before any
HTTP route is served. Any client certificate trusted by the configured CA and
valid for client authentication can connect: use a dedicated CA or issuing
intermediate to restrict which callers receive access. In `required` mode, certificate authentication protects the connection and
supplements UI OAuth/OIDC authorization if enabled.

## OIDC or certificate login on the same UI listener

To let browsers authenticate through OIDC while explicitly mapped certificate
clients use the same UI/API without an OIDC session, enable OIDC normally and set:

```yaml
server:
  tls:
    enabled: true
    certFile: /etc/ui-tls/tls.crt
    keyFile: /etc/ui-tls/tls.key
    clientAuth:
      enabled: true
      mode: oidc-or-certificate
      caFile: /etc/ui-clients/ca.crt
      identities:
      - fingerprint: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
        email: automation@example.org
        name: Automation
        groups: [report-readers]
```

`oidc-or-certificate` requires OIDC and at least one explicit identity mapping.
Use the lowercase hexadecimal SHA-256 hash of the DER-encoded leaf certificate
(without colons) as `fingerprint`. The configured email and groups are used by the
same cluster and custom-board access controls as an OIDC profile. No identity or
groups are inferred from a certificate's CN, SAN, or subject fields. Possession
of a CA-issued certificate alone does not grant application access in this mode.

When a TLS-terminating reverse proxy presents a mapped client certificate,
**every request forwarded on that connection uses the proxy's mapped service
identity**, including browser requests. For mixed end-user OIDC and certificate
identities, use direct client connections or TLS passthrough so the UI sees each
client's actual certificate choice. Forwarded identity/certificate headers are
not used for certificate login.

Clients without a certificate follow the existing OIDC login/session flow.
Clients presenting a certificate must pass TLS verification and match a configured
fingerprint; an invalid or unmapped certificate is rejected even if an OIDC
session cookie is present. The mapped certificate identity takes precedence over
an existing OIDC session, so that session cannot broaden its permissions.
Certificate validity is checked again for every request on a persistent connection.

Certificate identities exist only for the current request and are never written
into session cookies. Withdrawing a client certificate removes certificate access;
logging out of OIDC cannot revoke a certificate. Disable its mapping or trust and
restart the application to revoke future certificate-authenticated requests. Leaf
certificate renewal changes the fingerprint: deploy the new fingerprint mapping
and restart before switching clients (both old and new mappings may coexist during
rotation), then remove the old mapping and restart. There is no automatic CN/SAN
mapping, CRL, or OCSP lookup.

Metrics does not use OIDC or identity mapping: its client authentication remains
`required`, independently configured. In OIDC-or-certificate mode, `/healthz`
accepts a connection with no client certificate, while an unmapped presented
certificate is rejected; health checks still require HTTPS.

## Deployment and renewal

Certificates can be provisioned by cert-manager or any other external process.
Mount the server key/certificate and client CA bundle read-only, using separate
Secrets for the listeners when their identities or trust differ. Configure the
reverse proxy and scraper to verify each server certificate's CA and DNS name;
for mutual TLS they must also present a client key and certificate issued by the
corresponding client CA.

Certificates, keys, and client CA bundles are loaded once at startup. After any
renewal or trust change, restart/roll out the application to load the new files.
Plan overlapping CA trust during CA rotation and roll the server and clients in
an order that retains a valid chain. There is no automatic certificate reload or
revocation-list/OCSP enforcement.

The `/healthz` route shares the UI listener and its TLS/client-authentication
requirements. With TLS only, configure Kubernetes HTTP probes with `scheme:
HTTPS`. Standard kubelet HTTP probes cannot present client certificates. With
required mutual TLS, use an authenticated probe client or a TCP probe; a TCP probe checks
only the listening socket, not application readiness. There is no unauthenticated
health-route exemption. Metrics scrapes likewise need independent HTTPS and
client-certificate configuration when metrics TLS/client authentication is enabled.

The Helm chart is maintained separately from this repository. Enabling these
settings also requires matching chart Secret mounts, probes, ingress backend
protocol, and ServiceMonitor/scraper settings; the application does not modify
those resources.
