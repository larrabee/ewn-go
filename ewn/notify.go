package ewn

import (
	"fmt"
	"github.com/spf13/viper"
)

// Notify send notifications over all chanels
func Notify(msg *Message, config *viper.Viper) () {
	fmt.Printf("%+v", msg)
	return
}
