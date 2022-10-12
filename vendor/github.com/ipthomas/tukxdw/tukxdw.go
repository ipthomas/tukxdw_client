package tukxdw

import (
	"encoding/json"
	"encoding/xml"
	"log"

	"github.com/ipthomas/tukdbint"

	"github.com/ipthomas/tukcnst"

	"github.com/ipthomas/tukdsub"
)

type DSUBSubscribeResponse struct {
	XMLName        xml.Name `xml:"Envelope"`
	Text           string   `xml:",chardata"`
	S              string   `xml:"s,attr"`
	A              string   `xml:"a,attr"`
	Xsi            string   `xml:"xsi,attr"`
	Wsnt           string   `xml:"wsnt,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Header         struct {
		Text   string `xml:",chardata"`
		Action string `xml:"Action"`
	} `xml:"Header"`
	Body struct {
		Text              string `xml:",chardata"`
		SubscribeResponse struct {
			Text                  string `xml:",chardata"`
			SubscriptionReference struct {
				Text    string `xml:",chardata"`
				Address string `xml:"Address"`
			} `xml:"SubscriptionReference"`
		} `xml:"SubscribeResponse"`
	} `xml:"Body"`
}
type DSUBSubscribe struct {
	BrokerUrl   string
	ConsumerUrl string
	Topic       string
	Expression  string
	Request     []byte
	BrokerRef   string
	UUID        string
}
type WorkflowDefinition struct {
	Ref                 string `json:"ref"`
	Name                string `json:"name"`
	Confidentialitycode string `json:"confidentialitycode"`
	CompleteByTime      string `json:"completebytime"`
	CompletionBehavior  []struct {
		Completion struct {
			Condition string `json:"condition"`
		} `json:"completion"`
	} `json:"completionBehavior"`
	Tasks []struct {
		ID                 string `json:"id"`
		Tasktype           string `json:"tasktype"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		Owner              string `json:"owner"`
		ExpirationTime     string `json:"expirationtime"`
		StartByTime        string `json:"startbytime"`
		CompleteByTime     string `json:"completebytime"`
		IsSkipable         bool   `json:"isskipable"`
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
type XDWTransaction struct {
	Action             string
	Pathway            string
	NHS_ID             string
	WorkflowDefinition WorkflowDefinition
	DSUB_BrokerURL     string
	DSUB_ConsumerURL   string
	Response           []byte
}

// XDW Workflow Document Structs

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

type DocumentEvents []DocumentEvent

type XDW_Int interface {
	processRequest() error
}

func (e DocumentEvents) Len() int {
	return len(e)
}
func (e DocumentEvents) Less(i, j int) bool {
	return e[i].EventTime > e[j].EventTime
}
func (e DocumentEvents) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Convienance method to obtain initialised WorkflowDefinition struct which will be needed in the your tukxdw.New_Transaction(&XDWTransaction) if i.Action='register'
func NewWorkflowDefinition(wfdefstring string, wfdefbytes []byte) (WorkflowDefinition, error) {
	wfdefstruct := WorkflowDefinition{}
	if wfdefbytes == nil {
		wfdefbytes = []byte(wfdefstring)
	}
	err := json.Unmarshal(wfdefbytes, &wfdefstruct)
	if err != nil {
		log.Println(err.Error())
	}
	return wfdefstruct, err
}
func New_Transaction(i XDW_Int) error {
	return i.processRequest()
}

func (i *XDWTransaction) processRequest() error {
	switch i.Action {
	case tukcnst.REGISTER:
		return i.registerWorkflowDefinition()
	case tukcnst.XDW_CONTENT_CREATOR:
	case tukcnst.XDW_CONTENT_CONSUMER:
		return i.newXDWContentConsumer()
	case tukcnst.XDW_CONTENT_UPDATER:
	}
	return nil
}
func (i *XDWTransaction) newXDWContentConsumer() error {
	var err error
	wfs := tukdbint.Workflows{Action: tukcnst.SELECT}
	wf := tukdbint.Workflow{XDW_Key: i.Pathway + i.NHS_ID}
	wfs.Workflows = append(wfs.Workflows, wf)
	err = tukdbint.NewDBEvent(&wfs)

	return err
}
func (i *XDWTransaction) registerWorkflowDefinition() error {
	var err error
	pwyExpressions := make(map[string]string)
	if err = i.persistXDWDefinition(); err == nil {
		log.Println("Parsing XDW Tasks for potential DSUB Broker Subscriptions")
		for _, task := range i.WorkflowDefinition.Tasks {
			for _, inp := range task.Input {
				log.Printf("Checking Input Task %s", inp.Name)
				if inp.AccessType == tukcnst.XDS_REGISTERED {
					pwyExpressions[inp.Name] = i.WorkflowDefinition.Ref
					log.Printf("Task %v %s task input %s included in potential DSUB Broker subscriptions", task.ID, task.Name, inp.Name)
				} else {
					log.Printf("Input Task %s does not require a dsub broker subscription", inp.Name)
				}
			}
			for _, out := range task.Output {
				log.Printf("Checking Output Task %s", out.Name)
				if out.AccessType == tukcnst.XDS_REGISTERED {
					pwyExpressions[out.Name] = i.WorkflowDefinition.Ref
					log.Printf("Task %v %s task output %s included in potential DSUB Broker subscriptions", task.ID, task.Name, out.Name)
				} else {
					log.Printf("Output Task %s does not require a dsub broker subscription", out.Name)
				}
			}
		}
	}
	log.Printf("Found %v potential DSUB Broker Subscriptions - %s", len(pwyExpressions), pwyExpressions)
	if len(pwyExpressions) > 0 {
		event := tukdsub.DSUBEvent{Action: tukcnst.CANCEL, Pathway: i.WorkflowDefinition.Ref}
		tukdsub.New_Transaction(&event)
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
func (i *XDWTransaction) persistXDWDefinition() error {
	log.Println("Processing WF Def for Pathway : " + i.WorkflowDefinition.Ref)
	xdw := tukdbint.XDW{Name: i.WorkflowDefinition.Ref}
	xdws := tukdbint.XDWS{Action: tukcnst.DELETE}
	xdws.XDW = append(xdws.XDW, xdw)
	if err := tukdbint.NewDBEvent(&xdws); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Deleted Existing XDW Definition for Pathway %s", i.WorkflowDefinition.Ref)

	xdwBytes, _ := json.Marshal(i)
	xdw = tukdbint.XDW{Name: i.WorkflowDefinition.Ref, IsXDSMeta: false, XDW: string(xdwBytes)}
	xdws = tukdbint.XDWS{Action: tukcnst.INSERT}
	xdws.XDW = append(xdws.XDW, xdw)
	if err := tukdbint.NewDBEvent(&xdws); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Persisted New XDW Definition for Pathway %s", i.WorkflowDefinition.Ref)
	return nil
}
