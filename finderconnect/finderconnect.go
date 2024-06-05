package finderconnect

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var URLgetIp string = "http://172.17.0.56:8701/findip?find=%s"
var URLgetLtm string = "http://172.17.0.56:8701/findltm?find=%s"
var URLgetDisp string = "http://172.17.0.56:8701/finddispositivo?find=%s"

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

type Finder struct {
}

func (f *Finder) GetIp(find string) (string, error) {
	URL := fmt.Sprintf(URLgetIp, find)
	result, err := FetchURL(URL)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (f *Finder) GetLTM(find string) (string, error) {
	URL := fmt.Sprintf(URLgetLtm, find)
	result, err := FetchURL(URL)
	if err != nil {
		return "", err
	}
	return result, nil
}
func (f *Finder) GetDisp(find string) (string, error) {
	URL := fmt.Sprintf(URLgetDisp, find)
	result, err := FetchURL(URL)
	if err != nil {
		return "", err
	}
	return result, nil
}
