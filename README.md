# AI Studio Monolith (Wails + Go + Clean Architecture)

> **Aplicación monolítica de escritorio de alto rendimiento construida con Go, Wails v2 y Frontend Web, diseñada bajo los principios de Clean Architecture (Ports & Adapters).**

---

## 📋 Descripción

**AI Studio** es una estación de trabajo de escritorio orientada a desarrolladores e inteligencia artificial. Su propósito es brindar una plataforma extensible y desacoplada para interactuar con múltiples modelos de lenguaje (LLMs), orquestar tareas en segundo plano mediante un Message Broker sin congelar la interfaz gráfica, inspeccionar repositorios Git y exponer métricas de telemetría para dashboards de Prometheus y Grafana.

### Principios Arquitectónicos
- **Clean Architecture & Hexagonal (Ports & Adapters):** El núcleo de negocio (`domain` y `usecase`) es 100% independiente de frameworks, APIs de LLMs y motores de renderizado.
- **Patrones de Diseño Implementados:**
  - **Strategy & Factory:** Para intercambiar proveedores LLM (OpenAI, Gemini) en tiempo de ejecución.
  - **Decorator:** Para inyectar observabilidad y conteo de tokens de forma transparente.
  - **Publisher/Subscriber & Worker Pool:** Para el procesamiento asíncrono de tareas en background.
- **Documentación SDD Completa:** Consulta la carpeta [`SDD/`](./SDD) para especificaciones técnicas exhaustivas de arquitectura, puertos, casos de uso e infraestructura.

---

## 🚀 Requisitos Previos

- **Go:** 1.22+ (probado y verificado en Go 1.25).
- **Node.js:** v18+ y npm (para empaquetar el frontend con Vite).
- **Wails CLI v2:** Instalado automáticamente con `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.
- **Sistema Operativo:** Windows 10/11 (utiliza Microsoft Edge WebView2 nativo).

---

## 🛠️ Instrucciones de Ejecución

### 1. Ejecución Rápida (Binario Precompilado)
Si ya compilaste el proyecto, simplemente ejecuta en PowerShell:

```powershell
.\proyecto_ia.exe
```
*(O haz doble clic sobre `proyecto_ia.exe` en el explorador de archivos).*

---

### 2. Modo Desarrollo con Recarga en Vivo (`wails dev`)
Para desarrollar con hot-reload tanto en el código Go como en el frontend:

```powershell
# Asegúrate de que $env:USERPROFILE\go\bin esté en tu PATH o usa:
& "$env:USERPROFILE\go\bin\wails.exe" dev
```

---

### 3. Compilación Oficial para Producción (`wails build`)
Para compilar un binario nativo optimizado con iconos y bindings de Wails:

```powershell
& "$env:USERPROFILE\go\bin\wails.exe" build
```
El ejecutable resultante se ubicará en `build/bin/proyecto_ia.exe`.

---

### 4. Variables de Entorno Opcionales (APIs LLM)
Por defecto, si no configuras claves de API, la aplicación funciona en **modo simulación limpia** para pruebas locales y offline. Si deseas conectar tus cuentas oficiales:

```powershell
$env:OPENAI_API_KEY = "sk-..."
$env:GEMINI_API_KEY = "AIza..."
.\proyecto_ia.exe
```

---

## 💡 Ejemplo de Uso Paso a Paso

Al abrir la aplicación, verás una interfaz oscura con 4 pestañas de navegación:

### Paso 1: Interactuar en el LLM Studio
1. Selecciona la pestaña **LLM Studio**.
2. Elige el proveedor (`OpenAI` o `Google Gemini`) y el modelo deseado (`gpt-4o`, `gemini-1.5-pro`).
3. Escribe un mensaje en la caja inferior (ej. *"Explica la regla de dependencias en Clean Architecture"*) y haz clic en **Enviar Prompt**.
4. Verás la respuesta generada por el adaptador correspondiente a través del puerto de dominio `LLMProvider`.

### Paso 2: Procesamiento Asíncrono en Background (Message Broker)
1. Cambia a la pestaña **Message Broker Tasks**.
2. Haz clic en el botón **+ Encolar Nueva Tarea**.
3. Observarás cómo se crea un nuevo `Job` con estado `pending`, pasa a `running` (procesado por una goroutine del `WorkerPool` en segundo plano) y finaliza en `completed` sin congelar ni un solo milisegundo la interfaz gráfica.

### Paso 3: Inspección de Repositorios Git
1. Ve a la pestaña **Git Repositories**.
2. Ingresa la ruta de tu proyecto (ej. `d:\workspace\proyecto_ia`) y presiona **Inspeccionar Repositorio**.
3. Se listará el estado de los archivos (`Working Tree`) y el historial de commits recientes recuperados mediante el adaptador nativo de Git.

### Paso 4: Monitoreo con Grafana / Prometheus
1. Pestaña **Metrics & Grafana**: Visualiza los KPIs en tiempo real de peticiones acumuladas, tokens procesados, latencia promedio y workers activos.
2. Abre tu navegador web e ingresa a:
   👉 **http://localhost:2112/metrics**
3. Verás las métricas exportadas en formato estándar de Prometheus listas para ser consumidas por un agente de Prometheus o Grafana.

---

## 🧪 Pruebas Automatizadas

Para validar que los contratos de dominio y adaptadores se cumplen:

```powershell
go test -v ./...
```

---

## 📂 Especificaciones Técnicas (Spec-Driven Development)

Para continuar trabajando y extendiendo la arquitectura en un nuevo repositorio Git, revisa la documentación detallada en [`SDD/`](./SDD):

- [SDD-01: Arquitectura del Sistema](./SDD/01_system_architecture_spec.md)
- [SDD-02: Dominio y Puertos](./SDD/02_domain_and_ports_spec.md)
- [SDD-03: Casos de Uso y Flujos](./SDD/03_usecases_and_workflows_spec.md)
- [SDD-04: Infraestructura y Adaptadores](./SDD/04_infrastructure_and_adapters_spec.md)
- [SDD-05: Presentación y Bindings](./SDD/05_presentation_and_wails_bindings_spec.md)
- [SDD-06: Guía de Extensión y Próximos Pasos](./SDD/06_roadmap_and_extension_guide.md)
