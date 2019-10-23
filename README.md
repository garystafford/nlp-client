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
