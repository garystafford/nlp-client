#!/bin/bash
#
# author: Gary A. Stafford
# site: https://programmaticponderings.com
# license: MIT License
# purpose: Build Go microservices for nlp demo
# date: 2021-06-14

readonly projectRoot="/Users/garystafford/Documents/projects"
readonly -a projects=(nlp-client rake-app prose-app lang-app dynamo-app)
readonly tag=1.2.1

pushd $projectRoot || exit

for project in "${projects[@]}"
do
  pushd "$project" || exit
  docker build -t "garystafford/$project:$tag" . --no-cache
  docker push "garystafford/$project:$tag"
  popd || exit
done