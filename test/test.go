package main

import (
  "bcgo/utils"
  //"encoding/json"
  "fmt"
  "time"
  "strings"
  "strconv"
  "os"
  "math/rand"
)

func subBlock(){
    handle,err:=ipfs.PubSubSub("block")
    if err!=nil {
       fmt.Println("error:%v",err)
       return
    }
    for {
        msg,_:=handle.Next()
        data:=string(msg.Data)
        fmt.Println(data,msg.Seqno)
        ipfs.PubSubPub("recieved",data)
        if strings.Contains(data,"0") {
            ipfs.PubSubPub("confirm",data)
        }
    }
}

func subTrans(){
    handle,err:=ipfs.PubSubSub("trans")
    if err!=nil{
      fmt.Println("error:%v",err)
      return
    }
    for {
        msg,_:=handle.Next()
        fmt.Println(string(msg.Data),msg.Seqno)        
    }
}

func subRecieve(){
    recieve:=map[string]int{}
    handle,_:=ipfs.PubSubSub("recieved")
    for {
        msg,_:=handle.Next()
        data:=string(msg.Data)
        recieve[data]+=1
        fmt.Println("recieve receipt:",data,recieve[data])
    }
}
func subConfirm(){
      confirm:=map[string]int{}
      handle,_:=ipfs.PubSubSub("confirm")
      for {
        msg,_:=handle.Next()
        data:=string(msg.Data)
        confirm[data]+=1
        fmt.Println("confirm receipt:",data,confirm[data])        
    }
}

func send(key,msg string){
    var i int
    for{
        time.Sleep(5*time.Second)
        i = rand.Intn(10)
        ipfs.PubSubPub(key,msg+"->"+strconv.Itoa(i))
    }
}

var ipfs utils.Ipfs

func main(){
  if len(os.Args)!=4 {
     fmt.Printf("%v\n",os.Args)
     fmt.Print("sample: ./test ipfs5:5001 <block|trans> msg")
     panic("参数错误")
  }
  addr:=os.Args[1]
  key:=os.Args[2]
  msg:=os.Args[3]
  ipfs = utils.NewIpfs(addr)
  go subBlock()
  go subTrans()
  go subRecieve()
  go subConfirm()
  go send(key,msg)
  //handle1.Cancel()
  //handle2.Cancel()
  for{
      time.Sleep(time.Nanosecond)
  }
}