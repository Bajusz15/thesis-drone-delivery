package outbound

import (
	"bytes"
	"drone-delivery/server/pkg/config"
	"drone-delivery/server/pkg/domain/models"
	"encoding/json"
	"errors"
	"net/http"
)

type JSONAdapter struct{} //just a wrapper

func NewJSONAdapter() *JSONAdapter {
	a := new(JSONAdapter)
	return a

}

type ProvisionData struct {
	Drone     models.Drone
	Warehouse models.Warehouse
}

func (a *JSONAdapter) FetchProvisionDroneEndpoint(wh models.Warehouse, d models.Drone) (success bool, err error) {
	payload := ProvisionData{
		Drone:     d,
		Warehouse: wh,
	}
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(payload)
	resp, err := http.Post(config.DroneSwarmURL+"/provision", "application/json", buf)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		//var r interface{}
		//_ = json.Unmarshal(body, &r)
		//log.Print(body)
		return false, errors.New("failed to start drone")
	}

	return true, nil
}

//func (a *JSONAdapter) GetDrones() ([]models.Drone, error) {
//	resp, err := http.Get(config.DroneSwarmURL + "/drones")
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	var drones []models.Drone
//	json.Unmarshal(body, &drones)
//	return drones, nil
//}
//
//func (a *JSONAdapter) GetParcels() ([]models.Drone, error) {
//	resp, err := http.Get(config.DroneSwarmURL + "/drones")
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	var drones []models.Drone
//	json.Unmarshal(body, &drones)
//	return drones, nil
//}
