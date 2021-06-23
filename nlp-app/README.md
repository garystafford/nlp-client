# NLP Application API Helm Chart

This Helm 3 chart will install all Kubernetes resources to the `nlp-app` namespace for the NLP Application API. First, place your environment specific values in the chart's `values.yaml`.

```shell
# perform dry run
helm install nlp-app ./nlp-app --namespace nlp-app --debug --dry-run

# apply chart resources
helm install nlp-app ./nlp-app --namespace nlp-app | kubectl apply -f
helm upgrade nlp-app ./nlp-app --namespace nlp-app --create-namespace
```

Resources included in Helm Chart:

```text
.
├── destination-rules.yaml
├── gateway.yaml
├── nlp-app-deployment.yaml
├── nlp-app-hpa.yaml
├── nlp-app-secret.yaml
├── nlp-app-service.yaml
└── virtual-services.yaml
```