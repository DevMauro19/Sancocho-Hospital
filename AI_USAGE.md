# AI_USAGE.md — Hospital de los Costeños con Narcolepsia

Este documento describe con honestidad cómo usamos IA en el proyecto, según lo
pedido en la sección 8.1 del enunciado. Recoge **solo lo que ha pasado hasta
ahora**; cada vez que usemos IA de nuevo, se agrega aquí.

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

**ChatGPT**, como apoyo para revisar la estructura del proyecto,
    diagnosticar problemas de Go y del servidor web, discutir decisiones
    de implementación y preparar documentación y explicaciones técnicas.

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

1. **ID de episodios con un contador global.** La idea inicial de un contador
   global era útil para producir identificadores únicos, y en la implementación
   real se conservó con `atomic.Int64` en `episode.go`. La ventaja es que cada
   episodio recibe un ID reproducible en cada ejecución del programa, mientras la
   atomicidad evita condiciones de carrera si dos goroutines generan episodios
   al mismo tiempo. Este diseño quedó en el paquete `hospital` y no en el
   `Hospital` porque el contador es un detalle técnico del dominio de episodios,
   no de la gestión del hospital.

2. **`Doctor.Attend` mezcla dos responsabilidades.** En la práctica, la opción
   adoptada fue dejar que `Attend` también ejecute `DiagnosePatient` cuando el
   paciente no tenía doctor asignado. Eso mantiene la regla de negocio de que un
   episodio atendido por un doctor queda asociado a ese doctor y la lista de
   pacientes de ese médico se actualiza. El `Orderly` no tiene diagnóstico, porque
   un auxiliar de enfermería mueve a un paciente dormido pero no lo atiende como
   especialista. La política quedó clara: el doctor diagnostica al atender, y el
   auxiliar solo registra el episodio sin reasignar la atención principal.

3. **El episodio guarda una "foto" del cuarto.** La decisión correcta fue
   registrar la habitación en el momento en que se despacha al atendente, no al
   final del flujo. En `Hospital.HandleSleepAttack` se intenta primero asignar la
   habitación, después se despacha al `Attender` y sólo entonces se crea el
   `EpisodeRecord`. Así, si no había cama disponible, el campo `room` queda en
   `nil`, y si sí la había, el episodio conserva la habitación exacta donde
   ocurrió el episodio.

4. **Suposiciones sobre archivos que la IA no había visto.** El borrador inicial
   de la IA usó campos que parecían razonables, pero el modelo real del proyecto
   usa nombres privados y una estructura más estricta: `Person` tiene `id`, `name`
   y `age`; `Patient` añade `state`, `level`, `currentLocation` y
   `assignedRoom`; `EpisodeRecord` guarda `id`, `dateTime`, `patient`, `attender`,
   `location` y `room`. La corrección fue ajustar el diseño al código real y no
   al pseudocódigo de un ejemplo aislado.

5. **La confirmación real debe ser la ejecución del proyecto.** La IA ayudó a
   plantear la estructura y la lógica, pero la verificación definitiva se hizo con
   el código real del repositorio: `gofmt`, `go vet`, `go test ./...` y la
   ejecución de la interfaz web con `go run ./cmd/web`, probando el endpoint
   `/api/state` en el puerto configurado por la variable `PORT` (predeterminado
   9090).

## 4. Qué aprendió cada integrante

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


- **[Santiago Perez]:** Me encargué de `hospital/room.go`. `Room` representa una
  habitación con número, capacidad, estado (`RoomState`) y un slice
  `occupants []*Patient`. Todos los campos son privados, así que desde afuera
  solo se usa a través de sus métodos: `NewRoom` crea el cuarto y rechaza una
  capacidad menor o igual a 0; `IsAvailable` compara `len(occupants)` con la
  capacidad; `Occupy` agrega al paciente con `append` si hay cama libre y no
  está repetido; `Release` lo saca con `append(occupants[:i], occupants[i+1:]...)`.
  Después de cada cambio se llama a `updateState`, que es el único lugar donde
  cambia el estado, para que el cuarto nunca quede "Occupied" con camas libres.
  `Room` solo maneja el cuarto: el estado del paciente lo cambia el `Hospital`.
  Con esto aprendí cómo funcionan los slices en Go: por dentro apuntan a un
  arreglo, por eso `Occupants()` devuelve una copia con `copy` y así nadie puede
  modificar el slice interno desde afuera. También aprendí a usar `iota` para
  definir estados como constantes tipadas (`Available`, `Occupied`) con su
  método `String()` en vez de usar strings sueltos. Además repasé cómo funcionan
  los arreglos de punteros: como `occupants` guarda `*Patient`, el cuarto apunta
  a los mismos pacientes del hospital y no a copias, y por eso `contains` y
  `Release` pueden encontrar al paciente comparando punteros.


- **[Pablo Escudero]:** Con el attender se maneja el despacho del polimorfismo, 
    nuestro código no tiene concurrencia por lo que lo ideal sería que cuando goroutine
    llama para asignar una habitación, que espere o que salga una excepción de habitación
    no disponible.


## 5. Verificación hecha con ayuda de la IA

- Validación del archivo de diseño con un esquema mínimo antes de integrar la
  solución real, para comprobar que el patrón de `Attender`, `Doctor`,
  `Orderly`, `Patient` y `EpisodeRecord` era coherente.
- Verificación real del proyecto actual:
  - `gofmt -w .` en la raíz y en `cmd/web` para dejar todos los archivos con
    el formato estándar de Go.
  - `go vet ./...` sobre el módulo, sin advertencias relevantes.
  - `go test ./...`, resultado correcto: todas las pruebas del paquete
    `hospital` pasan.
  - `go run ./cmd/web` con la variable de entorno `PORT=9090` para levantar la
    API web y comprobar que `/api/state` responde con JSON válido.
  - Prueba del endpoint `http://127.0.0.1:9090/api/state`, que devolvió el
    estado inicial del hospital sin errores de tipo `null` en `hallway` ni en
    `episodes`.

La comprobación final del proyecto no se basa únicamente en el esqueleto mínimo:
la evidencia real proviene del módulo completo y de la ejecución de la API web
sobre el código que quedó en el repositorio.