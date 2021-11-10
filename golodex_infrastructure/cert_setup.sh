#!/bin/bash
echo "setting up the cert-manager"
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
sleep 15
echo "waiting for the cert-manager to be running"
p=0
export timedOut="false"
export notRunning="true"
while [[ $notRunning == "true" ]]
do
  kubectl -n cert-manager get pods
  export pending=$(kubectl -n cert-manager get pods | grep -E 'Pending|ContainerCreating|Failed|Crash|0/1')
  echo "pending: $pending"
  export running=$(kubectl -n cert-manager get pods | grep -E 'Running')
  if [[ -z "$running" ]]; then
    echo "cert manager not running yet"
  else
    if [[ -z "$pending" ]]; then
      echo "cert manager appears to be up and running"
      export notRunning="false"
    else
      echo "some of the pods are running, but some are not"
    fi
  fi
  ((p++))
  if [[ $p -gt 8 ]]; then
    echo "timed out waiting for cert-manager to start"
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
echo "just a bit more sleep just to make sure..."
sleep 10
echo "creating local self signed tls cert/secret, because this must exist in the namespace cert is being created."
cat <<EOF > local-selfsigned-issuer.yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned-issuer
  namespace: golodex
spec:
  selfSigned: {}
EOF

cat <<EOF > local-selfsigned-cert.yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: golodex-couchdb-cert
  namespace: golodex
spec:
  commonName: '*.thehangingpen.com'
  secretName: golodex-couchdb-cert
  issuerRef:
    name: selfsigned-issuer
    kind: Issuer
  duration: 2160h0m0s
  renewBefore: 360h0m0s
  dnsNames:
  - '*.thehangingpen.com'
  - thehangingpen.com
  - localhost
EOF
echo "applying the issuer and the cert..."
kubectl apply -f local-selfsigned-issuer.yaml
sleep 5
kubectl apply -f local-selfsigned-cert.yaml
export found="false"
i=0
while [[ $found == "false" ]]
do
  export found="true"
  kubectl -n golodex get secret golodex-couchdb-cert >/dev/null 2>&1 || { \
      echo "no secret yet" & export found="false" ; }
  ((i++))
  if [[ $i -gt 8 ]]; then
    echo "timed out waiting for secret"
    export found="true"
    export timedOut="true"
  fi
  if [[ $found == "false" ]]; then
    echo "sleeping for a bit"
    sleep 15
  else
    echo "I was true...$found"
  fi
done
echo "cleaning templates"
rm local-*.yaml
if [[ $timedOut == "true" ]]; then
  echo "something went wrong, the secret never showed up, so I am bailing on the rest right now."
  exit 1
fi

