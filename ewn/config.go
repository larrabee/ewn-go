package ewn

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// GetConfig return config struct
func GetConfig(configPath string) (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigType("json")

	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	err = vp.ReadConfig(file)
	if err != nil {
		return nil, err
	}
	return vp, nil
}

// InitConfig create default config file in given path
func InitConfig(configPath string) error {
	config := `{
  "email": {
    "enabled": true,
    "host": "mail.example.com",
    "user": "notify@example.com",
    "from": null,
    "secure": false,
    "pass": "password",
    "recipients": [
      "user1@example.com", 
      "user2@example.com"
    ]
  },
  "graylog": {
    "enabled": true,
    "host": "graylog.example.com",
    "port": 12206,
    "tag": "ewn",
    "mtu": 1400
  }
}
	`
	err := ioutil.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		return err
	}
	return nil
}
