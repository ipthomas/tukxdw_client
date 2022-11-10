package tukdbint

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukhttp"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB_URL = ""
	DBConn *sql.DB
)

type TukDBConnection struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	DBTimeout     string
	DBReadTimeout string
	DB_URL        string
	DBReader_Only bool
}
type ServiceStates struct {
	Action       string         `json:"action"`
	LastInsertId int64          `json:"lastinsertid"`
	Count        int            `json:"count"`
	ServiceState []ServiceState `json:"servicestate"`
}
type ServiceState struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
}
type Templates struct {
	Action       string     `json:"action"`
	LastInsertId int64      `json:"lastinsertid"`
	Count        int        `json:"count"`
	Templates    []Template `json:"templates"`
}
type Template struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	IsXML    bool   `json:"isxml"`
	Template string `json:"template"`
}
type Subscription struct {
	Id         int64  `json:"id"`
	Created    string `json:"created"`
	BrokerRef  string `json:"brokerref"`
	Pathway    string `json:"pathway"`
	Topic      string `json:"topic"`
	Expression string `json:"expression"`
}
type Subscriptions struct {
	Action        string         `json:"action"`
	LastInsertId  int64          `json:"lastinsertid"`
	Count         int            `json:"count"`
	Subscriptions []Subscription `json:"subscriptions"`
}
type Event struct {
	Id                 int64  `json:"id"`
	Creationtime       string `json:"creationtime"`
	DocName            string `json:"docname"`
	ClassCode          string `json:"classcode"`
	ConfCode           string `json:"confcode"`
	FormatCode         string `json:"formatcode"`
	FacilityCode       string `json:"facilitycode"`
	PracticeCode       string `json:"practicecode"`
	Expression         string `json:"expression"`
	Authors            string `json:"authors"`
	XdsPid             string `json:"xdspid"`
	XdsDocEntryUid     string `json:"xdsdocentryuid"`
	RepositoryUniqueId string `json:"repositoryuniqueid"`
	NhsId              string `json:"nhsid"`
	User               string `json:"user"`
	Org                string `json:"org"`
	Role               string `json:"role"`
	Topic              string `json:"topic"`
	Pathway            string `json:"pathway"`
	Comments           string `json:"comments"`
	Version            int    `json:"ver"`
	TaskId             int    `json:"taskid"`
	BrokerRef          string `json:"brokerref"`
}
type Events struct {
	Action       string  `json:"action"`
	LastInsertId int64   `json:"lastinsertid"`
	Count        int     `json:"count"`
	Events       []Event `json:"events"`
}
type Workflow struct {
	Id        int64  `json:"id"`
	Created   string `json:"created"`
	Pathway   string `json:"pathway"`
	NHSId     string `json:"nhsid"`
	XDW_Key   string `json:"xdw_key"`
	XDW_UID   string `json:"xdw_uid"`
	XDW_Doc   string `json:"xdw_doc"`
	XDW_Def   string `json:"xdw_def"`
	Version   int    `json:"version"`
	Published bool   `json:"published"`
	Status    string `json:"status"`
}
type Workflows struct {
	Action       string     `json:"action"`
	LastInsertId int64      `json:"lastinsertid"`
	Count        int        `json:"count"`
	Workflows    []Workflow `json:"workflows"`
}
type XDWS struct {
	Action       string `json:"action"`
	LastInsertId int64  `json:"lastinsertid"`
	Count        int    `json:"count"`
	XDW          []XDW  `json:"xdws"`
}
type XDW struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	IsXDSMeta bool   `json:"isxdsmeta"`
	XDW       string `json:"xdw"`
}
type IdMaps struct {
	Action       string
	LastInsertId int64
	Where        string
	Value        string
	Cnt          int
	LidMap       []IdMap
}
type IdMap struct {
	Id  int64  `json:"id"`
	Lid string `json:"lid"`
	Mid string `json:"mid"`
}
type EventAcks struct {
	Action       string
	LastInsertId int64
	Where        string
	Value        string
	Cnt          int
	EventAck     []EventAck
}
type EventAck struct {
	Id           int64  `json:"id"`
	CreationTime string `json:"creationtime"`
	EventID      int64  `json:"eventid"`
	User         string `json:"user"`
	Org          string `json:"org"`
	Role         string `json:"role"`
}

// sort interface for events
type EventsList []Event

func (e EventsList) Len() int {
	return len(e)
}
func (e EventsList) Less(i, j int) bool {
	return e[i].Id > e[j].Id
}
func (e EventsList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type TUK_DB_Interface interface {
	newEvent() error
}

// NewDBEvent takes an Interface (struct Events, Workflows, Subscriptions, XDWS, TukDBConnection) and executes a mysql request. The response is poputlated into the interface struct.
// If DB_URL = "" the sql query is executed against the DBConn established using a DSN, if a value is present in DB_URL, the query is sent to the API Gateway URL (DB_URL) as a json string of the provided interface struct
func NewDBEvent(i TUK_DB_Interface) error {
	return i.newEvent()
}

// functions for mysql DB access via DSN

func (i *TukDBConnection) newEvent() error {
	var err error
	if i.DB_URL != "" {
		log.Println("Database API URL provided. Will connect to mysql instance via AWS API Gateway url " + i.DB_URL)
		DB_URL = i.DB_URL
	} else {
		if i.DBUser == "" {
			i.DBUser = "root"
		}
		if i.DBPassword == "" {
			i.DBPassword = "rootPass"
		}
		if i.DBHost == "" {
			i.DBHost = "localhost"
		}
		if i.DBPort == "" {
			i.DBPort = ":3306"
		} else {
			if !strings.HasPrefix(i.DBPort, ":") {
				i.DBPort = ":" + i.DBPort
			}
		}
		if i.DBName == "" {
			i.DBName = "tuk"
		}
		if i.DBTimeout == "" {
			i.DBTimeout = "60s"
		} else {
			if !strings.HasSuffix(i.DBTimeout, "s") {
				i.DBTimeout = i.DBTimeout + "s"
			}
		}
		if i.DBReadTimeout == "" {
			i.DBReadTimeout = "2s"
		} else {
			if !strings.HasSuffix(i.DBReadTimeout, "s") {
				i.DBReadTimeout = i.DBReadTimeout + "s"
			}
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=%s&readTimeout=%s",
			i.DBUser,
			i.DBPassword,
			i.DBHost+i.DBPort,
			i.DBName,
			i.DBTimeout,
			i.DBReadTimeout)
		log.Println("No Database API URL provided. Opening DB Connection to mysql instance via DSN - " + dsn)
		DBConn, err = sql.Open(tukcnst.MYSQL, dsn)
	}

	return err
}
func GetSubscriptions(brokerref string, pathway string, expression string) Subscriptions {
	subs := Subscriptions{Action: tukcnst.SELECT}
	sub := Subscription{BrokerRef: brokerref, Pathway: pathway, Expression: expression}
	subs.Subscriptions = append(subs.Subscriptions, sub)
	subs.newEvent()
	return subs
}
func (i *Subscriptions) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_SUBSCRIPTIONS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Subscriptions) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.SUBSCRIPTIONS, reflectStruct(reflect.ValueOf(i.Subscriptions[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		for rows.Next() {
			sub := Subscription{}
			if err := rows.Scan(&sub.Id, &sub.Created, &sub.BrokerRef, &sub.Pathway, &sub.Topic, &sub.Expression); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Subscriptions = append(i.Subscriptions, sub)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func GetEvents(user string, pathway string, nhsid string, expression string, taskid int, version int) Events {
	events := Events{Action: tukcnst.SELECT}
	event := Event{User: user, Pathway: pathway, NhsId: nhsid, Expression: expression, TaskId: taskid, Version: version}
	events.Events = append(events.Events, event)
	events.newEvent()
	return events
}
func (i *Events) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_EVENTS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Events) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.EVENTS, reflectStruct(reflect.ValueOf(i.Events[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		for rows.Next() {
			ev := Event{}
			if err := rows.Scan(&ev.Id, &ev.Creationtime, &ev.DocName, &ev.ClassCode, &ev.ConfCode, &ev.FormatCode, &ev.FacilityCode, &ev.PracticeCode, &ev.Expression, &ev.Authors, &ev.XdsPid, &ev.XdsDocEntryUid, &ev.RepositoryUniqueId, &ev.NhsId, &ev.User, &ev.Org, &ev.Role, &ev.Topic, &ev.Pathway, &ev.Comments, &ev.Version, &ev.TaskId); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Events = append(i.Events, ev)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func GetAllWorkflows() Workflows {
	wfs := Workflows{Action: tukcnst.SELECT}
	wfs.newEvent()
	return wfs
}
func GetPathwayWorkflows(pathway string) Workflows {
	wfs := Workflows{Action: tukcnst.SELECT}
	wf := Workflow{Pathway: pathway}
	wfs.Workflows = append(wfs.Workflows, wf)
	wfs.newEvent()
	return wfs
}
func GetActiveWorkflowNames() []string {
	var activewfs []string
	wfs := GetWorkflows("", "", "", "", 0, false, "")
	log.Printf("Active Workflow Count %v", wfs.Count)
	for _, v := range wfs.Workflows {
		if v.Id != 0 {
			activewfs = append(activewfs, v.Pathway)
		}
	}
	log.Printf("Set %v Active Pathways - %s", len(activewfs), activewfs)
	return activewfs
}
func GetWorkflows(pathway string, nhsid string, xdwkey string, xdwuid string, version int, published bool, status string) Workflows {
	wfs := Workflows{Action: tukcnst.SELECT}
	wf := Workflow{Pathway: pathway, NHSId: nhsid, XDW_Key: xdwkey, XDW_UID: xdwuid, Version: version, Published: published, Status: status}
	wfs.Workflows = append(wfs.Workflows, wf)
	wfs.newEvent()
	return wfs
}
func (i *Workflows) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_WORKFLOWS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Workflows) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.WORKFLOWS, reflectStruct(reflect.ValueOf(i.Workflows[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			workflow := Workflow{}
			if err := rows.Scan(&workflow.Id, &workflow.Pathway, &workflow.NHSId, &workflow.Created, &workflow.XDW_Key, &workflow.XDW_UID, &workflow.XDW_Doc, &workflow.XDW_Def, &workflow.Version, &workflow.Published, &workflow.Status); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Workflows = append(i.Workflows, workflow)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func GetWorkflowDefinitionNames() []string {
	var xdwdefs []string
	xdws := XDWS{Action: tukcnst.SELECT}
	xdw := XDW{IsXDSMeta: false}
	xdws.XDW = append(xdws.XDW, xdw)
	if err := xdws.newEvent(); err == nil {
		for _, xdw := range xdws.XDW {
			if xdw.Id > 0 {
				xdwdefs = append(xdwdefs, xdw.Name)
			}
		}
	}
	log.Printf("Returning %v XDW Config files", len(xdwdefs))
	return xdwdefs
}
func GetWorkflowXDSMetaNames() []string {
	var xdwdefs []string
	xdws := XDWS{Action: tukcnst.SELECT}
	xdw := XDW{IsXDSMeta: true}
	xdws.XDW = append(xdws.XDW, xdw)
	if err := xdws.newEvent(); err == nil {
		for _, xdw := range xdws.XDW {
			if xdw.Id > 0 {
				xdwdefs = append(xdwdefs, xdw.Name)
			}
		}
	}
	log.Printf("Returning %v XDS Meta files", len(xdwdefs))
	return xdwdefs
}
func GetWorkflowDefinitions(name string) (XDWS, error) {
	xdws := XDWS{Action: tukcnst.SELECT}
	err := xdws.newEvent()
	return xdws, err
}
func GetWorkflowDefinition(name string) (XDW, error) {
	var err error
	xdws := XDWS{Action: tukcnst.SELECT}
	xdw := XDW{Name: name}
	xdws.XDW = append(xdws.XDW, xdw)
	if err = xdws.newEvent(); err == nil && xdws.Count == 1 {
		return xdws.XDW[1], nil
	}
	return xdw, err
}
func GetWorkflowXDSMeta(name string) (XDW, error) {
	var err error
	xdws := XDWS{Action: tukcnst.SELECT}
	xdw := XDW{Name: name, IsXDSMeta: true}
	xdws.XDW = append(xdws.XDW, xdw)
	if err = xdws.newEvent(); err == nil && xdws.Count == 1 {
		return xdws.XDW[1], nil
	}
	return xdw, err
}
func SetWorkflowDefinition(name string, config string, isxdsmeta bool) error {
	xdws := XDWS{Action: tukcnst.DELETE}
	xdw := XDW{Name: name, IsXDSMeta: isxdsmeta}
	xdws.XDW = append(xdws.XDW, xdw)
	xdws.newEvent()
	xdws = XDWS{Action: tukcnst.INSERT}
	xdw = XDW{Name: name, IsXDSMeta: isxdsmeta, XDW: config}
	xdws.XDW = append(xdws.XDW, xdw)
	return xdws.newEvent()
}
func (i *XDWS) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_XDWS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.XDW) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.XDWS, reflectStruct(reflect.ValueOf(i.XDW[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			xdw := XDW{}
			if err := rows.Scan(&xdw.Id, &xdw.Name, &xdw.IsXDSMeta, &xdw.XDW); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.XDW = append(i.XDW, xdw)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func GetTemplate(templatename string, isxml bool) (Template, error) {
	var err error
	tmplts := Templates{Action: tukcnst.SELECT}
	tmplt := Template{Name: templatename, IsXML: isxml}
	tmplts.Templates = append(tmplts.Templates, tmplt)
	if err = tmplts.newEvent(); err == nil && tmplts.Count == 1 {
		return tmplts.Templates[1], nil
	}
	return tmplt, err
}
func SetTemplate(templatename string, isxml bool, templatestr string) error {
	tmplts := Templates{Action: tukcnst.DELETE}
	tmplt := Template{Name: templatename, IsXML: isxml}
	tmplts.Templates = append(tmplts.Templates, tmplt)
	tmplts.newEvent()
	tmplts = Templates{Action: tukcnst.INSERT}
	tmplt = Template{Name: templatename, IsXML: isxml, Template: templatestr}
	tmplts.Templates = append(tmplts.Templates, tmplt)
	return tmplts.newEvent()
}
func (i *Templates) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_TEMPLATES
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Templates) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.TEMPLATES, reflectStruct(reflect.ValueOf(i.Templates[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			tmplt := Template{}
			if err := rows.Scan(&tmplt.Id, &tmplt.Name, &tmplt.IsXML, &tmplt.Template); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Templates = append(i.Templates, tmplt)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func GetIDMaps() IdMaps {
	idmaps := IdMaps{Action: tukcnst.SELECT}
	idmap := IdMap{}
	idmaps.LidMap = append(idmaps.LidMap, idmap)
	if err := idmaps.newEvent(); err != nil {
		log.Println(err.Error())
	}
	return idmaps
}
func GetIDMapsMappedId(localid string) string {
	idmaps := IdMaps{Action: tukcnst.SELECT}
	idmap := IdMap{Lid: localid}
	idmaps.LidMap = append(idmaps.LidMap, idmap)
	if err := idmaps.newEvent(); err != nil {
		log.Println(err.Error())
		return localid
	}
	if idmaps.Cnt == 1 {
		return idmaps.LidMap[1].Mid
	}
	return localid
}
func GetIDMapsLocalId(mid string) string {
	idmaps := IdMaps{Action: tukcnst.SELECT}
	idmap := IdMap{Mid: mid}
	idmaps.LidMap = append(idmaps.LidMap, idmap)
	if err := idmaps.newEvent(); err != nil {
		log.Println(err.Error())
		return mid
	}
	if idmaps.Cnt == 1 {
		return idmaps.LidMap[1].Lid
	}
	return mid
}
func (i *IdMaps) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_IDMAPS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.LidMap) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.ID_MAPS, reflectStruct(reflect.ValueOf(i.LidMap[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			idmap := IdMap{}
			if err := rows.Scan(&idmap.Id, &idmap.Lid, &idmap.Mid); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.LidMap = append(i.LidMap, idmap)
			i.Cnt = i.Cnt + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func GetServiceState(servicename string) (ServiceState, error) {
	var err error
	if !strings.HasSuffix(servicename, "srvc") {
		servicename = servicename + "srvc"
	}
	srvcs := ServiceStates{Action: tukcnst.SELECT}
	srvc := ServiceState{Name: servicename}
	srvcs.ServiceState = append(srvcs.ServiceState, srvc)
	err = NewDBEvent(&srvcs)
	if err == nil && srvcs.Count == 1 {
		return srvcs.ServiceState[1], nil
	}
	return srvc, err
}
func SetServiceState(servicename string, state string) error {
	srvcs := ServiceStates{Action: tukcnst.DELETE}
	srvc := ServiceState{Name: servicename}
	srvcs.ServiceState = append(srvcs.ServiceState, srvc)
	srvcs.newEvent()
	srvcs = ServiceStates{Action: tukcnst.INSERT}
	srvc = ServiceState{Name: servicename, Service: state}
	srvcs.ServiceState = append(srvcs.ServiceState, srvc)
	return srvcs.newEvent()
}
func (i *ServiceStates) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_SERVICESTATES
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.ServiceState) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.SERVICE_STATES, reflectStruct(reflect.ValueOf(i.ServiceState[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			srvc := ServiceState{}
			if err := rows.Scan(&srvc.Id, &srvc.Name, &srvc.Service); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.ServiceState = append(i.ServiceState, srvc)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func HasEventAck(eventid int64) bool {
	evacks := EventAcks{Action: tukcnst.SELECT}
	evack := EventAck{EventID: eventid}
	evacks.EventAck = append(evacks.EventAck, evack)
	if err := evacks.newEvent(); err != nil {
		log.Println(err.Error())
	}
	return evacks.Cnt > 0
}
func GetTaskNotes(pwy string, nhsid string, taskid int, ver int) string {
	notes := ""
	evs := Events{Action: tukcnst.SELECT}
	ev := Event{Pathway: pwy, NhsId: nhsid, TaskId: taskid, Version: ver}
	evs.Events = append(evs.Events, ev)
	err := NewDBEvent(&evs)
	if err == nil && evs.Count > 0 {
		for _, note := range evs.Events {
			if note.Id != 0 {
				notes = notes + note.Comments + "\n"
			}
		}
		log.Printf("Found TaskId %v Notes %s", taskid, notes)
	}
	return notes
}
func (i *EventAcks) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_EVENT_ACKS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.EventAck) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.EVENT_ACKS, reflectStruct(reflect.ValueOf(i.EventAck[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			eventack := EventAck{}
			if err := rows.Scan(&eventack.Id, &eventack.CreationTime, &eventack.EventID, &eventack.User, &eventack.Org, &eventack.Role); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.EventAck = append(i.EventAck, eventack)
			i.Cnt = i.Cnt + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func reflectStruct(i reflect.Value) map[string]interface{} {
	params := make(map[string]interface{})
	structType := i.Type()
	for f := 0; f < i.NumField(); f++ {
		if structType.Field(f).Name == "Id" || structType.Field(f).Name == "EventID" || structType.Field(f).Name == "LastInsertId" {
			tint64 := i.Field(f).Interface().(int64)
			if tint64 > 0 {
				params[strings.ToLower(structType.Field(f).Name)] = tint64
			}
		} else {
			if structType.Field(f).Name == "Version" || structType.Field(f).Name == "TaskId" {
				tint := i.Field(f).Interface().(int)
				if tint != -1 {
					params[strings.ToLower(structType.Field(f).Name)] = tint
				}
			} else {
				if i.Field(f).Interface() != "" {
					params[strings.ToLower(structType.Field(f).Name)] = i.Field(f).Interface()
				}
			}
		}
	}
	return params
}
func createPreparedStmnt(action string, table string, params map[string]interface{}) (string, []interface{}, error) {
	var vals []interface{}
	stmntStr := "SELECT * FROM " + table
	if len(params) > 0 {
		switch action {
		case tukcnst.SELECT:
			var paramStr string
			stmntStr = stmntStr + " WHERE "
			for param, val := range params {
				paramStr = paramStr + param + "= ? AND "
				vals = append(vals, val)
			}
			paramStr = strings.TrimSuffix(paramStr, " AND ")
			stmntStr = stmntStr + paramStr
		case tukcnst.INSERT:
			var paramStr string
			var qStr string
			stmntStr = "INSERT INTO " + table + " ("
			for param, val := range params {
				paramStr = paramStr + param + ", "
				qStr = qStr + "?, "
				vals = append(vals, val)
			}
			paramStr = strings.TrimSuffix(paramStr, ", ") + ") VALUES ("
			qStr = strings.TrimSuffix(qStr, ", ")
			stmntStr = stmntStr + paramStr + qStr + ")"
		case tukcnst.DEPRECATE:
			switch table {
			case tukcnst.WORKFLOWS:
				stmntStr = "UPDATE workflows SET version = version + 1 WHERE xdw_key=?"
				vals = append(vals, params["xdw_key"])
			case tukcnst.EVENTS:
				stmntStr = "UPDATE events SET version = version + 1 WHERE pathway=? AND nhsid=?"
				vals = append(vals, params["pathway"])
				vals = append(vals, params["nhsid"])
			}
		case tukcnst.UPDATE:
			switch table {
			case tukcnst.WORKFLOWS:
				stmntStr = "UPDATE workflows SET xdw_doc = ?, published = ?, status = ? WHERE pathway = ? AND nhsid = ? AND version = ?"
				vals = append(vals, params["xdw_doc"])
				vals = append(vals, params["published"])
				vals = append(vals, params["status"])
				vals = append(vals, params["pathway"])
				vals = append(vals, params["nhsid"])
				vals = append(vals, params["version"])
			}
		case tukcnst.DELETE:
			stmntStr = "DELETE FROM " + table + " WHERE "
			var paramStr string
			for param, val := range params {
				paramStr = paramStr + param + "= ? AND "
				vals = append(vals, val)
			}
			paramStr = strings.TrimSuffix(paramStr, " AND ")
			stmntStr = stmntStr + paramStr
		}
		log.Printf("Created Prepared Statement %s - Values %s", stmntStr, vals)
	}
	return stmntStr, vals, nil
}
func setRows(ctx context.Context, sqlStmnt *sql.Stmt, vals []interface{}) (*sql.Rows, error) {
	if len(vals) > 0 {
		return sqlStmnt.QueryContext(ctx, vals...)
	} else {
		return sqlStmnt.QueryContext(ctx)
	}
}
func setLastID(ctx context.Context, sqlStmnt *sql.Stmt, vals []interface{}) (int64, error) {
	if len(vals) > 0 {
		sqlrslt, err := sqlStmnt.ExecContext(ctx, vals...)
		if err != nil {
			log.Println(err.Error())
			return 0, err
		}
		id, err := sqlrslt.LastInsertId()
		if err != nil {
			log.Println(err.Error())
			return 0, err
		} else {
			return id, nil
		}
	}
	return 0, nil
}

// functions for AWS Aurora DB Access via AWS API GW URL
func (i *ServiceStates) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.SERVICE_STATES, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *Subscriptions) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.SUBSCRIPTIONS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *Events) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.EVENTS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *Workflows) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.WORKFLOWS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *XDWS) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.XDWS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *IdMaps) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.ID_MAPS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *EventAcks) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.EVENT_ACKS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *Templates) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.TEMPLATES, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func aws_APIRequest(action string, resource string, body []byte) tukhttp.AWS_APIRequest {
	return tukhttp.AWS_APIRequest{
		URL:      DB_URL,
		Act:      action,
		Resource: resource,
		Timeout:  5,
		Body:     body,
	}
}
