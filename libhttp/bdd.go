package libhttp

import (
	"fmt"
	"log"
	"reflect"
	"runtime"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var PROJECT_NAME string = "---"

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
	lt.DB = db

	// Migración automática de la estructura EndpointSave
	lt.DB.AutoMigrate(&EndpointSave{})
}

type EndpointSave struct {
	ID          uint   `gorm:"column:id;primaryKey"`
	Route       string `gorm:"column:route;size:100"`
	Controller  string `gorm:"column:controller;size:100"`
	Port        string `gorm:"column:port;size:6"`
	Type        string `gorm:"column:type;size:15"`
	Roles       string `gorm:"column:roles;size:100"`
	Project     string `gorm:"column:project;size:100"`
	Description string `gorm:"column:desc;size:300"`
}

func (EndpointSave) TableName() string {
	return "endpoints"
}

// SaveEndpointLog guarda el registro del punto final en la base de datos
func (lt *lthttp) SaveEndpointLog(endpoint Endpoint) {
	logEntry := EndpointSave{
		Route:       endpoint.Name,
		Controller:  endpoint.Controller,
		Port:        lt.Port,
		Type:        endpoint.Method,
		Roles:       endpoint.Roles,
		Project:     PROJECT_NAME,
		Description: "---", // Por defecto
	}

	var existingLog EndpointSave
	result := lt.DB.Where("route = ? AND port = ?", logEntry.Route, logEntry.Port).First(&existingLog)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Registro no encontrado, cree uno nuevo
			if err := lt.DB.Create(&logEntry).Error; err != nil {
				log.Printf("Error creating new endpoint log: %v", err)
			} else {
				log.Println("New endpoint log created successfully")
			}
		} else {
			log.Printf("Error querying endpoint log: %v", result.Error)
		}
	} else {
		// Registro encontrado, actualícelo
		if err := lt.DB.Model(&existingLog).Updates(logEntry).Error; err != nil {
			log.Printf("Error updating endpoint log: %v", err)
		} else {
			log.Println("Endpoint log updated successfully")
		}
	}
}

// Set nombre del proyecto
func (lt *lthttp) SetProjectName(name string) {
	PROJECT_NAME = name
}

// getFunctionName obtener nombre del controller
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
