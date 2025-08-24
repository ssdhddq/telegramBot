package error

import "fmt"

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %v", msg, err)
}
