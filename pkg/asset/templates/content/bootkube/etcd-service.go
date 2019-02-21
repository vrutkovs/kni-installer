package bootkube

import (
	"os"
	"path/filepath"

	"github.com/metalkube/kni-installer/pkg/asset"
	"github.com/metalkube/kni-installer/pkg/asset/templates/content"
)

const (
	etcdServiceKubeSystemFileName = "etcd-service.yaml"
)

var _ asset.WritableAsset = (*EtcdServiceKubeSystem)(nil)

// EtcdServiceKubeSystem is the constant to represent contents of etcd-service.yaml file
type EtcdServiceKubeSystem struct {
	FileList []*asset.File
}

// Dependencies returns all of the dependencies directly needed by the asset
func (t *EtcdServiceKubeSystem) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Name returns the human-friendly name of the asset.
func (t *EtcdServiceKubeSystem) Name() string {
	return "EtcdServiceKubeSystem"
}

// Generate generates the actual files by this asset
func (t *EtcdServiceKubeSystem) Generate(parents asset.Parents) error {
	fileName := etcdServiceKubeSystemFileName
	data, err := content.GetBootkubeTemplate(fileName)
	if err != nil {
		return err
	}
	t.FileList = []*asset.File{
		{
			Filename: filepath.Join(content.TemplateDir, fileName),
			Data:     []byte(data),
		},
	}
	return nil
}

// Files returns the files generated by the asset.
func (t *EtcdServiceKubeSystem) Files() []*asset.File {
	return t.FileList
}

// Load returns the asset from disk.
func (t *EtcdServiceKubeSystem) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(filepath.Join(content.TemplateDir, etcdServiceKubeSystemFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	t.FileList = []*asset.File{file}
	return true, nil
}
