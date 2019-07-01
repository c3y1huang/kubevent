[![CircleCI](https://circleci.com/gh/innobead/kubevent.svg?style=svg)](https://circleci.com/gh/innobead/kubevent)

# Kubevent

Publish K8s events of builtin resource objects from K8s clusters to external event brokers.

## Goals

- Register resource events to subscribe
- Publish the registered resource events to external event brokers
- Support metrics of registered resource events and published events
- Support extensible event brokers

## Non-Goals
- Message queue implementation
- Event broker implementation


## Architecture

![Architecture](docs/arch.png)

### Kubevent Controller

**Kubevent Controller** is responsible for watching user registered resource events, managing **Kubevent Event Publisher** to serve publishing events with resource data to external connected event brokers.

### Kubevent Config

**Kubevent Config** configure registered events and event brokers to adopt in Kubevent controller.

### Kubevent Event Handler

**Kubevent Event Handler** is responsible for serving event publishing to external event brokers.


## Supported Items

### Resource
Builtin K8s API resource.

### Event

Resource operation event like create, delete, update or generic.

### Event Broker

- Apache Kafka
- NATS
- AMQP (Advanced Message Queuing Protocol)
- STOMP (Simple Text Oriented Messaging Protocol)
- MQTT (Message Queuing Telemetry Transport)
- STOMP over WebSocket 
- Webhook
- Cloud solutions
  - AWS MQ, SQS (SImple Queue Service), Kinesis
  - Azure Service Bus, Event Grid, Event Hubs
  - GCP Cloud Pub/Sub
