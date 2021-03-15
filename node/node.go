package node

import (
  "fmt"
  "os"
  "time" 
  "strings"
  "errors"
  "strconv"
  "encoding/json"
  "gopkg.in/mgo.v2/bson"
  "bcgo/utils"
  "bcgo/block"
  "bcgo/transaction"
  "bcgo/wallet"
  "bcgo/chain"
)
type Node struct{
  Config string
  HttpServer string
  Me string
  EntryNode string
  IpfsServer string
  Ipfs  utils.Ipfs
  DBurl     string
  Display   bool
  Peers     string
  Diffcult  int
  DiffcultIndex int
  Nodes      []string
  EntryNodes []string
  ClientNodesId []string
  Dbclient   interface{}
  Wallet     wallet.Wallet
  BlockSyncing bool
  Mining       bool 
  Blockchain   *chain.Chain
}
func New(args map[string]interface{}) *Node{
  peer := new(Node)
  for key,value := range args{
    switch key{
      case "Config":
        peer.Config = value.(string)
      case "HttpServer":
        peer.HttpServer = value.(string)
      case "Me":
        peer.Me = value.(string)
      case "EntryNode":
        peer.EntryNode = value.(string)
      case "IpfsServer":
        peer.IpfsServer = value.(string)
      case "DBurl":
        peer.DBurl = value.(string)
      case "Display":
        peer.Display = value.(bool)
      case "Peers":
        peer.Peers = value.(string)
    }
  }

  peer.Ipfs = utils.NewIpfs(peer.IpfsServer)
  peer.BlockSyncing = false
  peer.Mining       = false  

  peer.Diffcult = utils.ZERO_DIFF

  //this.isolateUTXO={}
  //this.isolatePool=[]
  //this.tradeUTXO = {}
  //this.isolateBlockPool = []
  //this.tradeUTXO = new UTXO("trade")
  //this.isolateUTXO = new UTXO("isolate")

  var err error
  //连接数据库    
  utils.Logger.Info("连接数据库...")
  arr:=strings.Split(peer.DBurl,"/")
  if len(arr)!=4 {
    utils.Logger.Dangerf("数据库连接串格式错误,%s",peer.DBurl)
  }
  err =utils.DB.Init(arr[0]+"//"+arr[2],arr[3])
  if err!=nil {
    utils.Logger.Danger(err)
    os.Exit(1)
  }
  utils.Logger.Successf("数据库已连接,%s",peer.DBurl)
  //创建钱包  
  utils.Logger.Info("初始化钱包...")
  mywallet := new(wallet.Wallet)
  err =mywallet.ChooseByName(peer.Me)
  if err!=nil {
      peer.Wallet = *wallet.Create(peer.Me)
      utils.Logger.Successf("钱包已创建,name:%s,address:%s",peer.Wallet.Name,peer.Wallet.Address)
  }else{
      peer.Wallet = *mywallet
      utils.Logger.Warningf("钱包已连接,name:%s,address:%s",peer.Wallet.Name,peer.Wallet.Address)
  }
  return peer
}

func (self *Node)Add(args ...interface{}) (interface{},error){
  sum:=0
  for _,item :=range args{
    temp,_:=strconv.Atoi(item.(string))
    sum+=temp
  }
  utils.Logger.Info(sum)
  return strconv.Itoa(sum),nil
}

func (self *Node)IpfsAdd(args ...interface{}) (interface{},error){
  cid,err := self.Ipfs.Add([]byte(args[0].(string)))

  return cid,err
}

func (self *Node)SubTransaction(){
    handle,err:=self.Ipfs.PubSubSub("transaction")
    if err!=nil{
      fmt.Println("error:%v",err)
      return
    }
    prvkey:= self.Wallet.Prvkey[0]
    for {
        msg,_:=handle.Next()
        //fmt.Println(string(msg.Data),msg.Seqno) 
        trans:=new(transaction.Trans)
        json.Unmarshal(msg.Data,trans)
        //校验地址
        address:=self.Wallet.Address
        if address!=trans.Address{
          utils.Logger.Warning("钱包地址不正确！")
          continue
        }
        //校验数据
        verify:=utils.Crypto.Verify(trans.Data,trans.Sign,trans.VPubkey)
        if !verify{
          utils.Logger.Warning("数据已被篡改！")
          continue
        }
        hcid:=trans.HCid
        cid:=utils.Crypto.Decrypt(trans.Data,prvkey)
        data,_:=self.Ipfs.Cat(cid)
        utils.Logger.Infof("cat %s:%s",cid,data)
        utils.DB.Upsert("trans",bson.M{"hcid":hcid},trans)      
    }
}

func (self *Node)TransAdd(args ...interface{}) (interface{},error){
  cid,err := self.Ipfs.Add([]byte(args[0].(string)))
  if err!=nil{
    return nil,err
  }
  pubkey:=self.Wallet.Pubkey[0]
  prvkey:=self.Wallet.Prvkey[0]
  address:=self.Wallet.Address
  hcid:=utils.Hashlib.Sha256(cid)
  ecid:=utils.Crypto.Encrypt(cid,pubkey)
  sign:=utils.Crypto.Sign(ecid,prvkey)
  trans:=transaction.Trans{hcid,ecid,pubkey,address,pubkey,sign}
  btrans,_:=json.Marshal(trans)
  err1:=self.Ipfs.PubSubPub("transaction",string(btrans))
  if err1!=nil{
    return nil,err
  }
  return hcid,nil
}

func (self *Node)TransGet(args ...interface{}) (interface{},error){
  hcid := args[0].(string)
  trans:=new(transaction.Trans)
  utils.DB.Find("trans",bson.M{"hcid":hcid}).One(trans)      
  utils.Logger.Info(trans)  
  prvkey:=self.Wallet.Prvkey[0]

  //json.Unmarshal([]byte(message),trans)
  //校验地址
  address:=wallet.Address([]string{trans.VPubkey})
  if address!=trans.Address{
     return nil,errors.New("钱包地址不正确！")
  }
  //校验数据
  verify:=utils.Crypto.Verify(trans.Data,trans.Sign,trans.VPubkey)
  if !verify{
     return nil,errors.New("数据已被篡改！")
  }else{
     cid:=utils.Crypto.Decrypt(trans.Data,prvkey)
     data,_:=self.Ipfs.Cat(cid)
     rst := map[string]interface{}{
       "cid":cid,
       "data":string(data),
     }
     return rst,nil
  }
}

func (self *Node) RegisteNode(peers string){
    //self.Nodes = utils.set.union(this.nodes,peers)
    //self.Peers = string.Join(self.Nodes,',')
    //utils.Fs.WriteFileSync("peers",self.Peers)
}

func (self *Node) GenesisBlock(coinbase *transaction.Transaction) *block.Block{
    defer func(){
      utils.Logger.Info("GenesisBlock end")
    }()
    ch :=make(chan *block.Block)
    blockHeader := block.BlockHeader{
          Index:0,
          PrevHash:"0",
          Timestamp: time.Now().UnixNano(),
          LockTime:0,
          Diffcult:utils.ZERO_DIFF,
    }
    blockHeader.UpdateMerkleRoot([]*transaction.Transaction{coinbase})
    go block.FindNonce(blockHeader,ch)
    block:= <-ch 
    block.Data = []*transaction.Transaction{coinbase}
    return block
}
func (self *Node) adjustDiffcult(endIndex int){
    if ( endIndex%utils.ADJUST_DIFF != 0) || endIndex < utils.ADJUST_DIFF {
      return
    } 
    startIndex:= endIndex - utils.ADJUST_DIFF + 1
    block1,_ :=self.Blockchain.FindBlockByIndex(startIndex)
    sTime := block1.Timestamp
    block2,_ := self.Blockchain.FindBlockByIndex(endIndex)
    eTime := block2.Timestamp
    block_per_hour := utils.ADJUST_DIFF/int(eTime - sTime)*3600000
    oldDiffcult := self.Diffcult
    if block_per_hour > int(float32(utils.BLOCK_PER_HOUR)*1.1){//速度太快，增加难度
      self.Diffcult++  
    }else if block_per_hour < int(float32(utils.BLOCK_PER_HOUR)*0.9){ //速度太慢，减少难度
      self.Diffcult--
      if (self.Diffcult<utils.ZERO_DIFF) {
        self.Diffcult=utils.ZERO_DIFF
      }
    }
    self.DiffcultIndex = endIndex+1
    utils.DiffcultIndex = self.DiffcultIndex
    utils.Diffcult = self.Diffcult
    utils.Logger.Warnf("index:%v-%v,每小时出块:%v,难度值由%v调整为%v",
        startIndex,endIndex,block_per_hour,oldDiffcult,self.Diffcult)
}
func (self *Node) SyncLocalChain(){
    localChain := new(chain.Chain)
    utils.DB.Pipe("blockchain",[]bson.M{{"$project":bson.M{"_id":0}},{"$sort":bson.M{"index":1}}}).All(&localChain.Blocks)
    self.Blockchain = localChain
    utils.Logger.Successf("localchain has %v blocks",len(self.Blockchain.Blocks))

    self.adjustDiffcult(len(localChain.Blocks)-1)
    //默认不检查本地blockchain，节省时间
    //if (!localChain.isValid()) throw new Error("本地chain不合法，重新下载正确的chain")
}

func (self *Node) Mine() *block.Block{
    self.Mining = true
    //sync transaction from txPool
    var txPool []*transaction.Transaction 
    var fee = 0.0
    // var (fee=0
    //      inamt=0
    //      outamt=0
    //     )
    // value = []
    // go func(){
    //   txs:= self.txPoolSync()
    //   for _,tx := range txs {
    //       txPool.push(tx)
    //       //处理交易费
    //       inAmount:=0
    //       outAmount:=0
    //       for _,txIn := range tx.ins){
    //         prevTx := self.Blockchain.FindTransaction(txIn.PrevHash)
    //         inAmount += prevTx.Outs[txIn.Index].Amount
    //       }
    //       for _,txOut := range tx.outs{
    //         outAmount += txOut.amount
    //       }
    //       value = append(value , {inAmount,outAmount} )
    //       fee   +=  inAmount - outAmount
    //       inamt +=  inAmount
    //       outamt += outAmount       
    //   }
    // }()
    //utils.Logger.Warnf("fee=%v,inamt=%v,outamt=%v",fee,inamt,outamt)
    coinbase := transaction.NewCoinbase(self.Wallet.Address,fee)
    txPool = append(txPool,coinbase)
    
    prevBlock := self.Blockchain.Lastblock()
    
    //mine a block with a valid nonce
    blockHeader := block.BlockHeader{
        Index    :prevBlock.Index + 1, 
        PrevHash :prevBlock.Hash,
        Timestamp:time.Now().UnixNano(),
        LockTime :0,
        Diffcult :self.Diffcult,
       }
    blockHeader.UpdateMerkleRoot(txPool)
    utils.Logger.Warnf("is mining block %v,diffcult:%v",blockHeader.Index,blockHeader.Diffcult)
    ch := make(chan *block.Block)
    go block.FindNonce(blockHeader,ch)
    newBlock := <- ch
    newBlock.Data = txPool
    utils.Logger.Successf(utils.Logger.Red("[end] mine %v-%v"),newBlock.Index,newBlock.Nonce)
    return newBlock
    //"other miner mined"
    //remove transaction from txPool
    //await this.txPoolRemove(newBlock) 
    
    //newBlockDict = utils.obj2json(newBlock)

    //push to blockPool
//     this.emitter.emit("mined",newBlockDict)
//     //broadcast newBlock
//     logger.info(`broadcast block ${newBlock.index}-${newBlock.nonce}`)
//     this.broadcast(newBlockDict,"newBlock")
      
//     logger.info("mine广播完成")
    
//     //以下由blockPoolSync处理
//     //newBlock.save()
//     //self.blockchain.add_block(newBlock)
//     //self.updateUTXO(newBlock)
    
//     if (cb)
//       return cb(null,newBlock)
//     return newBlock
}

func (self *Node) BlockProcess(){
    // try{
    //   let promiseArray
    //   let start,end
    //   let results
    //   let maxindex = this.blockchain.maxindex()
    //   let blocksDict = await global.db.findMany("blockpool",{"index":maxindex+1},{"projection":{_id:0}})
    //   if (blocksDict.length!=0){
    //     let done = await this.blockPoolSync(blocksDict)
    //       .catch((error)=>{logger.error(error.stack)})
    //     logger.warn("blockprocess",done)
    //     if (!done){
    //       start = (maxindex - global.NUM_FORK >1)?maxindex - global.NUM_FORK:1
    //       end = maxindex + 1
    //       promiseArray = this.getARpcData("getBlocks",{start,end})
    //       results = await Promise.all(promiseArray)
    //       logger.warn("blockprocess",start,end,results)
    //       if (results.length<=0) 
    //         return setTimeout(this.blockProcess.bind(this),100)
    //       for (let result of results){
    //         for (let blockDict of result.data){
    //           const block = new Block(blockDict)
    //           if (block.isValid())
    //             block.saveToPool()
    //         }
    //       }
    //     }
    //   }else{
    //     let endBlock = await global.db.findOne("blockpool",{"index":{"$gte":maxindex+1}},{"projection":{_id:0},"sort":["index","ascending"]})
    //     if (!endBlock) {
    //       return setTimeout(this.blockProcess.bind(this),100)
    //     }
    //     if (maxindex+1 > endBlock.index -1 ) {
    //       return setTimeout(this.blockProcess.bind(this),100)
    //     } 
    //     console.log("want block",maxindex+1,endBlock.index - 1)
    //     start = maxindex+1
    //     end   = endBlock.index - 1
    //     promiseArray = this.getARpcData("getBlocks",{start,end})
    //     results = await Promise.all(promiseArray)
        
    //     console.log("get block...",results)
    //     if (results.length<=0) 
    //       return setTimeout(this.blockProcess.bind(this),100)
    //     for (let result of results){
    //       for (let blockDict of result.data){
    //         const block = new Block(blockDict)
    //         if (block.isValid())
    //           block.saveToPool()
    //       }
    //     }
    //   }
    //   setTimeout(this.blockProcess.bind(this),100)
    // }catch(error){
    //   logger.fatal("blockProcess",error.stack)
    //   setTimeout(this.blockProcess.bind(this),100)
    // }
}
  
func (self *Node) MinerProcess(){
    ch := make(chan int)
    //timestamp := time.Now().UnixNano()
    //txPoolCount := utils.DB.Count("transaction",utils.M{"$or":utils.M{[{"lockTime":0},{"lockTime":utils.M{"$lte":timestamp}}]}})
    //if (txPoolCount < utils.TRANSACTION_TO_BLOCK) {
    //  return
    //}
    
    //mine
    for{
      go func(){
        block:=self.Mine()
        self.Blockchain.Blocks = append(self.Blockchain.Blocks,block)
        block.Save()
        ch <- 1
      }()
      <- ch
    }

}
