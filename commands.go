package main

import (
	"strconv"
)

type Command struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type ServiceMenuCommand struct {
	Command
	Payload ServiceMenuCommandPayload `json:"payload"`
}
type ServiceMenuCommandPayload struct {
	Id     string                          `json:"id"`
	Params ServiceMenuCommandPayloadParams `json:"params"`
}
type ServiceMenuCommandPayloadParams struct {
	Id    string `json:"id"`
	IrKey string `json:"irKey"`
}

func NewServiceMenuCommand(id int) *ServiceMenuCommand {
	cmd := &ServiceMenuCommand{
		Command: Command{
			Id:   "show_service_menu_" + strconv.Itoa(id),
			Type: "request",
			URI:  "ssap://com.webos.applicationManager/launch",
		},
		Payload: ServiceMenuCommandPayload{
			Id: "com.webos.app.factorywin",
			Params: ServiceMenuCommandPayloadParams{
				Id:    "executeFactory",
				IrKey: "inStart",
			},
		},
	}

	return cmd
}
