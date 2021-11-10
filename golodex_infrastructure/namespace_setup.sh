#!/bin/bash
echo "golodex namespace"

kubectl apply -f - <<END
apiVersion: v1
kind: Namespace
metadata:
  name: "golodex"
  labels:
    kubernetes.io/metadata.name: "golodex"
END