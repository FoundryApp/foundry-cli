package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foundry/cli/logger"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	getEnvsCmd = &cobra.Command{
		Use:   "get-envs",
		Short: "Prints all environment variables in your development environment",
		Long:  "",
		Run:   runGetEnvs,
	}
)

func init() {
	rootCmd.AddCommand(getEnvsCmd)
}

func runGetEnvs(cmd *cobra.Command, args []string) {
	url := "https://us-central1-foundryapp.cloudfunctions.net/getUserEnvs"

	reqBody := struct {
		Data interface{} `json:"data"`
	}{struct{}{}}
	marshaledBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.FdebuglnFatal("Error marshaling getenvs request body:", err)
		logger.FatalLogln("Error getting environment variables (1):", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshaledBody))
	if err != nil {
		logger.FdebuglnFatal("Error creating getenvs request:", err)
		logger.FatalLogln("Error getting environment variables (2):", err)
	}

	bearer := "Bearer " + authClient.IDToken
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 30}
	res, err := client.Do(req)
	if err != nil {
		logger.FdebuglnFatal("Error doing getenvs request:", err)
		logger.FatalLogln("Error getting environment variables (3):", err)
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.FdebuglnFatal("Error reading getenvs response body:", err)
		logger.FatalLogln("Error getting environment variables (4):", err)
	}

	type Error struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	var respBody struct {
		Error *Error            `json:"error"`
		Envs  map[string]string `json:"result"`
	}

	err = json.Unmarshal(bodyBytes, &respBody)
	if err != nil {
		logger.FdebuglnFatal("Error unmarshaling getenvs request:", err)
		logger.FatalLogln("Error getting environment variables (5):", err)
	}

	if respBody.Error != nil {
		logger.FdebuglnFatal("Error response getenvs:", respBody.Error)
		logger.FatalLogln("Error getting environment variables (6):", respBody.Error.Message)
	}

	if len(respBody.Envs) == 0 {
		logger.SuccessLogln("No environment variable has been set yet")
	} else {
		logger.SuccessLogln("Following environment variables are set:")
		logger.Logln("")
		for k, v := range respBody.Envs {
			s := fmt.Sprintf("%s=%s\n", k, v)
			logger.Log(s)
		}
	}

}
