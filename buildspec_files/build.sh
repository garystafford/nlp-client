#!/usr/bin/env bash

set -e

export ISV_ACCOUNT
export ISV_ECR_REGION
$(aws ecr get-login --no-include-email --region ${ISV_ECR_REGION})
docker build -t ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/nlp-client:1.0.0 . --no-cache
docker push ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/nlp-client:1.0.0
