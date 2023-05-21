package recordgenerator

import (
	"errors"
	"testing"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
	"github.com/Khatunzev-Anton/filesorter/internal/services/recordgenerator/mocks"
	"github.com/stretchr/testify/mock"
)

func TestNewNameSalaryGenerator(t *testing.T) {
	namesrepook := mocks.NewNamesRepository(t)
	namesrepook.
		On("Names", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Return([]string{"John Smith", "Bill Burr"}, nil)

	namesreposhort := mocks.NewNamesRepository(t)
	namesreposhort.
		On("Names", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Once().
		Return([]string{"Li"}, nil)

	namesrepoempty := mocks.NewNamesRepository(t)
	namesrepoempty.
		On("Names", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Once().
		Return([]string{}, nil)

	namesrepoerr := mocks.NewNamesRepository(t)
	namesrepoerr.
		On("Names", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Once().
		Return(nil, errors.New("any error"))

	type args struct {
		rd        models.RecordDescriptor
		namesrepo NamesRepository
		minsalary uint32
		maxsalary uint32
	}

	rd := models.RecordDescriptor{
		models.FieldDescriptor{Name: "name", Type: "string", Size: 64},
		models.FieldDescriptor{Name: "salary", Type: "uint32", Size: 4},
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "namesrepook",
			args:    args{rd: rd, namesrepo: namesrepook, minsalary: 100000, maxsalary: 250000},
			wantErr: false,
		},
		{
			name:    "namesreposhort",
			args:    args{rd: rd, namesrepo: namesreposhort, minsalary: 100000, maxsalary: 250000},
			wantErr: true,
		},
		{
			name:    "namesrepoempty",
			args:    args{rd: rd, namesrepo: namesrepoempty, minsalary: 100000, maxsalary: 250000},
			wantErr: true,
		},
		{
			name:    "namesrepoerr",
			args:    args{rd: rd, namesrepo: namesrepoerr, minsalary: 100000, maxsalary: 250000},
			wantErr: true,
		},
		{
			name: "unknownfield",
			args: args{
				rd:        models.RecordDescriptor{models.FieldDescriptor{Name: "unknown field", Type: "string", Size: 32}},
				namesrepo: namesrepook,
				minsalary: 100000,
				maxsalary: 250000,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewNameSalaryGenerator(tt.args.rd, tt.args.namesrepo, tt.args.minsalary, tt.args.maxsalary)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewNameSalaryGenerator() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			_, err = got.GenerateRecord()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
