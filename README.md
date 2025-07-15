# Rick & Morty API Service

A service that fetches, caches, and exposes Rick & Morty character data with observability and automated deployments.

## Features

* **Data Ingestion**: Retrieves "Human" & "Alive" characters from the Rick & Morty public API.
* **Caching**: In-memory caching with Redis. Configurable TTL (default: 5 minutes).
* **Persistence**: Stores character data in PostgreSQL.
* **REST API**:
  open API doc -: https://app.swaggerhub.com/apis/yelchuridinesh/RickMorty/1.0.0
  * `GET /characters` with pagination, sorting (by `name` or `id`), validation, and rate limiting.
  * `GET /healthcheck` to verify Redis and PostgreSQL connectivity.
  * `GET /metrics` to expose Prometheus metrics.
* **Observability**:

  * Prometheus metrics:

    * `api_cache_hits_total`
    * `api_cache_misses_total`
    * `api_characters_processed_total`
    * `api_request_duration_seconds`
  * Jaeger tracing integration for distributed traces.
  * Prometheus alerting rules for:

    * High latency (>500 ms for 5m)
    * Cache miss ratio (>20% for 10m)
    * Low throughput (<10 characters in 5m)
* **Deployments**:

  * Docker containers
  * Kubernetes (Helm charts)
  * GitHub Actions CI/CD pipelines

## Prerequisites

* Go
* Docker
* Kubernetes cluster (e.g., Minikube)
* Helm 3
* Redis and PostgreSQL (via Bitnami Helm charts)
* cert-manager
* ingress-nginx
* kube-prometheus-stack
* Jaeger Helm chart



1. **API Endpoints**:

   * Healthcheck:

     ```http
     GET /healthcheck
     ```
   * List Characters:

     ```http
     GET /characters?page=1&limit=10&sort=name
     GET /characters?page=1&limit=10
     ```
   * Metrics:

     ```http
     GET /metrics
     ```

## Observability

* Exported metrics are scraped by Prometheus.

## CI/CD Workflows

### Integration Test (`.github/workflows/unit-test.yaml`)

* Trigger: `pull` to `main` on event `opened`, `Synchronized`, `Reopened`
* * Steps:
  1. Unit Tests 
  2. CodeQL Security Scan

### Docker CI  (`.github/workflows/docker-ci.yaml`)
* Trigger: `push` to `main`
* Steps:
  1. Build & test Go code
  2. Build & push Docker images to registry

### Kubernetes CD (`.github/workflows/k8s-cd.yaml`)

* Trigger: successful Docker CI on `main`
* Steps:

  1. Spin up Minikube cluster
  2. Install dependencies via Helm:

     * cert-manager
     * bitnami/postgresql
     * bitnami/redis
     * jaegertracing ([https://jaegertracing.github.io/helm-charts](https://jaegertracing.github.io/helm-charts))
     * ingress-nginx ([https://kubernetes.github.io/ingress-nginx](https://kubernetes.github.io/ingress-nginx))
     * prometheus-community/kube-prometheus-stack
  3. Deploy application via your Helm chart repository:

     ```text
     https://github.com/yelchuridinesh/rickmorty-helm-chart.git
     ```

### Integration & Unit Tests

* Unit tests run on pull requests (`.github/workflows/unit-tests.yaml`).
* Integration tests and code coverage run after deployment.

## GitOps Consideration

> **Note**: This project uses CI/CD-based installations for demonstration. In a production environment, consider adopting a GitOps approach (e.g., ArgoCD) to manage Helm chart synchronization across multiple repositories and namespaces more cleanly.Also, I haven't used Umbrella structure as it creates lot of complexity 
>
> [![CD Pipeline](https://github.com/yelchuridinesh/rickmorty-api/actions/workflows/k8s-cd.yaml/badge.svg)](https://github.com/yelchuridinesh/rickmorty-api/actions/workflows/k8s-cd.yaml)


```mermaid
flowchart LR

  %% CI/CD Pipeline
  subgraph CICD["CI/CD Pipeline"]
    Repo[Code Repository]
    GA[GitHub Actions Runner]
    Repo --> GA
    GA --> BuildDocker[Build & Test Docker Images]
    BuildDocker --> PushDocker[Push to Container Registry]
    GA --> BuildHelm[Package Helm Chart]
    BuildHelm --> PushHelm[Push Helm Chart to OCI Registry]
    GA --> Deploy[Deploy via Helm]
  end

  %% Kubernetes Cluster
  subgraph K8s["Kubernetes Cluster"]
    
    subgraph Ingest["Ingest Pipeline"]
      IC[InitContainer / CronJob]
      IC --> Redis[Redis Cache]
      IC --> Postgres[(PostgreSQL)]
    end

    subgraph APIService["API Service"]
      API_Svc[Go Gin Server]
      API_Svc --> Redis
      Redis --> API_Svc
      API_Svc --> Postgres
      API_Svc -->|"/metrics"| Prom[Prometheus]
      API_Svc -->|traces| Jaeger[Jaeger]
    end

    subgraph Network["Service & Ingress"]
      Svc[Service: rickmorty-api]
      Ingress[Ingress → nginx]
      API_Svc --> Svc
      Svc --> Ingress
      Ingress -.-> Client[Client / Browser]
    end

    subgraph Obs["Observability Stack"]
      Prom --> Graf[Grafana Dashboard]
      Jaeger --> Graf
      Prom --> AM[Alertmanager]
      AM --> Pager[On-Call / PagerDuty]
    end

  end

  %% CI/CD → Kubernetes link
  Deploy -.->|"helm pulls charts & images"| K8s



