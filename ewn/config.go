package ewn

import (
	"os"
	"github.com/spf13/viper"
)


func GetConfig(configPath string) (*viper.Viper, error) {
	var v = viper.New()
	v.SetConfigType("json")
	file, err1 := os.Open(configPath)
	defer file.Close()
	if err1 != nil {
		return nil, err1
	}
	err2 := v.ReadConfig(file)
	return v, err2
}
