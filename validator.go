package opt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func IsFile[T ~string](value T) error {
	_, err := os.Stat(string(value))
	return err
}

func IsMinMax(min, max int) func(int) error {
	return func(value int) error {
		if value < min || value > max {
			v := strconv.Itoa(value)
			s := fmt.Sprint("value is not in range ", min, " <= '", v, "' <= ", max)
			return errors.New(s)
		}
		return nil
	}
}
