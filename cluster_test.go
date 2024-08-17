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

package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.bytebuilders.dev/client"
	clustermodel "go.bytebuilders.dev/resource-model/apis/cluster"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	kmapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

var (
	TestClusterDisplayName  = "emruz-test-3"
	TestClusterName         = "emruz-test-3-linode"
	TestClusterProvider     = "Linode"
	TestCloudCredentialName = "emruz-linode"
	TestClusterID           = "87873"
	FeatureSetOpscenterCore = "opscenter-core"
	FeatureKubeUIServer     = "kube-ui-server"
	FeatureLicenseServer    = "license-proxyserver"
)

func TestClient_CheckClusterAPIs(t *testing.T) {
	providerOptions := clustermodel.ProviderOptions{
		Credential: TestCloudCredentialName,
		Name:       TestClusterProvider,
		ClusterID:  TestClusterID,
	}
	basicInfo := clustermodel.BasicInfo{
		Name:        TestClusterName,
		DisplayName: TestClusterDisplayName,
	}
	c := client.NewClient(client.TestServerURL).WithBasicAuth(client.TestServerUser, client.TestServrPassword)

	t.Run("CheckClusterExistence() should return NotImported phase when cluster is not imported", func(t *testing.T) {
		cluster, err := c.CheckClusterExistence(clustermodel.CheckOptions{
			Provider: providerOptions,
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, kmapi.ClusterPhaseNotImported, cluster.Status.Phase)
	})

	t.Run("ImportCluster() should import the cluster", func(t *testing.T) {
		cluster, err := c.ImportCluster(clustermodel.ImportOptions{
			BasicInfo: basicInfo,
			Provider:  providerOptions,
			Components: clustermodel.ComponentOptions{
				FluxCD: true,
				FeatureSets: []clustermodel.FeatureSet{
					{
						Name: FeatureSetOpscenterCore,
						Features: []string{
							FeatureKubeUIServer,
							FeatureLicenseServer,
						},
					},
				},
			},
		}, "")
		if !assert.Nil(t, err) {
			return
		}
		if !assert.Nil(t, waitForClusterToBeReady(c, cluster.Spec.Name)) {
			return
		}
	})

	t.Run("CheckClusterExistence() should return cluster phase when the cluster already exist", func(t *testing.T) {
		cluster, err := c.CheckClusterExistence(clustermodel.CheckOptions{
			Provider: providerOptions,
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.NotEqual(t, kmapi.ClusterPhaseNotImported, cluster.Status.Phase)
	})

	t.Run("ListClusters() should return non empty cluster list when the cluster exist", func(t *testing.T) {
		clusters, err := c.ListClusters(clustermodel.ListOptions{})
		if !assert.Nil(t, err) {
			return
		}
		assert.True(t, len(clusters.Items) > 0)
	})

	t.Run("ListClusters() should return empty cluster list when list options doesn't select any cluster", func(t *testing.T) {
		clusters, err := c.ListClusters(clustermodel.ListOptions{Provider: "None"})
		if !assert.Nil(t, err) {
			return
		}
		assert.True(t, len(clusters.Items) == 0)
	})

	t.Run("GetCluster() should return the cluster status when the cluster exist", func(t *testing.T) {
		cluster, err := c.GetCluster(clustermodel.GetOptions{
			Name: TestClusterName,
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, TestClusterName, cluster.Spec.Name)
		assert.Equal(t, TestClusterDisplayName, cluster.Spec.DisplayName)
		assert.Equal(t, TestClusterProvider, string(cluster.Spec.Provider))
		assert.Equal(t, kmapi.ClusterPhaseActive, cluster.Status.Phase)
	})

	t.Run("GetClusterClientConfig() should return client-config for the cluster", func(t *testing.T) {
		clientConfig, err := c.GetClusterClientConfig(clustermodel.GetOptions{
			Name: TestClusterName,
		})

		assert.Nil(t, err)
		assert.NotNil(t, clientConfig)

		cfg, err := clientConfig.ClientConfig()
		assert.Nil(t, err)

		kc, err := kubernetes.NewForConfig(cfg)
		assert.Nil(t, err)

		pods, err := kc.CoreV1().Pods(v1.NamespaceSystem).List(context.Background(), v1.ListOptions{Limit: 1})
		assert.Nil(t, err)

		assert.Len(t, pods.Items, 1)
	})

	t.Run("ConnectCluster() should return Active status when the cluster is already connected", func(t *testing.T) {
		cluster, err := c.ConnectCluster(clustermodel.ConnectOptions{
			Name:       TestClusterName,
			Credential: TestCloudCredentialName,
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, kmapi.ClusterPhaseActive, cluster.Status.Phase)
	})

	// TODO: Remove cluster components, make the cluster NotReady, then run reconfigure.
	//t.Run("ReconfigureCluster() should re-install cluster components", func(t *testing.T) {
	//	cluster, err := c.ReconfigureCluster(TestClusterName, true)
	//	if !assert.Nil(t, err) {
	//		return
	//	}
	//	assert.Nil(t, waitForClusterToBeReady(c, cluster.Spec.Name))
	//})

	t.Run("RemoveCluster() should remove the imported cluster", func(t *testing.T) {
		err := c.RemoveCluster(clustermodel.RemovalOptions{
			Name: TestClusterName,
			Components: clustermodel.ComponentOptions{
				FluxCD: true,
				FeatureSets: []clustermodel.FeatureSet{
					{
						Name: FeatureSetOpscenterCore,
						Features: []string{
							FeatureKubeUIServer,
							FeatureLicenseServer,
						},
					},
				},
			},
		}, "")
		if !assert.Nil(t, err) {
			return
		}
		assert.Nil(t, waitForClusterToBeRemoved(c, providerOptions))
	})

	t.Run("ListClusters() should return empty cluster list when no cluster exist", func(t *testing.T) {
		clusters, err := c.ListClusters(clustermodel.ListOptions{})
		if !assert.Nil(t, err) {
			return
		}
		assert.True(t, len(clusters.Items) == 0)
	})

	t.Run("GetCluster() should return error when the cluster does not exist", func(t *testing.T) {
		_, err := c.GetCluster(clustermodel.GetOptions{
			Name: TestClusterName,
		})
		assert.NotNil(t, err)
	})
}

func waitForClusterToBeReady(c *client.Client, clusterName string) error {
	return wait.PollUntilContextTimeout(context.TODO(), 2*time.Second, 5*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		cluster, err := c.GetCluster(clustermodel.GetOptions{
			Name: clusterName,
		})
		if err != nil {
			return false, err
		}
		if cluster.Status.Phase == kmapi.ClusterPhaseActive {
			return true, nil
		}
		return false, nil
	})
}

func waitForClusterToBeRemoved(c *client.Client, opts clustermodel.ProviderOptions) error {
	return wait.PollUntilContextTimeout(context.TODO(), 2*time.Second, 5*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		cluster, err := c.CheckClusterExistence(clustermodel.CheckOptions{
			Provider: opts,
		})
		if err != nil {
			return false, err
		}
		if cluster.Status.Phase == kmapi.ClusterPhaseNotImported {
			return true, nil
		}
		return false, nil
	})
}
