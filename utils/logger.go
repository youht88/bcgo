package utils

import (
    "log"
    "fmt"
    "os"
)

type _Logger struct {
   logger *log.Logger
   redPrefix string
   bluePrefix string
   yellowPrefix string
   greenPrefix string
}
var Logger *_Logger
func init(){
  Logger = new(_Logger)
  Logger.logger = log.New( os.Stdout, "", log.Lshortfile|log.LstdFlags)
}
func (self *_Logger) Red(data string) string{
  return "\033[1;31m"+data+"\033[0m"
}
func (self *_Logger) Blue(data string) string{
  return "\033[1;34m"+data+"\033[0m"
}
func (self *_Logger) Yellow(data string) string{
  return "\033[1;33m"+data+"\033[0m"
}
func (self *_Logger) Green(data string) string{
  return "\033[1;32m"+data+"\033[0m"
}

func (self *_Logger) Danger(data ...interface{}){
    self.logger.SetPrefix(self.Red("[Danger]"))
    s:=fmt.Sprint(data...)
    self.logger.Output(3,s)
    panic(s)
}
func (self *_Logger) Error(data ...interface{}){
    self.logger.SetPrefix(self.Red("[Error]"))
    s:=fmt.Sprint(data...)
    self.logger.Output(3,s)
}
func (self *_Logger) Success(data ...interface{}){
    self.logger.SetPrefix(self.Green("[Success]"))
    s:=fmt.Sprint(data...)
    self.logger.Output(3,s)
}

func (self *_Logger) Warn(data ...interface{}){
  self.Warning(data...)
}
func (self *_Logger) Warning(data ...interface{}){
    self.logger.SetPrefix(self.Yellow("[Warn]"))
    s:=fmt.Sprint(data...)
    self.logger.Output(3,s)
}
func (self *_Logger) Info(data ...interface{}){
    self.logger.SetPrefix(self.Blue("[Info]"))
    s:=fmt.Sprint(data...)
    self.logger.Output(3,s)
}

func (self *_Logger) Dangerf(format string,data ...interface{}){
    self.logger.SetPrefix(self.Red("[Danger]"))
    s:=fmt.Sprintf(format,data...)
    self.logger.Output(3,s)
    panic(s)
}
func (self *_Logger) Errorf(format string,data ...interface{}){
    self.logger.SetPrefix(self.Red("[Error]"))
    s:=fmt.Sprintf(format,data...)
    self.logger.Output(3,s)
}
func (self *_Logger) Successf(format string,data ...interface{}){
    self.logger.SetPrefix(self.Green("[Success]"))
    s:=fmt.Sprintf(format,data...)
    self.logger.Output(3,s)
}
func (self *_Logger) Warnf(format string,data ...interface{}){
    self.Warningf(format,data...)
}
func (self *_Logger) Warningf(format string,data ...interface{}){
    self.logger.SetPrefix(self.Yellow("[Warn]"))
    s:=fmt.Sprintf(format,data...)
    self.logger.Output(3,s)
}
func (self *_Logger) Infof(format string,data ...interface{}){
    self.logger.SetPrefix(self.Blue("[Info]"))
    s:=fmt.Sprintf(format,data...)
    self.logger.Output(3,s)
}
