package schedule

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-keg/simple/data/ent"
	"github.com/go-keg/simple/data/ent/user"
)

type Daily struct {
	ent      *ent.Client
	producer sarama.SyncProducer
}

func NewDaily(ent *ent.Client, producer sarama.SyncProducer) *Daily {
	return &Daily{ent: ent, producer: producer}
}

// Run mock daily send stat data to accounts
func (r Daily) Run(ctx context.Context) error {
	startID := 0
	for {
		accounts, err := r.ent.User.Query().Where(user.IDGT(startID)).Limit(200).All(ctx)
		if err != nil {
			return err
		}
		if len(accounts) == 0 {
			break
		}
		var messages []*sarama.ProducerMessage
		for i := range accounts {
			messages = append(messages, &sarama.ProducerMessage{
				Topic: "send_daily_stat",
				Key:   sarama.StringEncoder(fmt.Sprintf("account:%d", accounts[i].ID)),
				Value: sarama.StringEncoder("stat data..."),
			})
		}
		err = r.producer.SendMessages(messages)
		if err != nil {
			return err
		}
		startID = accounts[len(accounts)-1].ID
	}
	return nil
}
