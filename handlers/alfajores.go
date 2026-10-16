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


type alfajorRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Relleno     string `json:"relleno"`
	Cobertura   string `json:"cobertura"`
	Precio      int32  `json:"precio"`
	Tapas       int32  `json:"tapas"`
	Marca       string `json:"marca"`
	Descripcion string `json:"descripcion"`
}

func validarAlfajor(req alfajorRequest, validarID bool) error {
	if validarID && strings.TrimSpace(req.ID) == "" {
		return errors.New("el id es obligatorio")
	}
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("el nombre (name) no puede estar vacío")
	}
	if strings.TrimSpace(req.Relleno) == "" {
		return errors.New("el relleno no puede estar vacío")
	}
	if strings.TrimSpace(req.Marca) == "" {
		return errors.New("la marca no puede estar vacía")
	}
	if req.Precio <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}
	if req.Tapas != 2 && req.Tapas != 3 {
		return errors.New("la cantidad de tapas debe ser 2 o 3")
	}
	return nil
}

func toNullString(s string) sql.NullString {
	if strings.TrimSpace(s) == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
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

//POST/alfajores
func(h *Handler) createAlfajor(w http.ResponseWriter, r *http.Request) {
	var alfreq alfajorRequest
	if err := json.NewDecoder(r.Body).Decode(&alfreq); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}
	if err := validarAlfajor(alfreq, true); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nuevoAlfajor, err := h.queries.CreateAlfajor(r.Context(), db.CreateAlfajorParams{
		ID:          alfreq.ID,
		Name:        alfreq.Name,
		Relleno:     alfreq.Relleno,
		Cobertura:   toNullString(alfreq.Cobertura),
		Precio:      alfreq.Precio,
		Tapas:       alfreq.Tapas,
		Marca:       alfreq.Marca,
		Descripcion: toNullString(alfreq.Descripcion),
	})
	if err != nil {
		log.Printf("Error al crear el alfajor: %v", err)
		http.Error(w, "Error al guardar el alfajor en la base de datos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoAlfajor)
}


func (h *Handler) updateAlfajor(w http.ResponseWriter, r *http.Request, id string) {
	var alfreq alfajorRequest
	if err := json.NewDecoder(r.Body).Decode(&alfreq); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if err := validarAlfajor(alfreq, false); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actualizado, err := h.queries.UpdateAlfajor(r.Context(), db.UpdateAlfajorParams{
		ID:          id,
		Name:        alfreq.Name,
		Relleno:     alfreq.Relleno,
		Cobertura:   toNullString(alfreq.Cobertura),
		Precio:      alfreq.Precio,
		Tapas:       alfreq.Tapas,
		Marca:       alfreq.Marca,
		Descripcion: toNullString(alfreq.Descripcion),
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actualizado)
}

// DELETE /alfajores/{id}
func (h *Handler) deleteAlfajor(w http.ResponseWriter, r *http.Request, id string) {
	_, err := h.queries.GetAlfajor(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Alfajor no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error al buscar alfajor para eliminar: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeleteAlfajor(r.Context(), id); err != nil {
		log.Printf("Error al eliminar alfajor: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
