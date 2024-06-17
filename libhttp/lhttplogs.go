package libhttp

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var PROJECT_NAME string = "---"

// Nombre de la tabla
func (EndpointSave) TableName() string {
	return "endpoints"
}

// Obtener la conexión con la base de datos
func (lt *lthttp) SetConnectionDBSwagger(config interface{}) {
	lt.initDB(config)
}

// initDB inicializa una conexión a la base de datos
func (lt *lthttp) initDB(config interface{}) {
	// Usar la reflexión para extraer campos de la estructura pasada
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var dbUser, dbPassword, dbHost, dbPort, dbName string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		switch fieldName {
		case "DBUser":
			dbUser = field.String()
		case "DBPassword":
			dbPassword = field.String()
		case "DBHost":
			dbHost = field.String()
		case "DBPort":
			dbPort = field.String()
		case "DBName":
			// dbName = field.String()
			dbName = "swagger"
		}
	}

	// Comprobando la presencia de todos los campos obligatorios
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" {
		log.Fatalf("Incomplete database configuration")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	lt.DBSwagger = db

	// Migración automática de la estructura EndpointSave
	lt.DBSwagger.AutoMigrate(&EndpointSave{})
}

type EndpointSave struct {
	ID          uint   `gorm:"column:id;primaryKey"`
	Route       string `gorm:"column:route;size:100"`
	Controller  string `gorm:"column:controller;size:100"`
	Port        string `gorm:"column:port;size:6"`
	QueryParams string `gorm:"column:queryparams;size:100"`
	Body        string `gorm:"column:body;size:300"`
	Type        string `gorm:"column:type;size:15"`
	Roles       string `gorm:"column:roles;size:100"`
	Project     string `gorm:"column:project;size:100"`
	Description string `gorm:"column:desc;size:300"`
}

// SaveEndpointLog guarda el registro del punto final en la base de datos
func (lt *lthttp) SaveEndpointLog(endpoint Endpoint) {

	logEntry := EndpointSave{
		Route:       endpoint.Name,
		Controller:  endpoint.Controller,
		Port:        lt.Port,
		QueryParams: endpoint.QueryParams,
		Body:        endpoint.Body,
		Type:        endpoint.Method,
		Roles:       endpoint.Roles,
		Project:     PROJECT_NAME,
		Description: "---", // Por defecto
	}

	var existingLog EndpointSave
	result := lt.DBSwagger.Where("route = ? AND port = ?", logEntry.Route, logEntry.Port).First(&existingLog)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Registro no encontrado, cree uno nuevo
			if err := lt.DBSwagger.Create(&logEntry).Error; err != nil {
				log.Printf("Error creating new endpoint log: %v", err)
			} else {
				log.Println("New endpoint log created successfully")
			}
		} else {
			log.Printf("Error querying endpoint log: %v", result.Error)
		}
	} else {
		// Registro encontrado, verificar identidad
		if !CompareIfEndpointLogsAreSame(existingLog, logEntry) {
			// Los valores son diferentes, actualiza el registro
			if err := lt.DBSwagger.Model(&existingLog).Updates(logEntry).Error; err != nil {
				log.Printf("Error updating endpoint log: %v", err)
			} else {
				log.Println("Endpoint log updated successfully")
			}
		}
	}
}

// Set nombre del proyecto
func (lt *lthttp) SetProjectName(name string) {
	PROJECT_NAME = name
}

// getFunctionName obtener nombre del controller
func GetFunctionName(i interface{}) string {
	//return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}

func CompareIfEndpointLogsAreSame(existingLog, logEntry EndpointSave) bool {
	return existingLog.Route == logEntry.Route &&
		existingLog.Controller == logEntry.Controller &&
		existingLog.Port == logEntry.Port &&
		existingLog.Type == logEntry.Type &&
		existingLog.Roles == logEntry.Roles &&
		existingLog.Project == logEntry.Project &&
		existingLog.Description == logEntry.Description &&
		existingLog.QueryParams == logEntry.QueryParams &&
		existingLog.Body == logEntry.Body

}
