package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	db "alfajorboxd/db/sqlc"
	"alfajorboxd/views"
)

func validarReview(calificacion int) error {
	if calificacion < 1 || calificacion > 10 {
		return errors.New("la calificacion debe estar entre media estrella y 5 estrellas")
	}
	return nil
}

// Handler para /reviews
func (h *Handler) ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getReviews(w, r)
	case http.MethodPost:
		h.createReview(w, r)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /reviews
func (h *Handler) getReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.queries.ListUltimasReviews(r.Context(), 50)
	if err != nil {
		log.Printf("Error al listar reviews: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := views.ReviewsPage(leerTema(r), reviews).Render(r.Context(), w); err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// POST /reviews
func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario invalido", http.StatusBadRequest)
		return
	}

	idAlfajor := r.FormValue("id_alfajor")
	idUser, _ := strconv.Atoi(r.FormValue("id_user"))
	calificacion, _ := strconv.Atoi(r.FormValue("calificacion"))

	if err := validarReview(calificacion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetReview(r.Context(), db.GetReviewParams{
		IDAlfajor: idAlfajor,
		IDUser:    int32(idUser),
	})
	if err == nil {
		http.Error(w, "El usuario ya hizo una review de este alfajor", http.StatusConflict)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Error al buscar review: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	_, err = h.queries.CreateReview(r.Context(), db.CreateReviewParams{
		IDAlfajor:    idAlfajor,
		IDUser:       int32(idUser),
		Calificacion: int32(calificacion),
		Comentario:   toNullString(r.FormValue("comentario")),
	})
	if err != nil {
		log.Printf("Error al crear la review: %v", err)
		http.Error(w, "Error al guardar la review en la base de datos", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/alfajores/"+idAlfajor, http.StatusSeeOther)
}
