# Go Microservice: nlp-client

Go-based microservice, part of a set of (3) microservices for upcoming post.

## Run from IDE

Run each of the (3) service from a different terminal window

```bash
export NLP_CLIENT_PORT=:8080
export RAKE_PORT=:8081
export PROSE_PORT=:8082
export RACK_ENDPOINT=http://localhost:8081
export PROSE_ENDPOINT=http://localhost:8082
export AUTH_KEY=DqiSyCzJgUY9kbxWiF1QA7NY

go run *.go
```

Runing (3) service stack from Docker Swarm

```bash
export PORT=8080
export NLP_CLIENT_PORT=:${PORT}
export RAKE_PORT=:${PORT}
export PROSE_PORT=:${PORT}
export RACK_ENDPOINT=http://rake-app:${PORT}
export PROSE_ENDPOINT=http://prose-app:${PORT}
export AUTH_KEY=DqiSyCzJgUY9kbxWiF1QA7NY

docker stack deploy --compose-file stack.yml nlp
docker stack ps nlp --no-trunc
docker container ls

docker stack rm nlp

docker image ls
docker rmi $(docker images --filter "dangling=true" -q --no-trunc)
```

Sample output from Docker Swarm stack deployment

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