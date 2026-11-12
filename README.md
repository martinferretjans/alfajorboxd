# Alfajorboxd

Aplicación web para registrar y calificar alfajores, inspirada en Letterboxd y Backloggd. Los usuarios pueden cargar alfajores, dejar reviews calificándolos de media estrella a 5 estrellas y ver el promedio de cada uno.

## Requisitos

- Make
- Docker Compose

## Cómo ejecutar

```bash
git clone -b tp5 https://github.com/martinferretjans/alfajorboxd.git
cd alfajorboxd
make run-app
```

La aplicación queda disponible en `http://localhost:8080`. Al cortar la ejecución con `Ctrl+C` se eliminan los contenedores y volúmenes.

Para correr las pruebas de la base de datos:

```bash
make test
```

Con la aplicación levantada, `request.sh` prueba todas las rutas y verifica tanto el código HTTP como los datos devueltos:

```bash
./request.sh
```

# Documentación

La base de datos tiene tres tablas: `alfajor`, `usuario` y `review`. Al crearse la base se cargan datos de ejemplo desde `db/schema/seed.sql`.

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
| `imagen`      | `VARCHAR(255)` | Opcional                         | Nombre del archivo en `static/img/`.|

### Tabla `usuario`

| Campo   | Tipo SQL       | Restricciones        | Descripción                    |
|---------|----------------|----------------------|--------------------------------|
| `id`    | `SERIAL`       | `PRIMARY KEY`        | Identificador del usuario.     |
| `name`  | `VARCHAR(255)` | `NOT NULL`           | Nombre del usuario.            |
| `email` | `VARCHAR(255)` | `UNIQUE`, `NOT NULL` | Email del usuario.             |

### Tabla `review`

| Campo          | Tipo SQL       | Restricciones                                | Descripción                                   |
|----------------|----------------|----------------------------------------------|-----------------------------------------------|
| `id_alfajor`   | `VARCHAR(255)` | `PRIMARY KEY`, `REFERENCES alfajor(id)`      | Alfajor reseñado.                             |
| `id_user`      | `INT`          | `PRIMARY KEY`, `REFERENCES usuario(id)`      | Usuario que escribió la review.               |
| `calificacion` | `INT`          | `NOT NULL`, `CHECK (entre 1 y 10)`           | Cada punto equivale a media estrella.         |
| `comentario`   | `TEXT`         | Opcional                                     | Texto de la review.                           |
| `fecha`        | `DATE`         | `NOT NULL`, por defecto la fecha actual      | Fecha de la review.                           |

Un usuario puede hacer una sola review por alfajor. Si se borra un alfajor o un usuario, se borran sus reviews (`ON DELETE CASCADE`).

La calificación de un alfajor es el promedio de sus reviews pasado a estrellas (por ejemplo, reviews de 10, 9 y 9 dan 4.67 estrellas). Se calcula en las consultas `GetPromedioAlfajor` y `ListAlfajoresPopulares`.

# Vistas

Las páginas se renderizan en el servidor con `templ` (directorio `views`). Los formularios usan el patrón Post/Redirect/Get: hacen un POST y el servidor redirige a la página que corresponde.

| Ruta                  | Método | Descripción                                                        |
|-----------------------|--------|--------------------------------------------------------------------|
| `/`                   | GET    | Alfajores populares, reviews recientes y formulario de carga.      |
| `/alfajores`          | POST   | Crea un alfajor y redirige a `/`.                                  |
| `/alfajores/{id}`     | GET    | Detalle del alfajor, su promedio, sus reviews y el formulario de review. |
| `/reviews`            | GET    | Listado de las últimas reviews.                                    |
| `/reviews`            | POST   | Crea una review y redirige al detalle del alfajor.                 |
| `/usuarios`           | GET    | Listado de usuarios y formulario de alta.                          |
| `/usuarios`           | POST   | Crea un usuario y redirige a `/usuarios`.                          |
| `/tema`               | POST   | Cambia el tema de colores y vuelve a la página anterior.           |
| `/static/...`         | GET    | Archivos estáticos (CSS e imágenes).                               |

### Componentes

- `layout.templ`: estructura base de todas las páginas (head, cabecera con navegación y pie). Recibe el contenido como `templ.Component`.
- `alfajor_list.templ` y `alfajor_form.templ`: grilla de alfajores y formulario de carga.
- `review_list.templ` y `review_form.templ`: listado y formulario de reviews.
- `usuario_list.templ` y `usuario_form.templ`: tabla y formulario de usuarios.
- `index.templ`, `alfajor.templ`, `reviews.templ` y `usuarios.templ`: páginas completas que usan el layout.

### Temas de colores

Hay dos temas: uno claro inspirado en el envoltorio del Havanna de chocolate (dorado `#d8c28a` con marrón `#514227`), que es el tema por defecto, y uno oscuro marrón (`#60341f` con `#fbd705`). La estructura de la página está en `static/css/styles.css` y los colores de cada tema en `static/css/havanna.css` y `static/css/marron.css`. El botón "Cambiar tema" guarda el tema elegido en una cookie y el layout carga la hoja de estilos que corresponde.

### Imágenes

Las imágenes de los alfajores van en `static/img/`, y en el campo `imagen` se guarda solo el nombre del archivo (por ejemplo `havanna-chocolate.webp`). Los pósters son cuadrados con fondo blanco y la imagen se muestra centrada sin deformarse, por lo que sirven fotos de cualquier proporción. Si un alfajor no tiene imagen, se muestra un póster con su nombre.

Las imágenes de los alfajores de ejemplo fueron obtenidas de [alfajor.com.ar](https://alfajor.com.ar/ranking).

# Ejecución

La ejecución del proyecto está automatizada mediante el Makefile.

### Organización de los contenedores

Docker Compose define tres servicios:

- **`db`**: ejecuta PostgreSQL 16 y contiene la base de datos `alfajorboxd`. Al iniciarse ejecuta `schema.sql` y `seed.sql`.
- **`test`**: contiene Go, `sqlc`, `templ`, las dependencias y el código fuente del proyecto. Se utiliza para generar el código, compilar la aplicación y ejecutar las pruebas.
- **`app`**: usa la misma imagen que `test` y ejecuta el servidor web en el puerto 8080.

### Flujo de ejecución

Al ejecutar `make run-app`, se realizan automáticamente los siguientes pasos:

1. Se eliminan los contenedores y volúmenes de ejecuciones anteriores.
2. Se construye la imagen del servicio `app` utilizando el `Dockerfile`.
3. Se inicia el contenedor de PostgreSQL y se espera hasta que la base de datos esté lista.
4. Se crea un contenedor para ejecutar:
   - `sqlc generate`, que genera el código Go a partir de las consultas SQL.
   - `templ generate`, que genera el código Go a partir de las vistas `.templ`.
   - `go build`, que compila la aplicación.
   - La aplicación compilada, que queda escuchando en el puerto 8080.
5. Al finalizar (`Ctrl+C`), se eliminan los contenedores y los volúmenes creados durante la ejecución.

Al ejecutar `make test`, se realizan automáticamente los siguientes pasos:

1. Se eliminan los contenedores y volúmenes de ejecuciones anteriores.
2. Se construye la imagen del servicio `test` utilizando el `Dockerfile`.
3. Se inicia el contenedor de PostgreSQL.
4. Se espera hasta que la base de datos esté lista.
5. Se crea un contenedor para ejecutar:
   - `sqlc generate`, que genera el código Go a partir de las consultas SQL.
   - `templ generate`, que genera el código Go a partir de las vistas `.templ`.
   - `go build ./...`, que comprueba que todos los paquetes del proyecto compilen.
   - `go test ./... -v`, que ejecuta las pruebas de las tres tablas.
6. Al finalizar, se eliminan los contenedores y los volúmenes creados durante la ejecución.
