package BotFunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pterm/pterm"
)

type RunCodeResponce struct {
	StandardOutput      string `json:"stdout"`
	Error               string `json:"error"`
	StandardOutputError string `json:"stderr"`
}

func RunGoAndPythonCodeByWebAPI(
	language string,
	code string,
) (RunCodeResponce, error) {
	defer func() {
		err := recover()
		if err != nil {
			pterm.Warning.Printf("RunGoAndPythonCodeByWebAPI: %v\n", err)
		}
	}()
	// if crashed
	var fileName string = "main.go"
	var url string = `https://glot.io/api/run/go/latest`
	if language == "python" {
		fileName = "main.py"
		url = `https://glot.io/api/run/python/latest`
	}
	// sure the file name and the current url
	single := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		Name:    fileName,
		Content: code,
	}
	part, _ := json.Marshal(single)
	data := fmt.Sprintf(`{"files": [%v]}`, string(part))
	// init json datas
	newRequest, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return RunCodeResponce{}, fmt.Errorf("RunGoAndPythonCodeByWebAPI: %v", err)
	}
	// get new post request
	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.Header.Add("Authorization", "Token a917b67e-f68d-46b1-8b23-ff2513bf7065")
	// set headers
	resp, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		return RunCodeResponce{}, fmt.Errorf("RunGoAndPythonCodeByWebAPI: %v", err)
	}
	// send request
	buffer := bytes.NewBuffer([]byte{})
	io.Copy(buffer, resp.Body)
	// get buffer
	var ans RunCodeResponce
	err = json.Unmarshal(buffer.Bytes(), &ans)
	if err != nil {
		return RunCodeResponce{}, fmt.Errorf("RunGoAndPythonCodeByWebAPI: %v", err)
	}
	// get responce
	return ans, nil
	// return
}
