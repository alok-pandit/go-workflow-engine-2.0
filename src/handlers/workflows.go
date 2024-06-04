package handlers

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alok-pandit/go-workflow-engine-2.0/src/models"
	workflows "github.com/alok-pandit/go-workflow-engine-2.0/src/workflow_definitions"
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

	bpmn, err := ParseBPMN(data)

	if err != nil {
		panic(err)
	}

	return bpmn.Processes[0]

}

func fileDirWalk() []string {

	var subDirs []string

	rootDir := os.Getenv("WORKFLOW_ROOT")

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != rootDir { // Ignore the starting directory
			subDirs = append(subDirs, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
		return nil
	}

	return subDirs

}
func Genesis() {

	dirList := fileDirWalk()

	for _, dir := range dirList {

		procName := strings.Split(dir, "/")

		process := Build(dir + "/" + procName[2] + ".bpmn")

		sm := buildStageList(process)

		count := 0

		go populateStageMap(sm, count, process, procName[2])

	}

}

// func writeToJson(sm map[string]*models.Stage, dir string) error {

// 	// Marshal the map to JSON
// 	jsonData, err := json.Marshal(sm)
// 	if err != nil {
// 		fmt.Println("Error marshalling map to JSON:", err)
// 		return err
// 	}

// 	// Open the file for writing (with create if doesn't exist)
// 	file, err := os.Create(os.Getenv("WORKFLOW_ROOT") + dir + "/" + dir + ".json")
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return err
// 	}
// 	defer file.Close() // Close the file on exit

// 	// Write the JSON data to the file
// 	_, err = file.Write(jsonData)
// 	if err != nil {
// 		fmt.Println("Error writing JSON to file:", err)
// 		return err
// 	}

// 	fmt.Println("Successfully wrote data to json")

// 	return nil

// }

func populateStageMap(sm map[string]*models.Stage, c int, process models.Process, workflowDir string) {

	if c <= len(sm) {

		s := process.SequenceFlows[c]

		if sm[s.Source].Type == "Gateway" {
			if sm[s.Target].Name == "Success" {
				sm[s.Source].Next = s.Target
			}
			if sm[s.Target].Name == "Failure" {
				sm[s.Source].Failure = s.Target
			}
		} else {
			sm[s.Source].Next = s.Target
		}

		c = c + 1

		populateStageMap(sm, c, process, workflowDir)

	} else {

		for k, v := range sm {
			nextStage, ok := sm[v.Next]
			if ok {
				if nextStage.Type == "Gateway" {
					sm[k].Next = sm[nextStage.Next].Next
					sm[k].Failure = sm[nextStage.Failure].Next
				}
			}
		}

		workflows.Workflows[workflowDir] = sm

		node := workflows.BeginWFExecution("login")

		fmt.Printf("%+v\n\n", node)

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
