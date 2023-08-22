/** kafka.Producer 
  * 
  *  An implementation of a thread-safe Kafka Producer based on the
  *  confluent-kafka-go library.
  *
  *  tcarland@gmail.com
  *  Aug 15, 2023
 **/
package kafka

import (
    "bytes"
    "context"
    "log"

    "github.com/tcarland/kafka-go/utils"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)


/** Producer uses a BufferPool to manage a pool of reusable
  * byte.Buffer objects to avoid the overhead of * allocating 
  * a new buffer for each message.
  * The Produce() method should be run as a goroutine.
 **/
type Producer struct {
    topic     string
    brokers   string
    buffers  *utils.BufferPool
    bpc       chan *bytes.Buffer
    active    bool
}

// -----------------------------------

func NewProducer(brokers string, topic string) *Producer {
    return new(Producer).InitProducer(brokers, topic)
}


func (p *Producer) InitProducer(brokers string, topic string) *Producer {
    p.topic   = topic
    p.brokers = brokers
    p.buffers = utils.NewBufferPool(100)
    p.bpc     = make(chan *bytes.Buffer)
    p.active  = false
    return p
}

// -----------------------------------

// Kafka Producer to be ran as a goroutine
func (p *Producer) Produce(ctx context.Context) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": p.brokers})

	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				log.Printf("Error: %v\n", ev)
			default:
				log.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

    log.Printf("kafka.Producer.Produce() run '%s'", p.topic)
    p.active = true

	for p.active {
        select {
        case <- ctx.Done():
            p.active = false
            continue
        default:
            if ! p.active {
                break
            }

            b := <-p.bpc

		    err := producer.Produce(&kafka.Message {
                TopicPartition:   kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
                Value:            b.Bytes(),
                Headers:        []kafka.Header{{Key: "ipflow-dashboard", Value: []byte(p.topic)}}, 
            }, nil)
		
            if err == nil {
                log.Printf("Produce() event: '%s' ", b.String())
            } else if err.(kafka.Error).Code() == kafka.ErrQueueFull {
                log.Println("Producer queue full")            
		    } else if err.(kafka.Error).Code() != kafka.ErrTimedOut { 
			    log.Printf("Producer error: %v", err)
		    }

            p.buffers.Put(b)
        }
	}

    for producer.Flush(10000) > 0 {
        log.Println("Producer Flush() ...")
    }

    log.Printf("Producer finished for '%s'", p.topic)
	producer.Close()
}

// -----------------------------------

func (p *Producer) SendMessage(msg string) {
    b := p.buffers.Get()
    b.Write([]byte(msg))
    p.bpc <- b
}


func (p *Producer) IsActive() bool {
    return p.active
}

// -----------------------------------

func (p *Producer) CreateTopic(numParts int, replFactor int) {
    admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": p.brokers})

    if err != nil {
        log.Fatal(err.Error())
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    results, err := admin.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             p.topic,
			NumPartitions:     numParts,
			ReplicationFactor: replFactor}})
		
    if err != nil {
        log.Fatal(err.Error())
    }

    for _, result := range results {
        log.Printf("Producer.CreateTopic '%s'\n", result)
    }

    admin.Close()
}

// -----------------------------------
// producer.go