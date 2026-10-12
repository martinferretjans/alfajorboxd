package handlers

import (
    "database/sql"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strings"

    db "alfajorboxd/db/sqlc"
)

type Handler struct {
    queries *db.Queries
}

func New(queries *db.Queries) *Handler {
    return &Handler{queries: queries}
}

// Handler para /alfajores
func (h *Handler) AlfajoresHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    	case http.MethodGet:
        	h.getAlfajores(w, r)
    	default:
        	http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
    }
}

// Handler para /alfajores/{id}
func (h *Handler) AlfajorHandler(w http.ResponseWriter, r *http.Request) {
    
    if r.URL.Path == "/alfajores/" {
        h.AlfajoresHandler(w, r)
        return
    }
    
    id := strings.TrimPrefix(r.URL.Path, "/alfajores/")
    
    if id == "" || strings.Contains(id, "/") {
        http.Error(w, "ID invalido", http.StatusBadRequest)
        return
    }

    switch r.Method {
    	case http.MethodGet:
        	h.getAlfajor(w, r, id)
    	default:
        	http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
    }
}

// GET /alfajores
func (h *Handler) getAlfajores(w http.ResponseWriter, r *http.Request) {
    alfajores, err := h.queries.ListAlfajores(r.Context())
    if err != nil {
        log.Println(err)
        http.Error(w, "Error interno", http.StatusInternalServerError)
        return
    }

    if alfajores == nil {
        alfajores = []db.Alfajor{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alfajores)
}


// GET /alfajores/id
func (h *Handler) getAlfajor(w http.ResponseWriter, r *http.Request, id string) {

    alfajor, err := h.queries.GetAlfajor(r.Context(), id)

    if errors.Is(err, sql.ErrNoRows) {
        http.Error(w, "Alfajor no encontrado", http.StatusNotFound)
        return
    }

    if err != nil {
        log.Printf("Error al buscar alfajor: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alfajor)
}
