// El paquete "main" es especial: es el punto de entrada de un programa ejecutable.
// Todo lo que está en cmd/<nombre>/ produce un binario con ese nombre (aquí, "api").
package main

import (
	"context"   // para avisar a las tareas en curso que deben terminar
	"errors"    // para comparar errores
	"log/slog"  // logs estructurados de la librería estándar
	"net/http"  // servidor HTTP
	"os"        // variables de entorno y salida del programa
	"os/signal" // capturar Ctrl+C y la señal de apagado de Docker
	"syscall"   // nombres de señales del sistema (SIGTERM)
	"time"      // tiempos de espera

	"github.com/bayalad7/cargo-3d/backend/internal/httpapi"
)

// main es la primera función que se ejecuta.
func main() {
	// Logger que escribe JSON: fácil de leer para Docker y servicios de logs.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// El puerto viene de la variable de entorno PORT (Docker la define).
	// Si no existe, usamos 8080. Nunca escribimos configuración fija en el código.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configuramos el servidor con tiempos límite: sin ellos, un cliente lento
	// podría dejar conexiones abiertas para siempre.
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// "go" lanza la función en una goroutine: una tarea ligera que corre en paralelo.
	// Así el servidor atiende peticiones mientras main espera la señal de apagado.
	go func() {
		logger.Info("servidor iniciado", "port", port)
		// ListenAndServe bloquea hasta que el servidor se detiene.
		// ErrServerClosed es normal al apagar; cualquier otro error es grave.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("el servidor falló", "error", err)
			os.Exit(1)
		}
	}()

	// Esperamos Ctrl+C (SIGINT) o la orden de apagado de Docker (SIGTERM).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer ejecuta stop() cuando main termina, ya sea con return o por un
	// panic. Ojo: os.Exit (arriba) termina el programa SIN correr los defer.
	defer stop()
	// <-ctx.Done() se queda esperando hasta que llegue una de esas señales.
	<-ctx.Done()

	// Apagado ordenado: damos hasta 10 segundos para terminar las peticiones en curso.
	logger.Info("apagando el servidor")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("apagado forzado", "error", err)
	}
}
