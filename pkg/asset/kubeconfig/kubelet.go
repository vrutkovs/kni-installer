package kubeconfig

import (
	"path/filepath"

	"github.com/metalkube/kni-installer/pkg/asset"
	"github.com/metalkube/kni-installer/pkg/asset/installconfig"
	"github.com/metalkube/kni-installer/pkg/asset/tls"
)

var (
	kubeconfigKubeletPath       = filepath.Join("auth", "kubeconfig-kubelet")
	kubeconfigKubeletClientPath = filepath.Join("auth", "kubeconfig-kubelet-client")
)

// Kubelet is the asset for the kubelet kubeconfig.
// [DEPRECATED]
type Kubelet struct {
	kubeconfig
}

var _ asset.WritableAsset = (*Kubelet)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Kubelet) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.KubeCA{},
		&tls.KubeletCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *Kubelet) Generate(parents asset.Parents) error {
	kubeCA := &tls.KubeCA{}
	kubeletCertKey := &tls.KubeletCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(kubeCA, kubeletCertKey, installConfig)

	return k.kubeconfig.generate(
		kubeCA,
		kubeletCertKey,
		installConfig.Config,
		"kubelet",
		kubeconfigKubeletPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *Kubelet) Name() string {
	return "Kubeconfig Kubelet"
}

// Load is a no-op because kubelet kubeconfig is not written to disk.
func (k *Kubelet) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}

// KubeletClient is the asset for the kubelet kubeconfig.
type KubeletClient struct {
	kubeconfig
}

var _ asset.WritableAsset = (*KubeletClient)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *KubeletClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.KubeAPIServerCompleteCABundle{},
		&tls.KubeletClientCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *KubeletClient) Generate(parents asset.Parents) error {
	ca := &tls.KubeAPIServerCompleteCABundle{}
	clientcertkey := &tls.KubeletClientCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientcertkey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientcertkey,
		installConfig.Config,
		"kubelet",
		kubeconfigKubeletClientPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *KubeletClient) Name() string {
	return "Kubeconfig Kubelet Client"
}

// Load is a no-op because kubelet kubeconfig is not written to disk.
func (k *KubeletClient) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}
