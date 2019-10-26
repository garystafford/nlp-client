# Go Microservice: nlp-client

Go-based microservice, part of a set of (3) microservices for the application used in the blog post, [Amazon Elastic Container Registry (ECR) Cross-Account Access](https://programmaticponderings.com/). Please read the post for complete instructions on how to use the files in this repository.

## Run from IDE

Run each of the (3) service from a different terminal window.

```bash
export NLP_CLIENT_PORT=:8080
export RAKE_PORT=:8081
export PROSE_PORT=:8082
export RACK_ENDPOINT=http://localhost:8081
export PROSE_ENDPOINT=http://localhost:8082
export AUTH_KEY=SuP3r5eCRetAutHK3y

go run *.go
```

## Build Required Images for Docker

```bash
# change me
export ISV_ACCOUNT=01234567890
export ISV_ECR_REGION=us-east-2

$(aws ecr get-login --no-include-email --region us-east-2)
docker build -t ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/rake-app:1.0.0 . --no-cache
docker push ${ISV_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/rake-app:1.0.0

# change me
export CUSTOMER_ACCOUNT=09876543210
export CUSTOMER_ECR_REGION=us-west-2

$(aws ecr get-login --no-include-email --region ${CUSTOMER_ECR_REGION})
docker build -t ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/nlp-client:1.0.0 . --no-cache
docker push ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/nlp-client:1.0.0

docker build -t ${CUSTOMER_ACCOUNT}.dkr.ecr.${ISV_ECR_REGION}.amazonaws.com/prose-app:1.0.0 . --no-cache
docker push ${CUSTOMER_ACCOUNT}.dkr.ecr.${CUSTOMER_ECR_REGION}.amazonaws.com/prose-app:1.0.0

# display images
docker image ls --filter=reference='*amazonaws.com/*'
```
Runing (3) service stack from Docker Swarm.

```bash
export NLP_CLIENT_PORT=:8080
export RAKE_PORT=:8080
export PROSE_PORT=:8080
export RACK_ENDPOINT=http://rake-app:8080
export PROSE_ENDPOINT=http://prose-app:8080
export AUTH_KEY=SuP3r5eCRetAutHK3y

docker stack deploy --compose-file stack.yml nlp

# display containers
docker stack ps nlp --no-trunc
docker container ls

# delete stack
docker stack rm nlp
```

Sample output from Docker Swarm stack deployment.

```text
Admin-01:~/environment/ecr-cross-accnt-demo (master) $ docker stack deploy --compose-file stack.yml nlp
Creating network nlp_nlp-demo
Creating service nlp_nlp-client
Creating service nlp_rake-app
Creating service nlp_prose-app

Admin-01:~/environment/ecr-cross-accnt-demo (master) $ docker container ls
CONTAINER ID        IMAGE                                                           COMMAND             CREATED             STATUS              PORTS               NAMES
ac5501bb9a79        864887685992.dkr.ecr.us-east-2.amazonaws.com/rake-app:1.0.0     "/go/bin/app"       14 seconds ago      Up 13 seconds                           nlp_rake-app.1.jpctxbvzhcseo8uwuldwlp7hp
7dc171f89f9f        378568318651.dkr.ecr.us-west-2.amazonaws.com/nlp-client:1.0.0   "/go/bin/app"       15 seconds ago      Up 12 seconds                           nlp_nlp-client.1.t96hg46g76uwsvr7i6bweluxz
7ae5369d4293        378568318651.dkr.ecr.us-west-2.amazonaws.com/prose-app:1.0.0    "/go/bin/app"       15 seconds ago      Up 13 seconds                           nlp_prose-app.1.6wkb8x6slva7t253ksucfshyu
Admin-01:~/environment/ecr-cross-accnt-demo (master) $ 
```
