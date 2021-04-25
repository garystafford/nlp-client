# Go Microservice: nlp-client

Go-based microservice, part of a set of (5) microservices for the application used in the blog
post, [Amazon ECR Cross-Account Access for Containerized Applications on ECS](https://wp.me/p1RD28-6vd). Please read the
post for complete instructions on how to use the files in this repository.

## Build and Push to Docker Hub

```shell
docker build -t garystafford/nlp-client:1.1.0 . --no-cache
docker push garystafford/nlp-client:1.1.0
```

## Run from IDE

Create [DynamoDB CloudFormation stack](https://github.com/garystafford/dynamo-app/blob/master/dynamodb-table.yml) from
the `dynamodb-table.yml` CloudFormation template. Creates the `NLPText` table.

```bash
aws cloudformation create-stack \
    --stack-name dynamodb-table-stack \
    --template-body file://dynamodb-table.yml
```

Run each of the (5) service from a different terminal window.

```bash
    export NLP_CLIENT_PORT=8080
    export RAKE_PORT=8081
    export PROSE_PORT=8082
    export LANG_PORT=8083
    export DYNAMO_PORT=8084
    export RAKE_ENDPOINT=http://localhost:${RAKE_PORT}
    export PROSE_ENDPOINT=http://localhost:${PROSE_PORT}
    export LANG_ENDPOINT=http://localhost:${LANG_PORT}
    export DYNAMO_ENDPOINT=http://localhost:${DYNAMO_PORT}
    export API_KEY=SuP3r5eCRetAutHK3y
    export TEXT="The Nobel Prize is regarded as the most prestigious award in the World. Notable winners have included Marie Curie, Theodore Roosevelt, Albert Einstein, George Bernard Shaw, and Winston Churchill."

## Run service locally

go mod init github.com/garystafford/nlp-client
go mod tidy -v
go run *.go

curl -s -X GET \
    "http://localhost:${NLP_CLIENT_PORT}/health" \
    -H "X-API-Key: ${API_KEY}" \
    -H "Content-Type: application/json"

curl -s -X GET \
    "http://localhost:${NLP_CLIENT_PORT}/routes" \
    -H "X-API-Key: ${API_KEY}" \
    -H "Content-Type: application/json"

curl -s -X POST \
    "http://localhost:${NLP_CLIENT_PORT}//keywords" \
    -H "X-API-Key: ${API_KEY}" \
    -H "Content-Type: application/json" \
    -d "{\"text\": \"${TEXT}\"}"
```

## Build Required Images for Docker

```bash
# change me
export ISV_ACCOUNT=111222333444
export ISV_ECR_REGION=us-east-2

$(aws ecr get-login --no-include-email --region ${ISV_ECR_REGION})
docker build -t ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/rake-app:1.1.0 . --no-cache
docker push ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/rake-app:1.1.0

# change me
export CUSTOMER_ACCOUNT=999888777666
export CUSTOMER_ECR_REGION=us-west-2

$(aws ecr get-login --no-include-email --region ${CUSTOMER_ECR_REGION})
docker build -t ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/nlp-client:1.1.0 . --no-cache
docker push ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/nlp-client:1.1.0

docker build -t ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/prose-app:1.1.0 . --no-cache
docker push ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/prose-app:1.1.0

docker build -t ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/lang-app:1.1.0 . --no-cache
docker push ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/lang-app:1.1.0

docker build -t ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/dynamo-app:1.1.0 . --no-cache
docker push ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/dynamo-app:1.1.0

# display images
docker image ls --filter=reference='*amazonaws.com/*'
```

## Deploy Docker Stack

Running (5) service Stack locally, from Docker Swarm.

```bash
# change me
export ISV_ACCOUNT=111222333444
export ISV_ECR_REGION=us-east-2
export CUSTOMER_ACCOUNT=999888777666
export CUSTOMER_ECR_REGION=us-west-2

# don't change me
export NLP_CLIENT_PORT=8080
export RAKE_PORT=8080
export PROSE_PORT=8080
export LANG_PORT=8080
export DYNAMO_PORT=8080
export RAKE_ENDPOINT=http://localhost:${RAKE_PORT}
export PROSE_ENDPOINT=http://localhost:${PROSE_PORT}
export LANG_ENDPOINT=http://localhost:${LANG_PORT}
export DYNAMO_ENDPOINT=http://localhost:${DYNAMO_PORT}
export API_KEY=SuP3r5eCRetAutHK3y

docker stack deploy --compose-file stack.yml nlp

# display containers
docker stack ps nlp --no-trunc
docker container ls

# delete stack
docker stack rm nlp
```

Sample output from Docker Swarm stack deployment.

```text
> ~/environment/ecr-cross-accnt-demo (master) $ docker stack deploy --compose-file stack.yml nlp
Creating network nlp_nlp-demo
Creating service nlp_nlp-client
Creating service nlp_rake-app
Creating service nlp_prose-app
Creating service nlp_lang-app
Creating service nlp_dynamo-app

> ~/environment/ecr-cross-accnt-demo (master) $ docker container ls
CONTAINER ID        IMAGE                                                             COMMAND             CREATED             STATUS              PORTS               NAMES
ac5501bb9a79        111222333444.dkr.ecr.us-east-2.amazonaws.com/rake-app:1.1.0       "/go/bin/app"       14 seconds ago      Up 13 seconds                           nlp_rake-app.1.jpctxbvzhcseo8uwuldwlp7hp
7dc171f89f9f        999888777666.dkr.ecr.us-west-2.amazonaws.com/nlp-client:1.1.0     "/go/bin/app"       15 seconds ago      Up 12 seconds                           nlp_nlp-client.1.t96hg46g76uwsvr7i6bweluxz
7ae5369d4293        999888777666.dkr.ecr.us-west-2.amazonaws.com/prose-app:1.1.0      "/go/bin/app"       15 seconds ago      Up 13 seconds                           nlp_prose-app.1.6wkb8x6slva7t253ksucfshyu
4ab51b9f4271        999888777666.dkr.ecr.us-west-2.amazonaws.com/lang-app:1.1.0       "/go/bin/app"       15 seconds ago      Up 13 seconds                           nlp_lang-app.1.hlczjvxpppecosuwdt8bhu7wl
4ab51b9f4271        999888777666.dkr.ecr.us-west-2.amazonaws.com/dynamo-app:1.1.0     "/go/bin/app"       18 seconds ago      Up 13 seconds                           nlp_dynamo-app.1.twe6izgw7xhg6ruvu96sl6b74
```