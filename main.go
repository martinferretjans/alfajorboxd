package main

import(
	"fmt"
	"net/http"
)

func main(){
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)
	
	port := ":8080"
	fmt.Printf("Server escuchando en %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar servidor: %s\n", err)
	}
}
