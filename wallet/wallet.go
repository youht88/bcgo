package wallet

import (
    //"fmt"
    "errors"
    "strings"
    "bcgo/utils"
    "gopkg.in/mgo.v2/bson"
)

type Wallet struct{
    Name string
    Address string
    utils.Key
}
func IsAddress(nameOrAddress string) bool {
    return false
}
func DeleteByName(name string) error{
  return utils.DB.Delete("wallet",bson.M{"name":name})
}
func DeleteByAddress(address string) error{
  return utils.DB.Delete("wallet",bson.M{"address":address})
}

func Create(name string,keys ...utils.Key) *Wallet{ //prvkey=null,pubkey=null){
    wallet := new(Wallet)
    if len(keys)!=0{
        wallet.Key = keys[0]
    }else{
        wallet.Key = utils.Crypto.GenerateKeys()
    }
    address := Address(wallet.Key.Pubkey)
    wallet.Name = name
    wallet.Address = address
    utils.DB.Insert("wallet",*wallet)
    return wallet
}

func (self *Wallet) IsAddress(nameOrAddress string) bool{
    return IsAddress(nameOrAddress)
}

func (self *Wallet) ChooseByName(name string) error{
  w := []Wallet{}
  utils.DB.Find("wallet",bson.M{"name":name}).All(&w)
  if len(w)!=1 {
     return errors.New("没有记录或有多条记录")
  }
  *self = w[0]
  return nil
}
func (self *Wallet) ChooseByAddress(address string) error{
  w := []Wallet{}
  utils.DB.Find("wallet",bson.M{"address":address}).All(&w)
  if len(w)!=1 {
     return errors.New("没有记录或有多条记录")
  }
  *self = w[0]
  return nil
}
func Address(pubkey []string) string{
    var version string
    if len(pubkey)==1{
      version="00"
    }else{
      version="05"
    }
    publickey := version + strings.Join(pubkey,"")
    address := utils.Hashlib.Sha1(publickey)
    address58 := utils.Bufferlib.B58encode(address)
    return address58
}
func (self *Wallet) Getall() *[]Wallet{
  w:=[]Wallet{}
  iter:=utils.DB.Pipe("wallet",[]bson.M{{"$project":bson.M{"_id":0,"name":1,"address":1}}})
  iter.All(&w)
  return &w
}

//   static  isAddress(nameOrAddress){
//     let p,q,r,s,t
//     try{
//       t=utils.bufferlib.toBin(nameOrAddress,"base58")
//       if (t.length!=25 || (t[0]!=0 && t[0]!=5)) return false
//       r = t.slice(1,21).toString('hex')
//       s = utils.hashlib.doubleSha256(r)
//       p = s.slice(0,8)
//       q = t.slice(21,25).toString('hex')
//       if (p != q) return false
//       return true
//     }catch(e){
//       return false
//     }
//   }
