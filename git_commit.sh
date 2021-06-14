#!/bin/bash
#
# author: Gary A. Stafford
# site: https://programmaticponderings.com
# license: MIT License
# purpose: Verified commits and pushes to GitHub
# date: 2021-06-14

readonly projectRoot="/Users/garystafford/Documents/projects"
readonly -a projects=(nlp-client rake-app prose-app lang-app dynamo-app)
readonly commitMessage="Updating logging and main/run functions"

pushd $projectRoot || exit

for project in "${projects[@]}"
do
  pushd "$project" || exit
    gofmt main.go
    
    git add -A
    git commit --gpg-sign=DCD0DCF7AE909D61 -m "$commitMessage"
    git push
  popd || exit
done