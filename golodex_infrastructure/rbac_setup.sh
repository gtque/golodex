#!/bin/bash
echo "create service account"
kubectl apply -f - <<END
apiVersion: v1
kind: ServiceAccount
metadata:
  name: "golodex"
  namespace: "golodex"
END

kubectl apply -f - <<END
apiVersion: v1
kind: ServiceAccount
metadata:
  name: "golodex"
  namespace: "golodex"
END

kubectl apply -f - <<END
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: "golodex"
rules:
- apiGroups: [""]
  resources:
  - nodes
  - persistentvolumes
  - namespaces
  verbs: ["list", "watch"]
- apiGroups: ["storage.k8s.io"]
  resources:
  - storageclasses
  verbs: ["list", "watch"]
END

kubectl apply -f - <<END
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: "golodex-couchdb"
  namespace: "golodex"
  labels:
    app.kubernetes.io/name: "golodex"
    app.kubernetes.io/managed-by: "bash"
    app.kubernetes.io/version: "1.0"
    app.kubernetes.io/component: "couchdb"
    app.kubernetes.io/part-of: "golodex"
spec:
  rules:
  - host: golodex-couchdb.thehangingpen.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: golodex-couchdb
            port:
             number: 443
  tls:
  - hosts:
    - golodex-couchdb.thehangingpen.com
    secretName: golodex-couchdb-cert
END