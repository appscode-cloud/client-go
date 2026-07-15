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

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testUA = "unit-test-agent"

func newTestClient(t *testing.T, h http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, "tok", testUA, 5*time.Second, false)
}

func TestNewClient_DefaultUserAgent(t *testing.T) {
	c := NewClient("https://example.com", "tok", "", time.Second, false)
	if c.userAgent != defaultUserAgent {
		t.Fatalf("userAgent = %q, want default %q", c.userAgent, defaultUserAgent)
	}
}

func TestGetResource_ParsesStatusAndSendsHeaders(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/clusters/acme/prod/proxy/kubedb.com/v1/namespaces/team-a/postgreses/db1",
		func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("Authorization"); got != "token tok" {
				t.Errorf("Authorization = %q, want %q", got, "token tok")
			}
			if got := r.Header.Get("User-Agent"); got != testUA {
				t.Errorf("User-Agent = %q, want %q", got, testUA)
			}
			_, _ = w.Write([]byte(`{"status":{"phase":"Ready"}}`))
		})

	c := newTestClient(t, mux)
	obj, err := c.GetResource(context.Background(), "acme", "prod", "kubedb.com", "v1", "postgreses", "team-a", "db1")
	if err != nil {
		t.Fatalf("GetResource: %v", err)
	}
	if obj.Status.Phase != "Ready" {
		t.Fatalf("phase = %q, want Ready", obj.Status.Phase)
	}
}

func TestGetResource_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/clusters/acme/prod/proxy/kubedb.com/v1/namespaces/team-a/postgreses/missing",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"not found"}`))
		})

	c := newTestClient(t, mux)
	_, err := c.GetResource(context.Background(), "acme", "prod", "kubedb.com", "v1", "postgreses", "team-a", "missing")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !apiErr.NotFound() {
		t.Fatalf("want APIError NotFound, got %v", err)
	}
}

func TestGetSecret_DecodesBase64(t *testing.T) {
	enc := base64.StdEncoding.EncodeToString
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/clusters/acme/prod/proxy/core/v1/namespaces/team-a/secrets/db1-auth",
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"data":{"username":"` + enc([]byte("pg-user")) + `","password":"` + enc([]byte("p@ss")) + `"}}`))
		})

	c := newTestClient(t, mux)
	creds, err := c.GetSecret(context.Background(), "acme", "prod", "team-a", "db1-auth")
	if err != nil {
		t.Fatalf("GetSecret: %v", err)
	}
	if creds["username"] != "pg-user" || creds["password"] != "p@ss" {
		t.Fatalf("decoded creds = %v", creds)
	}
}

func TestGenerateModel_APIErrorStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v1/clusters/acme/prod/helm/options/model",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"message":"bad options"}`))
		})

	c := newTestClient(t, mux)
	_, err := c.GenerateModel(context.Background(), "acme", "prod", OptionsModelRequest{})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 APIError, got %v", err)
	}
}

func TestApplyEditor_SendsModelAndReturnsTask(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v1/clusters/acme/prod/helm/editor/",
		func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"k":"v"}` {
				t.Errorf("editor body = %q, want the model verbatim", string(body))
			}
			_, _ = w.Write([]byte(`{"id":"task-1"}`))
		})

	c := newTestClient(t, mux)
	task, err := c.ApplyEditor(context.Background(), "acme", "prod", EditorModel(`{"k":"v"}`))
	if err != nil {
		t.Fatalf("ApplyEditor: %v", err)
	}
	if task.ID != "task-1" {
		t.Fatalf("task id = %q", task.ID)
	}
}
