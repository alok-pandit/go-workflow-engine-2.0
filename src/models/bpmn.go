package models

import "encoding/xml"

type Stage struct {
	Name    string
	Type    string
	Next    string
	Failure string
}

// BPMN represents a BPMN 2.0 process definition
type BPMN struct {
	XMLName       xml.Name `xml:"definitions"`
	Collaboration Collaboration
	Processes     []Process `xml:"process"`
}

// Collaboration represents a collaboration element in BPMN
type Collaboration struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name,attr"`
	Participant Participant
}

// Participant represents a participant element in BPMN
type Participant struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	// Ref  string `xml:"processRef,attr"` // Not used in this example
}

// Process represents a process element in BPMN
type Process struct {
	ID                      string                   `xml:"id,attr"`
	Name                    string                   `xml:"name,attr"`
	IsExecutable            bool                     `xml:"isExecutable,attr"`
	Tasks                   []Task                   `xml:"task"`
	Gateways                []Gateway                `xml:"exclusiveGateway"` // Assuming only exclusive gateways are used
	StartEvent              StartEvent               `xml:"startEvent"`
	EndEvent                EndEvent                 `xml:"endEvent"`
	SequenceFlows           []SequenceFlow           `xml:"sequenceFlow"`
	IntermediateCatchEvents []IntermediateCatchEvent `xml:"intermediateCatchEvent"`
}

// Task represents a task element in BPMN
type Task struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// Gateway represents a gateway element in BPMN
type Gateway struct {
	ID string `xml:"id,attr"`
}

// StartEvent represents a start event element in BPMN
type StartEvent struct {
	ID string `xml:"id,attr"`
}

// EndEvent represents an end event element in BPMN
type EndEvent struct {
	ID string `xml:"id,attr"`
}

// SequenceFlow represents a sequence flow element in BPMN
type SequenceFlow struct {
	ID     string `xml:"id,attr"`
	Source string `xml:"sourceRef,attr"`
	Target string `xml:"targetRef,attr"`
}

// IntermediateCatchEvent represents an intermediate catch event element in BPMN
type IntermediateCatchEvent struct {
	ID          string                     `xml:"id,attr"`
	Name        string                     `xml:"name,attr"`
	Conditional ConditionalEventDefinition `xml:"conditionalEventDefinition"`
}

// ConditionalEventDefinition represents a conditional event definition element in BPMN
type ConditionalEventDefinition struct {
	Condition Condition `xml:"condition,attr"`
}

// Condition represents a condition element in BPMN (simplified)
type Condition struct {
	Language string `xml:"language,attr"`
}

var BPMNElementTypes = newTypeRegistry()

func newTypeRegistry() *typeRegistry {
	return &typeRegistry{
		Gateway:                "Gateway",
		StartEvent:             "StartEvent",
		EndEvent:               "EndEvent",
		Task:                   "Task",
		IntermediateCatchEvent: "IntermediateCatchEvent",
	}
}

type typeRegistry struct {
	Gateway                string
	StartEvent             string
	EndEvent               string
	Task                   string
	IntermediateCatchEvent string
}
