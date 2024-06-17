package respondwithjson

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type JsonResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Constructor para JsonResponse
func NewJsonResponse(message string, data interface{}, err string) JsonResponse {
	return JsonResponse{
		Message: message,
		Data:    data,
		Error:   err,
	}
}

// Responder con JSON detallado
func RespondWithJSON(w http.ResponseWriter, statusCode int, response JsonResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Responder con JSON simple
func RespondWithJSONSimple(w http.ResponseWriter, statusCode int, data interface{}) {
	response := NewJsonResponse("", data, "")
	RespondWithJSON(w, statusCode, response)
}

// Verificar y responder con JSON correcto
func CheckAndRespondJSON(w http.ResponseWriter, r *http.Request, object interface{}) bool {
	if r.Body == nil {
		response := NewJsonResponse("No request body provided", nil, "")
		RespondWithJSON(w, http.StatusBadRequest, response)
		return false
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Evita la decodificación si JSON contiene campos que no están en la estructura
	err := decoder.Decode(object)
	if err != nil {
		response := NewJsonResponse("Error parsing JSON", nil, err.Error())
		RespondWithJSON(w, http.StatusBadRequest, response)
		return false
	}

	return true
}

// Esta función obtiene un objeto y devuelve este mismo objeto en formato json, como un string
func GetStructTypes(input interface{}) (string, error) {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typeOfS := val.Type()

	fields := []map[string]string{}
	for i := 0; i < val.NumField(); i++ {
		field := typeOfS.Field(i)
		fieldType := field.Type.String()

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = field.Name
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		fields = append(fields, map[string]string{jsonTag: fieldType})
	}

	fieldTypes := make(map[string]string)
	for _, field := range fields {
		for k, v := range field {
			fieldTypes[k] = v
		}
	}

	jsonData, err := json.MarshalIndent(fieldTypes, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
