package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	repoOwner   = "kanisterio"
	repoName    = "kanister"
	searchTerm  = "blueprint"
	targetPath  = "examples"
	accessToken = "your github token"
)

type GitHubContent struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

func searchFiles(ctx context.Context, path string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, path), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request to GitHub API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var contents []GitHubContent
	err = json.Unmarshal(body, &contents)
	if err != nil {
		fmt.Println("Error decoding GitHub API response:", err)
		return
	}

	for _, content := range contents {
		if content.Type == "file" && strings.HasSuffix(content.Name, ".yaml") && strings.Contains(content.Name, searchTerm) {
			fmt.Println("YAML file found:", content.HTMLURL)
		} else if content.Type == "dir" {
			searchFiles(ctx, path+"/"+content.Name)
		}
	}
}

func main() {
	ctx := context.Background()
	searchFiles(ctx, targetPath)
}
