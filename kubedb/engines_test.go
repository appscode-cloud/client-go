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

import "testing"

func TestEngines_Complete(t *testing.T) {
	if got := len(SupportedEngines()); got != 30 {
		t.Fatalf("expected 30 supported engines, got %d", got)
	}
	for _, name := range SupportedEngines() {
		e, err := Lookup(name)
		if err != nil {
			t.Fatalf("Lookup(%q): %v", name, err)
		}
		if e.Group != "kubedb.com" {
			t.Errorf("%s: group = %q, want kubedb.com", name, e.Group)
		}
		if e.Version != "v1" && e.Version != "v1alpha2" {
			t.Errorf("%s: version = %q, want v1 or v1alpha2", name, e.Version)
		}
		if e.Kind == "" || e.Resource == "" {
			t.Errorf("%s: missing kind/resource: %+v", name, e)
		}
		if len(e.Ports) == 0 {
			t.Errorf("%s: no ports", name)
		}
	}
}

func TestEngines_KnownCoordinates(t *testing.T) {
	// Spot-check a sample against kubedb.dev/apimachinery to catch registry typos.
	want := map[string]struct {
		version, kind, resource, auth string
	}{
		"postgres":      {"v1", "Postgres", "postgreses", "password"},
		"redis":         {"v1", "Redis", "redises", "password"},
		"mysql":         {"v1", "MySQL", "mysqls", "password"},
		"memcached":     {"v1", "Memcached", "memcacheds", "password"}, // has auth by default
		"mssqlserver":   {"v1alpha2", "MSSQLServer", "mssqlservers", "password"},
		"qdrant":        {"v1alpha2", "Qdrant", "qdrants", "apikey"},
		"weaviate":      {"v1alpha2", "Weaviate", "weaviates", "apikey"},
		"milvus":        {"v1alpha2", "Milvus", "milvuses", "password"}, // no API key
		"perconaxtradb": {"v1", "PerconaXtraDB", "perconaxtradbs", "password"},
	}
	for name, w := range want {
		e, err := Lookup(name)
		if err != nil {
			t.Fatalf("Lookup(%q): %v", name, err)
		}
		if e.Version != w.version || e.Kind != w.kind || e.Resource != w.resource || e.AuthType != w.auth {
			t.Errorf("%s = {%s %s %s %s}, want {%s %s %s %s}",
				name, e.Version, e.Kind, e.Resource, e.AuthType, w.version, w.kind, w.resource, w.auth)
		}
	}
}

func TestEngines_AuthTypeValues(t *testing.T) {
	for _, name := range SupportedEngines() {
		e, _ := Lookup(name)
		if e.AuthType != "password" && e.AuthType != "apikey" {
			t.Errorf("%s: unexpected AuthType %q", name, e.AuthType)
		}
	}
}

func TestLookup_CaseInsensitiveAndUnknown(t *testing.T) {
	if e, err := Lookup("PostGres"); err != nil || e.Kind != "Postgres" {
		t.Fatalf("case-insensitive lookup failed: %+v, %v", e, err)
	}
	if _, err := Lookup("cockroach"); err == nil {
		t.Fatalf("expected error for unknown engine")
	}
}
