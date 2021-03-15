package utils

import (
    "encoding/base64"
)
var Bufferlib *_Bufferlib
func init(){
    Bufferlib = new(_Bufferlib)
}
type _Bufferlib struct{
    CodeTypes []string  //{"ascii","base64","utf8","hex","binary","base58"}
}
func (self *_Bufferlib) B64encode(str string) string{
   return base64.StdEncoding.EncodeToString([]byte(str))
}
func (self *_Bufferlib) B64decode(str string) string{
   decodeBytes,_ := base64.StdEncoding.DecodeString(str)  
   return string(decodeBytes)
}
func (self *_Bufferlib) B58encode(str string) string{
   return b58encode([]byte(str))
}
func (self *_Bufferlib) B58decode(str string) string{
   decodeBytes:= b58decode(str)  
   return string(decodeBytes)
}

// class Bufferlib{
//   constructor(){
//     this.codeTypes = ['ascii','base64','utf8','hex','binary','base58']
//   }
//   b64encode(str){
//     //对字符串进行base64编码
    
//     Buffer.from(str).toString('base64')
//   }
//   b64decode(str){
//     //对base64编码的字符串进行解码
//     return Buffer.from(str,'base64').toString()
//   }
//   b58encode(str){
//     return b58.encode(Buffer.from(str))
//   }
//   b58decode(str){
//     return b58.decode(str).toString()
//   }
//   toBin(str,codeType='utf8'){
//     //将特定编码类型的字符串压缩为bin码
//     if (this.codeTypes.includes(codeType)){
//       if (typeof str !== "string"){ 
//         str = JSON.stringify(str)
//         return Buffer.from(str,codeType)
//       }else if (codeType == "base58"){
//         return b58.decode(str)
//       }else{
//         return Buffer.from(str,codeType)
//       }
//     }else{
//       throw new Error(`code type must be one of ${this.codeTypes}`)
//     }
//   }
//   toString(buffer,codeType='utf8'){
//     //将压缩的bin码转换为对应类型的string
//     if (!Buffer.isBuffer(buffer)) throw new Error("first arg type must be buffer")
//     if (this.codeTypes.includes(codeType)){
//       if (codeType == "base58"){
//         return b58.encode(buffer)
//       }else{
//         return buffer.toString(codeType)
//       }
//     }else{
//       throw new Error(`code type must be one of ${this.codeTypes}`)
//     }
//   }
//   transfer(str,fromCode,toCode){
//     if (!this.codeTypes.includes(fromCode) || !this.codeTypes.includes(toCode) )
//       throw new Error(`code type must be one of ${this.codeTypes}`)
//     if (typeof str !== "string") {
//        str = JSON.stringify(str)
//        return this.toString(Buffer.from(str,'utf8'),toCode)
//     }else if (fromCode=="base58"){
//       return this.toString(b58.decode(str),toCode)
//     }else{
//       return this.toString(Buffer.from(str,fromCode),toCode)
//     }
//   }
// }
