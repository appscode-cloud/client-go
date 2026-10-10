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

// Package kubedb is a thin, dependency-free client for the KubeDB Platform API Server,
// exposing the subset of endpoints needed to provision databases in Helm editor mode.
// It is shared by the BMC Helix and ServiceNow provisioning adapters.
package kubedb

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const defaultUserAgent = "bytebuilders-client-go"

// Client talks to the KubeDB Platform API using a service-account bearer token.
type Client struct {
	baseURL   string
	token     string
	userAgent string
	http      *http.Client
}

// NewClient builds a KubeDB Platform API client. userAgent identifies the caller
// (e.g. "bmchelix-adapter"); an empty value falls back to a default.
func NewClient(baseURL, token, userAgent string, timeout time.Duration, insecureSkipTLS bool) *Client {
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecureSkipTLS}, //nolint:gosec // opt-in dev/test only
		DialContext:         (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
	}
	return &Client{
		baseURL:   strings.TrimRight(baseURL, "/"),
		token:     token,
		userAgent: userAgent,
		http: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

// APIError is returned for non-2xx responses from the KubeDB API.
type APIError struct {
	Method     string
	URL        string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("kubedb api %s %s: status %d: %s", e.Method, e.URL, e.StatusCode, truncate(e.Body, 512))
}

// NotFound reports whether the error is a 404 from the KubeDB API.
func (e *APIError) NotFound() bool { return e.StatusCode == http.StatusNotFound }

// GenerateModel calls PUT /helm/options/model, expanding the given options into a
// full editor model (returned opaquely).
func (c *Client) GenerateModel(ctx context.Context, owner, cluster string, req OptionsModelRequest) (EditorModel, error) {
	var model json.RawMessage
	if err := c.do(ctx, http.MethodPut, c.clusterPath(owner, cluster, "/helm/options/model"), req, &model); err != nil {
		return nil, err
	}
	return EditorModel(model), nil
}

// ApplyEditor calls PUT /helm/editor/, submitting the model as an async apply task.
func (c *Client) ApplyEditor(ctx context.Context, owner, cluster string, model EditorModel) (*TaskResponse, error) {
	var task TaskResponse
	if err := c.do(ctx, http.MethodPut, c.clusterPath(owner, cluster, "/helm/editor/"), model, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteEditor calls DELETE /helm/editor/, removing the release and companion objects.
func (c *Client) DeleteEditor(ctx context.Context, owner, cluster string, meta OptionsMetadata) (*TaskResponse, error) {
	var task TaskResponse
	if err := c.do(ctx, http.MethodDelete, c.clusterPath(owner, cluster, "/helm/editor/"), DeleteEditorRequest{Metadata: meta}, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// GetResource reads a namespaced CR through the Kubernetes proxy.
func (c *Client) GetResource(ctx context.Context, owner, cluster, group, version, resource, namespace, name string) (*Unstructured, error) {
	suffix := fmt.Sprintf("/proxy/%s/%s/namespaces/%s/%s/%s", group, version, namespace, resource, name)
	var obj Unstructured
	if err := c.do(ctx, http.MethodGet, c.clusterPath(owner, cluster, suffix), nil, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// GetSecret reads a namespaced core/v1 Secret through the proxy and returns its
// decoded (base64) string values.
func (c *Client) GetSecret(ctx context.Context, owner, cluster, namespace, name string) (map[string]string, error) {
	suffix := fmt.Sprintf("/proxy/core/v1/namespaces/%s/secrets/%s", namespace, name)
	var s Secret
	if err := c.do(ctx, http.MethodGet, c.clusterPath(owner, cluster, suffix), nil, &s); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(s.Data))
	for k, v := range s.Data {
		if dec, err := base64.StdEncoding.DecodeString(v); err == nil {
			out[k] = string(dec)
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// AvailableTypes returns the raw JSON of GET /available-types for a cluster.
func (c *Client) AvailableTypes(ctx context.Context, owner, cluster string) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, c.clusterPath(owner, cluster, "/available-types"), nil, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *Client) clusterPath(owner, cluster, suffix string) string {
	return fmt.Sprintf("%s/api/v1/clusters/%s/%s%s", c.baseURL, owner, cluster, suffix)
}

func (c *Client) do(ctx context.Context, method, url string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{Method: method, URL: url, StatusCode: resp.StatusCode, Body: string(data)}
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
