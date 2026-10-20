# AI_USAGE.md — Hospital de los Costeños con Narcolepsia

Este documento describe con honestidad cómo usamos IA en el proyecto, según lo
pedido en la sección 8.1 del enunciado. Recoge **solo lo que ha pasado hasta
ahora**; cada vez que usemos IA de nuevo, se agrega aquí.

> Los textos entre `[COMPLETAR: ...]` los tiene que escribir el equipo con sus
> propias palabras. No dejar ninguno al entregar: este archivo se contrasta con
> lo que cada integrante diga en la defensa.

## 1. Herramientas usadas y para qué

**Integrantes:** Carlos Mauricio Velasco, Pablo Escudero, Santiago Perez

**Claude (Anthropic), modelo Sonnet 5.5** , en varias sesiones de chat, usado como tutor y asistente de diseño
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

- **[Carlos Mauricio Velasco]:** Entendí qué cambia al embeber `Person` por valor o
  por puntero. Por valor, los datos viven dentro del `Doctor` y cada uno tiene su
  propia identidad; por puntero, el struct solo guarda una dirección, que queda
  en `nil` si no se inicializa (y `d.ID()` falla en ejecución, no al compilar), y
  dos copias del struct comparten el mismo `Person`. Elegimos valor porque un
  paciente y un doctor son personas distintas. También entendí por qué los
  métodos de `Doctor` usan receptor puntero: `Attend` y `DiagnosePatient` hacen
  `append` sobre `episodes` y `patients`, y con receptor por valor trabajarían
  sobre una copia y esos cambios se perderían sin ningún error de compilación.
  Además, `Patient.assignedDoctor` es un `*Doctor` y se compara con `!=`, lo que
  exige que todos hablen del mismo objeto, y solo `*Doctor` cumple `Attender`.


- **[Pablo Escudero]:** Entendí que un campo en minúscula es *no
  exportado*: en Go la visibilidad se define por **paquete**, y el compilador
  rechaza cualquier acceso desde otro paquete (por ejemplo, `p.state` desde
  `main` da el error "cannot refer to unexported field"). No es una protección
  en ejecución, y tampoco es por tipo: cualquier archivo del paquete `hospital`
  sí puede tocar esos campos, por eso la disciplina es modificarlos solo a
  través de métodos que mantengan coherente el estado (por ejemplo, que
  `state` y `room` del paciente cambien juntos). `MyEpisodes` devuelve una copia
  porque un slice es una referencia a un arreglo interno: si devolviéramos
  `d.episodes`, el código externo podría sobrescribir sus elementos y alterar el
  historial del doctor. La copia (`make` + `copy`) lo evita, aunque es superficial:
  los registros copiados siguen apuntando a los mismos `*Patient` y `*Room`.


- **[Santiago Perez]:** Entendí que el hospital guarda su personal en un
  `[]Attender`, y cada elemento es un valor de interfaz que lleva dentro el tipo
  concreto (`*Doctor`, `*Orderly`…). El hospital recorre el slice y llama a
  `a.Attend(p, ubicación)` sin saber de qué tipo es; Go decide en ejecución qué
  método corresponde. Por eso agregar un `Nurse` solo exige implementar
  `ID()`, `Name()` y `Attend()`, sin tocar el hospital. Solo si hace falta algo
  propio de un tipo (como `DiagnosePatient`, que el camillero no tiene) se usa
  una aserción de tipo (`a.(*Doctor)`). Sobre las goroutines: nuestra versión
  base no tiene concurrencia, pero si dos goroutines llamaran a `AssignRoom` a
  la vez, ambas podrían ver el mismo cuarto como disponible antes de que
  cualquiera lo ocupe (revisar y luego actuar no es atómico) y asignarlo a dos
  pacientes; es una condición de carrera que `go run -race` reportaría. La
  solución es un `sync.Mutex` en `Hospital` que cubra todo `AssignRoom` (buscar
  y ocupar), de modo que la segunda goroutine reciba otro cuarto o el error de
  "no hay cuarto". Además, `Room.Occupy` debe devolver un error si el cuarto
  ya está ocupado, como segunda defensa.

## 5. Verificación hecha con ayuda de la IA

- Compilación y prueba del borrador de `doctor.go` con un esqueleto mínimo, tal
  como se describe en el caso 5 de la sección 3.
- [COMPLETAR: verificaciones sobre el proyecto real —`gofmt`, `go vet`,
  `go test ./...`, `go test -race ./...`, ejecución de `go run .` desde un clon
  limpio— con su resultado.]