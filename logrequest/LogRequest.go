package logrequest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"time"

	"github.com/dgrijalva/jwt-go"
	printinfo "github.com/jaimenetel/golibs/printcallerinfo"
)

type HttpRequestDetails struct {
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             io.ReadCloser
	GetBody          func() (io.ReadCloser, error)
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
	Response         *http.Response
	//Context          context.Context
}
type LogRequest struct {
	User        string
	Endpoint    string
	Roles       []string
	ValidoHasta time.Time
	Claims      jwt.MapClaims
	Momento     time.Time
	Request     HttpRequestDetails
}

func (lr *LogRequest) SetRequestDetails(r *http.Request) {
	lr.Request = HttpRequestDetails{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             r.Body,
		GetBody:          r.GetBody,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          r.Trailer,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
		Response:         r.Response,
		//Context:          r.Context,
	}
	lr.Momento = time.Now()
	lr.Endpoint = r.RequestURI
}

type SimplifiedHttpRequestDetails struct {
	Method           string
	URL              string
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             url.Values
	PostForm         url.Values
	RemoteAddr       string
	RequestURI       string
	// Omitir campos complejos o no serializables
}

// ToSimplifiedHttpRequestDetails convierte HttpRequestDetails a una versión simplificada.
func (hrd *HttpRequestDetails) ToSimplifiedHttpRequestDetails() SimplifiedHttpRequestDetails {
	return SimplifiedHttpRequestDetails{
		Method:           hrd.Method,
		URL:              hrd.URL.String(),
		Proto:            hrd.Proto,
		ProtoMajor:       hrd.ProtoMajor,
		ProtoMinor:       hrd.ProtoMinor,
		Header:           hrd.Header,
		ContentLength:    hrd.ContentLength,
		TransferEncoding: hrd.TransferEncoding,
		Close:            hrd.Close,
		Host:             hrd.Host,
		Form:             hrd.Form,
		PostForm:         hrd.PostForm,
		RemoteAddr:       hrd.RemoteAddr,
		RequestURI:       hrd.RequestURI,
	}
}

// ToJSON convierte LogRequest a una cadena JSON.
func (lr *LogRequest) ToJSON() (string, error) {
	// Convertir HttpRequestDetails a una versión simplificada
	simplifiedRequest := lr.Request.ToSimplifiedHttpRequestDetails()

	// Crear una versión de LogRequest con la versión simplificada de HttpRequestDetails
	logRequestForJSON := struct {
		User        string
		Endpoint    string
		Roles       []string
		ValidoHasta time.Time
		Claims      jwt.MapClaims
		Request     SimplifiedHttpRequestDetails
		Momento     time.Time
	}{
		User:        lr.User,
		Endpoint:    lr.Endpoint,
		Roles:       lr.Roles,
		ValidoHasta: lr.ValidoHasta,
		Claims:      lr.Claims,
		Request:     simplifiedRequest,
		Momento:     lr.Momento,
	}

	// Serializar a JSON
	jsonData, err := json.MarshalIndent(logRequestForJSON, "", " ")
	if err != nil {
		return "", fmt.Errorf("error al serializar LogRequest a JSON: %v", err)
	}

	return string(jsonData), nil
}
func init() {
	printinfo.PrintCallerInfo()
}
