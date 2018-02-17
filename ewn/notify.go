package ewn

import (
	"fmt"
)

// Notify send notifications over all chanels
func Notify(strs ...interface{}) error {
	for _, str := range strs {
		fmt.Println(str)
	}
	return nil
}
