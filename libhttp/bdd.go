package libhttp

import (
	"fmt"
	"log"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// initDB инициализирует соединение с базой данных
func (lt *lthttp) initDB(config interface{}) {
	// Используем reflection для извлечения полей из переданной структуры
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

	// Проверяем наличие всех необходимых полей
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

	// Автоматическая миграция структуры EndpointSave
	lt.DB.AutoMigrate(&EndpointSave{})
}

type EndpointSave struct {
	Url         string `gorm:"primaryKey;size:100"`
	Route       string `gorm:"size:100"`
	Type        string `gorm:"size:15"`
	Project     string `gorm:"size:100"`
	Description string `gorm:"size:300"`
}

// SaveEndpointLog сохраняет лог endpoint в базу данных
func (lt *lthttp) SaveEndpointLog(url, route, tipo, project, desc string) {
	var logEntry EndpointSave
	result := lt.DB.First(&logEntry, "url = ?", url)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Запись не найдена, создаем новую
			logEntry = EndpointSave{
				Url:         url,
				Route:       route,
				Type:        tipo,
				Project:     project,
				Description: desc,
			}
			if err := lt.DB.Create(&logEntry).Error; err != nil {
				log.Printf("Error saving new endpoint log: %v", err)
			}
		} else {
			log.Printf("Error querying endpoint log: %v", result.Error)
		}
	} else {
		// Запись найдена, обновляем
		logEntry.Route = route
		logEntry.Type = tipo
		logEntry.Project = project
		logEntry.Description = desc
		if err := lt.DB.Save(&logEntry).Error; err != nil {
			log.Printf("Error updating endpoint log: %v", err)
		}
	}
}
