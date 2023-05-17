package fsorter

import (
	"errors"
	"fmt"
	"os"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
)

type fsorterquick struct {
	rd     models.RecordDescriptor
	rdsize int
}

func NewFSorterQuick(rd models.RecordDescriptor) (FSorter, error) {
	return &fsorterquick{
		rd:     rd,
		rdsize: rd.Size(),
	}, nil
}

func (r *fsorterquick) Sort(fname string, field string) error {
	found := false
	for _, fd := range r.rd {
		if fd.Name == field {
			found = true
		}
	}
	if !found {
		return errors.New("invalid field name")
	}

	f, err := os.OpenFile(fname, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file statistics: %w", err)
	}
	left, size := 0, fi.Size()
	if size%int64(r.rdsize) != 0 {
		return fmt.Errorf("invalid file size")
	}
	if size < int64(2*r.rdsize) {
		return nil
	}
	err = r.sort(f, field, int64(left), size-int64(r.rdsize))
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

func (r *fsorterquick) sort(f *os.File, field string, left int64, right int64) error {
	if right-left < int64(2*r.rdsize) {
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

	rightrecord, err := r.readrecordat(f, right)
	if err != nil {
		return err
	}
	var i int64
	start := left
	for i = 0; i < recordscnt; i++ {
		currentoffset := i*int64(r.rdsize) + start
		currentrecord, err := r.readrecordat(f, currentoffset)
		if err != nil {
			return err
		}
		less, err := currentrecord.(models.ComparableByField).LessThan(rightrecord.(models.ComparableByField), field)
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

	err = r.sort(f, field, start, left)
	if err != nil {
		return err
	}
	err = r.sort(f, field, left+int64(r.rdsize), right)
	if err != nil {
		return err
	}

	return nil
}

func (r *fsorterquick) readrecordat(f *os.File, offset int64) (models.SerializableRecord, error) {
	buf := make([]byte, r.rdsize)
	_, err := f.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}
	return models.FromBytes(buf, r.rd)
}

func (r *fsorterquick) swaprecordbuffersat(f *os.File, left int64, right int64) error {
	leftbuf := make([]byte, r.rdsize)
	_, err := f.ReadAt(leftbuf, left)
	if err != nil {
		return err
	}

	rightbuf := make([]byte, r.rdsize)
	_, err = f.ReadAt(rightbuf, right)
	if err != nil {
		return err
	}

	_, err = f.WriteAt(rightbuf, left)
	if err != nil {
		return err
	}

	_, err = f.WriteAt(leftbuf, right)
	if err != nil {
		return err
	}
	return nil
}
