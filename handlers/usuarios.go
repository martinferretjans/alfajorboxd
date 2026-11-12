package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	db "alfajorboxd/db/sqlc"
	"alfajorboxd/views"
)

func validarUsuario(name string, email string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if !strings.Contains(email, "@") {
		return errors.New("el email no es valido")
	}
	return nil
}

// Handler para /usuarios
func (h *Handler) UsuariosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsuarios(w, r)
	case http.MethodPost:
		h.createUsuario(w, r)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /usuarios
func (h *Handler) getUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.queries.ListUsuarios(r.Context())
	if err != nil {
		log.Printf("Error al listar usuarios: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := views.UsuariosPage(leerTema(r), usuarios).Render(r.Context(), w); err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// POST /usuarios
func (h *Handler) createUsuario(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario invalido", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	if err := validarUsuario(name, email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.CreateUsuario(r.Context(), db.CreateUsuarioParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		log.Printf("Error al crear el usuario: %v", err)
		http.Error(w, "Error al guardar el usuario en la base de datos", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
