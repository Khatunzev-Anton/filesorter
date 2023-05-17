package models

type FieldDescriptor struct {
	Name string `yaml:"name"`
	Size int    `yaml:"size"`
	Type string `yaml:"type"`
}
