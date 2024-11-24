/**  
  *  Application config parameters
  *
 **/ 
package config

import ()

var Version string = "0.3.0"

type KafkaSite struct {
    Brokers      string `yaml:"brokers"`
    Topic        string `yaml:"topic"`
    GroupId      string `yaml:"gid"`
    DoReset      bool   `yaml:"streamreset"`
    Replicas     int    `yaml:"replicationfactor"`
    Partitions   int    `yaml:"partitions"`
    Active       bool
}

func NewKafkaSite(brokers string, topic string, gid string) *KafkaSite {
    return new(KafkaSite).InitKafkaSite(brokers, topic, gid)
}

func (k *KafkaSite) InitKafkaSite(brokers string, topic string, gid string) *KafkaSite {
    k.Brokers    = brokers
    k.Topic      = topic
    k.GroupId    = gid
    k.DoReset    = false
    k.Replicas   = 1
    k.Partitions = 1
    k.Active     = false
    return k
}