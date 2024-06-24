package respondwithjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// JsonResponse es la estructura de la respuesta en formato JSON
type JsonResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Constructor para la respuesta JsonResponse
func NewJsonResponse(message string, data interface{}, err string) JsonResponse {
	return JsonResponse{
		Message: message,
		Data:    data,
		Error:   err,
	}
}

// Responder con el formato JSON
func RespondWithJSON(w http.ResponseWriter, statusCode int, response JsonResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Responder con JSON simple (simplemente data)
func RespondWithJSONSimple(w http.ResponseWriter, statusCode int, data interface{}) {
	response := NewJsonResponse("", data, "")
	RespondWithJSON(w, statusCode, response)
}

// Función para enviar una respuesta exitosa
func RespondWithSuccess(w http.ResponseWriter, data interface{}) {
	response := NewJsonResponse("Success", data, "")
	RespondWithJSON(w, http.StatusOK, response)
}

// Función para enviar una respuesta con el error
func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	var errMsg, message string
	if err != nil {
		errMsg = err.Error()
		message = "ERROR"
	}
	response := NewJsonResponse(message, nil, errMsg)
	RespondWithJSON(w, statusCode, response)
}

// Responder con JSON simple (simplemente data)
func RespondWithJSONMessageError(w http.ResponseWriter, statusCode int, messageError string) {
	response := NewJsonResponse("", "", messageError)
	RespondWithJSON(w, statusCode, response)
}

// Verificar y responder con JSON correcto
func CheckAndRespondJSON(w http.ResponseWriter, r *http.Request, object interface{}) error {
	if r.Body == nil {
		err := errors.New("request body is empty")
		return err
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Evita la decodificación si JSON contiene campos que no están en la estructura
	err := decoder.Decode(object)
	if err != nil {
		return err
	}
	return nil
}

// Esta función obtiene un objeto y devuelve este mismo objeto en formato json, y los tipos de variables del objeto. Por ejemplo: "name": "string"
// Ejemplo de uso: var json := GetStructTypes(ExampleObject{})
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

// Esta función convierte un objeto (o un modelo de objeto: ej. ExampleModel{}) a un formato JSON
func ConvertObjectToJSON(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ValidateFields comprueba que todos los campos pasados ​​no estén vacíos ni contengan espacios. (string, int)
func ValidateFields(fields ...interface{}) error {
	for _, field := range fields {
		value := reflect.ValueOf(field)
		switch value.Kind() {
		case reflect.String:
			str := value.String()
			if strings.TrimSpace(str) == "" || value.IsZero() {
				return fmt.Errorf("fields cannot be empty or contain spaces")
			}
		case reflect.Int:
			if value.Int() == 0 || value.IsZero() {
				return fmt.Errorf("integer fields cannot be zero")
			}
		default:
			return fmt.Errorf("unsupported field type: %s", value.Kind())
		}
	}
	return nil
}
