package jwt_temas

import (
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
