// Package manifests deals with creating manifests for all manifests to be installed for the cluster
package manifests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/metalkube/kni-installer/pkg/asset"
	"github.com/metalkube/kni-installer/pkg/asset/installconfig"
	"github.com/metalkube/kni-installer/pkg/asset/templates/content/bootkube"
	"github.com/metalkube/kni-installer/pkg/asset/tls"
)

const (
	manifestDir = "manifests"
)

var (
	kubeSysConfigPath = filepath.Join(manifestDir, "cluster-config.yaml")

	_ asset.WritableAsset = (*Manifests)(nil)

	customTmplFuncs = template.FuncMap{
		"indent": indent,
		"add": func(i, j int) int {
			return i + j
		},
	}
)

// Manifests generates the dependent operator config.yaml files
type Manifests struct {
	KubeSysConfig *configurationObject
	FileList      []*asset.File
}

type genericData map[string]string

// Name returns a human friendly name for the operator
func (m *Manifests) Name() string {
	return "Common Manifests"
}

// Dependencies returns all of the dependencies directly needed by a
// Manifests asset.
func (m *Manifests) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		&Ingress{},
		&DNS{},
		&Infrastructure{},
		&Networking{},
		&tls.RootCA{},
		&tls.EtcdCA{},
		&tls.EtcdClientCertKey{},
		&tls.EtcdMetricsCABundle{},
		&tls.EtcdMetricsSignerClientCertKey{},
		&tls.MCSCertKey{},

		&bootkube.KubeCloudConfig{},
		&bootkube.MachineConfigServerTLSSecret{},
		&bootkube.Pull{},
		&bootkube.CVOOverrides{},
		&bootkube.HostEtcdServiceEndpointsKubeSystem{},
		&bootkube.KubeSystemConfigmapEtcdServingCA{},
		&bootkube.KubeSystemConfigmapRootCA{},
		&bootkube.KubeSystemSecretEtcdClient{},
		&bootkube.OpenshiftConfigSecretEtcdMetricsClient{},
		&bootkube.OpenshiftConfigConfigmapEtcdMetricsServingCA{},

		&bootkube.OpenshiftMachineConfigOperator{},
		&bootkube.EtcdServiceKubeSystem{},
		&bootkube.HostEtcdServiceKubeSystem{},
	}
}

// Generate generates the respective operator config.yml files
func (m *Manifests) Generate(dependencies asset.Parents) error {
	ingress := &Ingress{}
	dns := &DNS{}
	network := &Networking{}
	infra := &Infrastructure{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig, ingress, dns, network, infra)

	// mao go to kube-system config map
	m.KubeSysConfig = configMap("kube-system", "cluster-config-v1", genericData{
		"install-config": string(installConfig.Files()[0].Data),
	})
	kubeSysConfigData, err := yaml.Marshal(m.KubeSysConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create kube-system/cluster-config-v1 configmap")
	}

	m.FileList = []*asset.File{
		{
			Filename: kubeSysConfigPath,
			Data:     kubeSysConfigData,
		},
	}
	m.FileList = append(m.FileList, m.generateBootKubeManifests(dependencies)...)

	m.FileList = append(m.FileList, ingress.Files()...)
	m.FileList = append(m.FileList, dns.Files()...)
	m.FileList = append(m.FileList, network.Files()...)
	m.FileList = append(m.FileList, infra.Files()...)

	asset.SortFiles(m.FileList)

	return nil
}

// Files returns the files generated by the asset.
func (m *Manifests) Files() []*asset.File {
	return m.FileList
}

func (m *Manifests) generateBootKubeManifests(dependencies asset.Parents) []*asset.File {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	etcdCA := &tls.EtcdCA{}
	mcsCertKey := &tls.MCSCertKey{}
	etcdClientCertKey := &tls.EtcdClientCertKey{}
	etcdMetricsCABundle := &tls.EtcdMetricsCABundle{}
	etcdMetricsSignerClientCertKey := &tls.EtcdMetricsSignerClientCertKey{}
	rootCA := &tls.RootCA{}
	dependencies.Get(
		clusterID,
		installConfig,
		etcdCA,
		etcdClientCertKey,
		etcdMetricsCABundle,
		etcdMetricsSignerClientCertKey,
		mcsCertKey,
		rootCA,
	)

	etcdEndpointHostnames := make([]string, *installConfig.Config.ControlPlane.Replicas)
	for i := range etcdEndpointHostnames {
		etcdEndpointHostnames[i] = fmt.Sprintf("etcd-%d", i)
	}

	templateData := &bootkubeTemplateData{
		Base64encodeCloudProviderConfig: "", // FIXME
		EtcdCaCert:                      string(etcdCA.Cert()),
		EtcdClientCert:                  base64.StdEncoding.EncodeToString(etcdClientCertKey.Cert()),
		EtcdClientKey:                   base64.StdEncoding.EncodeToString(etcdClientCertKey.Key()),
		EtcdMetricsCaCert:               string(etcdMetricsCABundle.Cert()),
		EtcdMetricsClientCert:           base64.StdEncoding.EncodeToString(etcdMetricsSignerClientCertKey.Cert()),
		EtcdMetricsClientKey:            base64.StdEncoding.EncodeToString(etcdMetricsSignerClientCertKey.Key()),
		McsTLSCert:                      base64.StdEncoding.EncodeToString(mcsCertKey.Cert()),
		McsTLSKey:                       base64.StdEncoding.EncodeToString(mcsCertKey.Key()),
		PullSecretBase64:                base64.StdEncoding.EncodeToString([]byte(installConfig.Config.PullSecret)),
		RootCaCert:                      string(rootCA.Cert()),
		CVOClusterID:                    clusterID.UUID,
		EtcdEndpointHostnames:           etcdEndpointHostnames,
		EtcdEndpointDNSSuffix:           installConfig.Config.ClusterDomain(),
	}

	kubeCloudConfig := &bootkube.KubeCloudConfig{}
	machineConfigServerTLSSecret := &bootkube.MachineConfigServerTLSSecret{}
	pull := &bootkube.Pull{}
	cVOOverrides := &bootkube.CVOOverrides{}
	hostEtcdServiceEndpointsKubeSystem := &bootkube.HostEtcdServiceEndpointsKubeSystem{}
	kubeSystemConfigmapEtcdServingCA := &bootkube.KubeSystemConfigmapEtcdServingCA{}
	kubeSystemConfigmapRootCA := &bootkube.KubeSystemConfigmapRootCA{}
	kubeSystemSecretEtcdClient := &bootkube.KubeSystemSecretEtcdClient{}
	openshiftConfigSecretEtcdMetricsClient := &bootkube.OpenshiftConfigSecretEtcdMetricsClient{}
	openshiftConfigConfigmapEtcdMetricsServingCA := &bootkube.OpenshiftConfigConfigmapEtcdMetricsServingCA{}

	openshiftMachineConfigOperator := &bootkube.OpenshiftMachineConfigOperator{}
	etcdServiceKubeSystem := &bootkube.EtcdServiceKubeSystem{}
	hostEtcdServiceKubeSystem := &bootkube.HostEtcdServiceKubeSystem{}
	dependencies.Get(
		kubeCloudConfig,
		machineConfigServerTLSSecret,
		pull,
		cVOOverrides,
		hostEtcdServiceEndpointsKubeSystem,
		kubeSystemConfigmapEtcdServingCA,
		kubeSystemConfigmapRootCA,
		kubeSystemSecretEtcdClient,
		openshiftConfigSecretEtcdMetricsClient,
		openshiftConfigConfigmapEtcdMetricsServingCA,
		openshiftMachineConfigOperator,
		etcdServiceKubeSystem,
		hostEtcdServiceKubeSystem,
	)
	assetData := map[string][]byte{
		"kube-cloud-config.yaml":                                  applyTemplateData(kubeCloudConfig.Files()[0].Data, templateData),
		"machine-config-server-tls-secret.yaml":                   applyTemplateData(machineConfigServerTLSSecret.Files()[0].Data, templateData),
		"pull.json":                                               applyTemplateData(pull.Files()[0].Data, templateData),
		"cvo-overrides.yaml":                                      applyTemplateData(cVOOverrides.Files()[0].Data, templateData),
		"host-etcd-service-endpoints.yaml":                        applyTemplateData(hostEtcdServiceEndpointsKubeSystem.Files()[0].Data, templateData),
		"kube-system-configmap-etcd-serving-ca.yaml":              applyTemplateData(kubeSystemConfigmapEtcdServingCA.Files()[0].Data, templateData),
		"kube-system-configmap-root-ca.yaml":                      applyTemplateData(kubeSystemConfigmapRootCA.Files()[0].Data, templateData),
		"kube-system-secret-etcd-client.yaml":                     applyTemplateData(kubeSystemSecretEtcdClient.Files()[0].Data, templateData),
		"openshift-config-secret-etcd-metrics-client.yaml":        applyTemplateData(openshiftConfigSecretEtcdMetricsClient.Files()[0].Data, templateData),
		"openshift-config-configmap-etcd-metrics-serving-ca.yaml": applyTemplateData(openshiftConfigConfigmapEtcdMetricsServingCA.Files()[0].Data, templateData),

		"04-openshift-machine-config-operator.yaml": []byte(openshiftMachineConfigOperator.Files()[0].Data),
		"etcd-service.yaml":                         []byte(etcdServiceKubeSystem.Files()[0].Data),
		"host-etcd-service.yaml":                    []byte(hostEtcdServiceKubeSystem.Files()[0].Data),
	}

	files := make([]*asset.File, 0, len(assetData))
	for name, data := range assetData {
		files = append(files, &asset.File{
			Filename: filepath.Join(manifestDir, name),
			Data:     data,
		})
	}

	return files
}

func applyTemplateData(data []byte, templateData interface{}) []byte {
	template := template.Must(template.New("template").Funcs(customTmplFuncs).Parse(string(data)))
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Load returns the manifests asset from disk.
func (m *Manifests) Load(f asset.FileFetcher) (bool, error) {
	fileList, err := f.FetchByPattern(filepath.Join(manifestDir, "*"))
	if err != nil {
		return false, err
	}
	if len(fileList) == 0 {
		return false, nil
	}

	kubeSysConfig := &configurationObject{}
	var found bool
	for _, file := range fileList {
		if file.Filename == kubeSysConfigPath {
			if err := yaml.Unmarshal(file.Data, kubeSysConfig); err != nil {
				return false, errors.Wrap(err, "failed to unmarshal cluster-config.yaml")
			}
			found = true
		}
	}

	if !found {
		return false, nil

	}

	m.FileList, m.KubeSysConfig = fileList, kubeSysConfig

	asset.SortFiles(m.FileList)

	return true, nil
}

func indent(indention int, v string) string {
	newline := "\n" + strings.Repeat(" ", indention)
	return strings.Replace(v, "\n", newline, -1)
}
