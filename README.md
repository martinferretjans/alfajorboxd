# Alfajorboxd

Aplicación web para registrar y calificar alfajores.

## Requisitos

- Make
- Docker Compose
- curl
- bash
- Git

## Cómo ejecutar

```bash
git clone -b tp3 https://github.com/martinferretjans/alfajorboxd.git
cd alfajorboxd
```

Para ejecutar automáticamente las pruebas de la API:

```bash
make test-api
```

Este comando realiza las tareas de preparación, ejecución y limpieza necesarias sin requerir intervención manual.

Para levantar la aplicación y utilizar la API:

```bash
make run
```

El servidor HTTP queda disponible en:

`http://localhost:8080`

Para detener los servicios:

```bash
make stop
```

# Documentación

La tabla principal del proyecto es `alfajor`, la cual guarda los datos de cada alfajor en la plataforma. Estuvimos trabajando también en las tablas `usuario` y `review`, que formarán parte de futuras entregas.

### Tabla `alfajor`

| Campo         | Tipo SQL       | Restricciones                    | Descripción                         |
|---------------|----------------|----------------------------------|-------------------------------------|
| `id`          | `VARCHAR`      | `PRIMARY KEY`                    | Identificador único del alfajor.    |
| `name`        | `VARCHAR(255)` | `NOT NULL`                       | Nombre del alfajor.                 |
| `relleno`     | `VARCHAR(255)` | `NOT NULL`                       | Tipo de relleno.                    |
| `cobertura`   | `VARCHAR(255)` | Opcional                         | Tipo de cobertura.                  |
| `precio`      | `INT`          | `NOT NULL`                       | Precio del alfajor.                 |
| `tapas`       | `INT`          | `NOT NULL`, `CHECK (2 o 3)`      | Cantidad de tapas.                  |
| `marca`       | `VARCHAR(255)` | `NOT NULL`                       | Marca del fabricante.               |
| `descripcion` | `TEXT`         | Opcional                         | Descripción adicional del alfajor.  |

### Endpoints disponibles

| Método | Endpoint | Descripción | Respuesta exitosa |
|---|---|---|---|
| GET | `/alfajores` | Obtener todos los alfajores. | 200 OK |
| GET | `/alfajores/{id}` | Obtener un alfajor por ID. | 200 OK |
| POST | `/alfajores` | Crear un nuevo alfajor. | 201 Created |
| PUT | `/alfajores/{id}` | Actualizar un alfajor existente. | 200 OK |
| DELETE | `/alfajores/{id}` | Eliminar un alfajor. | 204 No Content |

Los endpoints devuelven los resultados en formato JSON, excepto DELETE, que responde sin contenido al eliminar correctamente.

### Validaciones

Antes de crear o actualizar un alfajor, se comprueba que:

- El nombre, el relleno y la marca no estén vacíos.
- El precio sea mayor que cero.
- La cantidad de tapas sea 2 o 3.
- En la creación, el ID no esté vacío.

Si los datos recibidos son inválidos, se devuelve `400 Bad Request`.

Cuando se intenta consultar, actualizar o eliminar un alfajor inexistente, se devuelve `404 Not Found`.

Los campos `cobertura` y `descripcion` son opcionales. Si se reciben vacíos, contienen solamente espacios o se omiten, se almacenan como `NULL` en PostgreSQL.

### Ejemplos de peticiones

**Listar todos los alfajores:**

```bash
curl -i http://localhost:8080/alfajores
```

**Obtener un alfajor por ID:**

```bash
curl -i http://localhost:8080/alfajores/guaymallen-blanco
```

**Crear un alfajor:**

```bash
curl -i -X POST http://localhost:8080/alfajores \
  -H "Content-Type: application/json" \
  -d '{
    "id": "guaymallen-blanco",
    "name": "Guaymallen Blanco",
    "relleno": "Dulce de leche",
    "cobertura": "Chocolate blanco",
    "precio": 1800,
    "tapas": 2,
    "marca": "Guaymallen",
    "descripcion": "Alfajor simple blanco"
  }'
```

**Actualizar un alfajor:**

```bash
curl -i -X PUT http://localhost:8080/alfajores/guaymallen-blanco \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Guaymallen Triple Blanco",
    "relleno": "Dulce de leche",
    "cobertura": "Chocolate blanco",
    "precio": 2200,
    "tapas": 3,
    "marca": "Guaymallen"
  }'
```

**Eliminar un alfajor:**

```bash
curl -i -X DELETE http://localhost:8080/alfajores/guaymallen-blanco
```

## Organización de los contenedores

Docker Compose define tres servicios:

- **`db`**: ejecuta PostgreSQL 16 y contiene la base de datos `alfajorboxd`.
- **`test`**: contiene Go, `sqlc`, las dependencias y el código fuente. Se utiliza para generar código, compilar y ejecutar los tests de Go.
- **`app`**: ejecuta el servidor HTTP de la aplicación, conectado a PostgreSQL, exponiendo el puerto 8080.

### Pruebas de capa de datos

Al ejecutar `make test`, se realizan automáticamente los siguientes pasos:

1. Se eliminan los contenedores y volúmenes de ejecuciones anteriores.
2. Se construye la imagen del servicio `test` utilizando el `Dockerfile`.
3. Se inicia el contenedor de PostgreSQL.
4. Se espera hasta que la base de datos esté lista.
5. Se crea un contenedor para ejecutar:
   - `sqlc generate`, que genera el código Go a partir de las consultas SQL.
   - `go build ./...`, que comprueba que todos los paquetes del proyecto compilen.
   - `go test ./... -v`, que ejecuta todas las pruebas.
6. Al finalizar, se eliminan los contenedores y los volúmenes creados durante la ejecución.

### Pruebas de la API REST

Para ejecutar las pruebas del TP3:

```bash
make test-api
```

Este comando automatiza el proceso completo:

1. Elimina los contenedores y volúmenes anteriores del entorno de pruebas de la API.
2. Construye la imagen de la aplicación.
3. Inicia PostgreSQL y espera a que se encuentre disponible.
4. Genera el código de `sqlc` e inicia el servidor HTTP con Go.
5. Espera hasta que la API responda a las peticiones HTTP.
6. Ejecuta `request.sh`, que realiza las pruebas de los endpoints.
7. Al finalizar, elimina los contenedores y volúmenes utilizados, incluso si alguna prueba falla.

Para estas pruebas se utiliza un proyecto independiente de Docker Compose, llamado `alfajorboxd-test-api`, con el objetivo de aislar sus contenedores y volúmenes de los utilizados en la ejecución normal.

El puerto 8080 debe estar disponible antes de ejecutar `make test-api`.

### Archivo `request.sh`

El script utiliza `curl` para enviar peticiones HTTP y comprobar los códigos de respuesta.

Se realizan las siguientes pruebas:

| Prueba | Código esperado |
|---|---|
| Crear un alfajor válido (POST) | 201 |
| Crear un alfajor con datos inválidos (POST) | 400 |
| Listar todos los alfajores (GET) | 200 |
| Obtener el alfajor creado (GET) | 200 |
| Actualizar el alfajor (PUT) | 200 |
| Eliminar el alfajor (DELETE) | 204 |
| Intentar obtener el alfajor eliminado (GET) | 404 |

## Comandos disponibles

| Comando | Descripción |
|---|---|
| `make run` | Levanta la aplicación y PostgreSQL. |
| `make test` | Ejecuta los tests de Go y de la capa de datos. |
| `make test-api` | Ejecuta las pruebas automatizadas de la API REST. |
| `make stop` | Detiene y elimina los contenedores de la ejecución normal, conservando los volúmenes. |
| `make clean` | Elimina los contenedores y volúmenes de la ejecución normal. |