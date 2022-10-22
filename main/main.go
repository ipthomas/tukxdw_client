package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"tukxdw-client/tukcnst"
	"tukxdw-client/tukxdw"

	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukutil"
)

var (
	pathway = "toc"
	//pathways        = "lac,toc,pathalert,radconsult"
	nhs             = "9999999468"
	registerXdws    = true
	initTmplts      = false
	contentCreator  = false
	contentConsumer = false
	contentUpdator  = false
	//TUK_DB_URL_AWS          = "https://5k2o64mwt5.execute-api.eu-west-1.amazonaws.com/beta/"
	DSUB_BROKER_URL         = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritXDSDsub/Dsub"
	DSUB_CONSUMER_URL_AWS   = "https://cjrvrddgdh.execute-api.eu-west-1.amazonaws.com/beta/"
	DSUB_CONSUMER_URL_Local = "http://tukeventserver.ddns.net:8081/eventservice/event"
	user                    = "pcrofts"
	role                    = "CPS"
	org                     = "lth"
	notes                   = "User " + user + " from " + org + " in the role of " + role + " created new workflow for identified looked after child"
	_, b, _, _              = runtime.Caller(0)
	Basepath                = filepath.Dir(b)
	LogFile                 *os.File
)

func main() {
	initLog()
	initDB()
	initTemplates()
	RegisterXdws()
	ContentCreator()
	ContentConsumer()
	ContentUpdator()
	tukdbint.DBConn.Close()
	LogFile.Close()
}
func ContentUpdator() {
	if contentUpdator {
		log.Printf("Updating %s Workflow for NHS ID %s", pathway, nhs)
	}
}
func ContentConsumer() {
	if contentConsumer {
		trans := tukxdw.Transaction{
			Action:  tukcnst.XDW_ACTOR_CONTENT_CONSUMER,
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
			Action:  tukcnst.XDW_ACTOR_CONTENT_CREATOR,
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
		log.Println(string(trans.Response))
	}
}
func initTemplates() {
	if initTmplts {
		if xmlTmplts, err := tukutil.GetFolderFiles("./config/templates/xml/"); err == nil {
			for _, file := range xmlTmplts {
				filebytes := loadXMLFile(file)
				log.Printf("Persisting XML Template %s", file.Name())
				tmplts := tukdbint.Templates{Action: tukcnst.DELETE}
				tmplt := tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".xml"), IsXML: true}
				tukdbint.NewDBEvent(&tmplts)
				tmplts = tukdbint.Templates{Action: tukcnst.INSERT}
				tmplt = tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".xml"), IsXML: true, Template: string(filebytes)}
				tmplts.Templates = append(tmplts.Templates, tmplt)
				tukdbint.NewDBEvent(&tmplts)
			}
		}
		if htmlTmplts, err := tukutil.GetFolderFiles("./config/templates/html/"); err == nil {
			for _, file := range htmlTmplts {
				filebytes := loadHTMLFile(file)
				log.Printf("Persisting HTML Template %s", file.Name())
				tmplts := tukdbint.Templates{Action: tukcnst.DELETE}
				tmplt := tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".html")}
				tukdbint.NewDBEvent(&tmplts)
				tmplts = tukdbint.Templates{Action: tukcnst.INSERT}
				tmplt = tukdbint.Template{Name: strings.TrimSuffix(file.Name(), ".html"), Template: string(filebytes)}
				tmplts.Templates = append(tmplts.Templates, tmplt)
				tukdbint.NewDBEvent(&tmplts)
			}
		}
	}
}
func RegisterXdws() {
	if registerXdws {
		if xdwconfigs, err := tukutil.GetFolderFiles("./config/xdwconfig/"); err == nil {
			for _, file := range xdwconfigs {
				if strings.HasPrefix(file.Name(), pathway) {
					if filebytes, ismeta := loadConfigFile(file); filebytes != nil {
						trans := tukxdw.Transaction{
							Action:           tukcnst.XDW_ADMIN_REGISTER_DEFINITION,
							Pathway:          pathway,
							DSUB_BrokerURL:   DSUB_BROKER_URL,
							DSUB_ConsumerURL: DSUB_CONSUMER_URL_Local,
							Request:          filebytes,
						}
						if ismeta {
							trans.Action = tukcnst.XDW_ADMIN_REGISTER_XDS_META
						}
						tukxdw.Execute(&trans)
					}
				}
			}
		}
	}
}
func loadConfigFile(file fs.DirEntry) ([]byte, bool) {
	log.Printf("Loading file %s", file.Name())
	var fileBytes []byte
	var err error
	if strings.EqualFold(file.Name(), pathway+"_def.json") {
		if fileBytes, err = os.ReadFile("./config/xdwconfig/" + file.Name()); err != nil {
			log.Println(err.Error())
			return fileBytes, true
		}
		log.Printf("Loaded WF Def %s", file.Name())
		return fileBytes, false
	} else {
		if strings.EqualFold(file.Name(), pathway+"_meta.json") {
			if fileBytes, err = os.ReadFile("./config/xdwconfig/" + file.Name()); err != nil {
				log.Println(err.Error())
				return fileBytes, true
			}
			log.Println("Loaded WF XDS Meta for Pathway : " + file.Name())
			return fileBytes, true
		}
	}
	log.Println("Error file is not a valid XDW Definition file or a XDW XDS Meta file.")
	return fileBytes, true
}
func loadXMLFile(file fs.DirEntry) []byte {
	if strings.HasSuffix(file.Name(), ".xml") {
		xmltmplt, err := os.ReadFile("./config/templates/xml/" + file.Name())
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		log.Println("Loaded XML Template : " + file.Name())
		return xmltmplt
	}
	return nil
}
func loadHTMLFile(file fs.DirEntry) []byte {
	if strings.HasSuffix(file.Name(), ".html") {
		xmltmplt, err := os.ReadFile("./config/templates/html/" + file.Name())
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		log.Println("Loaded HTML Template : " + file.Name())
		return xmltmplt
	}
	return nil
}
func initLog() {
	LogFile = tukutil.CreateLog("./logs")
	log.Println("Loaded log file - " + LogFile.Name())
}
func initDB() {
	dbconn := tukdbint.TukDBConnection{
		DBUser:     "root",
		DBPassword: "rootPass",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "tuk",
	}
	if err := tukdbint.NewDBEvent(&dbconn); err != nil {
		log.Println(err.Error())
	}
}
