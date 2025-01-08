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

    "github.com/tcarland/tca-kafka-go/config"
    "github.com/tcarland/tca-kafka-go/utils"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)


/** Producer uses a BufferPool to manage a pool of reusable
  * byte.Buffer objects to avoid the overhead of * allocating 
  * a new buffer for each message.
  * The Produce() method should be run as a goroutine.
 **/
type Producer struct {
    name      string
    site     *config.KafkaSite
    buffers  *utils.BufferPool
    bpc       chan *bytes.Buffer
    active    bool
}

// -----------------------------------

func NewProducer( name string, site *config.KafkaSite ) *Producer {
    return new(Producer).InitProducer(name, site)
}

func (p *Producer) InitProducer( name string, site *config.KafkaSite ) *Producer {
    p.name    = name
    p.site    = site
    p.buffers = utils.NewBufferPool(100)
    p.bpc     = make(chan *bytes.Buffer)
    p.active  = false
    return p
}


// -----------------------------------


func (p *Producer) Version() string{
    return config.Version
}


// Kafka Producer to be ran as a goroutine
func (p *Producer) Produce(ctx context.Context) {
    producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": p.site.Brokers})

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

    log.Printf("kafka.Producer.Produce() run '%s'", p.site.Topic)
    p.active = true
    p.site.Active = true

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
                TopicPartition:   kafka.TopicPartition{Topic: &p.site.Topic, Partition: kafka.PartitionAny},
                Value:            b.Bytes(),
                Headers:        []kafka.Header{{Key: "ipflow-dashboard", Value: []byte(p.site.Topic)}}, 
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
    p.site.Active = false

    log.Printf("Producer finished for '%s'", p.site.Topic)
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

func (p *Producer) CreateTopic() {
    admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": p.site.Brokers})

    if err != nil {
        log.Fatal(err.Error())
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    results, err := admin.CreateTopics(
        ctx,
        []kafka.TopicSpecification{ {
            Topic:             p.site.Topic,
            NumPartitions:     p.site.Partitions,
            ReplicationFactor: p.site.Replicas,
        } })
        
    if err != nil {
        log.Fatal(err.Error())
    }

    for _, result := range results {
        log.Printf("Producer.CreateTopic '%s'\n", result)
    }

    admin.Close()
}
