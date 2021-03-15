package utils

import (
  //"fmt"
  "bytes"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/rpc"
  "github.com/gorilla/rpc/json"
)

type _Http struct{
  Router *mux.Router
  RpcServer *rpc.Server
}
var Http *_Http

func init(){
  Http = new(_Http)
  s := rpc.NewServer()
  s.RegisterCodec(json.NewCodec(), "application/json")
  s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
  
  Http.Router = mux.NewRouter()
  Http.RpcServer = s
}

func (self *_Http) Listen(url string){
  //start HttpServer
  Logger.Successf("节点Http服务已启动:%s",url)
  err := http.ListenAndServe(url,self.Router)
  if err!=nil{
     Logger.Danger(err)
  }
}

func (self *_Http) RegisterRpcService(rpcService interface{},rpcPath string){
    err:=self.RpcServer.RegisterService(rpcService, "")
    if err!=nil{
      Logger.Warn(err)
    }else{
      self.Router.Handle(rpcPath, self.RpcServer)
    }
}


//Result of RPC call is of this type
func RpcCall(url string,method string,args interface{})(interface{},error){
    defer func(){
        if err:=recover();err!=nil{
            Logger.Error(err)
        }else{
            //Logger.Warn("!!!!!!!!")
        }
    }()
    job := NewJob(
        func(jobargs ...interface{})(interface{},error){
            url := jobargs[0].(string)
            method:= jobargs[1].(string)
            args := jobargs[2].(interface{})

            message, err := json.EncodeClientRequest(method, &args)
            if err != nil {
              Logger.Errorf("%s", err)
              return nil,err
            }
            req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
            if err != nil {
              Logger.Errorf("%s", err)
              return nil,err
            }
            req.Header.Set("Content-Type", "application/json")
            client := new(http.Client)
            resp, err := client.Do(req)
            if err != nil {
              Logger.Errorf("Error in sending request to %s. %s", url, err)
              return nil,err
            }
            defer resp.Body.Close()

            var result interface{}
            err = json.DecodeClientResponse(resp.Body, &result)
            if err != nil {
              Logger.Errorf("Couldn't decode response. %s", err)
              return nil,err
            }
            return result,nil
        },url,method,args)  
    rst := <- job.C
    return rst.Result,rst.Err
}