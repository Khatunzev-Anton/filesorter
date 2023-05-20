package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Khatunzev-Anton/filesorter/internal/app/abstract"
	"github.com/Khatunzev-Anton/filesorter/internal/mainconfig"
	"github.com/Khatunzev-Anton/filesorter/internal/models"
	"github.com/Khatunzev-Anton/filesorter/internal/repositories"
	"github.com/Khatunzev-Anton/filesorter/internal/services/fgenerator"
	"github.com/Khatunzev-Anton/filesorter/internal/services/freader"
	"github.com/Khatunzev-Anton/filesorter/internal/services/fsorter"
	"github.com/Khatunzev-Anton/filesorter/internal/services/recordgenerator"
)

type App interface {
	abstract.FGenerator
	abstract.FReader
	abstract.FSorter
	RecordSize() int64
}

type app struct {
	config     mainconfig.MainConfig
	fgenerator fgenerator.FGenerator
	freader    freader.FReader
	fsorter    fsorter.FSorter
}

func NewApp(cfgFile string) (App, error) {
	config, err := mainconfig.NewMainConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get current executable directory: %w", err)
	}
	namesrepo, err := repositories.NewNameRepository(fmt.Sprintf("%[1]s/internal/data/names.txt", filepath.Dir(ex))) //???
	if err != nil {
		return nil, fmt.Errorf("failed to initialize namesrepo: %w", err)
	}

	g, err := recordgenerator.NewNameSalaryGenerator(config.GetRecordDescriptor(), namesrepo, 100000, 250000)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize generator: %w", err)
	}

	filegenerator, err := fgenerator.NewFGenerator(g)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file generator: %w", err)
	}

	filereader, err := freader.NewFReader(config.GetRecordDescriptor())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate freader: %w", err)
	}

	filesorter, err := fsorter.NewFSorterQuick(config.GetRecordDescriptor(), filereader)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate fsorter: %w", err)
	}

	application := &app{
		config:     config,
		fgenerator: filegenerator,
		freader:    filereader,
		fsorter:    filesorter,
	}
	return application, nil
}

func (r *app) GenerateFile(fname string, cnt int) error {
	return r.fgenerator.GenerateFile(fname, cnt)
}

func (r *app) ReadSerializableRecord(f io.ReadSeeker, offset int64) (models.SerializableRecord, error) {
	return r.freader.ReadSerializableRecord(f, offset)
}

func (r *app) ReadComparableByField(f io.ReadSeeker, offset int64) (models.ComparableByField, error) {
	return r.freader.ReadComparableByField(f, offset)
}

func (r *app) Sort(f io.ReadWriteSeeker, field string, from int64, to int64) error {
	return r.fsorter.Sort(f, field, from, to)
}

func (r *app) RecordSize() int64 {
	rd := r.config.GetRecordDescriptor()
	return int64(rd.Size())
}
