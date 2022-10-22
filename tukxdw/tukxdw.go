package tukxdw

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"tukxdw-client/tukcnst"

	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukdsub"
	"github.com/ipthomas/tukutil"
)

type Interface interface {
	execute() error
}
type Transaction struct {
	Action           string
	User             string
	Org              string
	Role             string
	Pathway          string
	Expression       string
	NHS_ID           string
	Task_ID          int
	XDWVersion       int
	DSUB_BrokerURL   string
	DSUB_ConsumerURL string
	Request          []byte
	Response         []byte
	Dashboard        Dashboard
	XDWDefinition    WorkflowDefinition
	XDSDocumentMeta  XDSDocumentMeta
	XDWDocument      XDWWorkflowDocument
	XDWState         XDWState
	Workflows        tukdbint.Workflows
	XDWEvents        tukdbint.Events
	XDWTaskStates    []XDWTaskState
}
type XDWTaskState struct {
	TaskID              int
	Created             string
	CompleteBy          string
	Status              string
	IsOverdue           bool
	LatestTaskEventTime time.Time
	TaskDuration        time.Duration
	PrettyTaskDuration  string
}
type Dashboard struct {
	Total        int
	InProgress   int
	TargetMet    int
	TargetMissed int
	Escalated    int
	Complete     int
}
type XDWState struct {
	Created                 string
	CompleteBy              string
	Status                  string
	IsPublished             bool
	IsOverdue               bool
	LatestWorkflowEventTime time.Time
	LatestTaskEventTime     time.Time
	WorkflowDuration        time.Duration
	PrettyWorkflowDuration  string
}
type XDSDocumentMeta struct {
	ID                    string `json:"id"`
	Repositoryuniqueid    string `json:"repositoryuniqueid"`
	Registryoid           string `json:"registryoid"`
	Languagecode          string `json:"languagecode"`
	Docname               string `json:"docname"`
	Docdesc               string `json:"docdesc"`
	DocID                 string `json:"docid"`
	Authorinstitution     string `json:"authorinstitution"`
	Authorperson          string `json:"authorperson"`
	Classcode             string `json:"classcode"`
	Classcodescheme       string `json:"classcodescheme"`
	Classcodevalue        string `json:"classcodevalue"`
	Typecode              string `json:"typecode"`
	Typecodescheme        string `json:"typecodescheme"`
	Typecodevalue         string `json:"typecodevalue"`
	Practicesettingcode   string `json:"practicesettingcode"`
	Practicesettingscheme string `json:"practicesettingscheme"`
	Practicesettingvalue  string `json:"practicesettingvalue"`
	Confcode              string `json:"confcode"`
	Confcodescheme        string `json:"confcodescheme"`
	Confcodevalue         string `json:"confcodevalue"`
	Facilitycode          string `json:"facilitycode"`
	Facilitycodescheme    string `json:"facilitycodescheme"`
	Facilitycodevalue     string `json:"facilitycodevalue"`
	Formatcode            string `json:"formatcode"`
	Formatcodescheme      string `json:"formatcodescheme"`
	Formatcodevalue       string `json:"formatcodevalue"`
	Mimetype              string `json:"mimetype"`
	Objecttype            string `json:"objecttype"`
}
type WorkflowDefinition struct {
	Ref                 string `json:"ref"`
	Name                string `json:"name"`
	Confidentialitycode string `json:"confidentialitycode"`
	StartByTime         string `json:"startbytime"`
	CompleteByTime      string `json:"completebytime"`
	CompletionBehavior  []struct {
		Completion struct {
			Condition string `json:"condition"`
		} `json:"completion"`
	} `json:"completionBehavior"`
	Tasks []struct {
		ID              string `json:"id"`
		Tasktype        string `json:"tasktype"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		ActualOwner     string `json:"actualowner"`
		ExpirationTime  string `json:"expirationtime"`
		StartByTime     string `json:"startbytime"`
		CompleteByTime  string `json:"completebytime"`
		IsSkipable      bool   `json:"isskipable"`
		PotentialOwners []struct {
			OrganizationalEntity struct {
				User string `json:"user"`
			} `json:"organizationalEntity"`
		} `json:"potentialOwners"`
		CompletionBehavior []struct {
			Completion struct {
				Condition string `json:"condition"`
			} `json:"completion"`
		} `json:"completionBehavior"`
		Input []struct {
			Name        string `json:"name"`
			Contenttype string `json:"contenttype"`
			AccessType  string `json:"accesstype"`
		} `json:"input,omitempty"`
		Output []struct {
			Name        string `json:"name"`
			Contenttype string `json:"contenttype"`
			AccessType  string `json:"accesstype"`
		} `json:"output,omitempty"`
	} `json:"tasks"`
}

//<xsd:simpleType name="tPredefinedStatus"> for WorkflowStatus
//     <xsd:enumeration value="CREATED" /> Workflow Created

//     <xsd:enumeration value="READY" /> First Task assigned owner

//     <xsd:enumeration value="RESERVED" /> not used

//     <xsd:enumeration value="IN_PROGRESS" /> First Task started (or completed)

//     <xsd:enumeration value="SUSPENDED" /> not used

//     <xsd:enumeration value="COMPLETED" /> Workflow Complete

//     <xsd:enumeration value="FAILED" /> not used

//     <xsd:enumeration value="ERROR" /> not used

//     <xsd:enumeration value="EXITED" /> not used

//     <xsd:enumeration value="OBSOLETE" /> Deprecated Workflows that are not complete

type XDWWorkflowDocument struct {
	XMLName                        xml.Name              `xml:"XDW.WorkflowDocument"`
	Hl7                            string                `xml:"hl7,attr"`
	WsHt                           string                `xml:"ws-ht,attr"`
	Xdw                            string                `xml:"xdw,attr"`
	Xsi                            string                `xml:"xsi,attr"`
	SchemaLocation                 string                `xml:"schemaLocation,attr"`
	ID                             ID                    `xml:"id"`
	EffectiveTime                  EffectiveTime         `xml:"effectiveTime"`
	ConfidentialityCode            ConfidentialityCode   `xml:"confidentialityCode"`
	Patient                        PatientID             `xml:"patient"`
	Author                         Author                `xml:"author"`
	WorkflowInstanceId             string                `xml:"workflowInstanceId"`
	WorkflowDocumentSequenceNumber string                `xml:"workflowDocumentSequenceNumber"`
	WorkflowStatus                 string                `xml:"workflowStatus"`
	WorkflowStatusHistory          WorkflowStatusHistory `xml:"workflowStatusHistory"`
	WorkflowDefinitionReference    string                `xml:"workflowDefinitionReference"`
	TaskList                       TaskList              `xml:"TaskList"`
}
type ConfidentialityCode struct {
	Code string `xml:"code,attr"`
}
type EffectiveTime struct {
	Value string `xml:"value,attr"`
}
type PatientID struct {
	ID ID `xml:"id"`
}
type Author struct {
	AssignedAuthor AssignedAuthor `xml:"assignedAuthor"`
}
type AssignedAuthor struct {
	ID             ID             `xml:"id"`
	AssignedPerson AssignedPerson `xml:"assignedPerson"`
}
type ID struct {
	Root                   string `xml:"root,attr"`
	Extension              string `xml:"extension,attr"`
	AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
}
type AssignedPerson struct {
	Name Name `xml:"name"`
}
type Name struct {
	Family string `xml:"family"`
	Prefix string `xml:"prefix"`
}
type WorkflowStatusHistory struct {
	DocumentEvent []DocumentEvent `xml:"documentEvent"`
}
type TaskList struct {
	XDWTask []XDWTask `xml:"XDWTask"`
}
type XDWTask struct {
	TaskData         TaskData         `xml:"taskData"`
	TaskEventHistory TaskEventHistory `xml:"taskEventHistory"`
}
type TaskData struct {
	TaskDetails TaskDetails `xml:"taskDetails"`
	Description string      `xml:"description"`
	Input       []Input     `xml:"input"`
	Output      []Output    `xml:"output"`
}
type TaskDetails struct {
	ID                    string `xml:"id"`
	TaskType              string `xml:"taskType"`
	Name                  string `xml:"name"`
	Status                string `xml:"status"`
	ActualOwner           string `xml:"actualOwner"`
	CreatedTime           string `xml:"createdTime"`
	CreatedBy             string `xml:"createdBy"`
	ActivationTime        string `xml:"activationtime"`
	LastModifiedTime      string `xml:"lastModifiedTime"`
	RenderingMethodExists string `xml:"renderingMethodExists"`
}
type TaskEventHistory struct {
	TaskEvent []TaskEvent `xml:"taskEvent"`
}
type AttachmentInfo struct {
	Identifier      string `xml:"identifier"`
	Name            string `xml:"name"`
	AccessType      string `xml:"accessType"`
	ContentType     string `xml:"contentType"`
	ContentCategory string `xml:"contentCategory"`
	AttachedTime    string `xml:"attachedTime"`
	AttachedBy      string `xml:"attachedBy"`
	HomeCommunityId string `xml:"homeCommunityId"`
}
type Part struct {
	Name           string         `xml:"name,attr"`
	AttachmentInfo AttachmentInfo `xml:"attachmentInfo"`
}
type Output struct {
	Part Part `xml:"part"`
}
type Input struct {
	Part Part `xml:"part"`
}
type DocumentEvent struct {
	EventTime           string `xml:"eventTime"`
	EventType           string `xml:"eventType"`
	TaskEventIdentifier string `xml:"taskEventIdentifier"`
	Author              string `xml:"author"`
	PreviousStatus      string `xml:"previousStatus"`
	ActualStatus        string `xml:"actualStatus"`
}
type TaskEvent struct {
	ID         string `xml:"id"`
	EventTime  string `xml:"eventTime"`
	Identifier string `xml:"identifier"`
	EventType  string `xml:"eventType"`
	Status     string `xml:"status"`
}

func Execute(i Interface) error {
	return i.execute()
}

// IHE XDW Actors

// IHE XDW Content Creator
func (i *Transaction) execute() error {
	switch i.Action {
	case tukcnst.XDW_ADMIN_REGISTER_DEFINITION:
		return i.registerWorkflowDefinition(false)
	case tukcnst.XDW_ADMIN_REGISTER_XDS_META:
		return i.registerWorkflowDefinition(true)
	case tukcnst.XDW_ACTOR_CONTENT_CREATOR:
		return i.contentCreator()
	case tukcnst.XDW_ACTOR_CONTENT_CONSUMER:
		return i.contentConsumer()
	case tukcnst.XDW_ACTOR_CONTENT_UPDATER:
		return i.contentUpdater()
	}
	return nil
}
func (i *Transaction) contentUpdater() error {
	return nil
}

// IHE XDW Content Creator
func (i *Transaction) contentCreator() error {
	log.Printf("Creating New Workflow for %s", i.Pathway+i.NHS_ID)
	var err error
	if err = i.loadWorkflowConfig(); err == nil {
		if err = i.deprecateWorkflow(); err == nil {
			i.createWorkflow()
			i.persistWorkflow()
		}
	}
	return err
}
func (i *Transaction) loadWorkflowConfig() error {
	log.Printf("Obtaining XDS Meta for Pathway %s", i.Pathway)
	var err error
	xdwmeta := tukdbint.XDW{Name: i.Pathway + "_meta", IsXDSMeta: true}
	xdwsmeta := tukdbint.XDWS{Action: tukcnst.SELECT}
	xdwsmeta.XDW = append(xdwsmeta.XDW, xdwmeta)
	if err = tukdbint.NewDBEvent(&xdwsmeta); err == nil {
		if xdwsmeta.Count == 1 {
			log.Printf("Loaded XDS Meta for Pathway %s", i.Pathway)
			if err = json.Unmarshal([]byte(xdwsmeta.XDW[1].XDW), &i.XDSDocumentMeta); err == nil {
				xdwdef := tukdbint.XDW{Name: i.Pathway, IsXDSMeta: false}
				xdwsdef := tukdbint.XDWS{Action: tukcnst.SELECT}
				xdwsdef.XDW = append(xdwsdef.XDW, xdwdef)
				if err = tukdbint.NewDBEvent(&xdwsdef); err == nil {
					if xdwsdef.Count == 1 {
						log.Printf("Loaded XDW definition for Pathway %s", i.Pathway)
						err = json.Unmarshal([]byte(xdwsdef.XDW[1].XDW), &i.XDWDefinition)
					}
				} else {
					err = errors.New("no xdw definition config found")
				}
			}
		} else {
			err = errors.New("no xdw meta config found")
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func (i *Transaction) deprecateWorkflow() error {
	log.Printf("Deprecating any current %s Workflow", i.Pathway+i.NHS_ID)
	var err error
	wfs := tukdbint.Workflows{Action: tukcnst.DEPRECATE}
	wf := tukdbint.Workflow{XDW_Key: i.Pathway + i.NHS_ID}
	wfs.Workflows = append(wfs.Workflows, wf)
	if err = tukdbint.NewDBEvent(&wfs); err == nil {
		log.Printf("Deprecating any current %s Workflow events", i.Pathway+i.NHS_ID)
		evs := tukdbint.Events{Action: tukcnst.DEPRECATE}
		ev := tukdbint.Event{Pathway: i.Pathway, NhsId: i.NHS_ID}
		evs.Events = append(evs.Events, ev)
		if err = tukdbint.NewDBEvent(&evs); err != nil {
			log.Println(err.Error())
		}
	}
	return err
}
func (i *Transaction) createWorkflow() {
	var authoid = getLocalId(i.Org)
	var patoid = tukcnst.NHS_OID_DEFAULT
	var wfid = tukutil.Newid()
	var effectiveTime = tukutil.Time_Now()
	i.XDWDocument.Xdw = tukcnst.XDWNameSpace
	i.XDWDocument.Hl7 = tukcnst.HL7NameSpace
	i.XDWDocument.WsHt = tukcnst.WHTNameSpace
	i.XDWDocument.Xsi = tukcnst.XMLNS_XSI
	i.XDWDocument.XMLName.Local = tukcnst.XDWNameLocal
	i.XDWDocument.SchemaLocation = tukcnst.WorkflowDocumentSchemaLocation
	i.XDWDocument.ID.Root = strings.ReplaceAll(tukcnst.WorkflowInstanceId, "^", "")
	i.XDWDocument.ID.Extension = wfid
	i.XDWDocument.ID.AssigningAuthorityName = getMappedId("LCS")
	i.XDWDocument.EffectiveTime.Value = effectiveTime
	i.XDWDocument.ConfidentialityCode.Code = i.XDWDefinition.Confidentialitycode
	i.XDWDocument.Patient.ID.Root = patoid
	i.XDWDocument.Patient.ID.Extension = i.NHS_ID
	i.XDWDocument.Patient.ID.AssigningAuthorityName = "NHS"
	i.XDWDocument.Author.AssignedAuthor.ID.Root = authoid
	i.XDWDocument.Author.AssignedAuthor.ID.Extension = strings.ToUpper(i.Org)
	i.XDWDocument.Author.AssignedAuthor.ID.AssigningAuthorityName = getMappedId(i.Org)
	i.XDWDocument.Author.AssignedAuthor.AssignedPerson.Name.Family = i.User
	i.XDWDocument.Author.AssignedAuthor.AssignedPerson.Name.Prefix = i.Role
	i.XDWDocument.WorkflowInstanceId = wfid + tukcnst.WorkflowInstanceId
	i.XDWDocument.WorkflowDocumentSequenceNumber = "1"
	i.XDWDocument.WorkflowStatus = tukcnst.OPEN
	i.XDWDocument.WorkflowDefinitionReference = strings.ToUpper(i.Pathway)
	tevidstr := tukutil.GetStringFromInt(int(i.newEventID()))
	for _, t := range i.XDWDefinition.Tasks {
		log.Printf("Creating Workflow Task ID - %v Name - %s", t.ID, t.Name)
		task := XDWTask{}
		task.TaskData.TaskDetails.ID = t.ID
		task.TaskData.TaskDetails.TaskType = t.Tasktype
		task.TaskData.TaskDetails.Name = t.Name
		task.TaskData.TaskDetails.ActualOwner = t.ActualOwner
		task.TaskData.TaskDetails.CreatedBy = i.Role + " " + i.User
		task.TaskData.TaskDetails.CreatedTime = effectiveTime
		task.TaskData.TaskDetails.RenderingMethodExists = "false"
		task.TaskData.TaskDetails.LastModifiedTime = effectiveTime
		task.TaskData.Description = t.Description
		task.TaskData.TaskDetails.Status = tukcnst.CREATED
		for _, inp := range t.Input {
			docinput := Input{}
			docinput.Part.Name = inp.Name
			docinput.Part.AttachmentInfo.Name = inp.Name
			docinput.Part.AttachmentInfo.AccessType = inp.AccessType
			docinput.Part.AttachmentInfo.ContentType = inp.Contenttype
			docinput.Part.AttachmentInfo.ContentCategory = tukcnst.MEDIA_TYPES
			task.TaskData.Input = append(task.TaskData.Input, docinput)
			log.Printf("Created Input Part - %s", inp.Name)
		}
		for _, outp := range t.Output {
			docoutput := Output{}
			docoutput.Part.Name = outp.Name
			docoutput.Part.AttachmentInfo.Name = outp.Name
			docoutput.Part.AttachmentInfo.AccessType = outp.AccessType
			docoutput.Part.AttachmentInfo.ContentType = outp.Contenttype
			docoutput.Part.AttachmentInfo.ContentCategory = tukcnst.MEDIA_TYPES
			task.TaskData.Output = append(task.TaskData.Output, docoutput)
			log.Printf("Created Output Part - %s", outp.Name)
		}
		tev := TaskEvent{}
		tev.EventTime = effectiveTime
		tev.ID = tevidstr
		tev.Identifier = t.ID
		tev.EventType = tukcnst.XDW_TASKEVENTTYPE_CREATED
		tev.Status = tukcnst.COMPLETE
		task.TaskEventHistory.TaskEvent = append(task.TaskEventHistory.TaskEvent, tev)
		i.XDWDocument.TaskList.XDWTask = append(i.XDWDocument.TaskList.XDWTask, task)
		log.Printf("Set Workflow Task Event %s %s status to %s", t.ID, tev.EventType, tev.Status)
	}
	docevent := DocumentEvent{}
	docevent.Author = i.User + i.Role
	docevent.TaskEventIdentifier = tevidstr
	docevent.EventTime = effectiveTime
	docevent.EventType = tukcnst.CREATE
	docevent.PreviousStatus = tukcnst.CREATE
	docevent.ActualStatus = tukcnst.CREATED
	log.Printf("Created Workflow Document Event Set status to %s", tukcnst.CREATED)
	i.XDWDocument.WorkflowStatusHistory.DocumentEvent = append(i.XDWDocument.WorkflowStatusHistory.DocumentEvent, docevent)
	i.Response, _ = xml.MarshalIndent(i.XDWDocument, "", "  ")
	i.XDWVersion = 0
	log.Printf("Created new %s Workflow for Patient %s", i.XDWDocument.WorkflowDefinitionReference, i.NHS_ID)
}

// IHE XDW Content Consumer
func (i *Transaction) contentConsumer() error {
	if err := i.setXDWS(); err != nil {
		return err
	}
	if i.Workflows.Count == 1 {
		log.Printf("Setting %s Workflow state for Patient %s", i.XDWDocument.WorkflowDefinitionReference, i.XDWDocument.Patient.ID.Extension)
		i.XDWState.Created = i.XDWDocument.EffectiveTime.Value
		i.XDWState.Status = i.XDWDocument.WorkflowStatus
		i.XDWState.IsPublished = i.Workflows.Workflows[1].Published
		i.setWorkflowLatestEventTime()
		i.setWorkflowDuration()
		workflowStartTime := tukutil.GetTimeFromString(i.XDWState.Created)
		workflowCompleteByDate := workflowStartTime
		if i.XDWDefinition.CompleteByTime == "" {
			i.XDWState.CompleteBy = "Non Specified"
		} else {
			period := strings.Split(i.XDWDefinition.CompleteByTime, "(")[0]
			periodDuration := tukutil.GetIntFromString(strings.Split(strings.Split(i.XDWDefinition.CompleteByTime, "(")[1], ")")[0])
			switch period {
			case "month":
				workflowCompleteByDate = tukutil.GetFutureDate(workflowStartTime, 0, periodDuration, 0, 0, 0)
			case "day":
				workflowCompleteByDate = tukutil.GetFutureDate(workflowStartTime, 0, 0, periodDuration, 0, 0)
			case "hour":
				workflowCompleteByDate = tukutil.GetFutureDate(workflowStartTime, 0, 0, 0, periodDuration, 0)
			case "min":
				workflowCompleteByDate = tukutil.GetFutureDate(workflowStartTime, 0, 0, 0, 0, periodDuration)
			}
			if workflowCompleteByDate.Before(workflowStartTime) {
				i.XDWState.CompleteBy = "Non Specified"
			} else {
				i.XDWState.CompleteBy = strings.Split(workflowCompleteByDate.String(), " +")[0]
			}
			i.setIsWorkflowOverdueState()
		}

		for _, deftask := range i.XDWDefinition.Tasks {
			for _, doctask := range i.XDWDocument.TaskList.XDWTask {
				if doctask.TaskData.TaskDetails.ID == deftask.ID {
					tstate := XDWTaskState{}
					tstate.TaskID = tukutil.GetIntFromString(doctask.TaskData.TaskDetails.ID)
					tstate.Created = doctask.TaskData.TaskDetails.CreatedTime

					taskStartTime := tukutil.GetTimeFromString(tstate.Created)
					if deftask.CompleteByTime == "" {
						tstate.CompleteBy = "Non Specified"
					} else {
						period := strings.Split(deftask.CompleteByTime, "(")[0]
						periodDuration := tukutil.GetIntFromString(strings.Split(strings.Split(deftask.CompleteByTime, "(")[1], ")")[0])
						switch period {
						case "month":
							i.XDWState.CompleteBy = strings.Split(tukutil.GetFutureDate(taskStartTime, 0, periodDuration, 0, 0, 0).String(), " +0")[0]
						case "day":
							i.XDWState.CompleteBy = strings.Split(tukutil.GetFutureDate(taskStartTime, 0, 0, periodDuration, 0, 0).String(), " +0")[0]
						case "hour":
							i.XDWState.CompleteBy = strings.Split(tukutil.GetFutureDate(taskStartTime, 0, 0, 0, periodDuration, 0).String(), " +0")[0]
						case "min":
							i.XDWState.CompleteBy = strings.Split(tukutil.GetFutureDate(taskStartTime, 0, 0, 0, 0, periodDuration).String(), " +0")[0]
						}
					}
				}
			}
		}

	}
	return nil
}

func (i *Transaction) setXDWS() error {
	log.Println("Setting XDW States")
	i.Workflows = tukdbint.Workflows{Action: tukcnst.SELECT}
	wf := tukdbint.Workflow{Pathway: i.Pathway, NHSId: i.NHS_ID, Version: i.XDWVersion}
	i.Workflows.Workflows = append(i.Workflows.Workflows, wf)
	i.XDWEvents = tukdbint.Events{Action: tukcnst.SELECT}
	ev := tukdbint.Event{Pathway: i.Pathway, NhsId: i.NHS_ID, Version: i.XDWVersion}
	i.XDWEvents.Events = append(i.XDWEvents.Events, ev)

	if err := tukdbint.NewDBEvent(&i.Workflows); err != nil {
		log.Println(err.Error())
		return err
	}
	if err := tukdbint.NewDBEvent(&i.XDWEvents); err != nil {
		log.Println(err.Error())
		return err
	}
	i.Dashboard.Total = i.Workflows.Count
	for _, wf := range i.Workflows.Workflows {
		if len(wf.XDW_Doc) > 0 {
			if i.XDWVersion == -1 || wf.Version == i.XDWVersion {
				if err := xml.Unmarshal([]byte(wf.XDW_Doc), &i.XDWDocument); err != nil {
					log.Println(err.Error())
					return err
				}
				log.Printf("%s Workflow Status is %s", wf.XDW_Key, i.XDWDocument.WorkflowStatus)
				if err := xml.Unmarshal([]byte(wf.XDW_Def), &i.XDWDefinition); err != nil {
					log.Println(err.Error())
					return err
				}
				switch i.XDWDocument.WorkflowStatus {
				case tukcnst.READY:
					i.Dashboard.InProgress = i.Dashboard.InProgress + 1
				case tukcnst.IN_PROGRESS:
					i.Dashboard.InProgress = i.Dashboard.InProgress + 1
				case tukcnst.COMPLETE:
					i.Dashboard.Complete = i.Dashboard.Complete + 1
				}
				if i.setIsWorkflowOverdueState() {
					i.Dashboard.TargetMissed = i.Dashboard.TargetMissed + 1
				} else {
					if i.XDWDocument.WorkflowStatus == tukcnst.COMPLETE {
						i.Dashboard.TargetMet = i.Dashboard.TargetMet + 1
					}
				}
			}
		}
	}
	return nil
}

// XDW Admin

func (i *Transaction) registerWorkflowDefinition(ismeta bool) error {
	var err error
	if i.Pathway == "" {
		return errors.New("pathway is not set")
	}
	if i.Request == nil || string(i.Request) == "" {
		return errors.New("request bytes is not set")
	}
	if ismeta {
		log.Println("Persisting Workflow XDS Meta")
		err = i.registerWorkflowXDSMeta()
	} else {
		log.Println("Registering Workflow Definition")
		err = i.registerWorkflowDef()
	}
	return err
}
func (i *Transaction) registerWorkflowDef() error {
	var err error
	err = json.Unmarshal(i.Request, &i.XDWDefinition)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	event := tukdsub.DSUBEvent{Action: tukcnst.CANCEL, Pathway: i.XDWDefinition.Ref}
	tukdsub.New_Transaction(&event)
	log.Printf("Cleaned Event Service Subscriptions for Pathway %s", i.XDWDefinition.Ref)
	pwyExpressions := make(map[string]string)
	if err = i.persistXDWDefinition(); err == nil {
		log.Println("Parsing XDW Tasks for potential DSUB Broker Subscriptions")
		for _, task := range i.XDWDefinition.Tasks {
			for _, inp := range task.Input {
				log.Printf("Checking Input Task %s", inp.Name)
				if inp.AccessType == tukcnst.XDS_REGISTERED {
					pwyExpressions[inp.Name] = i.XDWDefinition.Ref
					log.Printf("Task %v %s task input %s included in potential DSUB Broker subscriptions", task.ID, task.Name, inp.Name)
				} else {
					log.Printf("Input Task %s does not require a dsub broker subscription", inp.Name)
				}
			}
			for _, out := range task.Output {
				log.Printf("Checking Output Task %s", out.Name)
				if out.AccessType == tukcnst.XDS_REGISTERED {
					pwyExpressions[out.Name] = i.XDWDefinition.Ref
					log.Printf("Task %v %s task output %s included in potential DSUB Broker subscriptions", task.ID, task.Name, out.Name)
				} else {
					log.Printf("Output Task %s does not require a dsub broker subscription", out.Name)
				}
			}
		}
	}
	log.Printf("Found %v potential DSUB Broker Subscriptions - %s", len(pwyExpressions), pwyExpressions)
	if len(pwyExpressions) > 0 {
		event.Action = tukcnst.CREATE
		event.BrokerURL = i.DSUB_BrokerURL
		event.ConsumerURL = i.DSUB_ConsumerURL
		for expression := range pwyExpressions {
			event.Expressions = append(event.Expressions, expression)
		}
		err = tukdsub.New_Transaction(&event)
	}
	return err
}
func (i *Transaction) registerWorkflowXDSMeta() error {
	var err error
	xdw := tukdbint.XDW{Name: i.Pathway + "_meta", IsXDSMeta: true}
	xdws := tukdbint.XDWS{Action: tukcnst.DELETE}
	xdws.XDW = append(xdws.XDW, xdw)
	if err = tukdbint.NewDBEvent(&xdws); err == nil {
		log.Printf("Deleted Existing XDS Meta for Pathway %s", i.Pathway)
		xdw = tukdbint.XDW{Name: i.Pathway + "_meta", IsXDSMeta: true, XDW: string(i.Request)}
		xdws = tukdbint.XDWS{Action: tukcnst.INSERT}
		xdws.XDW = append(xdws.XDW, xdw)
		if err = tukdbint.NewDBEvent(&xdws); err == nil {
			log.Printf("Persisted Workflow XDS Meta for Pathway %s", i.Pathway)
		}
	}
	return err
}
func (i *Transaction) persistXDWDefinition() error {
	var err error
	xdw := tukdbint.XDW{Name: i.Pathway, IsXDSMeta: false}
	xdws := tukdbint.XDWS{Action: tukcnst.DELETE}
	xdws.XDW = append(xdws.XDW, xdw)
	if err = tukdbint.NewDBEvent(&xdws); err == nil {
		log.Printf("Deleted Existing XDW Definition for Pathway %s", i.Pathway)
		xdw = tukdbint.XDW{Name: i.Pathway, IsXDSMeta: false, XDW: string(i.Request)}
		xdws = tukdbint.XDWS{Action: tukcnst.INSERT}
		xdws.XDW = append(xdws.XDW, xdw)
		if err = tukdbint.NewDBEvent(&xdws); err == nil {
			log.Printf("Persisted New XDW Definition for Pathway %s", i.Pathway)
		}
	}
	return err
}

func (i *Transaction) GetTaskLastModifiedTime() string {
	for _, task := range i.XDWDocument.TaskList.XDWTask {
		if task.TaskData.TaskDetails.ID == tukutil.GetStringFromInt(i.Task_ID) {
			return task.TaskData.TaskDetails.LastModifiedTime
		}
	}
	return ""
}
func (i *Transaction) setWorkflowDuration() {
	ws, err := time.Parse(time.RFC3339, i.XDWDocument.EffectiveTime.Value)
	if err != nil {
		fmt.Println(err)
	}
	we := time.Now()
	if i.XDWDocument.WorkflowStatus == tukcnst.COMPLETE {
		we = i.XDWState.LatestWorkflowEventTime
	}
	i.XDWState.WorkflowDuration = we.Sub(ws)
	log.Println("Duration - " + i.XDWState.WorkflowDuration.String())
	//25h15m32.877428s
	totmins := int(i.XDWState.WorkflowDuration.Minutes())
	var onehour = 60
	var oneday = onehour * 24
	if totmins < onehour {
		if totmins == 0 {
			i.XDWState.PrettyWorkflowDuration = "Less than a Minute"
		} else {
			if totmins == 1 {
				i.XDWState.PrettyWorkflowDuration = "1 Min"
			} else {
				i.XDWState.PrettyWorkflowDuration = tukutil.GetStringFromInt(totmins) + " Mins"
			}
		}
	} else {
		if totmins < oneday {
			hrs := totmins / onehour
			if hrs == 0 {
				i.XDWState.PrettyWorkflowDuration = "1 Hour"
			} else {
				i.XDWState.PrettyWorkflowDuration = tukutil.GetStringFromInt(hrs) + " Hours " + tukutil.GetStringFromInt(totmins-(hrs*onehour)) + " Mins"
			}

		} else {
			days := totmins / oneday
			if days == 0 {
				i.XDWState.PrettyWorkflowDuration = "1 Day"
			} else {
				hrs := totmins - (days * oneday)
				mins := totmins - (days * oneday) - (hrs * onehour)
				i.XDWState.PrettyWorkflowDuration = tukutil.GetStringFromInt(days) + " Days " + tukutil.GetStringFromInt(hrs) + " Hrs " + tukutil.GetStringFromInt(mins) + " Mins"
			}
		}
	}
}
func (i *Transaction) setWorkflowLatestEventTime() error {
	log.Printf("Setting Latest Workflow Event Time for Pathway %s NHS ID %s", i.Pathway, i.NHS_ID)
	i.XDWState.LatestWorkflowEventTime = tukutil.GetTimeFromString(i.XDWDocument.EffectiveTime.Value)
	for _, docevent := range i.XDWDocument.WorkflowStatusHistory.DocumentEvent {
		etime := tukutil.GetTimeFromString(docevent.EventTime)
		if etime.After(i.XDWState.LatestWorkflowEventTime) {
			i.XDWState.LatestWorkflowEventTime = etime
		}
	}
	log.Printf("Latest Event Time set to %s ", i.XDWState.LatestWorkflowEventTime.String())
	return nil
}
func (i *Transaction) persistWorkflow() error {
	var err error
	wfs := tukdbint.Workflows{Action: tukcnst.INSERT}
	wf := tukdbint.Workflow{
		XDW_Key: strings.ToUpper(i.Pathway) + i.NHS_ID,
		XDW_UID: i.XDWDocument.ID.Extension,
		Version: i.XDWVersion,
	}
	xdwDocBytes, _ := xml.MarshalIndent(i.XDWDocument, "", "  ")
	xdwDefBytes, _ := json.Marshal(i.XDWDefinition)
	wf.XDW_Doc = string(xdwDocBytes)
	wf.XDW_Def = string(xdwDefBytes)
	wfs.Workflows = append(wfs.Workflows, wf)
	if err = tukdbint.NewDBEvent(&wfs); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("Persisted Workflow Version %v for Pathway %s NHS ID %s", i.XDWVersion, i.Pathway, i.NHS_ID)
	}
	return err
}
func (i *Transaction) newEventID() int64 {
	ev := tukdbint.Event{
		DocName:            i.XDWDocument.WorkflowDefinitionReference + "-" + i.NHS_ID,
		ClassCode:          i.XDSDocumentMeta.Classcode,
		ConfCode:           i.XDSDocumentMeta.Confcode,
		FormatCode:         i.XDSDocumentMeta.Formatcode,
		FacilityCode:       i.XDSDocumentMeta.Facilitycode,
		PracticeCode:       i.XDSDocumentMeta.Practicesettingcode,
		Expression:         i.Expression,
		Authors:            i.XDWDocument.Author.AssignedAuthor.AssignedPerson.Name.Prefix + " " + i.XDWDocument.Author.AssignedAuthor.AssignedPerson.Name.Family,
		XdsPid:             "NA",
		XdsDocEntryUid:     i.XDWDocument.ID.Root,
		RepositoryUniqueId: "ODD",
		NhsId:              i.NHS_ID,
		User:               i.User,
		Org:                i.Org,
		Role:               i.Role,
		Topic:              tukcnst.DSUB_TOPIC_TYPE_CODE,
		Pathway:            i.Pathway,
		Notes:              string(i.Request),
		Version:            0,
		BrokerRef:          "NA",
	}
	evs := tukdbint.Events{Action: tukcnst.INSERT}
	evs.Events = append(evs.Events, ev)
	if err := tukdbint.NewDBEvent(&evs); err != nil {
		log.Println(err.Error())
		return 0
	}
	log.Printf("Created Event ID :  = %v", evs.LastInsertId)
	return evs.LastInsertId
}
func (i *Transaction) setIsWorkflowOverdueState() bool {
	if i.XDWDefinition.CompleteByTime != "" {
		completionDate := i.getWorkflowCompleteByDate()
		if time.Now().After(completionDate) {
			if i.XDWDocument.WorkflowStatus == tukcnst.COMPLETE {
				i.setWorkflowLatestEventTime()
				return i.XDWState.LatestWorkflowEventTime.After(completionDate)
			}
		}
	}
	return false
}
func (i *Transaction) getWorkflowCompleteByDate() time.Time {
	workflowStartTime := tukutil.GetTimeFromString(i.XDWDocument.EffectiveTime.Value)
	days := tukutil.GetIntFromString(strings.Split(strings.Split(i.XDWDefinition.CompleteByTime, "(")[1], ")")[0])
	completionDate := tukutil.GetFutueDaysDate(workflowStartTime, days)
	return completionDate
}
func getLocalId(mid string) string {
	idmaps := tukdbint.IdMaps{Action: tukcnst.SELECT}
	idmap := tukdbint.IdMap{Mid: mid}
	idmaps.LidMap = append(idmaps.LidMap, idmap)
	if err := tukdbint.NewDBEvent(&idmaps); err != nil {
		log.Println(err.Error())
		return mid
	}
	if idmaps.Cnt == 1 {
		return idmaps.LidMap[1].Lid
	}
	return mid
}
func getMappedId(lid string) string {
	idmaps := tukdbint.IdMaps{Action: tukcnst.SELECT}
	idmap := tukdbint.IdMap{Lid: lid}
	idmaps.LidMap = append(idmaps.LidMap, idmap)
	if err := tukdbint.NewDBEvent(&idmaps); err != nil {
		log.Println(err.Error())
		return lid
	}
	if idmaps.Cnt == 1 {
		return idmaps.LidMap[1].Mid
	}
	return lid
}

// sort interface
type DocumentEvents []DocumentEvent

func (e DocumentEvents) Len() int {
	return len(e)
}
func (e DocumentEvents) Less(i, j int) bool {
	return e[i].EventTime > e[j].EventTime
}
func (e DocumentEvents) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
