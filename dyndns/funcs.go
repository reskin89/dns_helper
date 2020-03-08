package dyndns

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
	"github.com/aws/aws-sdk-go/aws/session"
)

//SessionObj will automatically create an AWS session object looking for environment vars
var SessionObj *session.Session

func init() {
	SessionObj = session.Must(session.New())
}

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

//NewConfiguration returns a configuration to use
func NewConfiguration(fileName string) (*Configuration, error) {
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

	} else {
		config.ZoneID = os.Getenv("DYN_ZONE_ID")
		config.Record = os.Getenv("DYN_DNS_RECORD")
		notify, err := strconv.ParseBool(os.Getenv("DYN_SNS_NOTIFY"))
		if err != nil {
			notify = false
		}
		config.Notify.SNSNotify = notify
		config.Notify.SNSTopic = os.Getenv("SNSTopic")
		config.Notify.SNSMessage = os.Getenv("SNSMessage")
		return config, nil
	}

	return nil, fmt.Errorf("no config found")
}
