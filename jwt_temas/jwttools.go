package jwt_temas

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var SUPER_SECRET_KEY = "64ece9a47243209e7f8739bde3ff17b4ea815c777fe0a4bdfadb889db9900340"

func getUserName(tokenString string) (string, error) {
	if tokenString == "" {
		return "", nil
	}

	// Limpiar el token para remover 'Bearer' si está presente
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// Definir la clave secreta
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(SUPER_SECRET_KEY), nil
	}

	// Parsear y validar el token
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extraer el usuario del campo 'sub'
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", nil
}

func DecodificarJWT2(tokenString string) (jwt.MapClaims, error) {
	//secretKeySeed := SUPER_SECRET_KEY
	//secretKey := hmac.New(sha256.New, []byte(secretKeySeed))
	// Decodificar el token
	atoken := strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(atoken, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SUPER_SECRET_KEY), nil
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
func GetUserFromBearerToken(token string) string {
	// Split the token by the space character
	name, _ := getUserName(token)
	fmt.Println("name:", name)
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	parts := myClaims["sub"].(string)

	return parts
}
func GetVerifUserFromBearerToken(token string) string {

	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	parts := myClaims["cliuser"].(string)
	if parts == "" {
		parts = myClaims["sub"].(string)
	}

	return parts
}
func GetRolesFromBearerToken(token string) []string {

	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	parts := myClaims["cliuser"].(string)
	if parts == "" {
		parts = myClaims["sub"].(string)
	}

	role := myClaims["role"].(string)
	roles := strings.Split(role, ",")
	return roles
}
func GetRolesFromBearerTokenString(token string) string {

	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	parts := myClaims["cliuser"].(string)
	if parts == "" {
		parts = myClaims["sub"].(string)
	}

	role := "," + myClaims["role"].(string) + ","

	return role
}

// CompareRolesWithToken comprueba si hay intersecciones entre roles de una cadena y una lista de cadenas
func CompareRolesWithToken(roleString, token string) bool {
	roles := GetRolesFromBearerToken(token)
	// Divide la cadena de roles en una lista de cadenas
	roleList := strings.Split(roleString, ",")
	// Comprobar intersección de lista
	for _, role := range roleList {
		for _, r := range roles {
			if strings.TrimSpace(role) == strings.TrimSpace(r) {
				return true
			}
		}
	}
	return false
}

// Funcion recibe (token, user(payload string), roles...) y devuelve bool comprobando si tokenDecodificado.payload == user OOO token.roles == cualquier rol que acabamos de pasar --> True
// CheckUserWithTokenOrRoles checks if the user payload or roles match with the token's claims
// Ejemplo: user = `{"exp":1.112123417e+09,"hostname":"","iat":1.112123417e+09,"iss":"Liftel.es","role":"ROLE_ADMIN","sub":"raul"}`
// var token = "eyJhbGci2iJIUzUxMiJ9.eyJzdWIiOiJSQVVMQUQiLCJyb2xlIjoiUk9MRV9BRE1JTixST0xFX0dFU1RJT04sUk9MRV9NQVRJQyxST0xFX1JFR05FVEVMIiwiaG9zdG5hbWUiOiJQLU1BVElDIiwiaXNzIjoiTGlmdGVsLmVzIiwiaWF0IjoxNzE5NDk3Nzc3LCJleHAiOjE3MTk1ODQxNzd9.5rkX99zpjSWLXVjqOfs4fM-3cAoUSiq2xZ9vj6AUDmpHEu4TFZ8R5PUZxlsnMpkvxkI6CD60HqkIfNvczpWECA"
func CheckUserWithTokenOrRoles(token string, user string, roles ...string) (bool, error) {
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error decoding the token:", err)
		return false, err
	}

	fmt.Println("Decoded Claims:", myClaims)

	// Convert user (string) to jwt.MapClaims
	var userClaims jwt.MapClaims
	err = json.Unmarshal([]byte(user), &userClaims)
	if err != nil {
		fmt.Println("Error deserializing the user (FORMAT JSON):", err)
		// return false, err
	} else {
		fmt.Println("User Claims:", userClaims)

		// Check if user claims match the token claims
		match := true
		for key, value := range userClaims {
			if myClaims[key] != value {
				match = false
				break
			}
		}
		if match {
			return true, nil
		}
	}

	// Check if any of the roles match
	tokenRoles := strings.Split(myClaims["role"].(string), ",")
	for _, role := range roles {
		if contains(tokenRoles, role) {
			return true, nil
		}
	}

	return false, nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
