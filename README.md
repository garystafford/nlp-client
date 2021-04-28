# Go Microservice: NLP Client

Go-based microservice, part of a set of (5) microservices for the application used in the blog
post, [Amazon ECR Cross-Account Access for Containerized Applications on ECS](https://wp.me/p1RD28-6vd). Please read the
post for complete instructions on how to use the files in this repository.

1. [nlp-client](https://github.com/garystafford/nlp-client)
2. [dynamo-app](https://github.com/garystafford/dynamo-app)
3. [prose-app](https://github.com/garystafford/prose-app)
4. [lang-app](https://github.com/garystafford/lang-app)
5. [rake-app](https://github.com/garystafford/rake-app)

## Architecture

![Architecture](diagram/nlp_diagram2.png)

## Routes

The follow routes are available though the NLP Client.

```json
[
  {
    "method": "GET",
    "path": "/error",
    "name": "main.getError"
  },
  {
    "method": "POST",
    "path": "/keywords",
    "name": "main.getKeywords"
  },
  {
    "method": "POST",
    "path": "/language",
    "name": "main.getLanguage"
  },
  {
    "method": "GET",
    "path": "/health",
    "name": "main.getHealth"
  },
  {
    "method": "GET",
    "path": "/health/:app",
    "name": "main.getHealthUpstream"
  },
  {
    "method": "GET",
    "path": "/routes",
    "name": "main.getRoutes"
  },
  {
    "method": "POST",
    "path": "/tokens",
    "name": "main.getTokens"
  },
  {
    "method": "POST",
    "path": "/entities",
    "name": "main.getEntities"
  },
  {
    "method": "POST",
    "path": "/sentences",
    "name": "main.getSentences"
  },
  {
    "method": "POST",
    "path": "/record",
    "name": "main.putDynamo"
  }
]
```

## Build and Push Docker Image to Docker Hub

```shell
git tag v1.1.0 && git push origin v1.1.0
docker build -t garystafford/nlp-client:1.1.0 . --no-cache
docker push garystafford/nlp-client:1.1.0
```

## Run Services Locally

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
    "http://localhost:${NLP_CLIENT_PORT}/keywords" \
    -H "X-API-Key: ${API_KEY}" \
    -H "Content-Type: application/json" \
    -d "{\"text\": \"${TEXT}\"}"
```

## Build Required Images for ECR

```bash
# change me
export ISV_ACCOUNT=890966919088
export ISV_ECR_REGION=us-west-2
export CUSTOMER_ACCOUNT=676164205626
export CUSTOMER_ECR_REGION=us-east-2

aws ecr get-login-password \
    --region ${ISV_ECR_REGION} \
| docker login \
    --username AWS \
    --password-stdin ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com

docker build -t ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/rake-app:1.1.0 . --no-cache
docker push ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/rake-app:1.1.0

aws ecr get-login-password \
    --region ${CUSTOMER_ECR_REGION} \
| docker login \
    --username AWS \
    --password-stdin ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com

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
export RAKE_ENDPOINT=rake-app:${RAKE_PORT}
export PROSE_ENDPOINT=prose-app:${PROSE_PORT}
export LANG_ENDPOINT=lang-app:${LANG_PORT}
export DYNAMO_ENDPOINT=dynamo-app:${DYNAMO_PORT}
export API_KEY=SuP3r5eCRetAutHK3y
export TEXT="The Nobel Prize is regarded as the most prestigious award in the World. Notable winners have included Marie Curie, Theodore Roosevelt, Albert Einstein, George Bernard Shaw, and Winston Churchill."

docker stack deploy --compose-file stack.yml nlp

# display containers
docker stack ps nlp --no-trunc
docker container ls

# delete stack
docker stack rm nlp
```

Sample output from Docker Swarm stack deployment.

```text
> docker container ls
CONTAINER ID        IMAGE                                                             COMMAND             CREATED             STATUS              PORTS               NAMES
ac5501bb9a79        111222333444.dkr.ecr.us-east-2.amazonaws.com/rake-app:1.1.0       "/go/bin/app"       14 seconds ago      Up 13 seconds                           nlp_rake-app.1.jpctxbvzhcseo8uwuldwlp7hp
7dc171f89f9f        999888777666.dkr.ecr.us-west-2.amazonaws.com/nlp-client:1.1.0     "/go/bin/app"       15 seconds ago      Up 12 seconds                           nlp_nlp-client.1.t96hg46g76uwsvr7i6bweluxz
7ae5369d4293        999888777666.dkr.ecr.us-west-2.amazonaws.com/prose-app:1.1.0      "/go/bin/app"       15 seconds ago      Up 13 seconds                           nlp_prose-app.1.6wkb8x6slva7t253ksucfshyu
4ab51b9f4271        999888777666.dkr.ecr.us-west-2.amazonaws.com/lang-app:1.1.0       "/go/bin/app"       15 seconds ago      Up 13 seconds                           nlp_lang-app.1.hlczjvxpppecosuwdt8bhu7wl
4ab51b9f4271        999888777666.dkr.ecr.us-west-2.amazonaws.com/dynamo-app:1.1.0     "/go/bin/app"       18 seconds ago      Up 13 seconds                           nlp_dynamo-app.1.twe6izgw7xhg6ruvu96sl6b74
```