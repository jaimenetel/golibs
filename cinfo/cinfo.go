package cinfo

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func NowAsString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func Cinfo(print bool) string {
	pc, file, line, ok := runtime.Caller(1) // El argumento 1 obtiene la información del llamador
	if !ok {
		fmt.Println("No se pudo obtener la información del llamador")
		return ""
	}

	// Obtener los detalles del llamador
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Println("No se pudo obtener la función del llamador")
		return ""
	}

	// Imprimir el nombre de la función y el archivo fuente
	salida := fmt.Sprintf("%s -  %s | %s, Linea: %d", NowAsString(), filepath.Base(fn.Name()), file, line)
	if print {
		fmt.Println(salida)
	}
	return fmt.Sprintf("%s -  %s Archivo: %s, Linea: %d", NowAsString(), filepath.Base(fn.Name()), file, line)
}
