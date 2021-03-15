package utils

import (
    "fmt"
    //"log"
    "math/big"
    "encoding/hex"
    "encoding/json"
    "os"

    "crypto"
    "crypto/rand"
    "crypto/x509"
    "crypto/elliptic"
    "crypto/ecdsa"
    "crypto/rsa"
    "encoding/pem"

    "github.com/tjfoc/gmsm/sm2"
    //"github.com/tjfoc/gmsm/sm4"
)

var Crypto CryptoIface

func init(){
  Crypto = new(SM2)
  //Crypto = new(ECC)
}

//crypto
type CryptoIface interface{
   GenerateKeys(files ...string)(Key) //GenerateKeys("prvkey.pem","pubkey.pem")
   Sign(message, prvkey string) string 
   Verify(message,signStr,pubkey string) bool
   Encrypt(message , pubkey string) string
   Decrypt(message , prvkey string) string
   //Encipher(message,key string) string
   //Decipher(message,key string) string
   //ToPEM()
}

type Key struct{
    Prvkey []string
    Pubkey []string
}
type SignStruct struct{
    R *big.Int 
    S *big.Int
}

type SM2 struct{}
func (self *SM2) GenerateKeys(files ...string)(Key){
  privateKey,_ := sm2.GenerateKey()
  private,_ := sm2.MarshalSm2PrivateKey(privateKey,nil)
  publicKey := &privateKey.PublicKey
  public,_  := sm2.MarshalSm2PublicKey(publicKey)
  prvkey:=hex.EncodeToString(private[:])
  pubkey:=hex.EncodeToString(public[:])
  return Key{[]string{prvkey},[]string{pubkey}}
}
func (self *SM2) Sign(message, prvkey string) string{
    private , _:=hex.DecodeString(prvkey)
    privateKey , err:=sm2.ParsePKCS8UnecryptedPrivateKey(private)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    r, s, err := sm2.Sign(privateKey, []byte(message))
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    signStruct:=SignStruct{r,s}
    signByte , err :=json.Marshal(signStruct)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    signStr := hex.EncodeToString(signByte)
    return signStr
} 
func (self *SM2) Verify(message,signStr,pubkey string) bool{
    public , err :=hex.DecodeString(pubkey)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    publicKey , err := sm2.ParseSm2PublicKey(public)
    if err!=nil{
        fmt.Println(err)
        return false
    }

    signStruct := SignStruct{}
    signByte , err := hex.DecodeString(signStr)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    json.Unmarshal(signByte,&signStruct)
    var r,s *big.Int
    r = signStruct.R
    s = signStruct.S
    flag := sm2.Verify( publicKey, []byte(message), r,s)
    return flag
}
func (self *SM2) Encrypt(message , pubkey string) string{
    public , err :=hex.DecodeString(pubkey)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    publicKey , err := sm2.ParseSm2PublicKey(public)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    encryptByte, err := publicKey.Encrypt([]byte(message))
    encryptText := hex.EncodeToString(encryptByte)
    return string(encryptText)
}
func (self *SM2) Decrypt(message , prvkey string) string{
    private , _:=hex.DecodeString(prvkey)
    privateKey , err:=sm2.ParsePKCS8UnecryptedPrivateKey(private)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    decryptHexByte,_ := hex.DecodeString(message)
    decryptByte,err  :=  privateKey.Decrypt(decryptHexByte)
    return string(decryptByte)
}

type RSA struct{
    modulusLength string
    namedCipher string
    PEM_PRIVATE_BEGIN string 
    PEM_PRIVATE_END   string 
    PEM_PUBLIC_BEGIN  string 
    PEM_PUBLIC_END    string 
}
func NewRSA(modulusLength ,namedCipher string ) *RSA{
    rsa := new(RSA)
    rsa.modulusLength = modulusLength
    rsa.namedCipher = namedCipher
    rsa.PEM_PRIVATE_BEGIN = "-----BEGIN PRIVATE KEY-----\n"
    rsa.PEM_PRIVATE_END   ="\n-----END PRIVATE KEY-----"
    rsa.PEM_PUBLIC_BEGIN  ="-----BEGIN PUBLIC KEY-----\n"
    rsa.PEM_PUBLIC_END    ="\n-----END PUBLIC KEY-----"
    return rsa
}
func ToPEM(x509_PublicKey []byte) {

}
func (self *RSA) GenerateKeys(files ...string) Key{
    //得到私钥
    privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    //通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
    private := x509.MarshalPKCS1PrivateKey(privateKey) 
    privateBlock := pem.Block{
		Type:  "rsa private key",
		Bytes: private,
	}

    //处理公钥,公钥包含在私钥中
    publickKey := privateKey.PublicKey
    //接下来的处理方法同私钥
    //通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
    public, _ := x509.MarshalPKIXPublicKey(&publickKey)
    publicBlock := pem.Block{
		Type:  "rsa public key",
		Bytes: public,
	}
    //save to privatekey.pem
    for idx,fileName :=range files{
      if idx==0 && fileName!=""{
        privateFile, _ := os.Create(fileName)
        pem.Encode(privateFile, &privateBlock)
        privateFile.Close()
      }else if idx==1 && fileName!=""{
        publicFile, _ := os.Create(fileName)
        pem.Encode(publicFile, &publicBlock)
        publicFile.Close()
      }
    }

    prvkey:=hex.EncodeToString(private)
    pubkey:=hex.EncodeToString(public[:])
    return Key{[]string{prvkey},[]string{pubkey}}
}
func (self *RSA) Encrypt(message , pubkey string) string {
    public , _:=hex.DecodeString(pubkey)
    publicKey , _:=x509.ParsePKIXPublicKey(public)
    
    encryptByte, _ := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(message))
    encryptText := hex.EncodeToString(encryptByte)
    return encryptText
}
func (self *RSA) Decrypt(message , prvkey string) string{
    private , _:=hex.DecodeString(prvkey)
    privateKey , _:=x509.ParsePKCS1PrivateKey(private)

    decryptHexByte,_ :=hex.DecodeString(message)
    decryptByte, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptHexByte)
    
    return string(decryptByte)
}
func (self *RSA) Sign(message,prvkey string) string {  
    private , _:=hex.DecodeString(prvkey)
    privateKey , _:=x509.ParsePKCS1PrivateKey(private)
    messageHash := Hashlib.Sha256(message)
    messageHashByte,_:=hex.DecodeString(messageHash)
    signByte, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, messageHashByte)
    if err!=nil{
        fmt.Printf("Error from signing: %s\n", err)
    }
    signStr := hex.EncodeToString(signByte)
    return signStr
}
func (self *RSA) Verify(message , signStr , pubkey string) bool{
    public , _ :=hex.DecodeString(pubkey)
    publicKey , _ := x509.ParsePKIXPublicKey(public)

    signByte , _ := hex.DecodeString(signStr)
    messageHash := Hashlib.Sha256(message)
    messageHashByte,_:=hex.DecodeString(messageHash)
    //验证签名
    err := rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA256,messageHashByte, signByte) 
    if err != nil {
        return false
    } else {
        return true
    }
}

type ECC struct{
  namedCurve string 
  namedCipher string
  PEM_PRIVATE_BEGIN string 
  PEM_PRIVATE_END   string 
  PEM_PUBLIC_BEGIN  string 
  PEM_PUBLIC_END    string
}
func NewEcc(namedCurve,namedCipher string) *ECC{
  ecc := new(ECC)
  ecc.namedCurve = namedCurve
  ecc.namedCipher = namedCipher
  ecc.PEM_PRIVATE_BEGIN = "-----BEGIN EC PRIVATE KEY-----\n"
  ecc.PEM_PRIVATE_END   ="\n-----END EC PRIVATE KEY-----"
  ecc.PEM_PUBLIC_BEGIN  = "-----BEGIN PUBLIC KEY-----\n"
  ecc.PEM_PUBLIC_END    ="\n-----END PUBLIC KEY-----"   
  return ecc
}
func (self *ECC) GenerateKeys(files ...string)(Key){
    privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    private, _ := x509.MarshalECPrivateKey(privateKey)
    privateBlock := pem.Block{
		Type:  "esdsa private key",
		Bytes: private,
	}
    
    publicKey := privateKey.PublicKey
    //x509序列化
    public, _ := x509.MarshalPKIXPublicKey(&publicKey)
    publicBlock := pem.Block{
		Type:  "esdsa public key",
		Bytes: public,
	}

    //save to privatekey.pem
    for idx,fileName :=range files{
      if idx==0 && fileName!=""{
        privateFile, _ := os.Create(fileName)
        pem.Encode(privateFile, &privateBlock)
        privateFile.Close()
      }else if idx==1 && fileName!=""{
        publicFile, _ := os.Create(fileName)
        pem.Encode(publicFile, &publicBlock)
        publicFile.Close()
      }
    }
    //待研究
    pem.EncodeToMemory(&privateBlock)
    pem.EncodeToMemory(&publicBlock)
    
    prvkey:=hex.EncodeToString(private)
    pubkey:=hex.EncodeToString(public[:])
    return Key{[]string{prvkey},[]string{pubkey}}
}
func (self *ECC) Encrypt(message , pubkey string) string {
    return "尚不支持该功能"
}
func (self *ECC) Decrypt(message , prvkey string) string{
    return "尚不支持该功能"
}

func (self *ECC) Sign(message,prvkey string) string {  
    private , _:=hex.DecodeString(prvkey)
    privateKey , err:=x509.ParseECPrivateKey(private)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    messageHash := Hashlib.Sha256(message)
    r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte(messageHash))
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    //rText ,_:= r.MarshalText()
    //sText ,_:= s.MarshalText()
    //signStruct:=SignStruct{string(rText),string(sText)}
    signStruct:=SignStruct{r,s}
    signByte , err :=json.Marshal(signStruct)
    if err!=nil{
        fmt.Println(err)
        return ""
    }
    signStr := hex.EncodeToString(signByte)
    return signStr
}
func (self *ECC) Verify(message , signStr , pubkey string) bool{
    public , err :=hex.DecodeString(pubkey)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    publicStream , err := x509.ParsePKIXPublicKey(public)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    publicKey :=publicStream.(*ecdsa.PublicKey)
    signStruct := SignStruct{}
    signByte , err := hex.DecodeString(signStr)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    json.Unmarshal(signByte,&signStruct)
    var r,s *big.Int
    //r.UnmarshalText([]byte(signStruct.R))
    //s.UnmarshalText([]byte(signStruct.S))
    r = signStruct.R
    s = signStruct.S
    messageHash := Hashlib.Sha256(message)
    //flag := ecdsa.Verify( publicKey, []byte(messageHash), &r,&s)
    flag := ecdsa.Verify( publicKey, []byte(messageHash), r,s)
    return flag
}



// class Logger {
//   getLogger(name="default",confFile="log4js.json"){
//     var log4js_config = require("./"+confFile)
//     log4js.configure(log4js_config);
//     const logger = log4js.getLogger(name)
//     return logger
//   } 
// }

// logger=new Logger().getLogger()

// class Crypto{
//   toPEM(key,type){
//     type = type.toUpperCase()
//     if (type=="PUBLIC"){
//       return this.PEM_PUBLIC_BEGIN+key+this.PEM_PUBLIC_END
//     }else if (type =="PRIVATE"){
//       return this.PEM_PRIVATE_BEGIN+key+this.PEM_PRIVATE_END
//     }else {
//       return null
//     }
//   }
//   sign(message,prvkey=null,prvfile=null){
//     let signature=null
//     try{
//       if (prvfile)
//         prvkey = fs.readFileSync(prvfile,"utf8")
//       if (prvkey){
//         const signObj = crypto.createSign('sha256')
//         signObj.update(message)
//         const prvkeyPEM = this.toPEM(prvkey,"private")
//         const signStr = signObj.sign(prvkeyPEM).toString('base64');
//         return signStr  
//       }else{
//         return null
//       }
//     }catch(e){
//       throw e
//     }
//   }
//   verify(message,signStr,pubkey=null,pubfile=null){
//     let verify=null
//     try{
//       if (pubfile)
//         pubkey = fs.readFileSync(pubfile,"utf8")
//       if (pubkey){
//         const verifyObj = crypto.createVerify('sha256')
//         verifyObj.update(message)
//         const pubkeyPEM = this.toPEM(pubkey,"public")
//         const verifyBool = verifyObj.verify(pubkeyPEM,Buffer.from(signStr,"base64"));
//         return verifyBool  
//       }else{
//         return false
//       }
//     }catch(e){
//       console.log(`error ${e.name} with ${e.message}`)
//       return false
//     }
//   }
//   encrypt(message,pubkey=null,pubfile=null){
//     let encrypted=null
//     try{
//       if (pubfile)
//         pubkey = fs.readFileSync(pubfile,"utf8")
//       if (pubkey){
//         const pubkeyPEM = this.toPEM(pubkey,'public')
//         encrypted = crypto.publicEncrypt({key:pubkeyPEM},Buffer.from(message)).toString('base64')
//         return encrypted
//       }else {
//         return null
//       }  
//     }catch(e){
//       console.log(`error ${e.name} with ${e.message}`)
//       return null
//     }
//   }
//   decrypt(message,prvkey=null,prvfile=null){
//     let decrypted=null
//     try{
//       if (prvfile)
//         prvkey = fs.readFileSync(prvfile,"utf8")
//       if (prvkey){
//         const prvkeyPEM = this.toPEM(prvkey,'private')
//         decrypted = crypto.privateDecrypt({key:prvkeyPEM},Buffer.from(message,'base64')).toString()
//         return decrypted
//       }else {
//         return null
//       }  
//     }catch(e){
//       console.log(`error ${e.name} with ${e.message}`)
//       return null
//     }
//   }
//   enCipher(message,key){
//     let encipher = crypto.createCipher(this.namedCipher,key)
//     let encrypted = encipher.update(JSON.stringify(message),"utf8","base64")
//     encrypted += encipher.final('base64') 
//     return encrypted
//   }
//   deCipher(message,key){
//     let decipher = crypto.createDecipher(this.namedCipher,key)
//     let decrypted = decipher.update(message,"base64",'utf8')
//     decrypted += decipher.final('utf8')
//     return decrypted
//   }  
// }
// class ECC extends Crypto{
//   constructor(namedCurve='secp256k1',namedCipher='aes192'){
//     super()
//     this.namedCurve = namedCurve
//     this.namedCipher = namedCipher
//     this.PEM_PRIVATE_BEGIN = "-----BEGIN EC PRIVATE KEY-----\n"
//     this.PEM_PRIVATE_END="\n-----END EC PRIVATE KEY-----"
//     this.PEM_PUBLIC_BEGIN = "-----BEGIN PUBLIC KEY-----\n"
//     this.PEM_PUBLIC_END="\n-----END PUBLIC KEY-----"
//   }
//   generateKeys(prvfile="private",pubfile="public"){
//     try{
//      const key = crypto.generateKeyPairSync("ec",{
//         namedCurve       :this.namedCurve,
//         publicKeyEncoding:{
//           type  :"spki",
//           format:"der"
//         },
//         privateKeyEncoding:{
//           type  :"sec1",
//           format:"der"
//         }
//      })
//      const pubkey = key.publicKey.toString("base64")
//      const prvkey = key.privateKey.toString("base64")
//      if (prvfile){
//        fs.writeFileSync(prvfile,prvkey)
//      }
//      if (pubfile){
//        fs.writeFileSync(pubfile,pubkey)
//      }
//      return {"prvkey":prvkey,"pubkey":pubkey}
//     }catch(e){
//      console.log(`error ${e.name} with ${e.message}`)
//      return null
//     } 
//   }
//   encrypt(){
//     console.log("尚不支持改功能")
//   }
//   decrypt(){
//     console.log("尚不支持改功能")
//   }
//   genECDH(){
//     const ecdh = crypto.createECDH(this.namedCurve)
//     const pubkey = ecdh.generateKeys("base64")
//     const prvkey = ecdh.getPrivateKey("base64")
//     return {prvkey,pubkey}
//   }  
//   computeSecret(prvkey,pubkey){
//     const ecdh = crypto.createECDH(this.namedCurve)
//     ecdh.setPrivateKey(prvkey,"base64")
//     return ecdh.computeSecret(Buffer.from(pubkey,"base64")).toString("base64")
//   }
//   getKeys(prvkey){
//     const bufferlib = new Bufferlib()
//     const ecdh = crypto.createECDH(this.namedCurve)
//     prvkey = bufferlib.transfer(prvkey,"base64","hex")
//     ecdh.setPrivateKey(prvkey,"hex")
//     const pubkey = ecdh.getPublicKey("hex")
//     /*组装public key der,转换为base64
//     3056【sequence 类型 长度86】
//     3010【sequence 类型 长度16】
//     0607【OID类型 长度 07】
//     2a8648ce3d0201 【 OID value = "1.2.840.10045.2.1"=>{42,134,72,206,61,2,1}】
//     0605【OID类型 长度05】
//     2b8104000a【OID value = "1.3.132.0.10"=>{43,129,04,00,10}=>{0x 2b 81 04 00 0a}】
//     034200【bit string类型，长度66，前导00】
//     */
//     const pubkey_der="3056301006072a8648ce3d020106052b8104000a034200"+pubkey
//     /*组装private key der,转换为base64
//     3074【sequence类型，长度116】
//     0201【Integer类型，长度01】
//     01 【value=1 ，ecprivkeyVer1=1】
//     0420【byte类型，长度32】
//     ....【私钥】
//     a007【a0结构类型，长度07】
//     0605【OID类型，长度05】
//     2b8104000a【OID value named secp256k1 elliptic curve = 1.3.132.0.10 =>{43,129,04,00,10}=>{0x 2b 81 04 00 10}】
//     a144【a1结构类型，长度68】
//     034200【bitstring类型，长度66，前导00】
//    【0x 04开头的非压缩公钥】
//     */
//     const prvkey_der="30740201010420"+prvkey+
//                      "a00706052b8104000aa144034200"+pubkey
  
//     return {"prvkey":bufferlib.transfer(prvkey_der,"hex","base64"),
//             "pubkey":bufferlib.transfer(pubkey_der,"hex","base64")}
//   }


//   genKeys(keyStr,num){
//     if (!num) num=1
//     const hashlib=new Hashlib()
//     const bufferlib = new Bufferlib()
//     let seed  = hashlib.sha512(keyStr)
//     let keys=[]
//     for (let i=0 ;i<num;i++){
//       const temp = hashlib.sha512(seed)      
//       seed =temp.slice(64,128)
//       const prvkey=Buffer.from(temp.slice(0,64),'hex').toString('base64')
//       keys.push(this.getKeys(prvkey))
//     }
//     return keys   
//   }
  
// }
// class RSA extends Crypto{
//   constructor(modulusLength=1024,namedCipher='aes192'){
//     super()
//     this.modulusLength = modulusLength
//     this.namedCipher = namedCipher
//     this.PEM_PRIVATE_BEGIN = "-----BEGIN PRIVATE KEY-----\n"
//     this.PEM_PRIVATE_END="\n-----END PRIVATE KEY-----"
//     this.PEM_PUBLIC_BEGIN = "-----BEGIN PUBLIC KEY-----\n"
//     this.PEM_PUBLIC_END="\n-----END PUBLIC KEY-----"
//   }
//   generateKeys(prvfile="private",pubfile="public"){
//     try{
//      const key = crypto.generateKeyPairSync("rsa",{
//         modulusLength    :this.modulusLength,
//         publicKeyEncoding:{
//           type  :"spki",
//           format:"der"
//         },
//         privateKeyEncoding:{
//           type  :"pkcs8",
//           format:"der"
//         }
//      })
//      const pubkey = key.publicKey.toString("base64")
//      const prvkey = key.privateKey.toString("base64")
//      if (prvfile){
//        fs.writeFileSync(prvfile,prvkey)
//      }
//      if (pubfile){
//        fs.writeFileSync(pubfile,pubkey)
//      }
//      return {"prvkey":prvkey,"pubkey":pubkey}
//     }catch(e){
//      console.log(`error ${e.name} with ${e.message}`)
//      return null
//     } 
//   }
// }
    
// class MySet{
  
//   union(a,b){
//     let c = new Set([...a, ...b]);
//     return [...c];
//   }
//   difference(a,b){
//    let m = new Set([...a])
//    let n = new Set([...b])
//    let c = new Set([...m].filter(x => !n.has(x)));
//    return [...c];
//   }
//   intersect(a,b){
//     let m = new Set([...a])
//     let n = new Set([...b])
//     let c = new Set([...m].filter(x => n.has(x)));//ES6
//     return [...c];
//   }
//   removeRepeat(a){
//     let c = new Set([...a]);
//     return [...c];
//   }
// }

// class MyHttp{
//   async get(urls){
//     const promiseArray = urls.map(
//       url => this.httpGet(url))
//     return Promise.all(promiseArray)
//   }
//   async httpGet(url){
//     return new Promise((resolve,reject)=>{
//       var urlObj=new URL(url)
//       var options = { 
//           hostname: urlObj.hostname, 
//           port: urlObj.port, 
//           path: urlObj.pathname, 
//           method: 'GET',
//       };
//       var req = http.request(options, function (res) {
//         res.setEncoding('utf8');
//         let rawData = '';
//         res.on('data', (chunk) => { rawData += chunk; });
//         res.on('end', () => {
//           if (res.statusCode==200){
//             resolve(rawData);
//           }else{
//             reject(res.statusText);
//           }
//         });
//       }); 
         
//       req.on('error', function (e) { 
//           reject('problem with request: ' + e.message); 
//       }); 
         
//       req.end();
//     })
//   }
//   async httpPost(url,data){
//     return new Promise((resolve,reject)=>{
//       var postData = querystring.stringify(data)
//       var urlObj=new URL(url)
//       var options = { 
//           hostname: urlObj.hostname, 
//           port: urlObj.port, 
//           path: urlObj.pathname, 
//           method: 'POST',
//       };
//       var req = http.request(options, function (res) {
//         res.setEncoding('utf8');
//         let rawData = '';
//         res.on('data', (chunk) => { rawData += chunk; });
//         res.on('end', () => {
//           if (res.statusCode==200){
//             resolve(rawData);
//           }else{
//             reject(res.statusText);
//           }
//         });
//       }); 
         
//       req.on('error', function (e) { 
//           reject('problem with request: ' + e.message); 
//       }); 
//       console.log(postData)
//       req.write(postData)
//       req.end();
//     })
//   }
// }
// exports.obj2json = function(obj){
//   return JSON.parse(JSON.stringify(obj))
// }

// exports.ecc  = new ECC()  
// exports.rsa  = new RSA() 
// exports.hashlib = new Hashlib()
// exports.bufferlib  = new Bufferlib()
// exports.logger  = new Logger()
// exports.set     = new MySet()
// exports.db      = new DB()
// exports.http    = new MyHttp()
