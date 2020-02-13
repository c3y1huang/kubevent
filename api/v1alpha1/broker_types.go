/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BrokerSpec defines the desired state of Broker
type BrokerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Enum=kafka;amqp
	Provider          string              `json:"provider"` //TODO need to have validation for values in (kafka)
	Activate          bool                `json:"activate"`
	Kafka             KafkaBroker         `json:"kafka,omitempty"`
	AMQP              AMQPBroker          `json:"amqp,omitempty"`
	GroupVersionKinds []*GroupVersionKind `json:"resources,omitempty"`
}

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

//KafkaBroker defines the Kafka broker server info
type KafkaBroker struct {
	Addresses []string `json:"addresses"`
	Topic     string   `json:"topic"`
	// +optional
	TlsConfig *TlsConfig `json:"tls_config,omitempty"`
}

//AMQPBroker defines the AMQP broker server info
type AMQPBroker struct {
	Addresses []string `json:"addresses"`
	Exchange  string   `json:"exchange"`
	// +optional
	TlsConfig *TlsConfig `json:"tls_config,omitempty"`
}

//TlsConfig defines the TLS configurations
type TlsConfig struct {
	Insecure   bool   `json:"insecure,omitempty"`
	CACert     string `json:"ca_cert,omitempty"`
	ClientCert string `json:"client_cert,omitempty"`
	ClientKey  string `json:"client_key,omitempty"`
}

// BrokerStatus defines the observed state of Broker
type BrokerStatus struct {
	Provider string `json:"provider"` //TODO need to have validation for values in (kafka)
	Name     string `json:"name"`
	Active   bool   `json:"active"`
}

// +kubebuilder:object:root=true

// Broker is the Schema for the brokers API
type Broker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerSpec   `json:"spec,omitempty"`
	Status BrokerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BrokerList contains a list of Broker
type BrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Broker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Broker{}, &BrokerList{})
}
