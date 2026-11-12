package handlers

import (
    "database/sql"
    "errors"
    "log"
    "net/http"
    "strconv"
    "strings"

    db "alfajorboxd/db/sqlc"
    "alfajorboxd/views"
)

type Handler struct {
    queries *db.Queries
}

func New(queries *db.Queries) *Handler {
    return &Handler{queries: queries}
}


type alfajorRequest struct {
	ID          string
	Name        string
	Relleno     string
	Cobertura   string
	Precio      int32
	Tapas       int32
	Marca       string
	Descripcion string
	Imagen      string
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

// Handler para /
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	alfajores, err := h.queries.ListAlfajoresPopulares(r.Context())
	if err != nil {
		log.Printf("Error al listar alfajores: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	reviews, err := h.queries.ListUltimasReviews(r.Context(), 6)
	if err != nil {
		log.Printf("Error al listar reviews: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := views.IndexPage(leerTema(r), alfajores, reviews).Render(r.Context(), w); err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// Handler para /alfajores
func (h *Handler) AlfajoresHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    	case http.MethodGet:
        	http.Redirect(w, r, "/", http.StatusSeeOther)
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
    	default:
        	http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
    }
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

	promedio, err := h.queries.GetPromedioAlfajor(r.Context(), id)
	if err != nil {
		log.Printf("Error al calcular promedio: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	reviews, err := h.queries.ListReviewsByAlfajor(r.Context(), id)
	if err != nil {
		log.Printf("Error al listar reviews del alfajor: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	usuarios, err := h.queries.ListUsuarios(r.Context())
	if err != nil {
		log.Printf("Error al listar usuarios: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := views.AlfajorPage(leerTema(r), alfajor, promedio, reviews, usuarios).Render(r.Context(), w); err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

//POST/alfajores
func(h *Handler) createAlfajor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario invalido", http.StatusBadRequest)
		return
	}

	precio, _ := strconv.Atoi(r.FormValue("precio"))
	tapas, _ := strconv.Atoi(r.FormValue("tapas"))

	alfreq := alfajorRequest{
		ID:          r.FormValue("id"),
		Name:        r.FormValue("name"),
		Relleno:     r.FormValue("relleno"),
		Cobertura:   r.FormValue("cobertura"),
		Precio:      int32(precio),
		Tapas:       int32(tapas),
		Marca:       r.FormValue("marca"),
		Descripcion: r.FormValue("descripcion"),
		Imagen:      r.FormValue("imagen"),
	}
	if err := validarAlfajor(alfreq, true); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.CreateAlfajor(r.Context(), db.CreateAlfajorParams{
		ID:          alfreq.ID,
		Name:        alfreq.Name,
		Relleno:     alfreq.Relleno,
		Cobertura:   toNullString(alfreq.Cobertura),
		Precio:      alfreq.Precio,
		Tapas:       alfreq.Tapas,
		Marca:       alfreq.Marca,
		Descripcion: toNullString(alfreq.Descripcion),
		Imagen:      toNullString(alfreq.Imagen),
	})
	if err != nil {
		log.Printf("Error al crear el alfajor: %v", err)
		http.Error(w, "Error al guardar el alfajor en la base de datos", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
