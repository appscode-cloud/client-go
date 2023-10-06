/*
Copyright 2020 AppsCode Inc.

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
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ClusterAuthInfoTemplateSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            ClusterAuthInfoTemplateStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ClusterAuthInfoTemplateSpec struct {
	UID     string `json:"uid" protobuf:"bytes,1,opt,name=uid"`
	OwnerID int64  `json:"ownerID" protobuf:"bytes,2,opt,name=ownerID"`

	// CertificateAuthorityData contains PEM-encoded certificate authority certificates.
	// +optional
	CertificateAuthorityData []byte `json:"certificateAuthorityData,omitempty" protobuf:"bytes,3,opt,name=certificateAuthorityData"`

	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS.
	// +optional
	ClientCertificateData []byte `json:"clientCertificateData,omitempty" protobuf:"bytes,4,opt,name=clientCertificateData"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS.
	// +optional
	ClientKeyData []byte `json:"clientKeyData,omitempty" protobuf:"bytes,5,opt,name=clientKeyData"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	// +optional
	Token string `json:"token,omitempty" protobuf:"bytes,6,opt,name=token"`
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,7,opt,name=username"`
	// Password is the password for basic authentication to the kubernetes cluster.
	// +optional
	Password string `json:"password,omitempty" protobuf:"bytes,8,opt,name=password"`

	// Impersonate is the username to act-as.
	// +optional
	Impersonate string `json:"impersonate,omitempty" protobuf:"bytes,9,opt,name=impersonate"`
	// ImpersonateGroups is the groups to impersonate.
	// +optional
	ImpersonateGroups []string `json:"impersonateGroups,omitempty" protobuf:"bytes,10,rep,name=impersonateGroups"`
	// ImpersonateUserExtra contains additional information for impersonated user.
	// +optional
	ImpersonateUserExtra map[string]ExtraValue `json:"impersonateUserExtra,omitempty" protobuf:"bytes,11,rep,name=impersonateUserExtra"`

	// AuthProvider specifies a custom authentication plugin for the kubernetes cluster.
	// +optional
	AuthProvider *AuthProviderConfig `json:"authProvider,omitempty" protobuf:"bytes,13,opt,name=authProvider"`
}

// AuthProviderConfig holds the configuration for a specified auth provider.
type AuthProviderConfig struct {
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// +optional
	Config map[string]string `json:"config,omitempty" protobuf:"bytes,2,rep,name=config"`
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
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterAuthInfoTemplate `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

type ClusterAuthInfoTemplateStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`
}
