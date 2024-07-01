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

package cluster

type BasicInfo struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	OwnerID      int64  `json:"ownerID,omitempty"`
	UserID       int64  `json:"userID,omitempty"`
	ClusterUID   string `json:"clusterUID,omitempty"`
	HubClusterID string `json:"hubClusterID,omitempty"`
}

type ProviderOptions struct {
	Name          string `json:"name"`
	Credential    string `json:"credential,omitempty"`
	ClusterID     string `json:"clusterID,omitempty"`
	Project       string `json:"project,omitempty"`
	Region        string `json:"region,omitempty"`
	ResourceGroup string `json:"resourceGroup,omitempty"`
	KubeConfig    string `json:"kubeConfig,omitempty"`
}

type ComponentOptions struct {
	FluxCD         bool         `json:"fluxCD,omitempty"`
	FeatureSets    []FeatureSet `json:"featureSets,omitempty"`
	AllFeatures    bool         `json:"allFeatures,omitempty"`
	ClusterProfile string       `json:"clusterProfile,omitempty"`
}

type FeatureSet struct {
	Name     string   `json:"name"`
	Features []string `json:"features,omitempty"`
}

type ListOptions struct {
	Provider string `json:"provider,omitempty"`
}

type GetOptions struct {
	Name string `json:"name,omitempty"`
}

type CheckOptions struct {
	BasicInfo BasicInfo       `json:"basicInfo,omitempty"`
	Provider  ProviderOptions `json:"provider"`
}

type ImportOptions struct {
	BasicInfo  BasicInfo        `json:"basicInfo,omitempty"`
	Provider   ProviderOptions  `json:"provider,omitempty"`
	Components ComponentOptions `json:"components,omitempty"`
}

type ConnectOptions struct {
	Name       string `json:"name"`
	Credential string `json:"credential,omitempty"`
	KubeConfig string `json:"kubeConfig,omitempty"`
}

type ReconfigureOptions struct {
	BasicInfo  BasicInfo        `json:"basicInfo"`
	Components ComponentOptions `json:"components,omitempty"`
}

type RemovalOptions struct {
	Name       string           `json:"name"`
	Components ComponentOptions `json:"components,omitempty"`
}
