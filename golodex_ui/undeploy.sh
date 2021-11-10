#!/bin/bash

printf "deleting objects...\n"
kubectl --v=0 -n golodex delete secret,configmap,service,deployment,ingress,role,serviceaccount,rolebinding -l app.kubernetes.io/name=golodex-ui
sleep 1