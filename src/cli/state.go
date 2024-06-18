package main

import (
	"encoding/json"
	"log"
	"os"
)

type Filter struct {
	Tag       string `json:"tag"`
	ProjectId string `json:"project"`
	Priority  string `json:"priority"`
	Search    string
}

type State struct {
	Contexts      []string `json:"contexts"`
	ActiveContext string   `json:"activeContext"`
	Filter        Filter   `json:"filter"`
}

var AppState State = State{
	Contexts:      []string{},
	Filter:        Filter{},
	ActiveContext: "",
}

func ReadAppState() {
	if _, err := os.Stat(configPath + "/state.json"); os.IsNotExist(err) {
		initAppStateFile()
	}

	appConfigJson, err := os.ReadFile(configPath + "/state.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(appConfigJson, &AppState)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteAppState() {
	appStateJson, err := json.Marshal(AppState)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configPath+"/state.json", appStateJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func initAppStateFile() {
	appState := State{
		Contexts:      []string{},
		ActiveContext: "",
		Filter:        Filter{},
	}
	appStateJson, err := json.Marshal(appState)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configPath+"/state.json", appStateJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func AddContext(context string) {
	hasContext := false
	for _, c := range AppState.Contexts {
		if c == context {
			hasContext = true
		}
	}

	if !hasContext {
		AppState.Contexts = append(AppState.Contexts, context)
	}

	WriteAppState()
}

func SwitchToContext(context string) {
	hasContext := false
	for _, c := range AppState.Contexts {
		if c == context {
			hasContext = true
		}
	}

	if !hasContext {
		return
	}

	AppState.ActiveContext = context
	WriteAppState()
}
