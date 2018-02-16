package ewn

import (
  "fmt"
)

func Notify(strs ...interface{}) (error) {
  for _, str := range strs {
    fmt.Println(str)
  }
  return nil
}
