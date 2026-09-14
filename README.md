# Alfajorboxd

Aplicación web para registrar y calificar alfajores.

## Requisitos

- Make
- Docker Compose

## Cómo ejecutar

```bash
git clone -b tp2 https://github.com/martinferretjans/alfajorboxd.git
cd alfajorboxd
make test
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

# Ejecución

La ejecución del proyecto está automatizada mediante el Makefile.

### Organización de los contenedores

Docker Compose define dos servicios independientes:

- **`db`**: ejecuta PostgreSQL 16 y contiene la base de datos `alfajorboxd`.
- **`test`**: contiene Go, `sqlc`, las dependencias y el código fuente del proyecto. Se utiliza para generar el código, compilar la aplicación y ejecutar las pruebas.

### Flujo de ejecución

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
