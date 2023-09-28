package main

import (
	"blueprint"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/generateBlueprint", generateBlueprintHandler)
	fmt.Println("Server started on :8083")
	http.ListenAndServe(":8083", nil)
}

func generateBlueprintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type Response struct {
		BlueprintYAML string `json:"blueprintYAML"`
	}

	var requestBody struct {
		Name           string   `json:"name"`
		Actions        []string `json:"actions"`
		ConfigMapNames string   `json:"configMapNames"`
		ContainerName  string   `json:"containerName"`
		BackupCommand  string   `json:"backupCommand"`
		RestoreCommand string   `json:"restoreCommand"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blueprintYAML, err := blueprint.GenerateBlueprint(requestBody.Name, requestBody.Actions, requestBody.ConfigMapNames, requestBody.ContainerName, requestBody.BackupCommand, requestBody.RestoreCommand)
	if err != nil {
		http.Error(w, "Failed to generate blueprint", http.StatusInternalServerError)
		return
	}

	response := Response{
		BlueprintYAML: blueprintYAML,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
