package libhttp

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

// initDB инициализирует соединение с базой данных
func (lt *lthttp) initDB(config *DatabaseConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		// config.DBName,
		"swagger",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	lt.DB = db

	// Автоматическая миграция структуры EndpointLog
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
