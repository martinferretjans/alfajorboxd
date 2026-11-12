#!/bin/bash

#se corre con la aplicacion levantada con make run-app, usa los datos de db/schema/seed.sql
BASE_URL="http://localhost:8080"

verificar() {
  if echo "$2" | grep -q "$3"; then
    echo "OK: $1"
  else
    echo "FALLO: $1 (no se encontro '$3')"
  fi
}

echo "1. GET / (Pagina principal - Esperado: 200 OK con los alfajores del seed)"
RESP=$(curl -s -i "$BASE_URL/")
verificar "codigo 200" "$RESP" "200 OK"
verificar "lista de alfajores" "$RESP" "Havanna de chocolate"
verificar "reviews recientes" "$RESP" "Mucho chocolate, me encanto."
echo -e "\n----------------------------------------\n"

echo "2. POST /alfajores (Crear - Esperado: 303 See Other a /)"
RESP=$(curl -s -i -X POST "$BASE_URL/alfajores" \
  -d "id=guaymallen-blanco" \
  -d "name=Guaymallen blanco" \
  -d "relleno=Dulce de leche" \
  -d "cobertura=Chocolate blanco" \
  -d "precio=1800" \
  -d "tapas=2" \
  -d "marca=Guaymallen" \
  -d "descripcion=Afajor simple blanco")
verificar "codigo 303" "$RESP" "303 See Other"
verificar "redirige a /" "$RESP" "Location: /"
RESP=$(curl -s "$BASE_URL/")
verificar "el alfajor aparece en la lista" "$RESP" "Guaymallen blanco"
echo -e "\n----------------------------------------\n"

echo "3. POST /alfajores (Datos invalidos - Esperado: 400 Bad Request)"
RESP=$(curl -s -i -X POST "$BASE_URL/alfajores" \
  -d "id=invalido" \
  -d "name=" \
  -d "relleno=Dulce de leche" \
  -d "precio=-100" \
  -d "tapas=5" \
  -d "marca=X")
verificar "codigo 400" "$RESP" "400 Bad Request"
verificar "mensaje de error" "$RESP" "el nombre (name) no puede estar"
echo -e "\n----------------------------------------\n"

echo "4. GET /alfajores/guaymallen-blanco (Detalle - Esperado: 200 OK)"
RESP=$(curl -s -i "$BASE_URL/alfajores/guaymallen-blanco")
verificar "codigo 200" "$RESP" "200 OK"
verificar "cobertura" "$RESP" "Chocolate blanco"
verificar "precio" "$RESP" "\$1800"
verificar "sin reviews" "$RESP" "Sin reviews"
echo -e "\n----------------------------------------\n"

echo "5. GET /alfajores/no-existe (Esperado: 404 Not Found)"
RESP=$(curl -s -i "$BASE_URL/alfajores/no-existe")
verificar "codigo 404" "$RESP" "404 Not Found"
verificar "mensaje de error" "$RESP" "Alfajor no encontrado"
echo -e "\n----------------------------------------\n"

echo "6. GET /alfajores/havanna-chocolate (Promedio - Esperado: 4.67)"
RESP=$(curl -s "$BASE_URL/alfajores/havanna-chocolate")
verificar "promedio de 10, 9 y 9" "$RESP" "4.67"
verificar "cantidad de reviews" "$RESP" "3 reviews"
echo -e "\n----------------------------------------\n"

echo "7. POST /usuarios (Crear - Esperado: 303 See Other a /usuarios)"
RESP=$(curl -s -i -X POST "$BASE_URL/usuarios" -d "name=ana" -d "email=ana@alfajorboxd.com")
verificar "codigo 303" "$RESP" "303 See Other"
RESP=$(curl -s "$BASE_URL/usuarios")
verificar "el usuario aparece en la lista" "$RESP" "ana@alfajorboxd.com"
echo -e "\n----------------------------------------\n"

echo "8. POST /usuarios (Email invalido - Esperado: 400 Bad Request)"
RESP=$(curl -s -i -X POST "$BASE_URL/usuarios" -d "name=pepe" -d "email=pepe")
verificar "codigo 400" "$RESP" "400 Bad Request"
verificar "mensaje de error" "$RESP" "el email no es valido"
echo -e "\n----------------------------------------\n"

echo "9. POST /reviews (Crear - Esperado: 303 See Other al detalle del alfajor)"
RESP=$(curl -s -i -X POST "$BASE_URL/reviews" -d "id_alfajor=guaymallen-blanco" -d "id_user=1" -d "calificacion=9" -d "comentario=Muy rico el blanco")
verificar "codigo 303" "$RESP" "303 See Other"
verificar "redirige al alfajor" "$RESP" "Location: /alfajores/guaymallen-blanco"
RESP=$(curl -s "$BASE_URL/alfajores/guaymallen-blanco")
verificar "comentario de la review" "$RESP" "Muy rico el blanco"
verificar "promedio de una review de 9" "$RESP" "4.50"
echo -e "\n----------------------------------------\n"

echo "10. POST /reviews (Repetida - Esperado: 409 Conflict)"
RESP=$(curl -s -i -X POST "$BASE_URL/reviews" -d "id_alfajor=guaymallen-blanco" -d "id_user=1" -d "calificacion=3")
verificar "codigo 409" "$RESP" "409 Conflict"
verificar "mensaje de error" "$RESP" "ya hizo una review"
echo -e "\n----------------------------------------\n"

echo "11. POST /reviews (Calificacion invalida - Esperado: 400 Bad Request)"
RESP=$(curl -s -i -X POST "$BASE_URL/reviews" -d "id_alfajor=guaymallen-blanco" -d "id_user=2" -d "calificacion=11")
verificar "codigo 400" "$RESP" "400 Bad Request"
verificar "mensaje de error" "$RESP" "la calificacion debe estar"
echo -e "\n----------------------------------------\n"

echo "12. GET /reviews (Listar - Esperado: 200 OK con la review nueva)"
RESP=$(curl -s -i "$BASE_URL/reviews")
verificar "codigo 200" "$RESP" "200 OK"
verificar "review nueva" "$RESP" "Muy rico el blanco"
echo -e "\n----------------------------------------\n"

echo "13. POST /tema (Cambiar a marron - Esperado: 303 See Other y cookie tema=marron)"
RESP=$(curl -s -i -X POST "$BASE_URL/tema" -d "tema=marron")
verificar "codigo 303" "$RESP" "303 See Other"
verificar "cookie" "$RESP" "Set-Cookie: tema=marron"
RESP=$(curl -s --cookie "tema=marron" "$BASE_URL/")
verificar "hoja de estilos marron" "$RESP" "/static/css/marron.css"
RESP=$(curl -s "$BASE_URL/")
verificar "sin cookie usa havanna" "$RESP" "/static/css/havanna.css"
echo -e "\n"
