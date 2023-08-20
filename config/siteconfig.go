/**  
  * Parses application config parameters
  *
 **/ 
package config

import ()


type KafkaSite struct {
    Brokers      string `yaml:"brokers"`
    Topic        string `yaml:"topic"`
    GroupId      string `yaml:"gid"`
    DoReset      bool   `yaml:"streamreset"`
    Replica      int    `yaml:"replicationfactor"`
    Partitions   int    `yaml:"partitions"`
    Active       bool
}

