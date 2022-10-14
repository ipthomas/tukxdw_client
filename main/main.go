package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukutil"
	"github.com/ipthomas/tukxdw"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Dir(b)
	LogFile    *os.File
	//TUK_DB_URL_AWS          = "https://5k2o64mwt5.execute-api.eu-west-1.amazonaws.com/beta/"
	DSUB_BROKER_URL         = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritXDSDsub/Dsub"
	DSUB_CONSUMER_URL_AWS   = "https://cjrvrddgdh.execute-api.eu-west-1.amazonaws.com/beta/"
	DSUB_CONSUMER_URL_Local = "http://tukeventserver.ddns.net:8081/eventservice/event"
)

func main() {
	registerXDWs()
}
func registerXDWs() {
	initLog()
	initDB()
	if xdwDefFiles, err := loadWorkflowDefinitionFolder(); err == nil {
		var err error
		trans := tukxdw.XDWTransaction{Action: tukcnst.REGISTER, DSUB_BrokerURL: DSUB_BROKER_URL, DSUB_ConsumerURL: DSUB_CONSUMER_URL_Local}
		for _, file := range xdwDefFiles {
			if trans.WorkflowDefinition, err = tukxdw.NewWorkflowDefinition("", loadWorkflowDefinitionFile(file)); err == nil {
				tukxdw.New_Transaction(&trans)
			}
		}
	}
	tukdbint.DBConn.Close()
	LogFile.Close()
}
func loadWorkflowDefinitionFile(file fs.DirEntry) []byte {
	if strings.HasSuffix(file.Name(), ".json") && strings.Contains(file.Name(), "_xdwdef") {
		xdwfile, err := os.ReadFile(Basepath + "/xdwconfig/" + file.Name())
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		log.Println("Loaded WF Def for Pathway : " + file.Name())
		return xdwfile
	}
	return nil
}
func loadWorkflowDefinitionFolder() ([]fs.DirEntry, error) {
	var err error
	var f *os.File
	var fileInfo []fs.DirEntry
	f, err = os.Open(Basepath + "/xdwconfig/")
	if err != nil {
		log.Println(err)
		return fileInfo, err
	}
	fileInfo, err = f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}
	return fileInfo, err
}
func initLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	lf, err := os.OpenFile(Basepath+"/"+tukutil.Tuk_Month()+tukutil.Tuk_Day()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		LogFile = lf
		log.SetOutput(LogFile)
		log.Println("-----------------------------------------------------------------------------------")
		log.Println("Loaded log file - " + LogFile.Name())
	} else {
		log.Println(err.Error())
	}
}
func initDB() {
	dbconn := tukdbint.TukDBConnection{
		DBUser:     "root",
		DBPassword: "rootPass",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "tuk",
	}
	tukdbint.NewDBEvent(&dbconn)
}
