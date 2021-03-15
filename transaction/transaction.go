package transaction

import (
  //"strconv"
  //"strings"
  "time"
  "bcgo/utils"
  "encoding/json"
)
type Trans struct{
    HCid string
    Data string
    EPubkey string
    Address string
    VPubkey string
    Sign string
}
type TXin struct{
    PrevHash string
    Index int
    InAddr string
    Pubkey []string
    Sign   []string
}
func NewTXin(args map[string]interface{}) *TXin{
  txin := new(TXin)
  for key,value := range args{
    switch key{
      case "PrevHash":
        txin.PrevHash = value.(string)
      case "Index":
        txin.Index = value.(int)
      case "InAddr":
        txin.InAddr = value.(string)
      case "Pubkey":
        txin.Pubkey = value.([]string)
      case "Sign":
        txin.Sign = value.([]string)
    }
  }
  return txin
}

type TXout struct{
  Amount float32
  OutAddr string
  SignNum int
  ContractHash string
  Script string
  Assets interface{}
}
func NewTXout(args map[string]interface{}) *TXout{
  txout := new(TXout)
  for key,value := range args{
    switch key{
      case "Amount":
        txout.Amount = value.(float32)
      case "OutAddr":
        txout.OutAddr = value.(string)
      case "SignNum":
        signNum := value.(int)
        if signNum!=1 {
          txout.SignNum = signNum
        }else{
          txout.SignNum = 1
        }
      case "ContractHash":
        txout.ContractHash = value.(string)
      case "Script":
        txout.Script = value.(string)
      case "Assets":
        txout.Assets = value
    }
  }
  if (txout.Script!="" && txout.ContractHash==""){
    txout.ContractHash = utils.Hashlib.Md5(txout.Script)
  }  
  return txout
}

type Transaction struct{
  Ins []TXin
  InsLen int
  Outs []TXout
  OutsLen int
  Timestamp int64
  LockTime  int64
  SignType  string
  Hash string
} 
func New(args map[string]interface{}) *Transaction{
  trans := new(Transaction)
  var (
     timestamp  int64
     signType   string
     hash       string
  )
  for key,value := range args{
    switch key{
      case "Ins":
        trans.Ins = value.([]TXin)
      case "Outs":
        trans.Outs = value.([]TXout)
      case "Timestamp":
        timestamp = value.(int64)
      case "LockTime":
        trans.LockTime = value.(int64)
      case "SignType":
        signType = value.(string)
      case "Hash":
        hash = value.(string)
    }
  }
  
  trans.InsLen  = len(trans.Ins)
  trans.OutsLen = len(trans.Outs)
  if timestamp!=0 {
    trans.Timestamp = timestamp
  }else{
    trans.Timestamp = time.Now().UnixNano()
  }

  if signType!="" {
    trans.SignType = signType
  }else{
    trans.SignType = "all"
  }

  if hash!=""{
    trans.Hash = hash
  }else{
    trans.Hash = utils.Hashlib.Sha256(trans.preHeaderString())
  }

  return trans
}

func (self *Transaction) IsValid(txAmount []map[string]float32) bool{
    utils.Logger.Info("transaction begin verify ",self.Hash)
    if (self.IsCoinbase()){
      if (!(self.InsLen==1 && self.Outs[0].Amount<=utils.REWARD)){
        utils.Logger.Error("transaction verify","coinbase transaction reward error.")
        return false
      }
      //check fee is coded in function block.isValid
      return true
    }
    if (utils.Hashlib.Sha256(self.preHeaderString())!=self.Hash) {
      utils.Logger.Error("transaction verify","交易内容与hash不一致")
      return false
    }
    //验证lockTime合法性
    if (utils.LockTime!=0 && self.LockTime >= time.Now().UnixNano()){
      utils.Logger.Warn("transaction verify","lockTime时期尚未到来")
      return false
    }
    
    /*
    //验证每条输入
    signType := self.SignType
    outByte,_ :=json.Marshal(self.Outs)
    outsHash := utils.Hashlib.Sha256(string(outByte))
    //logger.error("isValid",outsHash)
    prevTxAmount:=[]
    for idx,vin := range self.Ins {
      if (!vin.canUnlockWith({signType,outsHash,prevTxAmount})) return false
    }
    txInAmount:=0
    if (prevTxAmount.length>0)
      txInAmount = prevTxAmount.map(x=>x.amount).reduce((x,y)=>x+y)
    let txOutAmount = this.outs.map(x=>x.amount).reduce((x,y)=>x+y)
    utils.Logger.Warn("transaction verify","txInAmount",txInAmount,"txOutAmount",txOutAmount)
    if ( txInAmount < txOutAmount ){
      utils.Logger.Error("transaction verify","输入的金额小于输出的金额","txInAmount",txInAmount,"txOutAmount",txOutAmount)
      return false
    }
    //txAmount的作用是供block.isValid函数判断coinbase交易的交易费是否合法
    txAmount.push({hash:this.hash,txInAmount:txInAmount,txOutAmount:txOutAmount})
    */
    return true
}
func (self *Transaction) preHeaderString() string{
  ins:=[]TXin{}
  for _,item := range self.Ins {
    item.Pubkey =[]string{}
    item.Sign =[]string{}
    ins = append(ins,item)
  }
  
  tran:=Transaction{
    Ins  : ins,
    Outs : self.Outs,
    Timestamp : self.Timestamp,
    LockTime  : self.LockTime,
    SignType  : self.SignType,
  }
  rst,_ := json.Marshal(tran)
  return string(rst)
}

func (self *Transaction) IsCoinbase() bool{
    return self.Ins[0].Index == -1
}

func NewCoinbase(outAddr string,fee float64) *Transaction{
    ins:=[]TXin{*NewTXin(map[string]interface{}{
                   "PrevHash":"",
                   "Index":-1,
                   "InAddr":"",
               })}
    outs:=[]TXout{*NewTXout(map[string]interface{}{
                     "Amount":utils.REWARD,
                     "OutAddr":outAddr,
                     "Script":"",
                  })}
    if fee > 0 {
      outs=append(outs,*NewTXout(map[string]interface{}{
              "Amount":fee,
              "OutAddr":outAddr}))}
    
    return New(map[string]interface{}{
                "Ins":ins,
                "Outs":outs})
}

func ParseTransaction(data Transaction) *Transaction{
    var ins []TXin
    var outs []TXout
    for i:=0;i<len(data.Ins);i++ {
      ins = append(ins,data.Ins[i])
    }
    for j:=0;j<len(data.Outs);j++ {
      outs = append(outs,data.Outs[j])
    }
    
    args:=map[string]interface{}{
       "Ins"     : ins,
       "Outs"    : outs,
       "Hash"     : data.Hash ,
       "Timestamp": data.Timestamp ,
       "LockTime" : data.LockTime  ,
       "SignType" : data.SignType  ,
    }
    return New(args)
}

