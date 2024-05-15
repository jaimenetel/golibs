package jwt_temas

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Handler define una interfaz para manejar acciones
type Handler interface {
	HandleAction(data string) string
}

// LibraryClass es una estructura que usa un Handler para realizar acciones
type KeyGetter struct {
	Handler
}

// PerformAction ejecuta una acción utilizando el Handler proporcionado
func (l *KeyGetter) PerformAction() string {
	fmt.Println("Performing action in library...")
	return l.HandleAction("Important data")
}

func (l *KeyGetter) CreateJWTToken(user string, roles []string) (string, error) {
	secretKey := l.HandleAction("getSecretKey")
	fmt.Println("secretKey: ", secretKey)
	// Crear el token
	token := jwt.New(jwt.SigningMethodHS256)

	// Establecer los claims del token
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user
	claims["sub"] = user
	claims["roles"] = roles
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expira en 72 horas

	// Firmar el token con una clave secreta
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (l *KeyGetter) DecodificarJWT(tokenString string) (jwt.MapClaims, error) {
	secretKey := l.HandleAction("getSecretKey")
	// Decodificar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error al decodificar el token: %v", err)
	}

	// Validar el token
	if !token.Valid {
		return nil, fmt.Errorf("Token no válido")
	}

	// Obtener el contenido del token (payload)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Error al obtener las claims del token")
	}
	fmt.Println(claims)
	return claims, nil
}

func (l *KeyGetter) DecodificarJWTVerbose(tokenString string) (jwt.MapClaims, string, time.Time, []string, error) {
	secretKey := l.HandleAction("getSecretKey")
	//secretKey := "64ece9a47243209e7f8739bde3ff17b4ea815c777fe0a4bdfadb889db9900340"
	// Decodificar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, "", time.Time{}, nil, fmt.Errorf("Error al decodificar el token: %v", err)
	}

	// Validar el token
	if !token.Valid {
		return nil, "", time.Time{}, nil, fmt.Errorf("Token no válido")
	}

	// Obtener el contenido del token (payload)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", time.Time{}, nil, fmt.Errorf("Error al obtener las claims del token")
	}

	// Obtener el usuario y la fecha de expiración del token
	usuario, usuarioOk := claims["sub"].(string)
	if !usuarioOk {
		return nil, "", time.Time{}, nil, fmt.Errorf("Error al obtener el usuario del token")
	}

	exp, expOk := claims["exp"].(float64)
	if !expOk {
		return nil, "", time.Time{}, nil, fmt.Errorf("Error al obtener la fecha de expiración del token")
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
