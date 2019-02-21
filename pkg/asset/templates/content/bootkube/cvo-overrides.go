package bootkube

import (
	"os"
	"path/filepath"

	"github.com/metalkube/kni-installer/pkg/asset"
	"github.com/metalkube/kni-installer/pkg/asset/templates/content"
)

const (
	cVOOverridesFileName = "cvo-overrides.yaml.template"
)

var _ asset.WritableAsset = (*CVOOverrides)(nil)

// CVOOverrides is an asset that generates the cvo-override.yaml.template file.
// This is a gate to prevent CVO from installing these operators which conflict
// with resources already owned by other operators.
// This files can be dropped when the overrides list becomes empty.
type CVOOverrides struct {
	FileList []*asset.File
}

// Dependencies returns all of the dependencies directly needed by the asset
func (t *CVOOverrides) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Name returns the human-friendly name of the asset.
func (t *CVOOverrides) Name() string {
	return "CVOOverrides"
}

// Generate generates the actual files by this asset
func (t *CVOOverrides) Generate(parents asset.Parents) error {
	fileName := cVOOverridesFileName
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
func (t *CVOOverrides) Files() []*asset.File {
	return t.FileList
}

// Load returns the asset from disk.
func (t *CVOOverrides) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(filepath.Join(content.TemplateDir, cVOOverridesFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	t.FileList = []*asset.File{file}
	return true, nil
}
