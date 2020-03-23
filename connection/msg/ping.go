package msg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"foundry/cli/logger"
)

type PingBody struct {
	Token string `json:"token"`
}

type PingMsg struct {
	URL		string
	Body 	PingBody
}

func NewPingMsg(url, t string) *PingMsg {
	return &PingMsg{
		URL: 	url,
		Body: PingBody{t},
	}
}

func (pm *PingMsg) Send() error {
	j, err := json.Marshal(pm.Body)
	if err != nil {
		return err
	}

	res, err := http.Post(pm.URL, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Debugln("<Non-OK ping response> Error reading ping response body: ", err)
			return err
		}

		bodyString := string(bodyBytes)
		logger.Debugln("Non-OK ping response: %s\n", bodyString)
	}

	return nil
}


// func ping(ticker *time.Ticker, token, url string) {
//   for {
//     select {
//     case <-ticker.C:
//       // Ping the server
//       var body = struct {
//         Token string `json:"token"`
//       }{token}

//       jBody, err := json.Marshal(body)
//       if err != nil {
//         logger.Debugln("Error marshaling ping body: ", err)
//         continue
//       }

//       res, err := http.Post(url, "application/json", bytes.NewBuffer(jBody))
//       if err != nil {
//         logger.Debugln("Error making ping post request: ", err)
//         continue
//       }

//       if res.StatusCode != http.StatusOK {
//         bodyBytes, err := ioutil.ReadAll(res.Body)
//         if err != nil {
//           logger.Debugln("Error reading ping response body: ", err)
//           continue
//         }

//         bodyString := string(bodyBytes)
//         logger.Debugln("Non-OK ping response: %s\n", bodyString)
//       }
//     }
//   }
// }