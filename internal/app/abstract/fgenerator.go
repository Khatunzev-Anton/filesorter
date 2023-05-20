package abstract

type FGenerator interface {
	GenerateFile(fname string, cnt int) error
}
