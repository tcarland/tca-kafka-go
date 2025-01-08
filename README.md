tca-kafka-go
=============

A library wrapper to the *confluentinc-kafka-go* package that provides added
thread-safe constructs for use with *goroutines*.


## Configuration

Consumer configuration is defined via the *KafkaSite* object,
which is typically provided from an input yaml configuration. The application
implementing this library would typically have its own configuration needs 
so there are no direct yaml parsing functions provided here. The application
config would simply include a top-level config object such as `kafka` to define 
the parameters like the following yaml snippet.
```yaml
kafka:
  uswest1:
    brokers: "foo1:9094,foo2:9094,foo3:9094"
    topic: "mytopic"
    gid: "grp1"
    streamreset: true
    replicationfactor: 3
    partitions: 3
  lab1:
    brokers: "localhost:9090"
    topic: "testtopic"
    gid: "test"
    streamreset: true
```

A *Config* object or struct would consist of various objects representing 
other parts of the application config as well as `map[string]*KafkaSite` 
for this library.
```golang
package config

import "github.com/tcarland/tca-kafka-go/config"

type Config struct {
    Sites   map[string]*config.KafkaSite `yaml:"kafka"` 
    [ ... ]
}
```

The config object can be generated manually via the *config.InitKafkaSite()* function
```golang
func (k *KafkaSite) InitKafkaSite(brokers string, topic string, gid string) *KafkaSite {}
```

The *Consumer* object takes the *KafkaSite* object at construction, though some 
options do not apply to the consumer. *Replicationfactor* and *Partitions* applies 
only to Producers, while the *GroupId(gid)* and *streamreset* options applies 
to Consumers