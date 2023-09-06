package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
)

func main() {
	spiceDbUrl := os.Getenv("SPICEDB_URL")
	spiceDbPSK := os.Getenv("SPICEDB_PSK")
	spiceDbClient, err := getSpiceDbClient(spiceDbUrl, spiceDbPSK)
	if err != nil {
		fmt.Printf("Failed to connect to SpiceDB!: %s", err)
		return
	}

	kafkaUrl := os.Getenv("KAFKA_URL")
	topic := os.Getenv("KAFKA_TOPIC")
	kafkaReader := getKafkaReader(kafkaUrl, topic)

	defer kafkaReader.Close()

	fmt.Println("Consuming!")
	for {
		m, err := kafkaReader.ReadMessage(context.TODO())
		if err != nil {
			fmt.Printf("Error reading from topic: %s\n", err)
			return
		}

		var evt DispatcherRunEvent
		if err = json.Unmarshal(m.Value, &evt); err != nil {
			fmt.Printf("Error unmarshalling message: %s, original: %s", err, string(m.Value))
			kafkaReader.CommitMessages(context.TODO(), m)
			continue
		}

		if evt.EventType == "create" {
			fmt.Printf("New run %s for org %s for service %s against host %s\n", evt.Payload.ID, evt.Payload.OrgID, evt.Payload.Service, evt.Payload.Recipient)

			_, err = spiceDbClient.WriteRelationships(context.TODO(), &v1.WriteRelationshipsRequest{
				Updates: []*v1.RelationshipUpdate{
					{
						Operation: v1.RelationshipUpdate_OPERATION_TOUCH,
						Relationship: &v1.Relationship{
							Resource: &v1.ObjectReference{
								ObjectType: "dispatcher/run",
								ObjectId:   evt.Payload.ID,
							},
							Relation: "host",
							Subject: &v1.SubjectReference{
								Object: &v1.ObjectReference{
									ObjectType: "inventory/host",
									ObjectId:   evt.Payload.Recipient,
								},
							},
						},
					},
					{
						Operation: v1.RelationshipUpdate_OPERATION_TOUCH,
						Relationship: &v1.Relationship{
							Resource: &v1.ObjectReference{
								ObjectType: "dispatcher/run",
								ObjectId:   evt.Payload.ID,
							},
							Relation: "service",
							Subject: &v1.SubjectReference{
								Object: &v1.ObjectReference{
									ObjectType: "dispatcher/service",
									ObjectId:   evt.Payload.Service,
								},
							},
						},
					},
				},
			})

			if err != nil {
				fmt.Printf("Failed to update SpiceDB: %s\n", err)
			}

		} else {
			fmt.Printf("Ignoring event of type %s\n", evt.EventType)
		}
		kafkaReader.CommitMessages(context.TODO(), m)

	}
}

func getSpiceDbClient(endpoint string, presharedKey string) (*authzed.Client, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithBlock())

	opts = append(opts, grpcutil.WithInsecureBearerToken(presharedKey))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return authzed.NewClient(
		endpoint,
		opts...,
	)
}

func getKafkaReader(kafkaURL, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

type DispatcherRunEvent struct {
	EventType string `json:"event_type"`
	Payload   struct {
		ID            string `json:"id"`
		OrgID         string `json:"org_id"`
		Recipient     string `json:"recipient"`
		CorrelationID string `json:"correlation_id"`
		Service       string `json:"service"`
		URL           string `json:"url"`
		Labels        struct {
			RemediationID string `json:"remediation_id"`
		} `json:"labels"`
		Name            string `json:"name"`
		WebConsoleURL   string `json:"web_console_url"`
		RecipientConfig struct {
			SatID    string `json:"sat_id"`
			SatOrgID string `json:"sat_org_id"`
		} `json:"recipient_config"`
		Status    string    `json:"status"`
		Timeout   int       `json:"timeout"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"payload"`
}
