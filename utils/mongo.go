package utils

import (
    //"fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type M bson.M

var DB *_DB

func init(){
    DB=new(_DB)
}

type _DB struct{
    session *mgo.Session
    db *mgo.Database
}

func (self *_DB) Init(url,database string) error{
  session, err := mgo.Dial(url)
  if err!=nil{
      return err
  }
  self.session = session
  self.db = session.DB(database)
  return err
}
func (self *_DB) Close(){
  self.session.Close()
}
func (self *_DB) Insert(collection string,doc ...interface{}) error {
  var table *mgo.Collection = self.db.C(collection)
  err := table.Insert(doc...)
  return err
}
func (self *_DB) Update(collection string,condition interface{},doc interface{}) error {
  var table *mgo.Collection = self.db.C(collection)
  err := table.Update(condition,doc)
  return err
}
func (self *_DB) UpdateId(collection,id string,doc interface{}) error {
  var table *mgo.Collection = self.db.C(collection)
  err := table.Update(bson.M{"_id":bson.ObjectIdHex(id)},doc)
  return err
}
func (self *_DB) Upsert(collection string,condition interface{}, doc interface{}) (info *mgo.ChangeInfo, err error){
  var table *mgo.Collection = self.db.C(collection)
  return table.Upsert(condition,doc)
}
func (self *_DB) UpsertId(collection,id string, doc interface{}) (info *mgo.ChangeInfo, err error){
  var table *mgo.Collection = self.db.C(collection)
  return table.Upsert(bson.M{"_id":bson.ObjectIdHex(id)},doc)
}
func (self *_DB) Delete(collection string,condition interface{}) error {
  var table *mgo.Collection = self.db.C(collection)
  err := table.Remove(condition)
  return err
}
func (self *_DB) DeleteId(collection,id string) error {
  var table *mgo.Collection = self.db.C(collection)
  err := table.Remove(bson.M{"_id":bson.ObjectIdHex(id)})
  return err
}

func (self *_DB) Find(collection string,condition interface{}) *mgo.Query {
  var table *mgo.Collection = self.db.C(collection)
  query := table.Find(condition)
  return query
}
func (self *_DB) FindId(collection,id string) *mgo.Query {
  var table *mgo.Collection = self.db.C(collection)
  query := table.Find(bson.M{"_id":bson.ObjectIdHex(id)})
  return query
}

func (self *_DB) Pipe(collection string,pipeline interface{}) *mgo.Iter {
  println(self.db)
  var table *mgo.Collection = self.db.C(collection)
  iter := table.Pipe(pipeline).Iter()
  return iter
}

// class DB{
//   constructor(){
//     this.client = {}    
//     this.db = null
//   }
  
//   async init(url){
//     return new Promise((resolve,reject)=>{
//       let conn = url.split("/")
//       let newUrl = "mongodb://"+conn[0]
//       let database = conn[1]
//       MongoClient.connect(newUrl,{useNewUrlParser:true},(err,client)=>{
//         if (err) reject(err)
//         this.client = client
//         this.db = client.db(database)
//         console.log("数据库已链接")
//         resolve(this.db)
//       })
//     })
//   }
//   close(){
//     if (!this.client) return
//     this.client.close()
//     console.log("数据库已关闭")
//   }
//   async deleteOne(collection,condition,cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         coll.deleteOne(condition,(err,result)=>{
//           if (typeof(cb)=="function"){
//             cb(err,result)
//             resolve(result)
//           }else if (err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })
//       })
//     })
//   }
//   async deleteMany(collection,condition,cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         coll.deleteMany(condition,(err,result)=>{
//           if (typeof(cb)=="function"){
//             cb(err,result)
//             resolve(result)
//           }else if (err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })
//       })
//     })
//   }
//   async insertOne(collection,doc,options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject (new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options={}
//         }
//         coll.insertOne(doc,options,(err,result)=>{
//           if (typeof(cb)=="function"){
//             cb(err,result)
//             resolve(result)
//           }else if (err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })
//       })
//     })
//   }
//   async insertMany(collection,docs,options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.insertMany(docs,options,(err,result)=>{
//           if (typeof(cb)=="function"){
//             cb(err,result)
//             resolve(result)
//           }else if (err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })
//       })
//     })
//   }
//   async findOne(collection,condition={},options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.findOne(condition,options,(err,result)=>{
//           if (typeof(cb) == "function"){
//             cb(err,result)
//             resolve(result)
//           }else if(err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })    
//       })
//     })
//   }
//   async findMany(collection,condition={},options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.find(condition,options,(err,result)=>{
//           if (typeof(cb) == "function"){
//             cb(err,result.toArray())
//             resolve(result)
//           }else if(err){
//             throw err
//           }else{
//             resolve(result.toArray())
//           }
//         })    
//       })
//     })
//   }
//   async updateOne(collection,condition={},update={},options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.updateOne(condition,update,options,(err,result)=>{
//           if (typeof(cb) == "function"){
//             cb(err,result)
//             resolve(result)
//           }else if(err){
//             reject(err)
//           }else{
//             resolve(result)
//           }
//         })    
//       })
//     })
//   }
//   async updateMany(collection,condition={},update={},options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.updateMany(condition,update,options,(err,result)=>{
//           if (typeof(cb) == "function"){
//             cb(err,result)
//             resolve(result)
//           }else if(err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })    
//       })
//     })
//   }
//   async count(collection,condition={},options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.countDocuments(condition,options,(err,result)=>{
//           if (typeof(cb) == "function"){
//             cb(err,result)
//             resolve(result)
//           }else if(err){
//             throw err
//           }else{
//             resolve(result)
//           }
//         })    
//       })
//     })
//   }
//   async aggregate(collection,pipeline=[],options={},cb=null){
//     return new Promise((resolve,reject)=>{
//       if (!this.db) reject(new Error("not init database"))
//       this.db.collection(collection,(err,coll)=>{
//         if (err) throw err
//         if (typeof(options)=="function"){
//           cb = options
//           options = {}
//         }
//         coll.aggregate(pipeline,options,(err,result)=>{
//           if (typeof(cb) == "function"){
//             cb(err,result.toArray())
//             resolve(result)
//           }else if(err){
//             throw err
//           }else{
//             resolve(result.toArray())
//           }
//         })    
//       })
//     })
//   }
// }
