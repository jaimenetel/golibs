package dbmanager

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBManagerMlt struct {
	db  *gorm.DB
	dsn string
}

type ManagerRegistry struct {
	managers map[string]*DBManagerMlt
	sync.Mutex
}

var (
	instanceMlt *ManagerRegistry
	onceMlt     sync.Once
)

// GetInstance retorna la única instancia del registro de DBManager.
func GetInstanceMlt() *ManagerRegistry {
	onceMlt.Do(func() {
		instanceMlt = &ManagerRegistry{
			managers: make(map[string]*DBManagerMlt),
		}
	})
	return instanceMlt
}

// AddDBManager añade o actualiza una instancia de DBManager en el registro, basada en el nombre.
func (r *ManagerRegistry) AddDBManager(nombre string) {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.managers[nombre]; exists {
		log.Printf("DBManager for '%s' already exists. Reusing the existing connection.", nombre)
		return
	}

	//viper.SetConfigName("config") // ajusta a tu archivo de configuración
	//viper.AddConfigPath(".")      // ajusta la ruta de tu archivo de configuración
	//err := viper.ReadInConfig()
	//if err != nil {
	//	log.Fatalf("Error reading config file: %v", err)
	//}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString(fmt.Sprintf("%s.user", nombre)),
		viper.GetString(fmt.Sprintf("%s.pass", nombre)),
		viper.GetString(fmt.Sprintf("%s.host", nombre)),
		viper.GetString(fmt.Sprintf("%s.port", nombre)),
		viper.GetString(fmt.Sprintf("%s.database", nombre)),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error opening database for '%s': %v", nombre, err)
	}

	r.managers[nombre] = &DBManagerMlt{db: db, dsn: dsn}
}

// GetDBManager retorna una instancia de DBManager por nombre del registro.
func (r *ManagerRegistry) GetDBManager(nombre string) *DBManagerMlt {
	r.Lock()
	defer r.Unlock()

	if manager, exists := r.managers[nombre]; exists {
		return manager
	}

	log.Fatalf("DBManager named '%s' does not exist.", nombre)
	return nil // Este punto no se alcanza debido al log.Fatalf anterior.
}

func init() {
	registry := GetInstanceMlt()
	registry.AddDBManager("screenfain")
	registry.AddDBManager("gateway")
	GetInstanceMlt().GetDBManager("screenfain").db.AutoMigrate(&FileEntry{})
	GetInstanceMlt().GetDBManager("screenfain").db.AutoMigrate(&FilesToSend{})
}
func Nomain() {
	//	registry := GetInstance()
	//	registry.AddDBManager("nombreBaseDeDatos") // Asume configuración para "nombreBaseDeDatos"

	//	dbManager := registry.GetDBManager("nombreBaseDeDatos")
	//	db := dbManager.db // Ahora puedes usar db para realizar operaciones en la base de datos.
}
