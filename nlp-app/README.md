# NLP Application API Helm Chart

This Helm 3 chart will install all Kubernetes resources to the `nlp-app` namespace for the NLP Application API. First,
update your environment specific values in the chart's `values.yaml` file.

Prerequisite: Metrics Server for HPA

<https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html>

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

kubectl get deployment metrics-server -n kube-system
```

Install Helm Chart

```shell
# perform dry run
helm install nlp-app ./nlp-app --namespace nlp-app --debug --dry-run

# apply chart resources
helm install nlp-app ./nlp-app --namespace nlp-app --create-namespace
```

Requirement: DynamoDB Service Account

```shell
aws iam create-policy --policy-name dynamodb-policy \
  --policy-document file://aws/dynamo-app-policy.json

eksctl create iamserviceaccount \
    --cluster $CLUSTER_NAME \
    --namespace nlp-app \
    --name dynamo-app-serviceaccount \
    --attach-policy-arn  arn:aws:iam::$AWS_ACCOUNT:policy/dynamodb-policy \
    --override-existing-serviceaccounts \
    --approve
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