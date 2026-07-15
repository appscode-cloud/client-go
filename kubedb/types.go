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

package kubedb

import "encoding/json"

// OptionsModelRequest is the body for PUT /helm/options/model — the small set of
// options the KubeDB editor expands into a full CR + companion objects.
type OptionsModelRequest struct {
	Metadata OptionsMetadata `json:"metadata"`
	Spec     map[string]any  `json:"spec,omitempty"`
}

// OptionsMetadata identifies the target resource kind and release for the editor.
type OptionsMetadata struct {
	Resource ResourceRef `json:"resource"`
	Release  ReleaseRef  `json:"release"`
}

// ResourceRef is a KubeDB group/version/kind reference.
type ResourceRef struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

// ReleaseRef names the installation (database) the editor operates on.
type ReleaseRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// EditorModel is the opaque model returned by options/model and consumed, unchanged,
// by the editor apply call.
type EditorModel = json.RawMessage

// DeleteEditorRequest is the body for DELETE /helm/editor/.
type DeleteEditorRequest struct {
	Metadata OptionsMetadata `json:"metadata"`
}

// TaskResponse is the async acknowledgement returned by editor apply/delete.
type TaskResponse struct {
	ID     string `json:"id,omitempty"`
	Task   string `json:"task,omitempty"`
	Status string `json:"status,omitempty"`
}

// Unstructured is a minimal view of a Kubernetes object — enough to read status and
// identify it. The full payload is not retained.
type Unstructured struct {
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Metadata   ObjectMeta     `json:"metadata"`
	Status     ResourceStatus `json:"status"`
}

// ObjectMeta is a minimal Kubernetes object metadata view.
type ObjectMeta struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// ResourceStatus captures the KubeDB database phase (e.g. "Ready", "Provisioning").
type ResourceStatus struct {
	Phase string `json:"phase"`
}

// Secret is a minimal core/v1 Secret view (base64-encoded data values).
type Secret struct {
	Data map[string]string `json:"data"`
}
