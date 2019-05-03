#!/bin/bash

cat <<'EOF' | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: hyper
spec:
  automountServiceAccountToken: false
  containers:
  - name: hyper
    image: k8s.gcr.io/hyperkube:v1.13.5
    env:
    - name: KUBECONFIG
      value: /etc/kubeconfig/kubeconfig
    command:
    - "/bin/bash"
    - "-c"
    - "--"
    args:
    - "while true; do sleep 30; done"
    volumeMounts:
    - name: kubeconfig
      mountPath: /etc/kubeconfig
  volumes:
  - name: kubeconfig
    secret:
      secretName: kubeconfig
EOF

while [[ "$(kubectl get pod hyper -o 'jsonpath={.status.conditions[?(@.type=="Ready")].status}')" != 'True' ]]; do
  echo "Pod is not ready, waiting"
  sleep 1
done

kubectl exec -it hyper -- bash

kubectl delete pod hyper
