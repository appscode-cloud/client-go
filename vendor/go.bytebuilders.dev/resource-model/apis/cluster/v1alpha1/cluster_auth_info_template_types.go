/*
Copyright AppsCode Inc. and Contributors

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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	ResourceKindClusterAuthInfoTemplate = "ClusterAuthInfoTemplate"
	ResourceClusterAuthInfoTemplate     = "clusterauthinfotemplate"
	ResourceClusterAuthInfoTemplates    = "clusterauthinfotemplates"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusterauthinfotemplates,singular=clusterauthinfotemplate,shortName=cauth,scope=Cluster,categories={kubernetes,resource-model,appscode}
// +kubebuilder:subresource:status
type ClusterAuthInfoTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterAuthInfoTemplateSpec   `json:"spec,omitempty"`
	Status            ClusterAuthInfoTemplateStatus `json:"status,omitempty"`
}

type ClusterAuthInfoTemplateSpec struct {
	UID     string `json:"uid"`
	OwnerID int64  `json:"ownerID"`

	// CertificateAuthorityData contains PEM-encoded certificate authority certificates.
	// +optional
	CertificateAuthorityData []byte `json:"certificateAuthorityData,omitempty"`

	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS.
	// +optional
	ClientCertificateData []byte `json:"clientCertificateData,omitempty"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS.
	// +optional
	ClientKeyData []byte `json:"clientKeyData,omitempty"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	// +optional
	Token string `json:"token,omitempty"`
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty"`
	// Password is the password for basic authentication to the kubernetes cluster.
	// +optional
	Password string `json:"password,omitempty"`

	// Impersonate is the username to act-as.
	// +optional
	Impersonate string `json:"impersonate,omitempty"`
	// ImpersonateGroups is the groups to impersonate.
	// +optional
	ImpersonateGroups []string `json:"impersonateGroups,omitempty"`
	// ImpersonateUserExtra contains additional information for impersonated user.
	// +optional
	ImpersonateUserExtra map[string]ExtraValue `json:"impersonateUserExtra,omitempty"`

	// AuthProvider specifies a custom authentication plugin for the kubernetes cluster.
	// +optional
	AuthProvider *AuthProviderConfig `json:"authProvider,omitempty"`
}

// AuthProviderConfig holds the configuration for a specified auth provider.
type AuthProviderConfig struct {
	Name string `json:"name"`
	// +optional
	Config map[string]string `json:"config,omitempty"`
}

func (auth *AuthProviderConfig) APIFormat() *api.AuthProviderConfig {
	if auth == nil {
		return nil
	}

	return &api.AuthProviderConfig{
		Name:   auth.Name,
		Config: auth.Config,
	}
}

func ToAuthProviderConfig(auth *api.AuthProviderConfig) *AuthProviderConfig {
	if auth == nil {
		return nil
	}

	return &AuthProviderConfig{
		Name:   auth.Name,
		Config: auth.Config,
	}
}

// ExtraValue masks the value so protobuf can generate
// +protobuf.nullable=true
// +protobuf.options.(gogoproto.goproto_stringer)=false
type ExtraValue []string

func (t ExtraValue) String() string {
	return fmt.Sprintf("%v", []string(t))
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type ClusterAuthInfoTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAuthInfoTemplate `json:"items,omitempty"`
}

type ClusterAuthInfoTemplateStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
