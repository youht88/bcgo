package utils

import (
    "encoding/hex"
    //"encoding/json"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"

    "github.com/tjfoc/gmsm/sm3"
)

var Hashlib hashlibIface 

func init(){
  Hashlib = new(HashGo)
}

type hashlibIface interface{
    Sha256(data string) string
    Sha256raw(data []byte) string
    Sha512(data string) string
    Sha1(data string) string
    Md5(data string) string
    DoubleSha256(data string) string
}

type HashGo struct{}

func (self *HashGo) Sha512(data string) string{
    sum := sha512.Sum512([]byte(data))
    return hex.EncodeToString(sum[:])
}
func (self *HashGo) Sha1(data string) string{
    sum := sha1.Sum([]byte(data))
    return hex.EncodeToString(sum[:])
}
func (self *HashGo) Md5(data string) string{
    sum := md5.Sum([]byte(data))
    return hex.EncodeToString(sum[:])
}
func (self *HashGo) Sha256(data string) string{
    sum := sha256.Sum256([]byte(data))
    return hex.EncodeToString(sum[:])
}
func (self *HashGo) Sha256raw(data []byte) string{
    sum := sha256.Sum256(data)
    return hex.EncodeToString(sum[:])
}
func (self *HashGo) DoubleSha256(data string) string{
    return self.Sha256(self.Sha256(data))
}

type HashGm struct{
    HashGo
}
func (self *HashGm) Sha256(data string) string{
    h := sm3.New()
    h.Write([]byte(data))
    sum := h.Sum(nil)
    return hex.EncodeToString(sum[:])
}
func (self *HashGm) Sha256raw(data []byte) string{
    h := sm3.New()
    h.Write(data)
    sum := h.Sum(nil)
    return hex.EncodeToString(sum[:])
}

func (self *HashGm) DoubleSha256(data string) string{
    return self.Sha256(self.Sha256(data))
}

