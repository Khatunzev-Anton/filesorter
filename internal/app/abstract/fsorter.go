package abstract

import "io"

type FSorter interface {
	Sort(f io.ReadWriteSeeker, field string, from int64, to int64) error
}
