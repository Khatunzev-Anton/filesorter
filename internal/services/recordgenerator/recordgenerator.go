package recordgenerator

import "github.com/Khatunzev-Anton/filesorter/internal/models"

//go:generate mockery --name RecordGenerator
type RecordGenerator interface {
	GenerateRecord() (models.SerializableRecord, error)
}
