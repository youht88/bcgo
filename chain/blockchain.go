package chain

import (
  "bcgo/utils"
  "bcgo/block"
)

type UTXO struct{
  utxoSet map[string]interface{}
  name string
}
func NewUTXO(name string) *UTXO{
    utxo:=new(UTXO)
    /* sample struct like this 
       {3a75be...:[{"index":0,"txout":TXout1},{"index":1,"txout":TXout2}],
        m09qf3...:[{"index":0,"txout":TXout3}]}
    */
    utxo.name = name
    return utxo
}

type Chain struct{
    Blocks []*block.Block
    UTXO     
}
func newChain(name string) *Chain{
    chain:=new(Chain)
    if name=="" {
      name = "main"
    }
    chain.UTXO = *NewUTXO(name)
    return chain
}

func (self *Chain) IsValid() bool{
    blocks := self.Blocks[:]
    utils.Logger.Info("verifing blockchain...",len(blocks))
    for index:=1 ;index < len(blocks);index ++{
      curBlock := blocks[index]
      prevBlock := blocks[index - 1]
      if (prevBlock.Index+1 != curBlock.Index){
        utils.Logger.Error("index error",prevBlock.Index,curBlock.Index)
        return false
      }
      if (!curBlock.IsValid()){
        //checks the hash
        utils.Logger.Errorf("curBlock %v-%v  false",index,curBlock.Nonce)
        return false
      }
      if (prevBlock.Hash != curBlock.PrevHash){
        utils.Logger.Error("block ",curBlock.Index," hash error",prevBlock.Hash,curBlock.PrevHash)
        return false
      }
    }
    return true
}
func (self *Chain) Save() bool{
  for _,block := range self.Blocks{
    block.Save()
  }
  return true
}
func (self *Chain) maxindex() int{
    if len(self.Blocks)==0 {
      return -1
    }
    return len(self.Blocks) - 1
  }

func (self *Chain) GetRangeBlocks(start,end int) []*block.Block{
    maxindex := self.maxindex()
    if (start>maxindex){
      end=maxindex
    }
    if (start<0 || start>maxindex){
      return []*block.Block{}
    }
    if (end<start  || end>maxindex){
      return []*block.Block{}
    }
    blocks := self.Blocks[start:end + 1]
    return blocks
  }

func (self *Chain) FindBlockByIndex(index int) (*block.Block,bool){
    emp := new(block.Block)
    if index<0 {
        return emp,false
    }
    if len(self.Blocks) >= index + 1 {
      return self.Blocks[index],true
    }
    return emp,false
}
func (self *Chain) FindBlockByHash(uhash string) (*block.Block,bool){
    emp := new(block.Block)
    for _,b := range self.Blocks {
      if b.Hash == uhash {
        return b,true
      }
    }
    return emp,false
}

func (self *Chain) Lastblock() *block.Block {
    return self.Blocks[len(self.Blocks)-1]
}

func (self *Chain) GetSPV() []map[string]interface{} {
    var chainSPV []map[string]interface{}
    for _,block := range self.Blocks {
       item := map[string]interface{}{
              "txCount":len(block.Data),
              "diffcult": block.Diffcult, 
              "hash":block.Hash, 
              "index": block.Index, 
              "merkleRoot":block.MerkleRoot,  
              "nonce": block.Nonce, 
              "prev_hash":block.PrevHash,  
              "timestamp": block.Timestamp}
      chainSPV = append(chainSPV,item)
    }
    return chainSPV
  }