package awssqs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/goutil"
)

func NewSubscriber(sess *session.Session, topic string, queueURL string) *Subscriber {

	return &Subscriber{
		client:   sqs.New(sess),
		queueURL: queueURL,
		timeout:  1,
		topic:    topic,
	}
}

type Subscriber struct {
	processor func(message string)
	client    *sqs.SQS
	queueURL  string
	stop      chan bool
	timeout   int64
	topic     string
}

func (sub *Subscriber) Unsubscribe() {
	//Todo: apply this function
	sub.stop <- true
}

func (sub *Subscriber) SetProcessor(processor gopubsubengine.MessageProcessor) {

	sub.processor = processor
}

func (sub *Subscriber) getMessage() {

	//fmt.Println("get Message", sub.queueURL)

	msgResult, err := sub.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &sub.queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   &sub.timeout,
	})

	if err != nil {

		fmt.Println("error receiving message ", err)
		return
	}
	//fmt.Println("receive ", len(msgResult.Messages), " messages")

	if len(msgResult.Messages) > 0 {

		//fmt.Println("Message ID:     " + *msgResult.Messages[0].MessageId)

		//fmt.Println("Message Handle: "+*msgResult.Messages[0].ReceiptHandle, *msgResult.Messages[0].Body)

		for _, msg := range msgResult.Messages {

			if sub.processor != nil {

				groupID := msg.Attributes[sqs.MessageSystemAttributeNameMessageGroupId]
				if *groupID == sub.topic {

					sub.processor(*msg.Body)

					_, err := sub.client.DeleteMessage(&sqs.DeleteMessageInput{

						QueueUrl:      &sub.queueURL,
						ReceiptHandle: msg.ReceiptHandle,
					})
					if err != nil {
						fmt.Println("error deleting message ", err)
					}
				} else {
					// fmt.Println("received on topic", *groupID)
				}
			}

		}
	}
}

func (sub *Subscriber) start() {
	sub.stop = goutil.Schedule(sub.getMessage, time.Microsecond*100)
	//TODO : Change the default schedule time
}
