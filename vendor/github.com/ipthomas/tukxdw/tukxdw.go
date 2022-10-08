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
	WorkflowDefinition WorkflowDefinition
	DSUB_BrokerURL     string
	DSUB_ConsumerURL   string
}
type XDW_Int interface {
	processRequest() error
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
	case tukcnst.XDW_CONTENT_UPDATER:
	}
	return nil
}
func (i *XDWTransaction) registerWorkflowDefinition() error {
	pwyExpressions := make(map[string]string)
	if err := i.persistXDWDefinition(); err == nil {
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
		tukdsub.New_Transaction(&event)
	}
	return nil
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
