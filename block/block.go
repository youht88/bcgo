package block

import (
  //"time"
  //"fmt"
  "strings"
  "strconv"
  "math"
  "encoding/json"
  
  "bcgo/utils"
  "bcgo/transaction"
  "gopkg.in/mgo.v2/bson"
)

type BlockHeader struct {
  Index int
  PrevHash string
  Timestamp int64
  LockTime int64
  Diffcult int
  MerkleRoot string
}

func (self *BlockHeader) HeaderString() string{
  headerByte,_ := json.Marshal(self)
  return string(headerByte)
}
func (self *BlockHeader) UpdateMerkleRoot(data []*transaction.Transaction) string{
  var txHash = []string{}
  for  _,item :=range data{
    txHash = append(txHash,item.Hash)
  }
  self.MerkleRoot=utils.Hashlib.Sha256(strings.Join(txHash,""))
  //merkleTree = merkle.Tree()
  //merkleRoot = merkleTree.makeTree(txHash)
  //self.merkleRoot = merkleRoot.value 
  return self.MerkleRoot
}

type Block struct{
  BlockHeader
  Data []*transaction.Transaction
  Nonce int
  Hash string 
}
func New(args map[string]interface{}) *Block{
  block := new(Block)
  for key,value := range args{
    switch key{
      case "Index":
        block.Index = value.(int)
      case "Nonce":
        block.Nonce = value.(int)
      case "PrevHash":
        block.PrevHash = value.(string)
      case "Timestamp":
        block.Timestamp = value.(int64)
      case "Diffcult":
        block.Diffcult = value.(int)
      case "MerkleRoot":
        block.MerkleRoot = value.(string)
      case "Data":
        block.Data = value.([]*transaction.Transaction)
      case "Hash":
        hash := value.(string)
        if hash!=""{
          block.Hash = value.(string)
        }else{
          block.Hash = block.UpdateHash("")
        }
    }
  }
  return block
}

func FindNonce(blockHeader BlockHeader,ch chan *Block) {
    defer func(){
      utils.Logger.Info("FindNonce end")
    }()
    utils.Logger.Info("FindNonce start")
    diffcult := blockHeader.Diffcult
    length  := diffcult / 4 
    mod    := diffcult % 4
    difNum := int64(math.Pow(2,float64(4-mod)))
    newBlock := new(Block)
    newBlock.BlockHeader = blockHeader
    preHeaderStr := blockHeader.HeaderString()
    newBlock.UpdateHash(preHeaderStr)
    
    for {
        subHash,_ := strconv.ParseInt(newBlock.Hash[:length+1],16,64)
        if  subHash < difNum {
           break
        }
        newBlock.Nonce +=1
        newBlock.UpdateHash(preHeaderStr)
        //time.Sleep(time.Nanosecond)
    }
    ch <- newBlock
}

func (self *Block) UpdateHash(headerStr string) string{
    if (headerStr!=""){
      self.Hash = utils.Hashlib.Sha256(headerStr+strconv.Itoa(self.Nonce))
    }else{
      self.Hash = utils.Hashlib.Sha256(self.BlockHeader.HeaderString()+strconv.Itoa(self.Nonce))
    }
    return self.Hash
}
func (self *Block) Save(){
    _,err:=utils.DB.Upsert("blockchain",bson.M{"index":self.Index},bson.M{"$set":self})
    if err!=nil{
      utils.Logger.Warning(err)
    }
}

func (self *Block) SaveToPool(){
  index := self.Index
  nonce := self.Nonce
  utils.Logger.Warnf("save block %v-%v to pool",index,nonce)
  utils.DB.Upsert("blockpool",utils.M{"hash":self.Hash},utils.M{"$set":self})
}
func (self *Block) RemoveFromPool(){
  err:=utils.DB.Delete("blockpool",utils.M{"hash":self.Hash})
  if err!=nil {
    utils.Logger.Warn("removeFromPool error",err)
  }
}
func (self *Block) IsValid() bool{
  if (self.Index == 0 ) {
    return true
  }
  utils.Logger.Infof(utils.Logger.Blue("verify block #%d-%d"),self.Index,self.Nonce)
  if (self.Index >= utils.DiffcultIndex && self.Index < utils.DiffcultIndex + utils.ADJUST_DIFF - 1 && self.Diffcult < utils.Diffcult) {
    utils.Logger.Errorf("%s is not worked because of diffcult is %d but little then %d",self.Hash,self.Diffcult,utils.Diffcult)
    return false
  }
  //logger.debug("verify proof of work")
  self.UpdateHash("")

  length := self.Diffcult / 4 
  mod := self.Diffcult % 4
  difNum := int64(math.Pow(2,float64(4-mod)))
  subHash , _ := strconv.ParseInt(self.Hash[:length+1],16,64) 
  if  subHash >= difNum{
    utils.Logger.Errorf("%v is not worked because of WOF is not valid",self.Hash)
    return false
  }
  
  //logger.debug(`${this.hash} is truly worked`)
  utils.Logger.Info("verify transaction data")
  txAmount :=[]map[string]float32{}
  
  for _,trans := range self.Data {
    if (!trans.IsValid(txAmount)) {
      utils.Logger.Errorf("%s is not worked because of transaction is not valid",self.Hash)
      return false
    }
  }
  //校验coinbase的交易费是否合法
  var fee float32
  for _,item := range txAmount{
    fee = fee + item["txInamount"] - item["txOutAmount"]
  }
  if (len(self.Data[0].Outs)>1 && self.Data[0].Outs[1].Amount > fee){
    utils.Logger.Error("矿工交易费设置不合法",self.Data[0].Outs[1].Amount,fee)
    return false
  }
  return true
}


/*
const fs=require('fs')
const async = require("async")
const utils = require("./utils.js")
const logger = utils.logger.getLogger()

const Transaction = require('./transaction.js').Transaction

class Block{
  constructor(args){
    this.index     = parseInt(args.index) || 0
    this.nonce     = parseInt(args.nonce) || 0
    this.prevHash  = args.prevHash ||""
    this.timestamp = parseInt(args.timestamp)
    this.diffcult  = parseInt(args.diffcult)
    this.merkleRoot= args.merkleRoot || ""
    this.data      = []
    for (var i=0 ;i < args.data.length;i++){
      this.data.push(Transaction.parseTransaction(args.data[i]))
    }
    this.hash     = args.hash     || this.updateHash()
  }  
  headerString(){
    return [this.index.toString(),
        this.prevHash,
        this.getMerkleRoot(),
        this.timestamp.toString(),
        this.diffcult.toString(),
        this.nonce.toString()].join("")
  }
  preHeaderString(){
    return [this.index.toString(),
        this.prevHash,
        this.getMerkleRoot(),
        this.timestamp.toString(),
        this.diffcult.toString()].join("")
  }
  getMerkleRoot(){
    let txHash=[]
    for (let item of this.data){
      txHash.push(item.hash)
    }
    this.merkleRoot=utils.hashlib.sha256(txHash.join(""))
    //merkleTree = merkle.Tree()
    //merkleRoot = merkleTree.makeTree(txHash)
    //self.merkleRoot = merkleRoot.value
    return this.merkleRoot

  }
  updateHash(preHeaderStr=null){
    if (preHeaderStr)
      this.hash = utils.hashlib.sha256(preHeaderStr+this.nonce.toString())
    else
      this.hash = utils.hashlib.sha256(this.headerString())
    return this.hash
  }
  dumps(){
    return {
      "index"      :this.index,
      "hash"       :this.hash,
      "prevHash"   :this.prevHash,
      "diffcult"   :this.diffcult,
      "nonce"      :this.nonce,
      "timestamp"  :this.timestamp,
      "merkleRoot" :this.merkleRoot,
      "data"       :this.data
    }
  }
  async save(){
    global.db.updateOne("blockchain",{"index":this.index},{"$set":this.dumps()},{"upsert":true})
    .catch(e=>console.log("save error:",e))
  }
  async saveToPool(){
    return new Promise(async (resolve,reject)=>{
      const {index,nonce} = this
      logger.warn(`save block ${index}-${nonce} to pool`)
      await global.db.updateOne("blockpool",{"hash":this.hash},{"$set":this.dumps()},{"upsert":true})
        .then(()=>resolve())
        .catch(e=>{console.log("saveToPool error:",e)
                   reject(e)   
              })
    })      
  }
  async removeFromPool(){
    global.db.deleteOne("blockpool",{"hash":this.hash})
    .catch(e=>console.log("removeFromPool error",e))
  }
  isValid(){
    if (this.index == 0 ) return true
    logger.debug(`verify block #${this.index}-${this.nonce}`)
    if (this.index >= global.diffcultIndex && this.index < parseInt(global.diffcultIndex) + parseInt(global.ADJUST_DIFFCULT) - 1 && this.diffcult < global.diffcult) {
      console.log(this.index,global.diffcultIndex,parseInt(global.diffcultIndex) + parseInt(global.ADJUST_DIFFCULT) - 1)
      logger.error(`${this.hash} is not worked because of diffcult is ${this.diffcult} but little then ${global.diffcult}`)
      return false
    }
    //logger.debug("verify proof of work")
    this.updateHash()

    let length = Math.floor(this.diffcult / 4 )
    let mod = this.diffcult % 4
         
    if ( parseInt(this.hash.slice(0,length+1),16) >= 2**(4-mod) ){
      logger.error(`${this.hash} is not worked because of WOF is not valid`)
      return false
    }
    
    //logger.debug(`${this.hash} is truly worked`)
    logger.debug("verify transaction data")
    let txAmount=[]
    for (let transaction of this.data){
      console.log("transaction hash:",transaction.hash)
      if (!transaction.isValid(txAmount)) {
        logger.error(`${this.hash} is not worked because of transaction is not valid`)
        return false
      }
    }
    //校验coinbase的交易费是否合法
    let fee=0
    if (txAmount.length>0)
      fee = txAmount.map(x=>x.txInAmount - x.txOutAmount).reduce((x,y)=>x+y)
    if (this.data[0].outs[1] && this.data[0].outs[1].amount > fee){
      logger.error('矿工交易费设置不合法',this.data[0].outs[1].amount,fee)
      return false
    }
    return true
  }
}
exports.Block = Block

*/