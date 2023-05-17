package repositories

import (
	"bufio"
	"os"
)

type NamesRepository interface {
	Names(page int, pagesize int) ([]string, error)
}

func NewNameRepository(fname string) (NamesRepository, error) {
	return &namesRepository{fname: fname}, nil
}

type namesRepository struct {
	fname string
}

func (r *namesRepository) Names(page int, pagesize int) ([]string, error) {
	file, err := os.Open(r.fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	recordnum := 0
	pagestart := page * pagesize
	pageend := (page + 1) * pagesize
	result := make([]string, 0, pagesize)
	for scanner.Scan() {

		if recordnum >= pagestart && recordnum < pageend {
			result = append(result, scanner.Text())
		}
		if recordnum >= pageend {
			break
		}
		recordnum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
