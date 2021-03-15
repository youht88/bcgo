package utils

import (
  "io/ioutil"
  "bytes"
  "encoding/gob"
)
var Fs =new(_Fs)

type _Fs struct{}

func (self _Fs) WriteFileSync(file string,data []byte) error{
    return ioutil.WriteFile(file,data,0666)
}

func (self _Fs) ReadFileSync(file string)([]byte,error){
    return ioutil.ReadFile(file)
}

func DeepCopy(dst, src interface{}) error {
    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(src); err != nil {
        return err
    }
    return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}