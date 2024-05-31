package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/alok-pandit/go-workflow-engine-2.0/src/models"
)

// ParseBPMN parses a BPMN file and returns a BPMN struct
func ParseBPMN(data []byte) (*models.BPMN, error) {
	var bpmn models.BPMN
	err := xml.Unmarshal(data, &bpmn)
	return &bpmn, err
}

func Build(filename string) models.Process {

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Parse the BPMN data
	bpmn, err := ParseBPMN(data)
	if err != nil {
		panic(err)
	}

	return bpmn.Processes[0]

}

func ExecuteProcess(process models.Process) {

	sm := buildStageList(process)

	count := 0

	populateStageMap(sm, count, process)

}

var StageMap map[string]*models.Stage

func populateStageMap(sm map[string]*models.Stage, c int, process models.Process) {
	if c <= len(sm) {
		s := process.SequenceFlows[c]
		sm[s.Source].Next = s.Target
		if sm[s.Source].Type == "Gateway" {
			if sm[s.Target].Name == "Success" {
				sm[s.Source].Next = s.Target
				if len(sm[s.Source].Failure) < 1 {
				a:
					for _, p := range process.IntermediateCatchEvents {
						if p.ID != sm[s.Source].Next {
							sm[s.Source].Failure = p.ID
							break a
						}
					}
				}
			}
			if sm[s.Target].Name == "Failure" {
				sm[s.Source].Failure = s.Target
				if len(sm[s.Source].Next) < 1 {
				b:
					for _, p := range process.IntermediateCatchEvents {
						if p.ID != sm[s.Source].Failure {
							sm[s.Source].Next = p.ID
							break b
						}
					}
				}
			}
		}
		c = c + 1
		populateStageMap(sm, c, process)
	} else {
		b, err := json.MarshalIndent(sm, "", "  ")

		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println("Stage Map Rec: ", string(b))

		StageMap = sm
	}
}

func buildStageList(process models.Process) map[string]*models.Stage {
	stageMap := make(map[string]*models.Stage)

	// Create stages for all tasks, gateways, start and end events
	for _, task := range process.Tasks {

		stageMap[task.ID] = &models.Stage{Name: task.Name, Type: models.BPMNElementTypes.Task}

	}

	for _, gateway := range process.Gateways {
		stageMap[gateway.ID] = &models.Stage{Name: "Gateway", Type: models.BPMNElementTypes.Gateway}
	}

	stageMap[process.StartEvent.ID] = &models.Stage{Name: "Start Event", Type: models.BPMNElementTypes.StartEvent}

	stageMap[process.EndEvent.ID] = &models.Stage{Name: "End Event", Type: models.BPMNElementTypes.EndEvent}

	for _, cond := range process.IntermediateCatchEvents {

		stageMap[cond.ID] = &models.Stage{Name: cond.Name, Type: models.BPMNElementTypes.IntermediateCatchEvent}

	}

	return stageMap

}
