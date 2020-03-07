package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"gopkg.in/yaml.v2"
)

type configuration struct {
	ZoneId string `yaml:"zone_id"`
	Record string `yaml:"record"`
}

var SessionObj *session.Session
var Config *configuration

func init() {
	log.Println("Establishing Session")
	SessionObj = session.Must(session.NewSession())
	log.Println("Checking for Configuration....")
	check := fileExists()
	if check {
		configBytes, err := ioutil.ReadFile("dnsConfig.yml")
		if err != nil {
			log.Fatal("Unable to Read Config!")
		}
		err = yaml.Unmarshal(configBytes, &Config)
		if err != nil {
			log.Fatal("Unable To Parse Config!")
		}
	} else {
		Config.ZoneId = os.Getenv("ZONE_ID")
		Config.Record = os.Getenv("DNS_RECORD")
	}
}

func main() {
	NewIP, err := getPublicIP()
	if err != nil {
		log.Fatal(err)
	}

	err = updateDNS(NewIP)
	if err != nil {
		log.Fatal(err)
	}
}

func fileExists() bool {
	exists, err := os.Stat("dnsConfig.yml")
	if os.IsNotExist(err) {
		return false
	}

	return exists.IsDir()
}

func getPublicIP() ([]byte, error) {
	resp, err := http.Get("ifconfig.me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func updateDNS(ipAddress []byte) error {
	updater := route53.New(SessionObj)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(Config.Record),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(string(ipAddress)),
							},
						},
					},
				},
			},
			Comment: aws.String("Updated Automagically!"),
		},
		HostedZoneId: aws.String(Config.ZoneId),
	}

	resp, err := updater.ChangeResourceRecordSets(params)
	if err != nil {
		return err
	}
	log.Println("Update Response: \n", resp)

	return nil
}
