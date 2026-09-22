# CustomBoard cluster visibility

A CustomBoard can optionally list the configured clusters where it is available:

```yaml
apiVersion: ui.policyreporter.kyverno.io/v1alpha1
kind: CustomBoard
metadata:
  name: production-board
spec:
  title: Production Board
  clusters: [production]
  namespaces:
    list: [team-a]
```

Use the `slug` values in `/api/config`'s `clusters` array, not cluster display names.
Matching is exact and case-sensitive. Multiple entries allow any listed cluster;
unknown IDs do not match. Omitting `clusters` or setting `clusters: []` leaves the
board available on all clusters, subject to existing access controls.

Static UI configuration supports the same field:

```yaml
customBoards:
  - name: Production Board
    clusters: [production]
    namespaces:
      list: [team-a]
```

On other clusters, the board is hidden from navigation and its board-specific
API routes return 404. Boards remain defined in the UI's configuration or its
own Kubernetes cluster; switching clusters does not discover remote CustomBoard
CRDs. This controls board applicability, not access through otherwise authorized
generic APIs. It does not filter report properties or audit annotations.

Install the updated CustomBoard CRD before using `spec.clusters`, so Kubernetes
retains the new field.
