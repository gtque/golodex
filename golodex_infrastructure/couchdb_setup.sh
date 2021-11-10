#!/bin/bash
echo "in microk8s, rbac must be enabled to install olm"
microk8s enable rbac
sleep 15
echo "installing olm"
curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.18.3/install.sh | bash -s v0.18.3
echo "installing operator"
kubectl create -f https://operatorhub.io/install/couchdb-operator.yaml
echo "checking operators"
kubectl get csv -n operators

echo "deploying couchdb"
echo "but rbac has to be disabled, at least until a proper service account is defined."
microk8s disable rbac
kubectl apply -f - <<END
apiVersion: couchdb.databases.cloud.ibm.com/v1
kind: CouchDBCluster
metadata:
  name: golodex-couchdb
  namespace: golodex
spec:
  size: 1
  disk: 1Gi
  memory: 1Gi
  cpu: "1"
  storageClass: "microk8s-hostpath"
  environment:
    adminPassword: "Asupe7S3cretP@ssw0rd!"
    tls: false
    couch_peruser_enable: true
END

echo "waiting for couchdb to be running"
p=0
export timedOut="false"
export notRunning="true"
while [[ $notRunning == "true" ]]
do
  kubectl -n golodex get pods
  export pending=$(kubectl -n golodex get pods | grep -E 'Pending|ContainerCreating|Failed|Crash|0/1')
  echo "pending: $pending"
  export running=$(kubectl -n golodex get pods | grep -E '2/2')
  if [[ -z "$running" ]]; then
    echo "couchdb not running yet"
  else
    if [[ -z "$pending" ]]; then
      echo "couchdb appears to be up and running"
      export notRunning="false"
    else
      echo "some of the pods are running, but some are not"
    fi
  fi
  ((p++))
  if [[ $p -gt 80 ]]; then
    echo "timed out waiting for couchdb to start"
    export notRunning="false"
    export timedOut="true"
  else
    if [[ $notRunning == "true" ]]; then
      echo "sleeping a bit before checking again."
      sleep 15
    fi
  fi
done
if [[ $timedOut == "true" ]]; then
  echo "something went wrong, the cert-manager never started, so I am bailing on the rest right now."
  exit 1
fi
kubectl -n golodex exec -i c-golodex-couchdb-m-0 -c db -- curl -X PUT http://$COUCHDB_ADMIN_USER:$COUCHDB_ADMIN_PW@127.0.0.1:5984/_node/_local/_config/couch_peruser/enable -d '"true"'
kubectl -n golodex exec -i c-golodex-couchdb-m-0 -c db -- curl -X PUT http://$COUCHDB_ADMIN_USER:$COUCHDB_ADMIN_PW@127.0.0.1:5984/_node/_local/_config/couch_peruser/delete_dbs -d '"true"'
echo "to access outside of the cluster, you will want to run:"
echo "with tls: kubectl port-forward svc/golodex-couchdb 8443:443 -n golodex"
echo "without tls: kubectl port-forward svc/golodex-couchdb 5984:5984 -n golodex"
