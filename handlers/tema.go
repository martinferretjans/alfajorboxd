package handlers

import (
	"net/http"
)

//devuelve el tema guardado en la cookie, por defecto havanna
func leerTema(r *http.Request) string {
	cookie, err := r.Cookie("tema")
	if err != nil || cookie.Value != "marron" {
		return "havanna"
	}
	return cookie.Value
}

// POST /tema
func (h *Handler) TemaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	tema := r.FormValue("tema")
	if tema != "marron" {
		tema = "havanna"
	}
	http.SetCookie(w, &http.Cookie{Name: "tema", Value: tema, Path: "/", MaxAge: 31536000})

	volver := r.Header.Get("Referer")
	if volver == "" {
		volver = "/"
	}
	http.Redirect(w, r, volver, http.StatusSeeOther)
}
