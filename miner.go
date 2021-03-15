package main

import (
  "fmt"
  "flag"
  "time"
  "net"
  "strconv"
  "encoding/json"
  "net/http"
  "bcgo/utils"
  "bcgo/node"
  //"bcgo/block"
  "bcgo/transaction"

  "github.com/gorilla/mux"
)

func main(){
  //var err error

  var EntryNode = flag.String("entryNode","","indicate which node to entry,e.g. ip|host:port")
  var Me = flag.String("me","me","indicate who am I,e.g. ip|host:port")
  var HttpServer = flag.String("httpServer","0.0.0.0:5000","default httpServer is 0.0.0.0:5000")
  var IpfsServer = flag.String("ipfsServer","","default ipfsServer is 0.0.0.0:5001")
  var DBurl = flag.String("db","mongodb://mongo:27017/ipfs","db connect,ip:port/db")
  var Display = flag.Bool("display",false,"display of node")
  var SyncNode = flag.Bool("syncNode",false,"sync node")
  var Full = flag.Bool("full",false,"full sync")
  var Debug = flag.Bool("debug",false,"if debug mode")
  //var Logging = flag.String("logging","debug","logging level one of (trace|debug|info|warn|error|fatal)")

  flag.Parse()
  
  utils.Logger.Info("创建节点...")
  peer := node.New(map[string]interface{}{
    //"Config":config,
    "HttpServer":*HttpServer,
    "IpfsServer":*IpfsServer,
    "EntryNode":*EntryNode,
    "Me":*Me,
    "DBurl":*DBurl,
    "Display":*Display,
    "SyncNode":*SyncNode,
    "Full":*Full,
    "Debug":*Debug,
    //"IoServer":ioServer,
    //"IoClient":ioClient
  })
  utils.Logger.Successf("节点已经创建,%s",peer.HttpServer)
  utils.Logger.Successf("IPFS已经创建,%s",peer.IpfsServer)
  // //导入本地区块链
  // go node.SyncLocalChain()
  // log.debug(`localchain has ${node.blockchain.maxindex()} blocks`)
  // //链接网络
  // log.debug("socketioConnect...")
  // go node.socketioConnect()
  // if args.entryNode==args.me{
  //   node.emitter.emit("start")
  // }
  
  //导入本地区块链
  peer.SyncLocalChain()
  
  RegisterHttpService(peer)
  utils.Http.RegisterRpcService(new(RpcService),"/rpc")
  go utils.Http.Listen(peer.HttpServer)
  
  go startSocketServer(peer)

  //go peer.MinerProcess()
  go peer.SubTransaction()

  time.Sleep(time.Second)
  utils.Logger.Warnf(utils.Logger.Red("has RpcService.Add is %v"),utils.Http.RpcServer.HasMethod("RpcService.Add"))
  result,_:=utils.RpcCall("http://127.0.0.1:5000/rpc","RpcService.Hello","youht")
  utils.Logger.Success(result)
  result1,_:=utils.RpcCall("http://127.0.0.1:5000/rpc","RpcService.Add",map[string]interface{}{"a":2,"b":3})
  utils.Logger.Success(result1)
    
  for{
    time.Sleep(time.Hour)
  }
}

type RpcService struct{}
func (self *RpcService) Hello(req *http.Request,name *string,result *string)error{
   *result = "hello "+*name
   return nil
}
func (self *RpcService) Add(req *http.Request,args *map[string]float32,result *float32)error{
  *result = (*args)["a"]+(*args)["b"]
  return nil
}

func RegisterHttpService(peer *node.Node){
  r := utils.Http.Router
  r.HandleFunc("/testrpc",func(res http.ResponseWriter,req *http.Request){
    result,err:=utils.RpcCall("http://127.0.0.1:5000/rpc","RpcService.Add",map[string]interface{}{"a":2,"b":3.3})
    if err!=nil{
      utils.Logger.Error("[client]",err)
      fmt.Fprintf(res,"error:",err)
    }else{
      utils.Logger.Success(result)
      fmt.Fprintf(res,"[client] the value is : %0.2f",result)
    }
  })
  r.HandleFunc("/add/{a}/{b}",func(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    a:= vars["a"]
    b:= vars["b"]
    println(a,b)
    // fetch job
    job := utils.NewJob(peer.Add,a,b)
    rst := <- job.C
    fmt.Fprintf(res, rst.Result.(string))
  })
  r.HandleFunc("/hello",func(res http.ResponseWriter,req *http.Request){
     res.Write([]byte("hello"))
  })
  r.HandleFunc("/ipfs/add/{msg}",func(res http.ResponseWriter,req *http.Request){
    msg:=mux.Vars(req)["msg"]
    if msg==""{
      res.Write([]byte("msg must be defined!"))
      return
    }
    job :=utils.NewJob(peer.IpfsAdd,msg)
    rst0 := <- job.C
    rst := map[string]interface{}{
        "cid":rst0.Result,
        "pubkey":peer.Wallet.Pubkey,
    }
    b,_ :=json.MarshalIndent(rst,"","    ")
    var s interface{}
    json.Unmarshal(b,&s) 
    utils.Logger.Warning(s)
    res.Write(b)
  })
  r.HandleFunc("/transaction/add/{data}",func(res http.ResponseWriter,req *http.Request){
    msg:=mux.Vars(req)["data"]
    if msg==""{
      res.Write([]byte("data must be defined!"))
      return
    }
    job :=utils.NewJob(peer.TransAdd,msg)
    rst0 := <- job.C
    rst := map[string]interface{}{
        "hcid":rst0.Result,
        "pubkey":peer.Wallet.Pubkey,
    }
    b,_ :=json.MarshalIndent(rst,"","    ")
    var s interface{}
    json.Unmarshal(b,&s) 
    utils.Logger.Warning(s)
    res.Write(b)
  })
  r.HandleFunc("/transaction/get/{Hcid}",func(res http.ResponseWriter,req *http.Request){
    hcid:=mux.Vars(req)["Hcid"]
    if hcid==""{
      res.Write([]byte("Hcid must be defined!"))
      return
    }
    job :=utils.NewJob(peer.TransGet,hcid)
    rst0 := <- job.C
    rst := map[string]interface{}{
        "data":rst0.Result,
        "pubkey":peer.Wallet.Pubkey,
    }
    b,_ :=json.MarshalIndent(rst,"","    ")
    var s interface{}
    json.Unmarshal(b,&s) 
    utils.Logger.Warning(s)
    res.Write(b)
  })
  r.HandleFunc("/wallet/me",func(res http.ResponseWriter,req *http.Request){
    //balance := peer.blockchain.utxo.getBalance(peer.Wallet.Address)
    rst := map[string]interface{}{
        "address":peer.Wallet.Address,
        "pubkey":peer.Wallet.Pubkey,
        "balance":0}
    b,_ :=json.MarshalIndent(rst,"","    ")
    var s interface{}
    json.Unmarshal(b,&s) 
    utils.Logger.Warning(s)
    res.Write(b)
  })
  r.HandleFunc("/genesis",func(res http.ResponseWriter,req *http.Request){
    coinbase:=transaction.NewCoinbase(peer.Wallet.Address,0)
    block :=peer.GenesisBlock(coinbase)
    utils.Logger.Successf("the block hash is :%s",utils.Logger.Green(block.Hash))
    block.Save()
    rtn,_ := json.MarshalIndent(block,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/mine",func(res http.ResponseWriter,req *http.Request){
    block :=peer.Mine()
    utils.Logger.Successf("the block hash is :%s",utils.Logger.Green(block.Hash))
    block.Save()
    rtn,_ := json.MarshalIndent(block,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/blockchain",func(res http.ResponseWriter,req *http.Request){
    chain := peer.Blockchain.Blocks
    rtn,_ := json.MarshalIndent(chain,"","    ")  
    res.Write(rtn)
  })
  
  r.HandleFunc("/blockchain/index/{index}",func(res http.ResponseWriter,req *http.Request){
    index,err:=strconv.Atoi(mux.Vars(req)["index"])
    if err!=nil{
      res.Write([]byte("index must be number!"))
      return
    }
    block,_ := peer.Blockchain.FindBlockByIndex(index)
    rtn,_ := json.MarshalIndent(block,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/blockchain/hash/{hash}",func(res http.ResponseWriter,req *http.Request){
    hash:=mux.Vars(req)["hash"]
    block,_ := peer.Blockchain.FindBlockByHash(hash)
    rtn,_ := json.MarshalIndent(block,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/blockchain/{start}/{end}", func(res http.ResponseWriter,req *http.Request){
    vars:=mux.Vars(req)
    start,err:=strconv.Atoi(vars["start"])
    if err!=nil{
      res.Write([]byte("start must be number!"))
      return
    }
    end,err  :=strconv.Atoi(vars["end"])
    if err!=nil{
      res.Write([]byte("end must be number!"))
      return
    }
    blocks := peer.Blockchain.GetRangeBlocks(start,end)
    rtn,_ := json.MarshalIndent(blocks,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/blockchain/lastblock", func(res http.ResponseWriter,req *http.Request){
    block := peer.Blockchain.Lastblock()
    rtn,_ := json.MarshalIndent(block,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/blockchain/spv", func(res http.ResponseWriter,req *http.Request){
    chainSPV := peer.Blockchain.GetSPV()
    rtn,_ := json.MarshalIndent(chainSPV,"","    ")  
    res.Write(rtn)
  })

  r.HandleFunc("/blockchain/isvalid", func(res http.ResponseWriter,req *http.Request){
    valid := peer.Blockchain.IsValid()
    rtn,_ := json.MarshalIndent(valid,"","    ")  
    res.Write(rtn)
  })

  
}


func startSocketServer(peer *node.Node){
  listener,err := net.Listen("tcp",peer.EntryNode)
  if err != nil {
      utils.Logger.Danger(err)
  }
  defer listener.Close()
  for  {
        conn,err := listener.Accept() //用conn接收链接
        if err != nil {
            utils.Logger.Danger(err)
        }
        go HandleSocket(conn)  //开启多个协程。
    }
}

func HandleSocket(conn net.Conn) { //这个是在处理客户端会阻塞的代码。
    var buf [1024]byte
    n,_:=conn.Read(buf[:])
    println(n,string(buf[:]))
    conn.Write([]byte(time.Now().Local().String()))//通过conn的wirte方法将这些数据返回给客户端。
    conn.Close() //与客户端断开连接。
}

