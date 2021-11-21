package awssqs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewPublisher(sess *session.Session, queueURL string, topic string) *Publisher {
	return &Publisher{
		client:         sqs.New(sess),
		queueURL:       queueURL,
		messageGroupID: topic,
	}
}

type Publisher struct {
	client         *sqs.SQS
	queueURL       string
	messageGroupID string
}

func (publisher *Publisher) Unpublish() {
	//TODO: apply this function
}
func (publisher *Publisher) Publish(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	_, err = publisher.client.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:   aws.Int64(0),
		MessageBody:    aws.String(string(data)),
		QueueUrl:       &publisher.queueURL,
		MessageGroupId: &publisher.messageGroupID,
	})
	if err != nil {
		fmt.Println("send error ", err)
	}
}

func (publisher *Publisher) PublishAttributes(message interface{}, attributes map[string]string) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	attrs := map[string]*sqs.MessageAttributeValue{}
	for key, value := range attributes {
		attrs[key] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}
	_, err = publisher.client.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(0),
		MessageAttributes: attrs,
		MessageBody:       aws.String(string(data)),
		QueueUrl:          &publisher.queueURL,
		MessageGroupId:    &publisher.messageGroupID,
	})
	if err != nil {
		fmt.Println("send error ", err)
	}
}
