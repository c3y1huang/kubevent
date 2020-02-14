#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CONTROL_PLANE_IP=${CONTROL_PLANE_IP:-}
KAFKA_NODE_PORT=${KAFKA_NODE_PORT:-30092}

# ref: https://github.com/bitnami/charts/tree/master/bitnami/kafka (note: hostpath provisioner is not working)
function install_kafka() {
  helm repo add bitnami https://charts.bitnami.com/bitnami
  helm install kafka bitnami/kafka \
    --set persistence.enabled=false \
    --set zookeeper.persistence.enabled=false \
    --set externalAccess.service.domain="${CONTROL_PLANE_IP}" \
    --set externalAccess.enabled=true \
    --set externalAccess.service.type=NodePort \
    --set externalAccess.service.nodePort=\{"${KAFKA_NODE_PORT}"\} \
    "${@}"
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

function show_menu() {
  cat <<EOF
Usage:
  install|uninstall kafka|amqp

Note:
  Make sure you have a running K8s cluster w/ default storage class installed.
EOF
}

if [[ $# -lt 2 ]]; then
  show_menu
  exit 1
fi

cmd=$1
service=$2
shift 2

case $cmd in
  "install"|"uninstall") ;;
  *) 
    show_menu
    exit 1
    ;;
esac

case $service in
  "kafka"|"amqp") ;;
  *) 
    show_menu
    exit 1
    ;;
esac

"$cmd"_"$service" "${@}"
