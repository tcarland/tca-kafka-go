/** kafka.Consumer
  *   
  *  An implementation of a thread-safe Kafka Consumer based on the
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
    "time"

    "github.com/tcarland/kafka-go/config"
    "github.com/tcarland/kafka-go/utils"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)


type Consumer struct {
    name        string
    bpc         chan *bytes.Buffer
    buffers    *utils.BufferPool
    msglist    *utils.SyncList
    site       *config.KafkaSite
    reset       int
    active      bool
}

// -----------------------------------

func NewConsumer( name string, site *config.KafkaSite ) *Consumer {
    return new(Consumer).InitConsumer(name, site)
}


func (c *Consumer) InitConsumer ( name string, site *config.KafkaSite ) *Consumer {
    c.name    = name
    c.site    = site
    c.bpc     = make(chan *bytes.Buffer)
    c.buffers = utils.NewBufferPool(100)
    c.msglist = utils.NewSyncList()
    c.reset   = 0
    c.active  = false
    return c
}


// Kafka Consumer goroutine
func (c *Consumer) Consume(ctx context.Context) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     c.site.Brokers,
        "broker.address.family": "v4",
		"group.id":              c.site.GroupId,
		"auto.offset.reset":    "latest",
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	consumer.SubscribeTopics([]string{c.site.Topic}, nil)  
    // .SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

    log.Printf("kafka.Consumer.Consume() run '%s'", c.name)
    c.site.Active = true
    c.active      = true

	for c.active {
        select {
        case <- ctx.Done():
            log.Println("Consumer.Consume() Context done")
            c.active = false
            close(c.bpc)
            continue
        default:
            b := c.buffers.Get()

		    msg, err := consumer.ReadMessage(time.Second * 6)
		
            if err == nil {
                b.Write(msg.Value)
                c.bpc <- b
		    } else if err.(kafka.Error).Code() != kafka.ErrTimedOut { 
			    log.Printf("Consumer error: %v (%v)\n", err, msg)
		    } else {
                if c.site.DoReset {
                    c.reset++
                }
                c.buffers.Put(b)
            }
        }
	}
    c.site.Active = false

    log.Printf("Consumer.Consume() finished for '%s'", c.name)
	consumer.Close()
}


// Process goroutine
func (c *Consumer) Process(ctx context.Context) {
    log.Printf("kafka.Consumer.Process() run '%s'", c.name)
    
    c.active = true
    for c.active {
        select {
        case <- ctx.Done():
            log.Println("Consumer.Process() Context done")
            c.active = false
            continue
        default:
            b := <-c.bpc
            if ! c.active {
                break
            }

            c.msglist.Lock()
            if c.site.DoReset && c.reset > 2 {
                log.Printf("Consumer.Process() -> SyncList, of %d items, reset event.", c.msglist.Size())
                c.msglist.Clear()
                c.reset = 0
            }
            str := b.String()
            c.msglist.PushBack(str)
            c.msglist.Unlock()
            c.buffers.Put(b)
        }
    }
    log.Printf("Consumer.Process() finished for '%s'", c.name) 
}


func (c *Consumer) IsActive() bool {
    return c.active
}


func (c *Consumer) GetMsgList() *utils.SyncList {
    return c.msglist
}


func (c* Consumer) GetSiteConfig() *config.KafkaSite {
    return c.site
}
