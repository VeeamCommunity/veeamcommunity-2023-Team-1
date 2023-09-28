package blueprint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"gopkg.in/yaml.v2"
)

type Blueprint struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Actions struct {
		Backup  Action `yaml:"backup"`
		Restore Action `yaml:"restore"`
		Delete  Action `yaml:"delete"`
	} `yaml:"actions"`
}

type Action struct {
	SecretNames        []string        `yaml:"secretNames,omitempty"`
	InputArtifactNames []string        `yaml:"inputArtifactNames,omitempty"`
	OutputArtifacts    OutputArtifacts `yaml:"outputArtifacts,omitempty"`
	Phases             []struct {
		Func string `yaml:"func"`
		Name string `yaml:"name"`
		Args struct {
			Namespace   string     `yaml:"namespace"`
			Image       string     `yaml:"image"`
			PodOverride struct{}   `yaml:"podOverride"`
			Volumes     []struct{} `yaml:"volumes"`
			Command     []string   `yaml:"command"`
		} `yaml:"args"`
	} `yaml:"phases,omitempty"`
}

type OutputArtifacts struct {
	CockroachDBCloudDump struct {
		KeyValue struct {
			S3Path string `yaml:"s3path"`
		} `yaml:"keyValue"`
	} `yaml:"cockroachDBCloudDump,omitempty"`
}

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

	blueprintYAML, err := GenerateBlueprint(requestBody.Name, requestBody.Actions, requestBody.ConfigMapNames, requestBody.ContainerName, requestBody.BackupCommand, requestBody.RestoreCommand)
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

func GenerateBlueprint(name string, actions []string, configMapNames, containerName, backupCommand, restoreCommand string) (string, error) {
	blueprint := Blueprint{
		APIVersion: "cr.kanister.io/v1alpha1",
		Kind:       "Blueprint",
		Metadata: struct {
			Name string `yaml:"name"`
		}{
			Name: name,
		},
	}

	if contains(actions, "backup") {
		backupAction := Action{
			SecretNames: []string{"Secret"},
			OutputArtifacts: OutputArtifacts{
				CockroachDBCloudDump: struct {
					KeyValue struct {
						S3Path string `yaml:"s3path"`
					} `yaml:"keyValue"`
				}{
					KeyValue: struct {
						S3Path string `yaml:"s3path"`
					}{
						S3Path: "/path/{{ .StatefulSet.Namespace }}/{{ index .Object.metadata.labels \"app.kubernetes.io/instance\" }}/{{ toDate \"2006-01-02T15:04:05.999999999Z07:00\" .Time  | date \"2006-01-02T15-04-05\" }}",
					},
				},
			},
			Phases: []struct {
				Func string `yaml:"func"`
				Name string `yaml:"name"`
				Args struct {
					Namespace   string     `yaml:"namespace"`
					Image       string     `yaml:"image"`
					PodOverride struct{}   `yaml:"podOverride"`
					Volumes     []struct{} `yaml:"volumes"`
					Command     []string   `yaml:"command"`
				} `yaml:"args"`
			}{
				{
					Func: "KubeTask",
					Name: "backupToS3",
					Args: struct {
						Namespace   string     `yaml:"namespace"`
						Image       string     `yaml:"image"`
						PodOverride struct{}   `yaml:"podOverride"`
						Volumes     []struct{} `yaml:"volumes"`
						Command     []string   `yaml:"command"`
					}{
						Namespace:   "YourNamespace",
						Image:       "latest/image",
						PodOverride: struct{}{},
						Volumes:     []struct{}{},
						Command:     []string{backupCommand},
					},
				},
			},
		}

		blueprint.Actions.Backup = backupAction
	}

	if contains(actions, "restore") {
		restoreAction := Action{
			InputArtifactNames: []string{"database"},
			Phases: []struct {
				Func string `yaml:"func"`
				Name string `yaml:"name"`
				Args struct {
					Namespace   string     `yaml:"namespace"`
					Image       string     `yaml:"image"`
					PodOverride struct{}   `yaml:"podOverride"`
					Volumes     []struct{} `yaml:"volumes"`
					Command     []string   `yaml:"command"`
				} `yaml:"args"`
			}{
				{
					Func: "KubeTask",
					Name: "restoreFromS3",
					Args: struct {
						Namespace   string     `yaml:"namespace"`
						Image       string     `yaml:"image"`
						PodOverride struct{}   `yaml:"podOverride"`
						Volumes     []struct{} `yaml:"volumes"`
						Command     []string   `yaml:"command"`
					}{
						Namespace:   "YourNamespace",
						Image:       "latest/image",
						PodOverride: struct{}{},
						Volumes:     []struct{}{},
						Command:     []string{restoreCommand},
					},
				},
			},
		}

		blueprint.Actions.Restore = restoreAction
	}

	blueprintYAML, err := yaml.Marshal(blueprint)
	if err != nil {
		return "", err
	}

	return string(blueprintYAML), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
