package recordgenerator

import "github.com/Khatunzev-Anton/filesorter/internal/models"

type RecordGenerator interface {
	GenerateRecord() (models.SerializableRecord, error)
}
