#!/bin/bash

REPO_ROOT="$(git rev-parse --show-toplevel)"

cat <<EOT | ytt -f ${REPO_ROOT}/knative/ksvc-template.yaml --data-value name=$1 | kubectl apply -f -
#@ load("@ytt:data", "data")

apiVersion: serving.knative.dev/v1
kind: Service
metadata:
 name: #@ data.values.name
 namespace: default
spec:
 template:
  spec:
   containers:
    - image: gcr.io/knative-samples/helloworld-go
      env:
        - name: TARGET
          value: #@ "{}".format(data.values.name)

EOT
