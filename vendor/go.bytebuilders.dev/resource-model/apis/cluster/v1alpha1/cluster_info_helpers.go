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
	"strconv"

	"go.bytebuilders.dev/resource-model/apis/cluster"
	"go.bytebuilders.dev/resource-model/crds"

	"k8s.io/apimachinery/pkg/labels"
	"kmodules.xyz/client-go/apiextensions"
)

type ClusterOptions struct {
	ResourceName    string   `protobuf:"bytes,1,opt,name=resourceName"`
	Provider        string   `protobuf:"bytes,2,opt,name=provider"`
	UserID          int64    `protobuf:"varint,3,opt,name=userID"`
	CID             string   `protobuf:"bytes,4,opt,name=cID"`
	OwnerID         int64    `protobuf:"varint,5,opt,name=ownerID"`
	ImportType      string   `protobuf:"bytes,6,opt,name=importType"`
	ExternalID      string   `protobuf:"bytes,7,opt,name=externalID"`
	ClusterManagers []string `protobuf:"bytes,9,opt,name=clusterManagers"`
}

func (_ ClusterInfo) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceClusterInfos))
}

func (clusterInfo *ClusterInfo) ApplyLabels(opts ClusterOptions) {
	labelMap := map[string]string{
		cluster.LabelResourceName:      opts.ResourceName,
		cluster.LabelClusterUID:        opts.CID,
		cluster.LabelClusterOwnerID:    strconv.FormatInt(opts.OwnerID, 10),
		cluster.LabelClusterProvider:   opts.Provider,
		cluster.LabelClusterImportType: opts.ImportType,
		cluster.LabelClusterExternalID: opts.ExternalID,
	}

	clusterInfo.ObjectMeta.SetLabels(labelMap)
}

func (clusterInfo *ClusterInfo) ApplyClusterManagerLabels(opts ClusterOptions) {
	labelMap := clusterInfo.Labels

	setManagerLabels(opts, labelMap)
	clusterInfo.ObjectMeta.SetLabels(labelMap)
}

func (_ ClusterInfo) FormatLabels(opts ClusterOptions) labels.Selector {
	labelMap := make(map[string]string)
	if opts.ResourceName != "" {
		labelMap[cluster.LabelResourceName] = opts.ResourceName
	}
	if opts.CID != "" {
		labelMap[cluster.LabelClusterUID] = opts.CID
	}
	if opts.OwnerID != 0 {
		labelMap[cluster.LabelClusterOwnerID] = strconv.FormatInt(opts.OwnerID, 10)
	}
	if opts.Provider != "" {
		labelMap[cluster.LabelClusterProvider] = opts.Provider
	}
	if opts.ImportType != "" {
		labelMap[cluster.LabelClusterImportType] = opts.ImportType
	}
	if opts.ExternalID != "" {
		labelMap[cluster.LabelClusterExternalID] = opts.ExternalID
	}

	setManagerLabels(opts, labelMap)
	return labels.SelectorFromSet(labelMap)
}

func setManagerLabels(opts ClusterOptions, labelMap map[string]string) {
	for _, mng := range opts.ClusterManagers {
		switch mng {
		case "ACE":
			labelMap[cluster.LabelClusterManagerACE] = "true"
		case "OCMHub":
			labelMap[cluster.LabelClusterManagerOCMHub] = "true"
		case "OCMSpoke":
			labelMap[cluster.LabelClusterManagerOCMSpoke] = "true"
		case "OCMMulticlusterControlplane":
			labelMap[cluster.LabelClusterManagerOCMMulticlusterControlplane] = "true"
		case "Rancher":
			labelMap[cluster.LabelClusterManagerRancher] = "true"
		case "OpenShift":
			labelMap[cluster.LabelClusterManagerOpenShift] = "true"
		case "vcluster":
			labelMap[cluster.LabelClusterManagerVCluster] = "true"
		}
	}
}
