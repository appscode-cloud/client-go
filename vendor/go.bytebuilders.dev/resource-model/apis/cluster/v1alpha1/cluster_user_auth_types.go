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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindClusterUserAuth = "ClusterUserAuth"
	ResourceClusterUserAuth     = "clusteruserauth"
	ResourceClusterUserAuths    = "clusteruserauths"
)

type TokenProviderName string

const (
	TokenProviderGoogle TokenProviderName = "gcp"
	TokenProviderAWS    TokenProviderName = "aws"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusteruserauths,singular=clusteruserauth,shortName=uauth,scope=Cluster,categories={kubernetes,resource-model,appscode}
// +kubebuilder:subresource:status
type ClusterUserAuth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterUserAuthSpec   `json:"spec,omitempty"`
	Status            ClusterUserAuthStatus `json:"status,omitempty"`
}

type ClusterUserAuthSpec struct {
	ClusterUID string `json:"clusterUID"`
	UserID     int64  `json:"userID"`

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

	// Provider Access Token params
	// +optional
	Provider       TokenProviderName      `json:"provider,omitempty"`
	GoogleOAuth    *GoogleOAuthProvider   `json:"googleOAuth,omitempty"`
	GoogleCloud    *GoogleCloudCredential `json:"googleCloud,omitempty"`
	AWS            *AWSProvider           `json:"aws,omitempty"`
	CredentialName string                 `json:"credentialName,omitempty"`
}

type GoogleCloudCredential struct {
	ProjectID      string `json:"projectID"`
	ServiceAccount string `json:"serviceAccount"`
}

type GoogleOAuthProvider struct {
	AccessToken string `json:"accessToken"`
	// +optional
	RefreshToken string `json:"refreshToken,omitempty"`
	// +optional
	Expiry int64 `json:"expiry,omitempty"`
}

type AWSProvider struct {
	// +optional
	Region string `json:"region"`
	// +optional
	ClusterID string `json:"clusterID"`
	// +optional
	AssumeRoleARN string `json:"assumeRoleArn"`
	// +optional
	AssumeRoleExternalID string `json:"assumeRoleExternalID"`
	// +optional
	SessionName string `json:"sessionName"`
	// +optional
	ForwardSessionName bool `json:"forwardSessionName"`
	// +optional
	Cache bool `json:"cache"`

	// Temporary Token for 15 mins only, if expired or not set create a new one
	// +optional
	Token string `json:"token,omitempty"`
	// +optional
	Expiration int64 `json:"expiration,omitempty"`

	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type ClusterUserAuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterUserAuth `json:"items,omitempty"`
}

type ClusterUserAuthStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
