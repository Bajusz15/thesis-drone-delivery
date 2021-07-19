package json

import (
	"bytes"
	"drone-delivery/drone-swarm/pkg/config"
	"drone-delivery/server/pkg/domain/models"
	"drone-delivery/server/pkg/network/inbound/http/rest"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Adapter struct {
}

func NewOutBoundAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) SendTelemetryDataToServer(t models.Telemetry) error {
	td := rest.TelemetryData{
		Telemetry: t,
	}
	postBody, _ := json.Marshal(td)
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(config.ServerHTTPDomain+":"+config.ServerHTTPPort+"/api/delivery/telemetry", "application/json", responseBody)
	//Handle Error
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return errors.New("request failed, response status is " + strconv.Itoa(resp.StatusCode))
	}
	//Read the response body
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//sb := string(body)
	return nil
}
