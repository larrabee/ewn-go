# Execute With Notify
[![Go Report Card](https://goreportcard.com/badge/github.com/larrabee/ewn-go)](https://goreportcard.com/report/github.com/larrabee/ewn-go) [![Build Status](https://travis-ci.org/larrabee/ewn-go.svg?branch=master)](https://travis-ci.org/larrabee/ewn-go)   
This is wrapper script, what can notify you when your tasks failed.

## Features
1. Notify with email
2. Publish logs to Graylog/Logstash
3. Prevent duplicate processes
4. Retry failed commands


## Install
1. Download binary from [Release page](https://github.com/larrabee/ewn-go/releases) or build it with `go build main.go`
2. Run `ewn --initconfig`. This command create default config file
3. Update config file with your parameters
## Configure
You must modify file `/etc/ewn.conf`

## usage

* Exec command with notify for default recipients ('user1@example.com' and 'user2@example.com'):  
`ewn -c 'echo "This is test command" && /bin/false'`  
* Exec command with notify for custom recipients ('user3@example.com' and 'user4@example.com'):  
`ewn -c 'echo "This is test command" && /bin/false' -r "user3@example.com" "user4@example.com"`  
* Exec command with comment message. It will be attached to email and Graylog message:  
`ewn -c 'echo "This is test command" && /bin/false' --comment "This command will always fail"`    
* Exec command that prevent duplicate processes:  
`ewn -c 'echo "This is test command" && /bin/sleep 3600' --dont-duplicate`  
* Exec command with 3 retry on failure and 30 sec interval between retries:  
`ewn -c 'echo "This is test command" && /bin/false' --retry 3 --retry-sleep 30`  
* Exec command with custom valid exit codes (default is 0). If command exit with exit code 1, 3 or 255 email will not sent.  
`ewn -c 'echo "This is test command" && /bin/false' --valid-exitcodes 1 3 255`  

You can combine any keys like:  
`ewn -c 'echo "This is test command" && /bin/false' -d --dont-duplicate -r "user3@example.com" "user4@example.com" --retry 3 --retry-sleep 30 --valid-exitcodes 1  --comment "This command will never fail because 1 in valid exitcodes"`  

## License
GPL v3
