#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# ref: https://github.com/bitnami/charts/tree/master/bitnami/kafka (note: hostpath provisioner is not working)
function install_kafka() {
  helm repo add bitnami https://charts.bitnami.com/bitnami
  helm install kafka bitnami/kafka --set persistence.enabled=false,zookeeper.persistence.enabled=false "${@}"
}

function uninstall_kafka() {
  helm uninstall kafka "${@}"

  mapfile -t pvcs < <(kubectl get persistentvolumeclaims -o json | jq -r '.items[].metadata.name' | grep -E 'kafka|zookeeper')
  for pvc in "${pvcs[@]}"; do
    kubectl delete pvc "$pvc"
  done
}

# ref: https://github.com/helm/charts/tree/master/stable/rabbitmq
function install_amqp() {
  helm install rabbitmq stable/rabbitmq "${@}"
}

function uninstall_amqp() {
  helm uninstall rabbitmq "${@}"
}

if [[ $# -lt 2 ]]; then
  cat <<EOF
Usage:
  install|uninstall kafka|aqmp

Note:
  Make sure you have a running K8s cluster w/ default storage class installed.
EOF
  exit 1
fi

cmd=$1
service=$2
shift 2

"$cmd"_"$service" "${@}"
