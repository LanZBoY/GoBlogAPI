#!/bin/bash
image_path=go-blog-api:dev-$(date -u +"%Y%m%dT%H%M%SZ")
docker build -t $image_path .
minikube image load $image_path