package wails

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App gestiona el ciclo de vida de la aplicación Wails y la emisión de eventos hacia la UI
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	// Cleanup ordenado de recursos
}

// Emit envía eventos asíncronos hacia la ventana del navegador en Wails
func (a *App) Emit(eventName string, data interface{}) {
	if a.ctx != nil {
		wailsRuntime.EventsEmit(a.ctx, eventName, data)
	}
}

// Greet método de verificación de conectividad IPC
func (a *App) Greet(name string) string {
	return "Hola " + name + ", bienvenido a AI Studio (Wails Clean Architecture)!"
}
