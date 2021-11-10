#!/bin/bash

./undeploy.sh

printf "Dry running install...\n"
kubectl --v=0 --validate=true --dry-run=client create -f ./kubernetes
sleep 1

printf "Dry run succeeded.\nInstalling manifests...\n"
kubectl --v=0 apply -f ./kubernetes
sleep 1