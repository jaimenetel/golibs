package numvalidator

import (
	"log"
	"net/mail"

	"github.com/nyaruka/phonenumbers"
)

// isValidPhoneNumber comprueba si el numero de telefono es valido
func IsValidPhoneNumber(phoneNumber string) bool {
	defaultRegion := "ES"

	num, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		log.Println(" ¿Phone number? || Parsing phone number:", err)
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

// CheckPhoneNumberSMSType comprueba numero de telefono y devuelve SMS OR SMSEX
func CheckPhoneNumberSMSType(phoneNumber string) string {
	defaultRegion := "ES"

	num, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		log.Println("¿Phone number? || Parsing phone number:", err)
		return "INVALID"
	}

	if phonenumbers.IsValidNumber(num) {
		regionCode := phonenumbers.GetRegionCodeForNumber(num)
		if regionCode == "ES" {
			return "sms"
		} else {
			return "smsex"
		}
	}
	return "INVALID"
}

// CheckPhoneNumberRegion comprueba el telefono (pais) y devuelve el codigo del pais (es, it, etc...)
func CheckPhoneNumberRegion(phoneNumber string) string {
	defaultRegion := "ES"

	num, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		log.Println("Error al analizar el número de teléfono:", err)
		return "INVALID"
	}

	if phonenumbers.IsValidNumber(num) {
		regionCode := phonenumbers.GetRegionCodeForNumber(num)
		return regionCode
	}
	return "INVALID"
}

// IsValidEmail verifica si el correo electrónico es válido
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

/*
// EJEMPLOS
func main() {
	num := "+34611156612"
	test(num)

	num = "611156612"
	test(num)

	num = "611156612"
	test(num)

	num = "+393491234567" // ITALIA
	test(num)

	// Pruebas de correo electrónico
	email := "test@example.com"
	fmt.Printf("Email: %s, IsValid: %t\n", email, IsValidEmail(email))

	email = "invalid-email"
	fmt.Printf("Email: %s, IsValid: %t\n", email, IsValidEmail(email))
}

func test(num string) {
	p1 := IsValidPhoneNumber(num)
	typeSms := CheckPhoneNumberSMSType(num)
	region := CheckPhoneNumberRegion(num)
	fmt.Println(num, p1, region, typeSms)
}
*/
