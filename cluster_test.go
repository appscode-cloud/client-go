package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.bytebuilders.dev/client"
	clustermodel "go.bytebuilders.dev/resource-model/apis/cluster"
	"go.bytebuilders.dev/resource-model/apis/cluster/v1alpha1"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	TestClusterDisplayName  = "emruz-test-3"
	TestClusterName         = "emruz-test-3-linode"
	TestClusterProvider     = "Linode"
	TestCloudCredentialName = "emruz-linode"
	TestClusterID           = "87873"
	FeatureSetOpscenterCore = "opscenter-core"
	FeatureKubeUIServer     = "kube-ui-server"
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
		assert.Equal(t, v1alpha1.ClusterPhaseNotImported, cluster.Status.Phase)
	})

	t.Run("ImportCluster() should import the cluster", func(t *testing.T) {
		cluster, err := c.ImportCluster(clustermodel.ImportOptions{
			BasicInfo: basicInfo,
			Provider:  providerOptions,
			Components: clustermodel.ComponentOptions{
				FluxCD:        true,
				LicenseServer: true,
				FeatureSets: []clustermodel.FeatureSet{
					{
						Name:     FeatureSetOpscenterCore,
						Features: []string{FeatureKubeUIServer},
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
		assert.NotEqual(t, v1alpha1.ClusterPhaseNotImported, cluster.Status.Phase)
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
		assert.Equal(t, v1alpha1.ClusterPhaseActive, cluster.Status.Phase)
	})

	t.Run("ConnectCluster() should return Active status when the cluster is already connected", func(t *testing.T) {
		cluster, err := c.ConnectCluster(clustermodel.ConnectOptions{
			Name:       TestClusterName,
			Credential: TestCloudCredentialName,
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, v1alpha1.ClusterPhaseActive, cluster.Status.Phase)
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
				FluxCD:        true,
				LicenseServer: true,
				FeatureSets: []clustermodel.FeatureSet{
					{
						Name:     FeatureSetOpscenterCore,
						Features: []string{FeatureKubeUIServer},
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
	return wait.PollImmediate(2*time.Second, 5*time.Minute, func() (done bool, err error) {
		cluster, err := c.GetCluster(clustermodel.GetOptions{
			Name: clusterName,
		})
		if err != nil {
			return false, err
		}
		if cluster.Status.Phase == v1alpha1.ClusterPhaseActive {
			return true, nil
		}
		return false, nil
	})
}

func waitForClusterToBeRemoved(c *client.Client, opts clustermodel.ProviderOptions) error {
	return wait.PollImmediate(2*time.Second, 5*time.Minute, func() (done bool, err error) {
		cluster, err := c.CheckClusterExistence(clustermodel.CheckOptions{
			Provider: opts,
		})
		if err != nil {
			return false, err
		}
		if cluster.Status.Phase == v1alpha1.ClusterPhaseNotImported {
			return true, nil
		}
		return false, nil
	})
}
