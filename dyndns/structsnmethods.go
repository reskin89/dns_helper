package dyndns

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

//Configuration loads the config into memory
type Configuration struct {
	ZoneID string    `yaml:"zone_id"`
	Record string    `yaml:"record"`
	Notify SNSNotify `yaml:"sns_notify,omitempty"`
	IP     []byte
}

//SNSNotify is for optional sns notifcation on success or failure of dns updates.
type SNSNotify struct {
	SNSNotify  bool   `yaml:"sns_notify"`
	SNSTopic   string `yaml:"sns_topic"`
	SNSMessage string `yaml:"sns_message"`
}

func (c Configuration) GetPublicIP() error {
	resp, err := http.Get("ifconfig.me")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.IP = data
	return nil
}

func (c Configuration) UpdateDNS() error {
	updater := route53.New(SessionObj)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(c.Record),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(string(c.IP)),
							},
						},
					},
				},
			},
			Comment: aws.String("Updated Automagically!"),
		},
		HostedZoneId: aws.String(c.ZoneID),
	}

	resp, err := updater.ChangeResourceRecordSets(params)
	if err != nil {
		return err
	}
	log.Println("Update Response: \n", resp)

	return nil
}
