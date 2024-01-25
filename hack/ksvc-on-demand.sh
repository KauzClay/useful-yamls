#!/bin/bash

cat <<EOT | ytt -f - --data-value name="$1" --data-value namespace="${2:-default}" |  kubectl apply -f -
#@ load("@ytt:data", "data")

apiVersion: serving.knative.dev/v1
kind: Service
metadata:
 name: #@ data.values.name
 namespace: #@ data.values.namespace
spec:
 template:
  metadata:
    annotations:
      autoscaling.knative.dev/minScale: "1"
  spec:
   containers:
    - image: gcr.io/knative-samples/helloworld-go
      env:
        - name: TARGET
          value: #@ "{}".format(data.values.name)

EOT
