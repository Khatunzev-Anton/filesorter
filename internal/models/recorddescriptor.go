package models

type RecordDescriptor []FieldDescriptor

func (r *RecordDescriptor) Size() int {
	result := 0
	for _, fd := range *r {
		result += fd.Size
	}
	return result
}
