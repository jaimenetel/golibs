package printcallerinfo

import (
	"fmt"
	"runtime"
	"time"
)

func NowAsString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func PrintCallerInfo() {
	pc, file, line, ok := runtime.Caller(1) // El argumento 1 obtiene la información del llamador
	if !ok {
		fmt.Println("No se pudo obtener la información del llamador")
		return
	}

	// Obtener los detalles del llamador
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Println("No se pudo obtener la función del llamador")
		return
	}

	// Imprimir el nombre de la función y el archivo fuente
	fmt.Printf("%s - Llamado desde: %s\nArchivo: %s, Linea: %d\n", NowAsString(), fn.Name(), file, line)
}

func init() {
	PrintCallerInfo()

}
