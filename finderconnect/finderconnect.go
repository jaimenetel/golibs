package finderconnect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var URLgetIp string = "http://172.17.0.56:8701/findip?find=%s"
var URLgetLtm string = "http://172.17.0.56:8701/findltm?find=%s"
var URLgetDisp string = "http://172.17.0.56:8701/finddispositivo?find=%s"

type Findiccid struct {
	Iccid  string `gorm:"column:iccid"`
	Tofind string `gorm:"column:tofind"`
	Tipo   string `gorm:"column:tipo"`
}

type _Cliente struct {
	Iccid   string `gorm:"column:iccid"`
	Cliente string `gorm:"column:cliente"`
	Codigo  string `gorm:"column:codigo"`
	Name    string `gorm:"column:name"`
}

type Dispositivo struct {
	ICCID   string   `json:"iccid"`
	IP      string   `json:"ip"`
	IMEI    string   `json:"imei"`
	Phone   string   `json:"phone"`
	LTM     string   `gorm:"column:ltm"`
	LTC     string   `gorm:"column:ltc"`
	Cliente _Cliente `json:"cliente,omitempty"`
}

func FetchURL(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error al hacer la solicitud GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta fallida con código de estado: %d", resp.StatusCode)
	}

	// Lee el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
	}

	return string(body), nil
}

func GetIp(find string) (string, error) {
	URL := fmt.Sprintf(URLgetIp, find)
	result, err := FetchURL(URL)
	if err != nil {
		return "", err
	}
	return result, nil
}

func GetLTM(find string) ([]Findiccid, error) {
	URL := fmt.Sprintf(URLgetLtm, find)
	result, err := FetchURL(URL)
	if err != nil {
		return []Findiccid{}, err
	}
	fmt.Println("Result:", result)
	var resti []Findiccid

	err = json.Unmarshal([]byte(result), &resti)
	if err != nil {
		return []Findiccid{}, err
	}
	return resti, nil
}
func GetDisp(find string) (Dispositivo, error) {
	URL := fmt.Sprintf(URLgetDisp, find)
	result, err := FetchURL(URL)
	if err != nil {
		return Dispositivo{}, err
	}
	fmt.Println("Result:", result)
	var dispositivo Dispositivo

	err = json.Unmarshal([]byte(result), &dispositivo)
	if err != nil {
		return Dispositivo{}, err
	}
	return dispositivo, nil
}
