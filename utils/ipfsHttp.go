package utils

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "fmt"
)
//import "github.com/tidwall/gjson"


type IpfsClient struct{
    protocol string
    address string
    version string
    preUrl string
}

type ipfsID struct{
    ID string
    PublicKey string
    Addresses []string
    AgentVersion string
    ProtocolVersion string
}
type ipfsPeers struct{
    Peers []string
}

type ipfsKey struct{
    Name string
    Id   string
}
type ipfsKeys struct{
    Keys []ipfsKey
}
type ipfsSwarm struct{
    Strings []string
}
type ipfsSwarmPeer struct{
    Addr string
    Peer string
    Latency string
    Muxer   string
    Streams []map[string]string
}
type ipfsSwarmPeers struct{
    Peers []ipfsSwarmPeer
}
type ipfsPubsubPeers struct{
    Strings []string
}
type ipfsPubsubMessage struct{
    Form []uint8
    Data []uint8
    Seqno []uint8
    TopicIDs []string
    XXX_unrecongnized []uint8
}
type ipfsPubsubSub struct{
    Message ipfsPubsubMessage
}
type ipfsNameResolve struct{
    Path string
}
type ipfsNamePublish struct{
    Name string
    Value string
}
func NewIpfsClient(address string) *IpfsClient{
    ipfsClient := new(IpfsClient)
    ipfsClient.address=address
    ipfsClient.protocol="http://"
    ipfsClient.version="/api/v0"
    ipfsClient.preUrl=ipfsClient.protocol+ipfsClient.address+ipfsClient.version
    return ipfsClient
}

func (self *IpfsClient) getPreUrl()string{
    if self.preUrl==""{
      self.preUrl = "http://localhost:5001/api/v0"
    }
    return self.preUrl
}

func (self *IpfsClient) Shutdown()(string,error){
    //curl "http://localhost:5001/api/v0/shutdown"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/shutdown")
    if err!=nil{
        return "",err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return "",err1
    }
    return string(body),nil     
}

func (self *IpfsClient) Version()(map[string]string,error){
    //curl "http://localhost:5001/api/v0/version?number=false&commit=false&repo=false&all=false"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/version?number=true&commit=true&repo=true&all=true")
    result:=make(map[string]string)
    if err!=nil{
        return result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return result,err1
    }
    err2 := json.Unmarshal(body,&result)
    if err2!=nil{
        return result,err
    }
    return result,nil     
}

func (self *IpfsClient) Id()(ipfsID,error){
    //curl "http://localhost:5001/api/v0/id?arg=<peerid>&format=<value>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/id")
    result:=new(ipfsID)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) BootstrapList()(ipfsPeers,error){
    //curl "http://localhost:5001/api/v0/bootstrap/list"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/bootstrap/list")
    result:=new(ipfsPeers)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) BootstrapRmAll()(ipfsPeers,error){
    //curl "http://localhost:5001/api/v0/bootstrap/rm/all"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/bootstrap/rm/all")
    result:=new(ipfsPeers)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) KeyGen(name string,t string,size int)(ipfsKey,error){
    //curl "http://localhost:5001/api/v0/key/gen?arg=<name>&type=<value>&size=<value>"
    if t!="rsa" && t!="ed25519" {
        t="rsa"
    }
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/key/gen"+"?arg="+name+"&type="+t+"&size="+strconv.Itoa(size))
    result:=new(ipfsKey)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}
func (self *IpfsClient) KeyList()(ipfsKeys,error){
    //curl "http://localhost:5001/api/v0/key/list"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/key/list")
    result:=new(ipfsKeys)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}
func (self *IpfsClient) KeyRm(name string)(ipfsKeys,error){
    //curl "http://localhost:5001/api/v0/key/rm?arg=<name>&l=<value>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/key/rm"+"?arg="+name)
    result:=new(ipfsKeys)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) Cat(cid string)(string,error){
    //curl "http://localhost:5001/api/v0/cat?arg=<ipfs-path>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/cat?arg="+cid)
    if err!=nil{
        return "",err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return "",err1
    }
    return string(body),nil     
}
func (self *IpfsClient) DhtPut(key string,value string)([]byte,error){
    //curl "http://localhost:5001/api/v0/dht/put?arg=<key>&arg=<value>&verbose=false"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/dht/put?arg="+key+"&arg="+value+"&verbose=fale")
    if err!=nil{
        return nil,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return nil,err1
    }
    return body,nil     
}

func (self *IpfsClient) NamePublish(cid string,key string)(ipfsNamePublish,error){
    //curl "http://localhost:5001/api/v0/name/publish?arg=<ipfs-path>&resolve=true&lifetime=24h&ttl=<value>&key=self"
    url:=self.getPreUrl()
    if key==""{
        key="self"
    }
    resp,err:=http.Get(url+"/name/publish"+"?arg="+cid+"&key="+key)
    result:=new(ipfsNamePublish)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) NameResolve(nsid string)(ipfsNameResolve,error){
    //curl "http://localhost:5001/api/v0/name/resolve?arg=<name>&recursive=false&nocache=false"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/name/resolve"+"?arg="+nsid)
    result:=new(ipfsNameResolve)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,&result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) SwarmPeers()(ipfsSwarmPeers,error){
    //curl "http://localhost:5001/api/v0/swarm/peers?verbose=<value>&streams=<value>&latency=<value>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/swarm/peers")
    result:=new(ipfsSwarmPeers)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) SwarmConnect(pid string)(ipfsSwarm,error){
    //curl "http://localhost:5001/api/v0/swarm/connect?arg=<address>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/swarm/connect"+"?arg="+pid)
    result:=new(ipfsSwarm)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}
func (self *IpfsClient) SwarmDisconnect(pid string)(ipfsSwarm,error){
    //curl "http://localhost:5001/api/v0/swarm/disconnect?arg=<address>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/swarm/disconnect"+"?arg="+pid)
    result:=new(ipfsSwarm)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) PubsubLs()(ipfsPubsubPeers,error){
    //curl "http://localhost:5001/api/v0/pubsub/ls"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/pubsub/ls")
    result:=new(ipfsPubsubPeers)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}
func (self *IpfsClient) PubsubPeers(topic string)(ipfsPubsubPeers,error){
    //curl "http://localhost:5001/api/v0/pubsub/peers?arg=<topic>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/pubsub/peers"+"?arg="+topic)
    result:=new(ipfsPubsubPeers)
    if err!=nil{
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        return *result,err
    }
    return *result,nil     
}

func (self *IpfsClient) PubsubPub(topic string,value string)(string,error){
    //curl "http://localhost:5001/api/v0/pubsub/pub?arg=<topic>&arg=<data>"
    url:=self.getPreUrl()
    resp,err:=http.Get(url+"/pubsub/pub"+"?arg="+topic+"&arg="+value)
    if err!=nil{
        return "",err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        return "",err1
    }
    return string(body),nil    
}

func (self *IpfsClient) PubsubSub(topic string,discover bool)(ipfsPubsubSub,error){
    //curl "http://localhost:5001/api/v0/pubsub/sub?arg=<topic>&discover=<value>"
    url:=self.getPreUrl()
    var disc string
    if discover{
        disc="true"
    }else{
        disc="false"
    }
    resp,err:=http.Get(url+"/pubsub/sub"+"?arg="+topic+"&discover="+disc)
    fmt.Println(resp)
    result:=new(ipfsPubsubSub)
    if err!=nil{
        fmt.Println(err)
        return *result,err
    }
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1!=nil{
        fmt.Println(err)
        return *result,err1
    }
    err2 := json.Unmarshal(body,result)
    if err2!=nil{
        fmt.Println(err)
        return *result,err
    }
    fmt.Println("okokok")
    return *result,nil     
}
