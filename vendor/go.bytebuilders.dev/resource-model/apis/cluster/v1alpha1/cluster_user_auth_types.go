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
// +kubebuilder:resource:path=clusteruserauths,singular=clusteruserauth,scope=Cluster,categories={kubernetes,resource-model,appscode}
// +kubebuilder:subresource:status
type ClusterUserAuth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ClusterUserAuthSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            ClusterUserAuthStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ClusterUserAuthSpec struct {
	ClusterUID string `json:"clusterUID" protobuf:"bytes,1,opt,name=clusterUID"`
	UserID     int64  `json:"userID" protobuf:"bytes,2,opt,name=userID"`

	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS.
	// +optional
	ClientCertificateData []byte `json:"clientCertificateData,omitempty" protobuf:"bytes,3,opt,name=clientCertificateData"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS.
	// +optional
	ClientKeyData []byte `json:"clientKeyData,omitempty" protobuf:"bytes,4,opt,name=clientKeyData"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	// +optional
	Token string `json:"token,omitempty" protobuf:"bytes,5,opt,name=token"`
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,6,opt,name=username"`
	// Password is the password for basic authentication to the kubernetes cluster.
	// +optional
	Password string `json:"password,omitempty" protobuf:"bytes,7,opt,name=password"`

	// Impersonate is the username to act-as.
	// +optional
	Impersonate string `json:"impersonate,omitempty" protobuf:"bytes,8,opt,name=impersonate"`
	// ImpersonateGroups is the groups to impersonate.
	// +optional
	ImpersonateGroups []string `json:"impersonateGroups,omitempty" protobuf:"bytes,9,rep,name=impersonateGroups"`
	// ImpersonateUserExtra contains additional information for impersonated user.
	// +optional
	ImpersonateUserExtra map[string]ExtraValue `json:"impersonateUserExtra,omitempty" protobuf:"bytes,10,rep,name=impersonateUserExtra"`

	// AuthProvider specifies a custom authentication plugin for the kubernetes cluster.
	// +optional
	AuthProvider *AuthProviderConfig `json:"authProvider,omitempty" protobuf:"bytes,11,opt,name=authProvider"`

	// Provider Access Token params
	// +optional
	Provider    TokenProviderName    `json:"provider,omitempty" protobuf:"bytes,12,opt,name=provider,casttype=TokenProviderName"`
	GoogleOAuth *GoogleOAuthProvider `json:"googleOAuth,omitempty" protobuf:"bytes,13,opt,name=googleOAuth"`
	AWS         *AWSProvider         `json:"aws,omitempty" protobuf:"bytes,14,opt,name=aws"`
}

type GoogleOAuthProvider struct {
	ClientID     string `json:"clientID" protobuf:"bytes,1,opt,name=clientID"`
	ClientSecret string `json:"clientSecret" protobuf:"bytes,2,opt,name=clientSecret"`
	AccessToken  string `json:"accessToken" protobuf:"bytes,3,opt,name=accessToken"`
	// +optional
	RefreshToken string `json:"refreshToken,omitempty" protobuf:"bytes,4,opt,name=refreshToken"`
	// +optional
	Expiry int64 `json:"expiry,omitempty" protobuf:"bytes,5,opt,name=expiry"`
}

type AWSProvider struct {
	// +optional
	Region string `json:"region" protobuf:"bytes,1,opt,name=region"`
	// +optional
	ClusterID string `json:"clusterID" protobuf:"bytes,2,opt,name=clusterID"`
	// +optional
	AssumeRoleARN string `json:"assumeRoleArn" protobuf:"bytes,3,opt,name=assumeRoleArn"`
	// +optional
	AssumeRoleExternalID string `json:"assumeRoleExternalID" protobuf:"bytes,4,opt,name=assumeRoleExternalID"`
	// +optional
	SessionName string `json:"sessionName" protobuf:"bytes,5,opt,name=sessionName"`
	// +optional
	ForwardSessionName bool `json:"forwardSessionName" protobuf:"bytes,6,opt,name=forwardSessionName"`
	// +optional
	Cache bool `json:"cache" protobuf:"bytes,7,opt,name=cache"`

	// Temporary Token for 15 mins only, if expired or not set create a new one
	// +optional
	Token string `json:"token,omitempty" protobuf:"bytes,8,opt,name=token"`
	// +optional
	Expiration int64 `json:"expiration,omitempty" protobuf:"varint,9,opt,name=expiration"`

	AccessKeyID     string `json:"accessKeyID" protobuf:"bytes,10,opt,name=accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey" protobuf:"bytes,11,opt,name=secretAccessKey"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type ClusterUserAuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterUserAuth `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

type ClusterUserAuthStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`
}
