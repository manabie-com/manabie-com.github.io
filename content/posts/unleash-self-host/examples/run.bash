#!/bin/bash

git clone https://github.com/manabie-com/manabie-com.github.io
cd manabie-com.github.io/content/posts/unleash-self-host/examples

minikube start
helm upgrade --install unleash ./ -f values.yaml