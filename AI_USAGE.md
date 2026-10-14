# AI_USAGE.md — Hospital de los Costeños con Narcolepsia

Este documento describe con honestidad cómo usamos IA en el proyecto, según lo
pedido en la sección 8.1 del enunciado. Recoge **solo lo que ha pasado hasta
ahora**; cada vez que usemos IA de nuevo, se agrega aquí.

> Los textos entre `[COMPLETAR: ...]` los tiene que escribir el equipo con sus
> propias palabras. No dejar ninguno al entregar: este archivo se contrasta con
> lo que cada integrante diga en la defensa.

## 1. Herramientas usadas y para qué

**Integrantes:** [COMPLETAR: nombre 1], [COMPLETAR: nombre 2], [COMPLETAR: nombre 3]

**Claude (Anthropic), modelo Sonnet 5.5** [COMPLETAR: confirmar el modelo y quién
tuvo la sesión], en una sesión de chat, usado como tutor y asistente de diseño
para:

- Leer el enunciado y proponer una **ruta de trabajo**: preparar el repositorio,
  orden de implementación de los archivos (tipos enumerados → `Person`/`Patient`
  → `Room` → `EpisodeRecord` → `Doctor`/`Orderly` → `Hospital` → tests → `main`)
  y errores comunes que evitar.
- Discutir una decisión de diseño: **dónde debe vivir `EpisodeRecord`**. La
  conclusión fue guardarlo en `Hospital.history` (fuente de verdad global) y en
  `Doctor.episodes` (que el enunciado exige para `MyEpisodes()`), y no dentro
  de `Patient`, para no duplicar datos en tres sitios.
- Redactar un **primer borrador de `hospital/doctor.go`** (`Doctor`,
  `NewDoctor`, `DiagnosePatient`, `Attend`, `AttendSleepEmergency`,
  `MyEpisodes`) y un borrador de `SevereReport()`.
- Redactar el `README.md` y este `AI_USAGE.md`.

**Qué NO hizo la IA (hasta ahora):** `Hospital`, `Room`, `Orderly`, los tests,
`main.go` y la interfaz gráfica. [COMPLETAR: si alguien usó otra herramienta
—Copilot, ChatGPT, etc.— indicar cuál y para qué archivos.]

## 2. Prompts decisivos

1. **"me podrías decir cómo puedo empezar con este proyecto"** (adjuntando el
   PDF del enunciado) — produjo la ruta de trabajo y nos hizo ver desde el
   inicio las partes que más pesan en la nota: `Person` embebido, un
   `[]Attender` usado de forma polimórfica, consultas que **devuelven** datos
   en vez de imprimir, y `error` en lugar de `panic` cuando no hay cuarto.

2. **"El EpisodeRecord debería estar dentro del paciente, no?"** — en vez de
   darnos la razón, la IA explicó por qué no: tres copias del mismo episodio
   obligan a mantener tres sitios sincronizados, y la consulta 5.4
   (`SevereReport`) se resuelve cruzando `h.patients` con `h.history`. Esto
   definió dónde se guardan los episodios.

3. **"muéstrame cómo quedaría el doctor"** — dio el borrador de `doctor.go` con
   la justificación de cada decisión (embeber `Person` por valor, receptor
   puntero, `MyEpisodes` devolviendo una copia, campo `attender` de tipo
   `Attender` en `EpisodeRecord`).

## 3. Casos donde la salida de la IA fue incorrecta o subóptima

1. **ID de episodios con un contador global.** El borrador usó una variable de
   paquete (`episodeCounter`, un `atomic.Int64`) para generar `E-001`, `E-002`…
   Funciona, pero es estado global mutable: los tests comparten el contador, los
   IDs dependen del orden en que corran y no son reproducibles.
   **Corrección:** [COMPLETAR: pasar el contador a un campo de `Hospital`, o
   explicar por qué se dejó así].

2. **`Doctor.Attend` mezcla dos responsabilidades.** El borrador llama a
   `DiagnosePatient` dentro de `Attend`, así que atender un episodio también
   asigna al paciente como "bajo cuidado" de ese doctor. Además falla si el
   paciente ya tiene *otro* doctor, lo que puede hacer que el hospital, al
   despachar por `[]Attender`, descarte a un doctor válido. **Corrección:**
   [COMPLETAR: separar ambas acciones, o definir una política clara para
   pacientes con otro doctor asignado].

3. **El episodio guarda una "foto" del cuarto.** `EpisodeRecord.room` se llena
   con `p.room` en el momento de atender. Si el paciente estaba en el pasillo
   (`room == nil`) y luego recibe cuarto, el episodio ya guardado sigue sin
   cuarto. La IA lo señaló como límite, pero no lo resolvió. Esto afecta el
   paso 3 del escenario (el paciente del pasillo recibe el cuarto liberado).
   **Corrección:** [COMPLETAR: decisión tomada].

4. **Suposiciones sobre archivos que la IA no había visto.** El borrador de
   `doctor.go` asume que `Patient` tiene los campos `room`, `currentLocation` y
   `assignedDoctor`, y que `EpisodeRecord` usa campos como `dateTime` y
   `attender`. Son nombres inventados por la IA, no definidos por nosotros.
   **Corrección:** [COMPLETAR: ajustar los nombres al modelo real del grupo].

5. **El código llegó sin compilar.** La IA entregó `doctor.go` sin haberlo
   compilado. Después se probó en un entorno aparte con Go 1.22.2, usando un
   esqueleto mínimo de `Person`, `Patient`, `Room`, `EpisodeRecord` y `Attender`
   escrito solo para esa prueba: `gofmt -l` sin salida, `go vet ./...` sin
   advertencias y un test de `Attend`/`MyEpisodes` en verde. Eso **no**
   confirma que funcione con nuestros archivos reales. La confirmación real es
   la que haga el equipo con `go vet ./...` y `go test ./...` en su máquina.

## 4. Qué aprendió cada integrante

Cada uno lo escribe con sus palabras. Para que sea útil en la defensa, conviene
que respondan (sin copiar la respuesta de la IA) a las preguntas de la
sección 8.2 del enunciado:

- **[COMPLETAR: nombre 1]:** [COMPLETAR: ¿qué cambia si `Person` se embebe por
  valor o por puntero?, ¿por qué los métodos de `Doctor` usan receptor puntero?]
- **[COMPLETAR: nombre 2]:** [COMPLETAR: ¿por qué los campos van en minúscula y
  qué impide exactamente que otro paquete los toque?, ¿por qué `MyEpisodes`
  devuelve una copia?]
- **[COMPLETAR: nombre 3]:** [COMPLETAR: ¿cómo se usa el `[]Attender` sin saber
  el tipo concreto?, ¿qué pasa si dos goroutines asignan el mismo cuarto a la
  vez?]

## 5. Verificación hecha con ayuda de la IA

- Compilación y prueba del borrador de `doctor.go` con un esqueleto mínimo, tal
  como se describe en el caso 5 de la sección 3.
- [COMPLETAR: verificaciones sobre el proyecto real —`gofmt`, `go vet`,
  `go test ./...`, `go test -race ./...`, ejecución de `go run .` desde un clon
  limpio— con su resultado.]