#!/bin/bash

BASE_URL="http://localhost:8080"

echo "1. POST /alfajores (Crear - Esperado: 201 Created)"
curl -i -X POST "$BASE_URL/alfajores" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "guaymallen-blanco",
    "name": "Guaymallen blanco",
    "relleno": "Dulce de leche",
    "cobertura": "Chocolate blanco",
    "precio": 1800,
    "tapas": 2,
    "marca": "Guaymallen",
    "descripcion": "Afajor simple blanco"
  }'
echo -e "\n----------------------------------------\n"

echo "2. POST /alfajores (Datos invalidos - Esperado: 400 Bad Request)"
curl -i -X POST "$BASE_URL/alfajores" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "invalido",
    "name": "",
    "relleno": "Dulce de leche",
    "precio": -100,
    "tapas": 5,
    "marca": "X"
  }'
echo -e "\n----------------------------------------\n"

echo "3. GET /alfajores (Listar todos - Esperado: 200 OK)"
curl -i -X GET "$BASE_URL/alfajores"
echo -e "\n----------------------------------------\n"

echo "4. GET /alfajores/guaymallen-blanco (Buscar uno - Esperado: 200 OK)"
curl -i -X GET "$BASE_URL/alfajores/guaymallen-blanco"
echo -e "\n----------------------------------------\n"

echo "5. PUT /alfajores/guaymallen-blanco (Actualizar - Esperado: 200 OK)"
curl -i -X PUT "$BASE_URL/alfajores/guaymallen-blanco" \
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
echo -e "\n----------------------------------------\n"

echo "6. DELETE /alfajores/guaymallen-blanco (Borrar - Esperado: 204 No Content)"
curl -i -X DELETE "$BASE_URL/alfajores/guaymallen-blanco"
echo -e "\n----------------------------------------\n"

echo "7. GET /alfajores/guaymallen-blanco (Verificar borrado - Esperado: 404 Not Found)"
curl -i -X GET "$BASE_URL/alfajores/guaymallen-blanco"
echo -e "\n"
