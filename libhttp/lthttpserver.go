package libhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	//logrequest "github.com/jaimenetel/golibs/logrequest"
)

// var secretKey string = "64ece9a47243209e7f8739bde3ff17b4ea815c777fe0a4bdfadb889db9900340"

type Endpoint struct {
	Name        string
	Handler     http.HandlerFunc
	Controller  string
	QueryParams string
	Roles       string
	Method      string
}

type lthttp struct {
	Port      string
	Endpoints []Endpoint
	DB        *gorm.DB // Connect con base de datos
}

var instance *lthttp
var oncelt sync.Once

// DatabaseConfig = user, password, host, port, name
func Ltinstance(config interface{}) *lthttp {
	oncelt.Do(func() {
		instance = &lthttp{}
		instance.initDB(config) // Inicializar la bdd
	})
	return instance
}

// Método para añadir endpoints y roles a lthttp. "args" --> (roles, method) Requiere: nombre, controller, queryParams(o nil), args(OPCIONAL) (roles, method)
func (lt *lthttp) AddEndpoint(name string, handler http.HandlerFunc, queryParams string, args ...string) {

	roles, method := ParseRolesAndMethod(args...)
	controllerName := GetFunctionName(handler)

	if queryParams == "" {
		queryParams = "none"
	}

	endpoint := Endpoint{
		Name:        name,
		Handler:     handler,
		Controller:  controllerName,
		QueryParams: queryParams,
		Roles:       roles,
		Method:      method,
	}

	lt.Endpoints = append(lt.Endpoints, endpoint)
}

// "args" --> (roles, method) Requiere: nombre, controller, prehandler, queryParams(o nil), args(OPCIONAL) (roles, method)
func (lt *lthttp) AddEndpointPreHandler(name string, handler http.HandlerFunc, prehandler func(http.HandlerFunc) http.HandlerFunc, queryParams string, args ...string) {

	// El prehandler es un middleware que toma y devuelve un http.HandlerFunc
	ohandler := prehandler(handler)

	roles, method := ParseRolesAndMethod(args...)
	controllerName := GetFunctionName(handler)

	if queryParams == "" {
		queryParams = "none"
	}

	endpoint := Endpoint{
		Name:        name,
		Handler:     ohandler,
		Controller:  controllerName,
		QueryParams: queryParams,
		Roles:       roles,
		Method:      method,
	}

	lt.Endpoints = append(lt.Endpoints, endpoint)
}

func (lt *lthttp) StartSinCOrs() {
	for _, endpoint := range lt.Endpoints {
		fmt.Println(endpoint)

		http.Handle(endpoint.Name, authMiddlewareRoleLog(endpoint.Handler, endpoint.Roles))

		if endpoint.Method == "" || (endpoint.Method != "POST" && endpoint.Method != "GET") {
			endpoint.Method = "POST"
		}

		// Guardar endpoint en bdd
		lt.SaveEndpointLog(endpoint)
	}
}
func (lt *lthttp) Start() {
	for _, endpoint := range lt.Endpoints {
		fmt.Println(endpoint)

		// Envuelve el handler original con los middlewares de auth y log, y luego con el CORS middleware
		handlerWithMiddleware := corsMiddleware(authMiddlewareRoleLog(endpoint.Handler, endpoint.Roles))

		if endpoint.Method == "" || (endpoint.Method != "POST" && endpoint.Method != "GET") {
			endpoint.Method = "POST"
		}

		// Guardar endpoint en bdd
		lt.SaveEndpointLog(endpoint)

		handlerWithMiddleware = ConfigMethodType(handlerWithMiddleware, endpoint.Method)

		http.Handle(endpoint.Name, handlerWithMiddleware)
	}
}

func (lt *lthttp) StartRenovado() {
	for _, endpoint := range lt.Endpoints {
		fmt.Println(endpoint)
		// Encadena los middlewares aquí
		handler := authMiddlewareRoleLog(endpoint.Handler, endpoint.Roles)
		//handler = additionalChecksMiddleware(handler)
		http.Handle(endpoint.Name, handler)
	}
}

func prehandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Aquí puedes añadir tus comprobaciones adicionales
		fmt.Println("Ejecutando comprobaciones adicionales")

		// Llama al siguiente middleware o al manejador final
		next.ServeHTTP(w, r)
	})
}

func IsIn(cadena string, slice []string) bool {
	for _, valor := range slice {
		if valor == cadena {
			return true
		}
	}
	return false
}
func IsInIsIn(cadenas []string, slice []string) bool {
	for _, valor := range cadenas {
		if IsIn(valor, slice) {
			return true
		}
	}
	return false
}

func CheckRoles(roles string, slice string) bool {
	sroles := strings.Split(roles, ",")
	proles := strings.Split(slice, ",")
	fmt.Println(proles)
	if strings.TrimSpace(slice) != "" {
		fmt.Println("Comprobando roles", len(proles))
		return IsInIsIn(sroles, proles)
	}
	fmt.Println("No hay restriccion")
	return true

}

func requestBodyToString(r *http.Request) (string, error) {
	// Leer el cuerpo del request
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Retorna el error para que el llamador lo maneje.
		return "", err
	}

	// Convertir el cuerpo del request a una cadena
	bodyString := string(bodyBytes)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyString, nil
}

func imprimirDatosSolicitud(w http.ResponseWriter, r *http.Request) {
	// Imprime el método de solicitud (GET, POST, etc.)

	//r.Write(os.Stdout)
	fmt.Printf("Método de solicitud: %s\n", r.Method)

	// Imprime la URL solicitada
	fmt.Printf("URL solicitada: %s\n", r.URL)

	// Imprime los encabezados de la solicitud
	fmt.Println("Encabezados de la solicitud:")
	for nombre, valores := range r.Header {
		for _, valor := range valores {
			fmt.Printf("%s: %s\n", nombre, valor)
		}
	}
	cuerpo, err := requestBodyToString(r)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Error EOF:", err)
		}
		fmt.Println("Error al leer el cuerpo de la solicitud:", err)
		return
	}
	if false {
		fmt.Printf("Cuerpo de la solicitud: %s\n", string(cuerpo))
	}
	// // Lee y muestra el cuerpo de la solicitud (si lo hubiera)
	// cuerpo := make([]byte, 0)
	// for {
	// 	buffer := make([]byte, 1024)
	// 	fmt.Println(buffer)
	// 	n, err := r.Body.Read(buffer)
	// if err != nil {
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	fmt.Println("Error al leer el cuerpo de la solicitud:", err)
	// 	return
	// }
	// 	cuerpo = append(cuerpo, buffer[:n]...)
	// }

	queryParams := r.URL.Query()
	fmt.Println("Parámetros de consulta:")
	for nombre, valores := range queryParams {
		for _, valor := range valores {
			fmt.Printf("%s: %s\n", nombre, valor)
		}
	}
}
func DecodificarJWT(tokenString string) (jwt.MapClaims, error) {
	secretKey := "64ece9a47243209e7f8739bde3ff17b4ea815c777fe0a4bdfadb889db9900340"
	// Decodificar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al decodificar el token: %v", err)
	}

	// Validar el token
	if !token.Valid {
		return nil, fmt.Errorf("token no válido")
	}

	// Obtener el contenido del token (payload)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error al obtener las claims del token")
	}
	fmt.Println(claims)
	return claims, nil
}

func DecodificarJWTVerbose(tokenString string) (jwt.MapClaims, string, time.Time, []string, error) {
	secretKey := "64ece9a47243209e7f8739bde3ff17b4ea815c777fe0a4bdfadb889db9900340"
	// Decodificar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, "", time.Time{}, nil, fmt.Errorf("error al decodificar el token: %v", err)
	}

	// Validar el token
	if !token.Valid {
		return nil, "", time.Time{}, nil, fmt.Errorf("token no válido")
	}

	// Obtener el contenido del token (payload)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", time.Time{}, nil, fmt.Errorf("error al obtener las claims del token")
	}

	// Obtener el usuario y la fecha de expiración del token
	usuario, usuarioOk := claims["sub"].(string)
	if !usuarioOk {
		return nil, "", time.Time{}, nil, fmt.Errorf("error al obtener el usuario del token")
	}

	exp, expOk := claims["exp"].(float64)
	if !expOk {
		return nil, "", time.Time{}, nil, fmt.Errorf("error al obtener la fecha de expiración del token")
	}
	fechaExpiracion := time.Unix(int64(exp), 0)

	// Obtener los roles
	var roles []string
	if rolesClaim, rolesOk := claims["roles"].([]interface{}); rolesOk {
		for _, role := range rolesClaim {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return claims, usuario, fechaExpiracion, roles, nil
}
func CheckAPIKey(r *http.Request) bool {
	apikey := r.Header.Get("X-Api-Key")
	fmt.Println("APIKEY:", apikey)
	if apikey == "" {
		return false
	}
	claims, _ := DecodificarJWT(apikey)
	fmt.Println(claims)
	return true
}
func corsMiddleware(next http.Handler) http.Handler {
	fmt.Println("CORS Middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PrintResponse(w http.ResponseWriter) {
	for key, values := range w.Header() {
		for _, value := range values {
			fmt.Println(key + ": " + value)
		}

	}
}
func authMiddlewareRoleLog(next http.Handler, roles string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obten el token JWT del encabezado de autorización
		//logrequest := logrequest.LogRequest{}
		bearer_string := "Bearer"
		imprimirDatosSolicitud(w, r)
		fmt.Println("tras solicitud", roles)
		//CheckAPIKey(r)
		tokenString := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), bearer_string, "", -1))

		//logrequest.Claims, logrequest.User, logrequest.ValidoHasta, logrequest.Roles, err1 = DecodificarJWTVerbose(tokenString)
		// logrequest.SetRequestDetails(r)
		// json, _ := logrequest.ToJSON()
		// fmt.Println(json)
		//GetMongoSaver().SaveJSON(json, "requestlog")

		if roles != "---" {
			if tokenString == "" {
				// http.Error(w, "Token JWT no proporcionado", http.StatusUnauthorized)
				RespondWithError(w, http.StatusUnauthorized, "Token JWT no proporcionado")
				return
			}
			fmt.Println("con token:", tokenString)

			fmt.Println("token valido")
			// Verifica el rol del usuario
			myClaims, err := DecodificarJWT(tokenString)
			if err != nil {
				// http.Error(w, "Token JWT no válido", http.StatusUnauthorized)
				RespondWithError(w, http.StatusUnauthorized, "Token JWT no válido")
				return
			}

			fmt.Println("tras claims")
			//role := logrequest.Claims["role"].(string)
			role := myClaims["role"].(string)
			fmt.Println("Roles: ", role)
			if roles != "---" {
				if !CheckRoles(role, roles) {
					// http.Error(w, "Acceso no autorizado", http.StatusForbidden)
					RespondWithError(w, http.StatusForbidden, "Acceso no autorizado")
					return
				}
			}
			fmt.Println("tenemos roles")
		}
		// Si el usuario tiene el rol ROLE_ADMIN, permite el acceso al manejador del endpoint
		next.ServeHTTP(w, r)
	})
}

// Ejemplo de función handler
func init() {
	//	PrintCallerInfo()

}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Generar la respuesta en formato JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// ConfigMethodType
func ConfigMethodType(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method && method != "" {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only "+method+" is supported")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ParseMethodAndRoles parses input arguments and returns method and roles. Si todo está vacio, nos añade roles "---" Method por defecto es "POST"
func ParseRolesAndMethod(args ...string) (roles, method string) {

	if len(args) == 0 {
		roles = "---"
	} else if len(args) == 1 {
		roles = args[0]
		if roles == "" {
			roles = "---"
		}
	} else if len(args) == 2 {
		roles = args[0]
		method = args[1]
	}
	return
}
