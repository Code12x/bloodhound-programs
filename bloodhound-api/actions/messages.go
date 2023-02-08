package actions

import "example/bloodhound/models"

const (
	ERROR       string = "ERROR"
	SUCCESS     string = "SUCCESS"
	PROCCESSING string = "PROCESSING"
)

type serverResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Payload models.SmellyServer `json:"payload"`
}

func badRequestMessage() serverResponse {
	return serverResponse{Status: ERROR, Message: "Incorect paramaters"}
}

func unauthorized() serverResponse {
	return serverResponse{Status: ERROR, Message: "You are not authorized to call this api"}
}

func noResultsMessage() serverResponse {
	return serverResponse{Status: ERROR, Message: "No results"}
}

func serverNotFound() serverResponse {
	return serverResponse{Status: ERROR, Message: "No Minecraft server found at that address"}
}

func playerNotFound() serverResponse {
	return serverResponse{Status: ERROR, Message: "Player not found"}
}

func serverUpdateSuccess(server models.SmellyServer) serverResponse {
	return serverResponse{Status: SUCCESS, Message: "Server upated", Payload: server}
}

func operationReceived() serverResponse {
	return serverResponse{Status: PROCCESSING, Message: "Operation received"}
}
