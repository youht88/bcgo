package utils

import (
    //"fmt"
    //"os"
    "runtime"
    //"net/http"
)
var (
  //Max_Num = os.Getenv("MAX_NUM")
  MaxWorker = runtime.NumCPU()
  MaxQueue = 1000
  JobQueue chan Job
)

func init(){
  //runtime.GOMAXPROCS = MaxWorker
  JobQueue = make(chan Job, MaxQueue)
  dispatcher := NewDispatcher(MaxWorker)
  dispatcher.Run()
}

type JobResult struct{
    Result interface{}
    Err    error
}
type Job struct {
  Func func(...interface{}) (interface{},error)
  Args []interface{}
  C   chan JobResult
}
func NewJob (Func func(...interface{}) (interface{},error),Args ...interface{}) Job{
  job:=Job{
    Func : Func,
    Args : Args,
    C : make(chan JobResult),
  }
  JobQueue <- job
  return job
}
type Worker struct {
  WorkerPool chan chan Job
  JobChannel chan Job
  Quit       chan bool
}
func NewWorker(workPool chan chan Job) Worker {
  return Worker {
    WorkerPool:workPool,
    JobChannel:make(chan Job),
    Quit:make(chan bool),
  }
}
func (w Worker) Start() {
  go func() {
    for {
      w.WorkerPool <- w.JobChannel
      select {
        case job:= <- w.JobChannel:
          // excute job
          result,err:=job.Func(job.Args...)
          job.C <- JobResult{result,err}
        case <- w.Quit:
          return
      }
    }
  }()
}
func (w Worker) Stop() {
  go func() {
    w.Quit <- true
  }()
}
type Dispatcher struct {
  MaxWorkers int
  WorkerPool chan chan Job
  Quit       chan bool
}
func NewDispatcher(maxWorkers int) *Dispatcher {
  pool := make(chan chan Job, maxWorkers)
  return &Dispatcher{MaxWorkers:maxWorkers, WorkerPool:pool,Quit:make(chan bool)}
}
func (d *Dispatcher) Run() {
  for i:=0; i<d.MaxWorkers; i++ {
    worker := NewWorker(d.WorkerPool)
    worker.Start()
  }
  go d.Dispatch()
}
func (d *Dispatcher) Stop() {
  go func() {
    d.Quit <- true
  }()
}
func (d *Dispatcher) Dispatch() {
  for {
    select {
      case job:=<- JobQueue:
        go func(job Job) {
          jobChannel := <- d.WorkerPool
          jobChannel <- job
        }(job)
      case <- d.Quit:
        return
      }
  }
}
