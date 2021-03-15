package main1

import (
  "log"
  "fmt"
  "flag"
  "bcgo/utils"
  "bcgo/node"
  "bcgo/wallet"
)

func main(){

  node := node.New(map[string]interface{}{
    "config":config,
    "httpServer":args.httpServer,
    "entryNode":args.entryNode,
    "entryKad":args.entryKad,
    "me":args.me,
    "db":args.db,
    "display":args.display,
    "ioServer":ioServer,
    "ioClient":ioClient
  })
  
  //链接数据库
  logger.debug("dbConnect...")
  node.dbConnect()

  //创建钱包  
  mywallet := wallet.New()
  go mywallet.ChooseByName(args.me)
  /*  .catch(async e=>{
      logger.error(`尚没有钱包，准备创建${args.me}的密钥钱包`)
      mywallet.create(args.me)
        .then(()=>logger.info("钱包创建成功"))
        .catch(e=>console.log("error2",e))
    })
  */
  node.wallet = mywallet
  logger.debug("mywallet.address",mywallet.address)
   
  //导入本地区块链
  go node.SyncLocalChain()
  logger.debug(`localchain has ${node.blockchain.maxindex()} blocks`)
  //链接网络
  logger.debug("socketioConnect...")
  go node.socketioConnect()
  if args.entryNode==args.me{
    node.emitter.emit("start")
  }
}

/*
let node
const start= async ()=>{
  //make node 
  node = new Node({
    "config":config,
    "httpServer":args.httpServer,
    "entryNode":args.entryNode,
    "entryKad":args.entryKad,
    "me":args.me,
    "db":args.db,
    "display":args.display,
    "ioServer":ioServer,
    "ioClient":ioClient
  })
  node.initEvents()

  //链接数据库
  logger.debug("dbConnect...")
  await node.dbConnect()
  //创建钱包  
  const mywallet = new Wallet()
  await mywallet.chooseByName(args.me)
    .catch(async e=>{
      logger.error(`尚没有钱包，准备创建${args.me}的密钥钱包`)
      mywallet.create(args.me)
        .then(()=>logger.info("钱包创建成功"))
        .catch(e=>console.log("error2",e))
    })
  node.wallet = mywallet
  logger.debug("mywallet.address",mywallet.address)
  //导入本地区块链
  await node.syncLocalChain()
  logger.debug(`localchain has ${node.blockchain.maxindex()} blocks`)
  //链接网络
  logger.debug("socketioConnect...")
  await node.socketioConnect()
  if (args.entryNode==args.me)
    node.emitter.emit("start")
}
start()
  .then(()=>console.log("node started."))
  .catch(e=>console.log(e))
*/
