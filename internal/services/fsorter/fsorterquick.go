package fsorter

import (
	"fmt"
	"io"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
)

type FReader interface {
	ReadComparableByField(f io.ReadSeeker, offset int64) (models.ComparableByField, error)
}

type fsorterquick struct {
	rd      models.RecordDescriptor
	freader FReader
	rdsize  int
}

func NewFSorterQuick(rd models.RecordDescriptor, freader FReader) (FSorter, error) {
	return &fsorterquick{
		rd:      rd,
		rdsize:  rd.Size(),
		freader: freader,
	}, nil
}

func (r *fsorterquick) Sort(f io.ReadWriteSeeker, field string, from int64, to int64) error {
	if (to-from)%int64(r.rdsize) != 0 {
		return fmt.Errorf("invalid file size")
	}
	if (to - from) < int64(2*r.rdsize) {
		return nil
	}
	err := r.sort(f, field, from, to-int64(r.rdsize))
	if err != nil {
		return err
	}
	return nil
}

/****
QuickSort:
	func sortColors(nums []int)  {
		if len(nums) < 2 {
			return
		}
		left,right := 0,len(nums) - 1

		center := rand.Intn(right)

		nums[center], nums[right] = nums[right],nums[center]
		for i := range nums {
			if nums[i] < nums[right] {
				nums[left], nums[i] = nums[i], nums[left]
				left++
			}
		}
		nums[right],nums[left] = nums[left],nums[right]

		sortColors(nums[:left])
		sortColors(nums[left + 1:])
	}
*/

func (r *fsorterquick) sort(f io.ReadWriteSeeker, field string, left int64, right int64) error {
	if right-left < int64(r.rdsize) {
		return nil
	}
	center := (left + right + int64(r.rdsize)) / 2
	centeroffset := center - (center % int64(r.rdsize))
	recordscnt := (right - left + int64(r.rdsize)) / int64(r.rdsize)

	fmt.Printf("\r\n===SEED %[1]d===", center)
	err := r.swaprecordbuffersat(f, centeroffset, right)
	if err != nil {
		return err
	}

	rightrecord, err := r.freader.ReadComparableByField(f, right)
	if err != nil {
		return err
	}
	var i int64
	start := left
	for i = 0; i < recordscnt; i++ {
		currentoffset := i*int64(r.rdsize) + start
		currentrecord, err := r.freader.ReadComparableByField(f, currentoffset)
		if err != nil {
			return err
		}
		less, err := currentrecord.LessThan(rightrecord, field)
		if err != nil {
			return err
		}

		if less {
			err = r.swaprecordbuffersat(f, left, currentoffset)
			if err != nil {
				return err
			}
			left += int64(r.rdsize)
		}
	}
	err = r.swaprecordbuffersat(f, left, right)
	if err != nil {
		return err
	}

	err = r.sort(f, field, start, left-int64(r.rdsize))
	if err != nil {
		return err
	}
	err = r.sort(f, field, left+int64(r.rdsize), right)
	if err != nil {
		return err
	}

	return nil
}

func (r *fsorterquick) swaprecordbuffersat(f io.ReadWriteSeeker, left int64, right int64) error {

	_, err := f.Seek(left, 0)
	if err != nil {
		return err
	}
	leftbuf := make([]byte, r.rdsize)
	_, err = f.Read(leftbuf)
	if err != nil {
		return err
	}

	_, err = f.Seek(right, 0)
	if err != nil {
		return err
	}
	rightbuf := make([]byte, r.rdsize)
	_, err = f.Read(rightbuf)
	if err != nil {
		return err
	}

	_, err = f.Seek(left, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(rightbuf)
	if err != nil {
		return err
	}

	_, err = f.Seek(right, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(leftbuf)
	if err != nil {
		return err
	}
	return nil
}
