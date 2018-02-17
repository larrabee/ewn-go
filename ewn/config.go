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

func InitConfig(configPath string) (error) {
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
  },
  "zabbix": {
    "enabled": true,
    "config": "/etc/zabbix/zabbix_agentd.conf",
    "key": "linux.cron.status"
  }
}
	`
	f, err1 := os.Create(configPath)
	if err1 != nil {
		return err1
	}
	defer f.Close()
	_, err2 := f.Write([]byte(config))
	if err2 != nil {
		return err2
	}
	return nil
}
