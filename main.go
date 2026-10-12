package main

import(
	"fmt"
	"os"
	"context"
	"log"
	"net/http"
	"database/sql"
	"alfajorboxd/handlers"
	db "alfajorboxd/db/sqlc"
	_ "github.com/lib/pq"
)

func main(){

	conn, err := abrirDB()
	
	if err != nil {
		log.Fatal("Error al conectar con db: ", err)
	}
	defer conn.Close()
	
	fmt.Println("Conexion con DB (Postgresql) exitosa")
	
	queries := db.New(conn)
	
	h := handlers.New(queries)
	
	
	
	//Prueba de conexion a la base
	alfajores, err := queries.ListAlfajores(context.Background())
	if err != nil{
		log.Fatal("Error al listar alfajores: ", err)
	}
	fmt.Printf("Alfajores encontrados: %d\n", len(alfajores))

	//Server http
	
	http.HandleFunc("/alfajores", h.AlfajoresHandler)
	http.HandleFunc("/alfajores/", h.AlfajorHandler)
	
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)
	
	port := ":8080"
	fmt.Printf("Server escuchando en %s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar servidor: %s\n", err)
	}
}

func abrirDB() (*sql.DB, error){
	dsn := os.Getenv("DATABASE_URL")
	
	db, err := sql.Open("postgres",dsn)
	
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	
	db.SetMaxOpenConns(25)
	
	return db, nil
}
