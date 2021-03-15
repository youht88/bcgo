package utils

import (
    "io/ioutil"
    //"encoding/json"
    //"strconv"
    //"fmt"
    //import "github.com/tidwall/gjson"
    "bytes"
    shell "github.com/ipfs/go-ipfs-api"
)

type Ipfs struct{
    shell *shell.Shell
    address string
}    

func NewIpfs(address string) Ipfs{
    ipfs := new(Ipfs)
    ipfs.shell   = shell.NewShell(address)
    ipfs.address = address
    return *ipfs
}

func (self Ipfs) Id()(*shell.IdOutput,error){
    result,err := self.shell.ID()
    return result,err 
}

func (self Ipfs) Add(data []byte)(string,error){
    cid,err := self.shell.Add(bytes.NewReader(data))
    return  cid,err
}

func (self Ipfs) Cat(cid string)([]byte,error){
    result,err1 := self.shell.Cat(cid)
    if err1!=nil{
        return nil,err1
    }
    data,err2 := ioutil.ReadAll(result)
    return data,err2 
}

func (self Ipfs)  PubSubPub(topic,data string) error{
     err := self.shell.PubSubPublish(topic, data )
     return err
} 
func (self Ipfs) PubSubSub(topic string)  (*shell.PubSubSubscription, error) {
    result,err := self.shell.PubSubSubscribe(topic)
    return result,err
}
func (self Ipfs) Publish(cid string,keys ...string) error{
    var key string
    if len(keys)!=0{
      key=keys[0]
    }
    err := self.shell.Publish(key,cid)
    return err
}
