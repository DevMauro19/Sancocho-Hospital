# Hospital de los Costeños con Narcolepsia

Proyecto en Go para simular la gestión de pacientes, habitaciones y episodios de sueño en un hospital. El repositorio combina una capa de dominio en `hospital/` con una API REST ligera en `cmd/web/` y una interfaz simple en `cmd/web/static/index.html`.

## Estado actual del proyecto

El proyecto ya incluye:

- lógica del dominio del hospital en Go
- gestión de habitaciones, pacientes, médicos y personal auxiliar
- API REST para consultar el estado del hospital
- registro de ataques de sueño y asignación de camas
- reinicio del escenario
- flujo de alta y eliminación de pacientes personalizados
- interfaz web mínima para operar el hospital desde el navegador
- pruebas automáticas del paquete `hospital`

## Estructura principal

```text
costenitos/
├── go.mod
├── main.go
├── Readme.md
├── AI_USAGE.md
├── cmd/
│   └── web/
│       ├── main.go
│       ├── server.go
│       └── static/
│           └── index.html
├── hospital/
│   ├── doctor.go
│   ├── episode.go
│   ├── hospital.go
│   ├── hospital_test.go
│   ├── orderly.go
│   ├── orderly_test.go
│   ├── patient.go
│   └── room.go
└── out.txt
```

## Cómo ejecutar

### 1) Escenario de consola

Desde la raíz del proyecto:

```bash
go run .
```

Esto ejecuta la demostración del dominio en consola.

### 2) API web

```bash
go run ./cmd/web
```

La aplicación escucha en el puerto `9090` por defecto. Si ese puerto está ocupado, puedes sobrescribirlo con la variable `PORT`:

```bash
$env:PORT = "9091"
go run ./cmd/web
```

En Linux/macOS:

```bash
PORT=9091 go run ./cmd/web
```

Luego abre en el navegador:

```text
http://127.0.0.1:9090
```

o el puerto que hayas configurado.

## API REST disponible

La web app expone estas rutas:

- `GET /api/state` — devuelve el estado completo del hospital
- `POST /api/attack` — registra un ataque de sueño para un paciente
- `POST /api/wake` — despierta a un paciente
- `POST /api/assign` — asigna una habitación disponible
- `POST /api/reset` — reinicia el escenario por defecto
- `POST /api/patients` — agrega un paciente personalizado
- `DELETE /api/patients/{id}` — elimina un paciente admitido

Ejemplo de creación de paciente:

```bash
curl -X POST http://127.0.0.1:9090/api/patients \
  -H "Content-Type: application/json" \
  -d '{"name":"Ana Gómez","age":31,"severity":"Moderate"}'
```

Ejemplo de eliminación:

```bash
curl -X DELETE http://127.0.0.1:9090/api/patients/P-006
```

## Interfaz web

La interfaz sirve para:

- ver pacientes admitidos
- ver habitaciones y ocupación
- revisar pacientes en el pasillo
- enviar ataques de sueño
- despertar pacientes
- asignar camas
- agregar pacientes personalizados
- eliminar pacientes

La UI está en `cmd/web/static/index.html` y se sirve automáticamente desde `cmd/web/main.go`.

## Reglas de negocio y diseño

El dominio incluye:

- `Hospital`: encapsula pacientes, habitaciones, personal y historial
- `Person`: información básica de una persona
- `Patient`: edad, severidad, ubicación, habitación asignada y estado del sueño
- `Doctor` y `Orderly`: implementan la interfaz `Attender`
- `Room`: controla disponibilidad y ocupación de camas
- `EpisodeRecord`: registra los episodios generados por ataques de sueño

El proyecto usa composición, interfaces polimórficas y `error` para manejar situaciones no válidas.

## Pruebas

Para ejecutar todas las pruebas:

```bash
go test ./...
```

En el estado actual, las pruebas del paquete `hospital` pasan correctamente.

## Observaciones importantes

- Si `go run ./cmd/web` falla con `bind: Only one usage of each socket address`, significa que el puerto ya está ocupado por otra instancia del proyecto o por otro proceso.
- La solución recomendada es cerrar el proceso que está usando ese puerto o arrancar la app con otra variable `PORT`.
- La app mantiene los datos en memoria; cuando se reinicia la app o el escenario, el estado vuelve al valor inicial definido por `reset()`.

## Conclusión

Este repositorio ya no es solo un ejercicio de dominio: también funciona como una pequeña aplicación web operativa con datos simulados y capacidad de personalizar la lista de pacientes sin depender únicamente de los casos predefinidos.
