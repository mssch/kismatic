package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	linuxArtifact  string = "./artifact/linux/kismatic.tar.gz"
	darwinArtifact string = "./artifact/darwin/kismatic.tar.gz"
)

var tag = flag.String("tag", "", "the name of the git tag")

func main() {
	// 0. Validate/get all the pre-reqs we need for creating the release
	flag.Parse()
	if *tag == "" {
		exit("tag cannot be empty")
	}
	fmt.Printf("Creating release for tag %q\n", *tag)
	authToken := os.Getenv("GITHUB_TOKEN")
	if authToken == "" {
		exit("GITHUB_TOKEN environment variable is empty.")
	}
	// Look for the linux and darwin tarballs
	if _, err := os.Stat(linuxArtifact); err != nil {
		exit(fmt.Sprintf("%v", err))
	}
	if _, err := os.Stat(darwinArtifact); err != nil {
		exit(fmt.Sprintf("%v", err))
	}

	// 1. Create the release with the given tag name
	url := "https://api.github.com/repos/apprenda/kismatic/releases"
	createBody := struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		Draft   bool   `json:"draft"`
	}{
		TagName: *tag,
		Name:    *tag,
		Draft:   true,
	}
	createData, err := json.Marshal(createBody)
	if err != nil {
		exit(fmt.Sprintf("error marshaling body: %s", err))
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(createData))
	if err != nil {
		exit(fmt.Sprintf("error creating request: %s", err))
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		exit(fmt.Sprintf("error doing HTTP request: %v", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		exit(fmt.Sprintf("got unexpected response code from server: %v", resp.Status))
	}

	createResp := struct {
		WebURL    string `json:"html_url"`
		UploadURL string `json:"upload_url"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		exit(fmt.Sprintf("error unmarshaling response: %v", err))
	}

	// 2. Upload the artifacts that are part of the release
	//
	// Gotta clean up the URL we get from GitHub
	// e.g. https://uploads.github.com/repos/alexbrand/kurl/releases/5014109/assets{?name,label}
	uploadURL := createResp.UploadURL[0:strings.Index(createResp.UploadURL, "{")]

	fmt.Println("Uploading linux artifact")
	artifactName := fmt.Sprintf("kismatic-%s-linux-amd64.tar.gz", *tag)
	if err := uploadArtifact(artifactName, linuxArtifact, uploadURL, authToken); err != nil {
		exit(fmt.Sprintf("error uploading linux artifact: %v", err))
	}
	fmt.Println("Uploading darwin artifact")
	artifactName = fmt.Sprintf("kismatic-%s-darwin-amd64.tar.gz", *tag)
	if err := uploadArtifact(artifactName, darwinArtifact, uploadURL, authToken); err != nil {
		exit(fmt.Sprintf("error uploading darwin artifact: %v", err))
	}
	fmt.Println("New release draft is up!")
	fmt.Println(createResp.WebURL)
}

func uploadArtifact(name, file, uploadURL, authToken string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	// Create request
	req, err := http.NewRequest(http.MethodPost, uploadURL, f)
	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-gzip")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))
	req.ContentLength = stat.Size()
	// Issue request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		return fmt.Errorf("got unexpected response code from server: %s", resp.Status)
	}
	return nil
}

func exit(msg string) {
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}
