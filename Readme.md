# 🏥 Programación Orientada a Objetos — Ejercicio Aplicado

# 🌴 Hospital de los Costeños con Narcolepsia 😴

| 📌 Información | 📝 Detalles                  |
|---|------------------------------|
| 🎓 Asignatura | Lenguajes y Compiladores     |
| 🏫 Universidad | Universidad EIA              |
| 👨‍🏫 Docente | Sebastián Zapata Ramírez     |
| 💻 Lenguaje | Go (Golang) 1.21 o superior  |
| 👥 Integrantes por grupo | Máximo 3 estudiantes         |
| 📅 Fecha de entrega | 19 de octubre de 2026, 23:59 |
| 🎤 Defensa oral | 20 de octubre de 2026        |
| 📊 Porcentaje de la nota final | 15 %                         |

---

## 🌞 1. Introducción: El origen del nombre

Cuenta la leyenda que la brisa cálida de la costa, combinada con un buen sancocho de pescado, produce un efecto de somnolencia imparable. Pero ¿qué sucede cuando el espíritu alegre, ruidoso y enérgico de un costeño colombiano se encuentra con la narcolepsia clínica?

¡Un caos absoluto e impredecible! 😂

En un momento están bailando champeta a todo volumen o discutiendo apasionadamente sobre béisbol y, al siguiente, se quedan dormidos de pie, apoyados contra una palmera. 🌴💤

El **Hospital de los Costeños con Narcolepsia** nació para resolver esta paradoja tan particular: pacientes cuyas almas quieren seguir de fiesta hasta las 6:00 a. m., pero cuyos sistemas nerviosos se apagan a las 2:15 p. m., ¡justo en medio de una conversación!

Este proyecto busca poner orden digital en ese carnaval somnoliento mediante la implementación de un sistema básico de gestión hospitalaria en Go.

### 💤 Los dos estados naturales del paciente

- 🍲 Comer sancocho.
- 🏥 Llegar al hospital.

## 🎯 2. Objetivo del proyecto

El sistema debe permitir realizar un seguimiento de la ubicación de cada paciente, identificar qué médico lo está atendiendo y registrar cada episodio repentino de sueño.

De esta manera, será posible asignar rápidamente una cama disponible antes de que el paciente termine bloqueando un pasillo. 🛏️😴

El entregable final será un programa en Go cuya función `main` ejecute el escenario de demostración e imprima los resultados de las consultas requeridas.

## ⚙️ 3. Objetivos y reglas técnicas

- 🐹 El proyecto debe estar desarrollado completamente en Go (Golang), versión 1.21 o superior.
- 📚 Solo se permite utilizar la biblioteca estándar de Go, a menos que se implemente la interfaz gráfica opcional.
- 🧱 El código debe reflejar principios de programación orientada a objetos mediante estructuras (`structs`), métodos, composición e interfaces.
- 🔒 La encapsulación se implementa mediante el sistema de paquetes de Go. Los campos deben ser privados (no exportados) y accederse a través de métodos.
- ⚠️ Los métodos que puedan fallar deben devolver un `error` como último valor de retorno.
- ✨ El código debe estar formateado con `gofmt` y superar las comprobaciones de `go vet` sin advertencias.

## 🏗️ 4. Modelo de entidades

El repositorio implementa las siguientes entidades principales:

| 🧩 Entidad | 📋 Responsabilidad principal |
|---|---|
| 🏥 `Hospital` | Gestiona médicos, pacientes, habitaciones e historial de episodios de sueño. |
| 🧑 `Person` | Representa los datos comunes de una persona. |
| 🛌 `Patient` | Almacena información del paciente, su estado y su ubicación actual. |
| 🩺 `Doctor` | Atiende pacientes y conserva el registro de sus propios episodios. |
| 🛏️ `Room` | Gestiona la disponibilidad y ocupación de las habitaciones. |
| 📝 `EpisodeRecord` | Registra cada episodio de sueño y la acción realizada. |

### 📌 Responsabilidades principales

- **🏥 `Hospital`:** administra los médicos, pacientes, habitaciones e historial de episodios.
- **🛌 `Patient`:** almacena datos como el identificador, nombre, edad, nivel de narcolepsia, estado y ubicación actual.
- **🩺 `Doctor`:** atiende a los pacientes y registra sus propios episodios.
- **🛏️ `Room`:** controla la disponibilidad y la ocupación de las habitaciones.
- **📝 `EpisodeRecord`:** documenta el episodio de sueño y la acción tomada para atenderlo.

El modelo utiliza enumeraciones tipadas en lugar de cadenas de texto convencionales para representar el estado del paciente y los niveles de gravedad de la narcolepsia.

Además, debe respetar el diseño requerido de composición e interfaces.

## 🔍 5. Consultas requeridas

Las siguientes consultas deben estar implementadas y devolver datos para que puedan utilizarse en otras partes del programa, en lugar de limitarse a imprimir resultados.

| 🔎 Consulta | 🎯 Propósito |
|---|---|
| `PatientsInHallway()` | Obtiene los pacientes que se encuentran en los pasillos. |
| `Doctor.MyEpisodes()` | Devuelve los episodios registrados por un médico. |
| `Hospital.AssignRoom(p *Patient)` | Intenta asignar una habitación disponible a un paciente. |
| `Hospital.SevereReport()` | Genera un informe de los pacientes o episodios clasificados como graves, según el modelo implementado. |

Estas consultas deben permitir responder al escenario de demostración de forma limpia, organizada y fácil de probar.

## 🎬 6. Escenario de demostración

El programa debe ejecutar el siguiente flujo, respetando el orden establecido:

1. 🏥 Crear el hospital.
2. 🩺 Contratar a **2 médicos y 1 auxiliar de enfermería**.
3. 🛏️ Registrar **3 habitaciones**, cada una con capacidad para un paciente.
4. 📝 Admitir **5 pacientes**, de los cuales al menos 2 deben estar clasificados como graves.
5. 😴 Simular episodios repentinos de sueño en **4 pacientes**, cada uno en una ubicación diferente.
6. 🛏️ Asignar una habitación a 3 de los pacientes. El cuarto no consigue una cama y debe permanecer en el pasillo.
7. 🌅 Hacer que uno de los pacientes despierte y libere su habitación.
8. 🔄 Asignar la habitación liberada al paciente que permanece en el pasillo.
9. 📊 Imprimir los resultados de las cuatro consultas requeridas.

Este escenario permite que el evaluador compare los resultados entre los diferentes grupos de manera consistente.

## ⭐ 7. Puntos adicionales (bonus)

Se pueden implementar funcionalidades opcionales para obtener puntos adicionales.

- 🖥️ **Interfaz gráfica:** utilizar Fyne, Wails o una API REST ligera con un frontend.
- ⚡ **Simulación concurrente:** utilizar goroutines y `sync.Mutex` para gestionar operaciones concurrentes de manera segura.
- 🧪 **Seguridad ante condiciones de carrera:** garantizar que el programa se ejecute correctamente mediante `go run -race .`.

## 📂 8. Estructura del proyecto

El repositorio debe tener la siguiente estructura:

```text
costenitos/
├── go.mod
├── Readme.md
├── AI_USAGE.md
├── main.go
├── cmd/
│   └── web/
│       ├── main.go
│       ├── server.go
│       ├── scenario.go
│       └── static/
│           └── index.html
├── hospital/
│   ├── hospital.go
│   ├── patient.go
│   ├── doctor.go
│   ├── orderly.go
│   ├── room.go
│   ├── episode.go
│   ├── hospital_test.go
│   └── orderly_test.go
```

### 📁 Descripción de los archivos

- `go.mod`: define el módulo y la versión de Go.
- `Readme.md`: contiene la documentación del proyecto.
- `AI_USAGE.md`: documenta el uso de herramientas de inteligencia artificial.
- `main.go`: construye el escenario de demostración y muestra los resultados del dominio en consola.
- `cmd/web/main.go`: lanza la API web y sirve la interfaz estática.
- `cmd/web/server.go`: implementa el backend REST del hospital y el estado en memoria.
- `cmd/web/scenario.go`: archivo auxiliar del paquete web; se mantiene en el mismo módulo.
- `cmd/web/static/index.html`: interfaz de usuario ligera para operar el hospital desde el navegador.
- `hospital/hospital.go`: implementa la lógica principal del hospital.
- `hospital/patient.go`: define el modelo y el comportamiento de los pacientes.
- `hospital/doctor.go`: define los médicos y sus funcionalidades.
- `hospital/orderly.go`: define el personal auxiliar de enfermería y su comportamiento.
- `hospital/room.go`: gestiona las habitaciones.
- `hospital/episode.go`: representa los registros de episodios de sueño.
- `hospital/hospital_test.go`: contiene las pruebas automatizadas del hospital.
- `hospital/orderly_test.go`: verifica el comportamiento del personal auxiliar.

## ▶️ 9. Cómo ejecutar el proyecto

### Consola

Desde la carpeta raíz del repositorio, ejecuta:

```bash
go run .
```

Este comando ejecuta el escenario de demostración de consola y muestra los resultados de las consultas requeridas.

### API web

También se incluye una interfaz ligera con API REST para interactuar desde el navegador:

```bash
go run ./cmd/web
```

La app usa el puerto `9090` por defecto. Si quieres cambiarlo, puedes definir la variable de entorno `PORT`:

```bash
$env:PORT = "9091"
go run ./cmd/web
```

Luego abre:

```text
http://127.0.0.1:9090
```

o el puerto que hayas configurado.

## 🧪 10. Cómo ejecutar las pruebas

Para ejecutar todas las pruebas automatizadas del proyecto, utiliza:

```bash
go test ./...
```

Este comando ejecutará las pruebas de todos los paquetes del módulo.

## 👥 11. Integrantes del equipo

- 👨‍💻 Carlos Mauricio Velasco
- 👨‍💻 Santiago Perez
- 👨‍💻 Pablo Escudero

## 🧠 12. Decisiones de diseño

El proyecto aplica los principios de programación orientada a objetos mediante las siguientes decisiones:

- 🧱 **Composición:** `Patient` y `Doctor` incorporan la estructura `Person` mediante campos embebidos.
- 🔌 **Polimorfismo:** la interfaz `Attender` permite trabajar con diferentes implementaciones concretas que cumplan el mismo contrato.
- 🏷️ **Enumeraciones tipadas:** se utilizan constantes tipadas para representar los estados de los pacientes y los niveles de gravedad.
- 🔒 **Encapsulación:** los campos privados se gestionan mediante métodos de acceso y operaciones definidas por cada entidad.
- 📦 **Separación de responsabilidades:** la lógica del dominio se concentra en el paquete `hospital`, mientras que `main.go` se encarga de construir el escenario y presentar los resultados.

Estas decisiones buscan que el código sea modular, mantenible, fácil de probar y coherente con el enfoque de programación orientada a objetos de Go.

## 🏁 13. Consideraciones finales

El objetivo de este proyecto no es únicamente modelar un hospital, sino también aplicar los conceptos de programación orientada a objetos en una solución práctica, organizada y verificable.

La implementación final debe ser clara, seguir las convenciones de Go y permitir que cada integrante explique y defienda las decisiones de diseño durante la evaluación individual.

**🌴 ¡Que nunca falten las camas, que los médicos permanezcan despiertos y que el sancocho no interrumpa las consultas!** 🍲🏥💤
