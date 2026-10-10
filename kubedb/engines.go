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
	"fmt"
	"sort"
	"strings"
)

// Engine describes the KubeDB coordinates for a single database kind. Values here
// are derived from kubedb.dev/apimachinery (apis/kubedb): a kind served under the GA
// group version uses "v1"; the remaining kinds use "v1alpha2". Confirm availability
// for a given cluster via `available-types`.
type Engine struct {
	// Name is the canonical lower-case engine name, e.g. "postgres".
	Name string
	// Group/Version/Kind/Resource identify the primary CR for the proxy path.
	Group    string
	Version  string
	Kind     string
	Resource string
	// Ports are the engine's default client service ports (informational).
	Ports []int
	// AuthType hints how connection credentials in the {name}-auth secret are shaped:
	// "password" (username/password) or "apikey" (Qdrant/Weaviate).
	AuthType string
}

const group = "kubedb.com"

// engines is the full set of databases KubeDB supports. Kinds served under the GA
// group version use "v1"; the rest use "v1alpha2".
var engines = map[string]Engine{
	// --- kubedb.com/v1 (GA) ---
	"elasticsearch": {Name: "elasticsearch", Group: group, Version: "v1", Kind: "Elasticsearch", Resource: "elasticsearches", Ports: []int{9200}, AuthType: "password"},
	"kafka":         {Name: "kafka", Group: group, Version: "v1", Kind: "Kafka", Resource: "kafkas", Ports: []int{9092}, AuthType: "password"},
	"mariadb":       {Name: "mariadb", Group: group, Version: "v1", Kind: "MariaDB", Resource: "mariadbs", Ports: []int{3306}, AuthType: "password"},
	"memcached":     {Name: "memcached", Group: group, Version: "v1", Kind: "Memcached", Resource: "memcacheds", Ports: []int{11211}, AuthType: "password"},
	"mongodb":       {Name: "mongodb", Group: group, Version: "v1", Kind: "MongoDB", Resource: "mongodbs", Ports: []int{27017}, AuthType: "password"},
	"mysql":         {Name: "mysql", Group: group, Version: "v1", Kind: "MySQL", Resource: "mysqls", Ports: []int{3306}, AuthType: "password"},
	"perconaxtradb": {Name: "perconaxtradb", Group: group, Version: "v1", Kind: "PerconaXtraDB", Resource: "perconaxtradbs", Ports: []int{3306}, AuthType: "password"},
	"pgbouncer":     {Name: "pgbouncer", Group: group, Version: "v1", Kind: "PgBouncer", Resource: "pgbouncers", Ports: []int{5432}, AuthType: "password"},
	"postgres":      {Name: "postgres", Group: group, Version: "v1", Kind: "Postgres", Resource: "postgreses", Ports: []int{5432}, AuthType: "password"},
	"proxysql":      {Name: "proxysql", Group: group, Version: "v1", Kind: "ProxySQL", Resource: "proxysqls", Ports: []int{6033}, AuthType: "password"},
	"redis":         {Name: "redis", Group: group, Version: "v1", Kind: "Redis", Resource: "redises", Ports: []int{6379}, AuthType: "password"},
	"redissentinel": {Name: "redissentinel", Group: group, Version: "v1", Kind: "RedisSentinel", Resource: "redissentinels", Ports: []int{26379}, AuthType: "password"},

	// --- kubedb.com/v1alpha2 ---
	"aerospike":   {Name: "aerospike", Group: group, Version: "v1alpha2", Kind: "Aerospike", Resource: "aerospikes", Ports: []int{3000}, AuthType: "password"},
	"cassandra":   {Name: "cassandra", Group: group, Version: "v1alpha2", Kind: "Cassandra", Resource: "cassandras", Ports: []int{9042}, AuthType: "password"},
	"clickhouse":  {Name: "clickhouse", Group: group, Version: "v1alpha2", Kind: "ClickHouse", Resource: "clickhouses", Ports: []int{9000, 8123}, AuthType: "password"},
	"documentdb":  {Name: "documentdb", Group: group, Version: "v1alpha2", Kind: "DocumentDB", Resource: "documentdbs", Ports: []int{10260}, AuthType: "password"},
	"druid":       {Name: "druid", Group: group, Version: "v1alpha2", Kind: "Druid", Resource: "druids", Ports: []int{8888}, AuthType: "password"},
	"hanadb":      {Name: "hanadb", Group: group, Version: "v1alpha2", Kind: "HanaDB", Resource: "hanadbs", Ports: []int{39017}, AuthType: "password"},
	"hazelcast":   {Name: "hazelcast", Group: group, Version: "v1alpha2", Kind: "Hazelcast", Resource: "hazelcasts", Ports: []int{5701}, AuthType: "password"},
	"ignite":      {Name: "ignite", Group: group, Version: "v1alpha2", Kind: "Ignite", Resource: "ignites", Ports: []int{10800}, AuthType: "password"},
	"milvus":      {Name: "milvus", Group: group, Version: "v1alpha2", Kind: "Milvus", Resource: "milvuses", Ports: []int{19530}, AuthType: "password"},
	"mssqlserver": {Name: "mssqlserver", Group: group, Version: "v1alpha2", Kind: "MSSQLServer", Resource: "mssqlservers", Ports: []int{1433}, AuthType: "password"},
	"oracle":      {Name: "oracle", Group: group, Version: "v1alpha2", Kind: "Oracle", Resource: "oracles", Ports: []int{1521}, AuthType: "password"},
	"pgpool":      {Name: "pgpool", Group: group, Version: "v1alpha2", Kind: "Pgpool", Resource: "pgpools", Ports: []int{9999}, AuthType: "password"},
	"qdrant":      {Name: "qdrant", Group: group, Version: "v1alpha2", Kind: "Qdrant", Resource: "qdrants", Ports: []int{6333, 6334}, AuthType: "apikey"},
	"rabbitmq":    {Name: "rabbitmq", Group: group, Version: "v1alpha2", Kind: "RabbitMQ", Resource: "rabbitmqs", Ports: []int{5672, 15672}, AuthType: "password"},
	"singlestore": {Name: "singlestore", Group: group, Version: "v1alpha2", Kind: "Singlestore", Resource: "singlestores", Ports: []int{3306}, AuthType: "password"},
	"solr":        {Name: "solr", Group: group, Version: "v1alpha2", Kind: "Solr", Resource: "solrs", Ports: []int{8983}, AuthType: "password"},
	"weaviate":    {Name: "weaviate", Group: group, Version: "v1alpha2", Kind: "Weaviate", Resource: "weaviates", Ports: []int{8080, 50051}, AuthType: "apikey"},
	"zookeeper":   {Name: "zookeeper", Group: group, Version: "v1alpha2", Kind: "ZooKeeper", Resource: "zookeepers", Ports: []int{2181}, AuthType: "password"},
}

// Lookup returns the Engine for a name (case-insensitive) or an error listing
// the supported engines.
func Lookup(name string) (Engine, error) {
	e, ok := engines[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return Engine{}, fmt.Errorf("unsupported engine %q (supported: %s)", name, strings.Join(SupportedEngines(), ", "))
	}
	return e, nil
}

// SupportedEngines returns the sorted list of engine names the client knows about.
func SupportedEngines() []string {
	out := make([]string, 0, len(engines))
	for k := range engines {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
