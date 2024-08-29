# Running Policy Reporter UI

## Requirements

### Tools Backend

* Go >= v1.22.4

### Tools Frontend

* [Bun](https://bun.sh/)

### Services

* Local Kubernetes Cluster
* Installed PolicyReport and ClusterPolicyReport CRDs
* Running Policy Reporter Core App with REST APIs enabled

## Running Steps

Prerequired Steps (without existing Cluster)

1. Create a local cluster (e.g. with Kind or Minikube)

```bash
kind create cluster -n kyverno
```

2. Install Kyverno + Kyverno PSS Policies

Kyverno installs the required CRDS and you get some sample PolicyReports by installing the PSS policies.

Add Helm chart

```bash
helm repo add kyverno https://kyverno.github.io/kyverno/
helm repo update
```

install Kyverno + CRDs

```bash
helm upgrade --install kyverno kyverno/kyverno -n kyverno --create-namespace
```

Add Pod Security Standard Policies

```bash
helm upgrade --install kyverno-policies kyverno/kyverno-policies -n kyverno --set podSecurityStandard=restricted
```

3. Checkout and Running Policy Reporter Core App v3.x

Clone Repository

```
git clone https://github.com/kyverno/policy-reporter.git
git checkout 3.x 
```

Install dependencies

```bash
go get ./...
```

Run the core app locally against your local cluster with enabled REST APIs.

```bash
go run main.go run -k /Users/<user>/.kube/config -p 8080 -r
```

4. Start the Policy Reporter UI Backend

Create an `config.yaml` file in the `backend` folder. You can copy and rename the `config.example.yaml`.

```bash
cp backend/config.example.yaml backend/config.yaml
```

Within the `/backend` folder of the `policy-reporter-ui` repository.

Install dependencies

```bash
go get ./...
```

Start the UI Backend on Port 8082 in local mode

```bash
go run main.go run --local --port 8082
```

5. Start Policy Reporter UI frontend

Create an `.env` file in the `frontend` folder. You can copy and rename the `.env.example`.

```bash
cp frontend/.env.example frontend/.env
```

Within the `/frontend` folder of the `policy-reporter-ui` repository.

Install dependencies

```bash
bun install
```

Start the UI Frontend on Port 3000

```bash
bun run dev
```

Access the UI under `http://localhost:3000`