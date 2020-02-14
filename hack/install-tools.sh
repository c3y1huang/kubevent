#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

function install_kafka() {
  GO111MODULE=off go get -v github.com/Shopify/sarama/tools/...

  cat <<EOF
Usage:
  kafka-console-producer -brokers <ADDRESSES> -topic kubevent -value hellofromconsole --verbose
  kafka-console-consumer -brokers <ADDRESSES> -topic kubevent -verbose -offset oldest
EOF
}

function show_menu() {
  cat <<EOF
Usage:
  install kafka
EOF
}

if [[ $# -lt 2 ]]; then
  show_menu
  exit 1
fi

cmd=$1
tool=$2
shift 2

case $cmd in
  "install") ;;
  *) 
    show_menu
    exit 1
    ;;
esac

case $tool in
  "kafka") ;;
  *) 
    show_menu
    exit 1
    ;;
esac

"$cmd"_"$tool" "${@}"
