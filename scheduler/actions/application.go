package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/google/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func ApplicationProcessor(commandFlag configuration.CommandFlag,
	opt configuration.Options, clusterName string, application kinds.Application) {

	executeDeployment(application, opt, clusterName, commandFlag)
}

func loadActionList(path string) kinds.ActionList {

	logger.Info(fmt.Sprintf("Loading %s", path))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}
	appc := kinds.ActionList{}
	err = yaml.Unmarshal(bytes, &appc)
	if err != nil {
		logger.Error(err)
	}
	return appc
}

func selectAndExecute(execute kinds.Execute, deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag, repoName string) {

	switch strings.ToLower(execute.Kind) {
	case "shell":
		ExecuteShellAction(execute, opt, repoName)

	case "kubernetes":
		ExecuteKubernetesAction(execute, clusterName, deployment, commandFlag, opt, repoName)

	case "runactionlist":
		tp, ok := execute.Configuration["Path"]
		if !ok {
			logger.Error("Could not find RunActionList Path")
			return
		}
		al := loadActionList(tp)

		executeActionList(al, deployment, opt, clusterName, commandFlag, repoName)
	}
}

func executeActionList(actionList kinds.ActionList, deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag, remoteVCSRepoName string) {
	for _, a := range actionList.Executions {

		selectAndExecute(a, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
	}
}

func executeDeployment(deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag) {
	logger.Info(fmt.Sprintf("Loading deployment %s\n", deployment.Name))

	remoteVCSRepoName, err := vcs.FetchRepo(deployment.Git, opt)
	if err != nil {
		logger.Error(err.Error())

		return
	}

	executeActionList(deployment.ActionList, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
}
