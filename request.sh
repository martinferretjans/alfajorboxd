
#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:8080"

RESPUESTA=$(mktemp)
trap 'rm -f "$RESPUESTA"' EXIT

# Ejecuta una petición y comprueba el código HTTP
verificar() {
    local esperado="$1"
    local descripcion="$2"
    shift 2

    local obtenido

    if ! obtenido=$(curl -sS -o "$RESPUESTA" -w "%{http_code}" "$@"); then
        echo "ERROR: $descripcion - No se pudo realizar la petición"
        exit 1
    fi

    if [ "$obtenido" != "$esperado" ]; then
        echo "ERROR: $descripcion"
        echo "Esperado: $esperado | Obtenido: $obtenido"
        cat "$RESPUESTA"
        exit 1
    fi

    echo "OK: $descripcion ($obtenido)"
}

# 1. Crear alfajor
verificar 201 "POST crear alfajor" \
    -X POST "$BASE_URL/alfajores" \
    -H "Content-Type: application/json" \
    -d '{
        "id": "guaymallen-blanco",
        "name": "Guaymallen blanco",
        "relleno": "Dulce de leche",
        "cobertura": "Chocolate blanco",
        "precio": 1800,
        "tapas": 2,
        "marca": "Guaymallen",
        "descripcion": "Alfajor simple blanco"
    }'

# 2. Intentar crear un alfajor inválido
verificar 400 "POST datos invalidos" \
    -X POST "$BASE_URL/alfajores" \
    -H "Content-Type: application/json" \
    -d '{
        "id": "invalido",
        "name": "",
        "relleno": "Dulce de leche",
        "precio": -100,
        "tapas": 5,
        "marca": "X"
    }'

# 3. Listar alfajores
verificar 200 "GET listar alfajores" \
    "$BASE_URL/alfajores"

# 4. Buscar el alfajor creado
verificar 200 "GET buscar alfajor" \
    "$BASE_URL/alfajores/guaymallen-blanco"

# 5. Actualizar alfajor
verificar 200 "PUT actualizar alfajor" \
    -X PUT "$BASE_URL/alfajores/guaymallen-blanco" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Guaymallen triple blanco",
        "relleno": "Dulce de leche",
        "cobertura": "Merengue",
        "precio": 2200,
        "tapas": 3,
        "marca": "Guaymallen",
        "descripcion": "Actualizado"
    }'

# Verificar que los cambios se guardaron
verificar 200 "GET verificar actualizacion" \
    "$BASE_URL/alfajores/guaymallen-blanco"

if ! grep -q '"Precio":2200' "$RESPUESTA"; then
    echo "ERROR: El precio no se actualizo correctamente"
    exit 1
fi

if ! grep -q '"Tapas":3' "$RESPUESTA"; then
    echo "ERROR: Las tapas no se actualizaron correctamente"
    exit 1
fi

echo "OK: Datos actualizados correctamente"

# 6. Eliminar alfajor
verificar 204 "DELETE eliminar alfajor" \
    -X DELETE "$BASE_URL/alfajores/guaymallen-blanco"

# 7. Comprobar que fue eliminado
verificar 404 "GET alfajor eliminado" \
    "$BASE_URL/alfajores/guaymallen-blanco"

echo ""
echo "Todas las pruebas de la API finalizaron correctamente."
