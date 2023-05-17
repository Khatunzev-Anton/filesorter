package mainconfig

import (
	"fmt"
	"os"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
	"gopkg.in/yaml.v2"
)

type MainConfig interface {
	GetRecordDescriptor() models.RecordDescriptor
}

type mainconfig struct {
	Record models.RecordDescriptor `yaml:"record"`
}

func NewMainConfig(fname string) (MainConfig, error) {
	mc := &mainconfig{}

	bytes, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("filed to open mainconfig file: %w", err)
	}

	err = yaml.Unmarshal(bytes, &mc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mainconfig file: %w", err)
	}

	return mc, err
}

func (r *mainconfig) GetRecordDescriptor() models.RecordDescriptor {
	return r.Record
}
