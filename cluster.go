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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	clustermodel "go.bytebuilders.dev/resource-model/apis/cluster"
	"go.bytebuilders.dev/resource-model/apis/cluster/v1alpha1"
)

func (c *Client) CheckClusterExistence(opts clustermodel.CheckOptions) (*v1alpha1.ClusterInfo, error) {
	org, err := c.getOrganization()
	if err != nil {
		return nil, err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s/check", org)

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider options. err: %w", err)
	}

	var cluster v1alpha1.ClusterInfo
	err = c.getParsedResponse(http.MethodPost, apiPath, jsonHeader, bytes.NewReader(body), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Client) ImportCluster(opts clustermodel.ImportOptions, responseID string) (*v1alpha1.ClusterInfo, error) {
	org, err := c.getOrganization()
	if err != nil {
		return nil, err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s/import", org)

	params := make([]queryParams, 0)
	if responseID != "" {
		params = append(params, queryParams{key: "response-id", value: responseID})
	}
	apiPath, err = setQueryParams(apiPath, params)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cluster basic info. err: %w", err)
	}

	var cluster v1alpha1.ClusterInfo
	err = c.getParsedResponse(http.MethodPost, apiPath, jsonHeader, bytes.NewReader(body), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Client) ListClusters(opts clustermodel.ListOptions) (*v1alpha1.ClusterInfoList, error) {
	org, err := c.getOrganization()
	if err != nil {
		return nil, err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s", org)

	params := make([]queryParams, 0)
	if opts.Provider != "" {
		params = append(params, queryParams{key: "provider", value: opts.Provider})
	}

	apiPath, err = setQueryParams(apiPath, params)
	if err != nil {
		return nil, err
	}
	var clusters v1alpha1.ClusterInfoList
	err = c.getParsedResponse(http.MethodGet, apiPath, jsonHeader, nil, &clusters)
	if err != nil {
		return nil, err
	}
	return &clusters, nil
}

func (c *Client) GetCluster(opts clustermodel.GetOptions) (*v1alpha1.ClusterInfo, error) {
	org, err := c.getOrganization()
	if err != nil {
		return nil, err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s/%s/status", org, opts.Name)

	var cluster v1alpha1.ClusterInfo
	err = c.getParsedResponse(http.MethodGet, apiPath, jsonHeader, nil, &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Client) ConnectCluster(opts clustermodel.ConnectOptions) (*v1alpha1.ClusterInfo, error) {
	org, err := c.getOrganization()
	if err != nil {
		return nil, err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s/%s/connect", org, opts.Name)

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal connect options. err: %w", err)
	}

	var cluster v1alpha1.ClusterInfo
	err = c.getParsedResponse(http.MethodPost, apiPath, jsonHeader, bytes.NewReader(body), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Client) ReconfigureCluster(opts clustermodel.ReconfigureOptions, responseID string) (*v1alpha1.ClusterInfo, error) {
	org, err := c.getOrganization()
	if err != nil {
		return nil, err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s/%s/reconfigure", org, opts.Name)

	params := make([]queryParams, 0)
	if responseID != "" {
		params = append(params, queryParams{key: "response-id", value: responseID})
	}
	apiPath, err = setQueryParams(apiPath, params)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal reconfigure options. err: %w", err)
	}

	var cluster v1alpha1.ClusterInfo
	err = c.getParsedResponse(http.MethodPost, apiPath, jsonHeader, bytes.NewReader(body), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Client) RemoveCluster(opts clustermodel.RemovalOptions, responseID string) error {
	org, err := c.getOrganization()
	if err != nil {
		return err
	}
	apiPath := fmt.Sprintf("/clustersv2/%s/%s/remove", org, opts.Name)

	params := make([]queryParams, 0)
	if responseID != "" {
		params = append(params, queryParams{key: "response-id", value: responseID})
	}
	apiPath, err = setQueryParams(apiPath, params)
	if err != nil {
		return err
	}

	body, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("failed to unmarshal remove options. err: %w", err)
	}
	_, err = c.getResponse(http.MethodPost, apiPath, jsonHeader, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}
