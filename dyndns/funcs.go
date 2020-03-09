package dyndns

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/yaml.v2"
)

//SessionObj will automatically create an AWS session object looking for environment vars
var SessionObj *session.Session

func init() {
	SessionObj = session.Must(session.New())
}

//fileExists checks to see if the filename provided (if provided) exists, and if not, returns false
func fileExists(fileName string) bool {
	if len(fileName) < 1 {
		return false
	}

	exists, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}

	if exists.IsDir() {
		return false
	}

	return true
}

//NewConfigurationFromEnvironment returns a configuration to use from environment variables
func NewConfigurationFromEnvironment(fileName string) (*Configuration, error) {
	var config *Configuration

	if len(os.Getenv("D2_ZONE_ID")) < 1 {
		return nil, fmt.Errorf("no route53 zoneid defined in env, expect D2_ZONE_ID")
	}
	config.ZoneID = os.Getenv("D2_ZONE_ID")

	if len(os.Getenv("D2_DNS_RECORD")) < 1 {
		return nil, fmt.Errorf("no dns record to modify specified, expect D2_DNS_RECORD")
	}
	config.Record = os.Getenv("D2_DNS_RECORD")
	notify, err := strconv.ParseBool(os.Getenv("D2_SNS_NOTIFY"))
	if err != nil {
		notify = false
	}
	config.Notify.SNSNotify = notify
	config.Notify.SNSTopic = os.Getenv("SNSTopic")
	config.Notify.SNSMessage = os.Getenv("SNSMessage")

	return config, nil
}

//NewConfigurationFromFile returns a pointer to a new configuration loaded from a yaml file
func NewConfigurationFromFile(fileName string) (*Configuration, error) {
	var config *Configuration

	loadFile := fileExists(fileName)
	if loadFile {
		configBytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal("Unable to Read Config!")
			return nil, err
		}
		err = yaml.Unmarshal(configBytes, &config)
		if err != nil {
			log.Fatal("Unable To Parse Config!")
			return nil, err
		}
	}
	return config, nil
}
