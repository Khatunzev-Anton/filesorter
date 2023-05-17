package fsorter

type FSorter interface {
	Sort(fname string, field string) error
}
