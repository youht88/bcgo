package utils

const REWARD float32 = 2.0
const BLOCK_PER_HOUR int = 3*60  //每小时出块数限制
const ADJUST_DIFF int =100   //每多少块调整一次难度
const ZERO_DIFF int = 6*4
const NUM_FORK int = 6
const TRANSACTION_TO_BLOCK int = 0
const SYNC_BLOCKCHAIN int = 10*1000*60  //多少毫秒同步blockchain
const CHECK_NODE int =  1000*60 //多少毫秒检查节点连接情况
const contractTimeout int = 5000
//global.emitter = new EventEmitter()

var DiffcultIndex int
var Diffcult int
var LockTime int64