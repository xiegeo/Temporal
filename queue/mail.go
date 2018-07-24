package queue

import (
	"encoding/json"
	"fmt"

	"github.com/RTradeLtd/Temporal/config"
	"github.com/RTradeLtd/Temporal/mail"
	"github.com/streadway/amqp"
)

var (
	// IpfsPinFailedContent is a to-be formatted message sent on IPFS pin failures
	IpfsPinFailedContent = "Pin failed for content hash %s on IPFS network %s, for reason %s"
	// IpfsPinFailedSubject is a subject for IPFS pin failed messages
	IpfsPinFailedSubject = "IPFS Pin Failed"
	// IpfsPrivateNetworkUnauthorizedSubject is a subject whenever someone tries to access a bad private network
	IpfsPrivateNetworkUnauthorizedSubject = "Unauthorized access to IPFS private network"
	// IpfsInitializationFailedSubject is a subject used when connecting to ipfs fails
	IpfsInitializationFailedSubject = "Connection to IPFS failed"
)

// EmailSend is a helper struct used to contained formatted content ot send as an email
type EmailSend struct {
	Subject      string   `json:"subject"`
	Content      string   `json:"content"`
	ContentType  string   `json:"content_type"`
	EthAddresses []string `json:"eth_addresses"`
}

// ProcessMailSends is a function used to process mail send queue messages
func ProcessMailSends(msgs <-chan amqp.Delivery, tCfg *config.TemporalConfig) error {
	mm, err := mail.GenerateMailManager(tCfg)
	if err != nil {
		return err
	}
	for d := range msgs {
		fmt.Println("email send notification detected")
		es := EmailSend{}
		err = json.Unmarshal(d.Body, &es)
		if err != nil {
			fmt.Println("error unmarshaling", err)
			d.Ack(false)
			continue
		}
		emails := make(map[string]string)
		for _, v := range es.EthAddresses {
			resp, err := mm.UserManager.FindEmailByAddress(v)
			if err != nil {
				//TODO: decide on how this should be handled
				fmt.Println("error finding email address for user", err)
				d.Ack(false)
				continue
			}
			emails[v] = resp[v]
		}
		for k, v := range emails {
			_, err = mm.SendEmail(es.Subject, es.Content, es.ContentType, k, v)
			if err != nil {
				//TODO: decide on how this should be handled
				fmt.Printf("error sending message to %s for error %s", v, err)
			}
		}
		fmt.Println("emails sent")
	}
	return nil
}
