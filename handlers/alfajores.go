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

type AlfajorInput struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Relleno     string  `json:"relleno"`
    Cobertura   *string `json:"cobertura"`
    Precio      int32   `json:"precio"`
    Tapas       int32   `json:"tapas"`
    Marca       string  `json:"marca"`
    Descripcion *string `json:"descripcion"`
}

func nullString(s *string) sql.NullString {
    if s == nil || strings.TrimSpace(*s) == "" {
        return sql.NullString{}
    }

    return sql.NullString{
        String: *s,
        Valid:  true,
    }
}

func validarAlfajor(a AlfajorInput, validarID bool) error {
    if validarID && strings.TrimSpace(a.ID) == "" {
        return errors.New("el id es obligatorio")
    }
    if strings.TrimSpace(a.Name) == "" {
        return errors.New("el nombre no puede estar vacio")
    }
    if strings.TrimSpace(a.Relleno) == "" {
        return errors.New("el relleno no puede estar vacio")
    }
    if strings.TrimSpace(a.Marca) == "" {
        return errors.New("la marca no puede estar vacia")
    }
    if a.Precio <= 0 {
        return errors.New("el precio debe ser mayor a 0")
    }
    if a.Tapas != 2 && a.Tapas != 3 {
        return errors.New("la cantidad de tapas debe ser 2 o 3")
    }

    return nil
}

func New(queries *db.Queries) *Handler {
    return &Handler{queries: queries}
}

// Handler para /alfajores
func (h *Handler) AlfajoresHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    	case http.MethodGet:
        	h.getAlfajores(w, r)
    	case http.MethodPost:
        	h.createAlfajor(w, r)
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
	case http.MethodPut:
    		h.updateAlfajor(w, r, id)
	case http.MethodDelete:
    		h.deleteAlfajor(w, r, id)
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

// POST /alfajores
func (h *Handler) createAlfajor(w http.ResponseWriter, r *http.Request) {
    var datos AlfajorInput

    // Decodificar JSON
    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        http.Error(w, "JSON invalido", http.StatusBadRequest)
        return
    }

    if err := validarAlfajor(datos, true); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
    }

    // Crear alfajor en la base
    alfajor, err := h.queries.CreateAlfajor(r.Context(), db.CreateAlfajorParams{
        ID:          datos.ID,
        Name:        datos.Name,
        Relleno:     datos.Relleno,
        Cobertura:   nullString(datos.Cobertura),
        Precio:      datos.Precio,
        Tapas:       datos.Tapas,
        Marca:       datos.Marca,
        Descripcion: nullString(datos.Descripcion),
    })

    if err != nil {
        log.Printf("Error al crear alfajor: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(alfajor)
}

// PUT /alfajores/{id}
func (h *Handler) updateAlfajor(w http.ResponseWriter, r *http.Request, id string) {
    var datos AlfajorInput

    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        http.Error(w, "JSON invalido", http.StatusBadRequest)
        return
    }

    if err := validarAlfajor(datos, false); err != nil {
    	http.Error(w, err.Error(), http.StatusBadRequest)
    	return
    }

    alfajor, err := h.queries.UpdateAlfajor(r.Context(), db.UpdateAlfajorParams{
        ID:          id,
        Name:        datos.Name,
        Relleno:     datos.Relleno,
        Cobertura:   nullString(datos.Cobertura),
        Precio:      datos.Precio,
        Tapas:       datos.Tapas,
        Marca:       datos.Marca,
        Descripcion: nullString(datos.Descripcion),
    })

    if errors.Is(err, sql.ErrNoRows) {
        http.Error(w, "Alfajor no encontrado", http.StatusNotFound)
        return
    }

    if err != nil {
        log.Printf("Error al actualizar alfajor: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alfajor)
}

// DELETE /alfajores/{id}
func (h *Handler) deleteAlfajor(w http.ResponseWriter, r *http.Request, id string) {
    _, err := h.queries.GetAlfajor(r.Context(), id)

    if errors.Is(err, sql.ErrNoRows) {
        http.Error(w, "Alfajor no encontrado", http.StatusNotFound)
        return
    }

    if err != nil {
        log.Printf("Error al buscar alfajor: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    err = h.queries.DeleteAlfajor(r.Context(), id)

    if err != nil {
        log.Printf("Error al eliminar alfajor: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
