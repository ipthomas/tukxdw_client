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
	pathway               = "pathalert"
	nhs                   = "9999999468"
	initlogging           = true
	initdb                = true
	registerEventServices = true
	registerXdws          = false
	initTmplts            = false
	initStatic            = false
	contentCreator        = false
	contentConsumer       = false
	contentUpdator        = false
	//TUK_DB_URL_AWS          = "https://5k2o64mwt5.execute-api.eu-west-1.amazonaws.com/beta/"
	DSUB_BROKER_URL       = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritXDSDsub/Dsub"
	DSUB_CONSUMER_URL_AWS = "https://fwa7l2kp71.execute-api.eu-west-1.amazonaws.com/beta/eventservice/event"
	//DSUB_CONSUMER_URL_Local = "http://tukeventserver.ddns.net:8081/eventservice/event"
	user       = "pbradley"
	role       = "Clinical"
	org        = "lth"
	notes      = "User " + user + " from " + org + " in the role of " + role + " created new " + pathway + " Workflow"
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Dir(b)
	LogFile    *os.File
)

func main() {
	initVars()
	initLog()
	initDB()
	initTemplates()
	initServices()
	initStaticFiles()
	RegisterXDWs()
	ContentCreator()
	ContentConsumer()
	ContentUpdator()
	tukdbint.DBConn.Close()
	LogFile.Close()
}
func initVars() {
	log.Println("Base Folder " + os.Getenv(tukcnst.ENV_TUK_CONFIG))
	log.Println("Config file " + os.Getenv(tukcnst.ENV_TUK_CONFIG_FILE) + ".json")
}
func initServices() {
	if registerEventServices {
		log.Println("Processing Event Service Config Files")
		if srvcs, err := tukutil.GetFolderFiles("./config/services/"); err == nil {
			for _, file := range srvcs {
				if strings.HasSuffix(file.Name(), ".json") {
					if filebytes := loadFile(file, "./config/services/"); filebytes != nil {
						srvcs := tukdbint.ServiceStates{Action: tukcnst.DELETE}
						srvc := tukdbint.ServiceState{Name: strings.TrimSuffix(file.Name(), ".json")}
						srvcs.ServiceState = append(srvcs.ServiceState, srvc)
						tukdbint.NewDBEvent(&srvcs)
						srvcs = tukdbint.ServiceStates{Action: tukcnst.INSERT}
						srvc = tukdbint.ServiceState{Name: strings.TrimSuffix(file.Name(), ".json"), Service: string(filebytes)}
						srvcs.ServiceState = append(srvcs.ServiceState, srvc)
						tukdbint.NewDBEvent(&srvcs)
					}
				}
			}
		}
	}
}
func ContentUpdator() {
	if contentUpdator {
		log.Printf("Updating %s Workflow for NHS ID %s", pathway, nhs)
	}
}
func ContentConsumer() {
	if contentConsumer {
		trans := tukxdw.Transaction{
			Actor:   tukcnst.XDW_ACTOR_CONTENT_CONSUMER,
			Pathway: pathway,
			NHS_ID:  nhs,
			User:    user,
			Org:     org,
			Role:    role,
		}
		if err := tukxdw.Execute(&trans); err != nil {
			log.Println(err.Error())
		}
		log.Printf("Consumed Workflow %s, current status %s - Is Overdue %v - Complete by %s - Workflow duration to date %s - Total Events to Date %v", trans.Pathway+trans.NHS_ID, trans.XDWState.Status, trans.XDWState.IsOverdue, trans.XDWState.CompleteBy, trans.XDWState.PrettyWorkflowDuration, trans.XDWEvents.Count)
	}
}
func ContentCreator() {
	if contentCreator {
		trans := tukxdw.Transaction{
			Actor:   tukcnst.XDW_ACTOR_CONTENT_CREATOR,
			Pathway: pathway,
			NHS_ID:  nhs,
			Request: []byte(notes),
			User:    user,
			Org:     org,
			Role:    role,
		}
		if err := tukxdw.Execute(&trans); err != nil {
			log.Println(err.Error())
		}
	}
}
func RegisterXDWs() {
	if registerXdws {
		log.Println("Processing XDW Config Files")
		if xdwconfigs, err := tukutil.GetFolderFiles("./config/xdwconfig/"); err == nil {
			for _, file := range xdwconfigs {
				suffix := ""
				if strings.HasPrefix(file.Name(), pathway) {
					if strings.HasSuffix(file.Name(), pathway+"_def.json") {
						suffix = "_def.json"
					} else {
						if strings.HasSuffix(file.Name(), pathway+"_meta.json") {
							suffix = "_meta.json"
						}
					}
					if suffix == "" {
						log.Printf("Skipping Invalid configuration file. Check name! %s", file.Name())
						continue
					}
					if filebytes := loadFile(file, "./config/xdwconfig/"); filebytes != nil {
						trans := tukxdw.Transaction{
							Actor:            tukcnst.XDW_ADMIN_REGISTER_DEFINITION,
							Pathway:          pathway,
							DSUB_BrokerURL:   DSUB_BROKER_URL,
							DSUB_ConsumerURL: DSUB_CONSUMER_URL_AWS,
							Request:          filebytes,
						}
						if suffix == "_meta.json" {
							trans.Actor = tukcnst.XDW_ADMIN_REGISTER_XDS_META
						}
						tukxdw.Execute(&trans)
					}
				}
			}
		}
	}
}
func initStaticFiles() {
	if initStatic {
		if staticfiles, err := tukutil.GetFolderFiles("./config/static/"); err == nil {
			for _, file := range staticfiles {
				filebytes := loadFile(file, "./config/static/")
				if filebytes != nil {
					log.Printf("Persisting static file %s", file.Name())
					statics := tukdbint.Statics{Action: tukcnst.DELETE}
					static := tukdbint.Static{Name: file.Name()}
					statics.Static = append(statics.Static, static)
					tukdbint.NewDBEvent(&statics)
					statics = tukdbint.Statics{Action: tukcnst.INSERT}
					static = tukdbint.Static{Name: file.Name(), Content: filebytes}
					statics.Static = append(statics.Static, static)
					tukdbint.NewDBEvent(&statics)
				}
			}
		}
	}
}
func initTemplates() {
	if initTmplts {
		if xmlTmplts, err := tukutil.GetFolderFiles("./config/templates/xml/"); err == nil {
			for _, file := range xmlTmplts {
				if strings.HasSuffix(file.Name(), ".xml") {
					filebytes := loadFile(file, "./config/templates/xml/")
					if filebytes != nil {
						log.Printf("Persisting XML Template %s", file.Name())
						tmplts := tukdbint.Templates{Action: tukcnst.DELETE}
						tmplt := tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".xml"), IsXML: true}
						tmplts.Templates = append(tmplts.Templates, tmplt)
						tukdbint.NewDBEvent(&tmplts)
						tmplts = tukdbint.Templates{Action: tukcnst.INSERT}
						tmplt = tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".xml"), IsXML: true, Template: string(filebytes)}
						tmplts.Templates = append(tmplts.Templates, tmplt)
						tukdbint.NewDBEvent(&tmplts)
					}
				}
			}
		}
		if htmlTmplts, err := tukutil.GetFolderFiles("./config/templates/html/"); err == nil {
			for _, file := range htmlTmplts {
				if strings.HasSuffix(file.Name(), ".html") {
					filebytes := loadFile(file, "./config/templates/html/")
					if filebytes != nil {
						log.Printf("Persisting HTML Template %s", file.Name())
						tmplts := tukdbint.Templates{Action: tukcnst.DELETE}
						tmplt := tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".html")}
						tmplts.Templates = append(tmplts.Templates, tmplt)
						tukdbint.NewDBEvent(&tmplts)
						tmplts = tukdbint.Templates{Action: tukcnst.INSERT}
						tmplt = tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".html"), Template: string(filebytes)}
						tmplts.Templates = append(tmplts.Templates, tmplt)
						tukdbint.NewDBEvent(&tmplts)
					}
				}
			}
		}
	}
}
func loadFile(file fs.DirEntry, folder string) []byte {
	var fileBytes []byte
	var err error
	fileBytes, err = os.ReadFile(folder + file.Name())
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("Loaded %s ", file.Name())
	}
	return fileBytes
}
func initLog() {
	if initlogging {
		LogFile = tukutil.CreateLog("./logs")
		log.Println("Loaded log file - " + LogFile.Name())
	}
}
func initDB() {
	if initdb {
		dbconn := tukdbint.TukDBConnection{
			DBUser:     "root",
			DBPassword: "rootPass",
			DBHost:     "tuk.coil1nnpqdlr.eu-west-1.rds.amazonaws.com",
			DBPort:     "3306",
			DBName:     "tuk",
		}
		if err := tukdbint.NewDBEvent(&dbconn); err != nil {
			log.Println(err.Error())
		}
	}
}
