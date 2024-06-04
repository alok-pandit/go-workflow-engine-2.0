package workflows

import (
	"fmt"

	"github.com/alok-pandit/go-workflow-engine-2.0/src/models"
)

var Workflows = make(map[string]map[string]*models.Stage)

func getStartNode(wf string) string {
	sn := ""

	for k1, v1 := range Workflows[wf] {
		if v1.Name == "Start Event" {
			sn = k1
			break
		}
	}

	return sn

}
func BeginWFExecution(wf_name string) *models.Stage {

	startNode := getStartNode(wf_name)

	return ExecuteNextStep(Workflows[wf_name][startNode], wf_name, true)

}

func ExecuteNextStep(node *models.Stage, wf_name string, isSuccess bool) *models.Stage {

	if isSuccess {
		node = Workflows[wf_name][node.Next]
	} else {
		node = Workflows[wf_name][node.Failure]
	}

	if node.Type != "EndEvent" {
		fmt.Printf("%+v\n\n %+v\n", node, node.Next)
	} else {
		fmt.Println("Workflow Completed")
	}

	return node
}
