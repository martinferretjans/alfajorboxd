package db
import (
	"context"
	"database/sql"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *Queries{
	dsn := os.Getenv("DATABASE_URL")
	if dsn == ""{
		t.Skip("el DATABASE_URL no esta seteado")
	}

	conexion, err:= sql.Open("postgres", dsn)
	if err!= nil{
		t.Fatalf("no se puedo abrir la conexion: %v", err)

	}
	t.Cleanup(func() {conexion.Close()})
	return New(conexion)
}

func crearUsuarioDePrueba(t *testing.T, q *Queries, ctx context.Context, email string) Usuario {
	t.Helper()

	usu, err := q.CreateUsuario(ctx, CreateUsuarioParams{
		Name:  "Usuario de Prueba",
		Email: email,
	})
	if err != nil {
		t.Fatalf("no se pudo crear el usuario de prueba: %v", err)
	}
	t.Cleanup(func() {
		_ = q.DeleteUsuario(ctx, usu.ID)
	})
	return usu
}

func crearAlfajorDePrueba(t *testing.T, q *Queries, ctx context.Context, id string) Alfajor {
	t.Helper()

	alf, err := q.CreateAlfajor(ctx, CreateAlfajorParams{
		ID:        id,
		Name:      "Alfajor de Prueba",
		Relleno:   "Dulce de leche",
		Cobertura: sql.NullString{String: "negro", Valid: true},
		Precio:    1500,
		Tapas:     2,
		Marca:     "guaymallen",
	})
	if err != nil {
		t.Fatalf("no se pudo crear el alfajor de prueba: %v", err)
	}
	t.Cleanup(func() {
		_ = q.DeleteAlfajor(ctx, alf.ID)
	})
	return alf
}

func TestCreateAlfajor(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	nuevoAlf, err := q.CreateAlfajor(ctx, CreateAlfajorParams{
		ID:        "prueba-create",
		Name:      "Alfajor de Prueba",
		Relleno:   "Dulce de leche",
		Cobertura: sql.NullString{String: "negro", Valid: true},
		Precio:    1500,
		Tapas:     2,
		Marca:     "guaymallen",
	})
	if err != nil {
		t.Fatalf("CreateAlfajor fallo: %v", err)
	}
	t.Cleanup(func() {
		_ = q.DeleteAlfajor(ctx, nuevoAlf.ID)
	})

	if nuevoAlf.Name != "Alfajor de Prueba" {
		t.Errorf("Name esperado 'Alfajor de Prueba', obtuve %q", nuevoAlf.Name)
	}
}

func TestGetAlfajor(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearAlfajorDePrueba(t, q, ctx, "prueba-get")

	alfajor, err := q.GetAlfajor(ctx, prueba.ID)
	if err != nil {
		t.Fatalf("GetAlfajor fallo: %v", err)
	}
	if alfajor.ID != prueba.ID {
		t.Errorf("ID esperado %q, obtuve %q", prueba.ID, alfajor.ID)
	}
	if alfajor.Name != prueba.Name {
		t.Errorf("Name esperado %q, obtuve %q", prueba.Name, alfajor.Name)
	}
}

func TestUpdateAlfajor(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearAlfajorDePrueba(t, q, ctx, "prueba-update")
 
	actualizado, err := q.UpdateAlfajor(ctx, UpdateAlfajorParams{
		ID:          prueba.ID,
		Name:        prueba.Name,
		Relleno:     prueba.Relleno,
		Cobertura:   prueba.Cobertura,
		Precio:      prueba.Precio,
		Tapas:       prueba.Tapas,
		Marca:       "fantoche",
		Descripcion: prueba.Descripcion,
		Imagen:      sql.NullString{String: "fantoche.jpg", Valid: true},
	})
	if err != nil {
		t.Fatalf("UpdateAlfajor fallo: %v", err)
	}
	if actualizado.Marca != "fantoche" {
		t.Errorf("Marca esperada %q, obtuve %q", "fantoche", actualizado.Marca)
	}
	if actualizado.Imagen.String != "fantoche.jpg" {
		t.Errorf("Imagen esperada %q, obtuve %q", "fantoche.jpg", actualizado.Imagen.String)
	}
}
func TestListAlfajores(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearAlfajorDePrueba(t, q, ctx, "prueba-list")

	lista, err := q.ListAlfajores(ctx)
	if err != nil {
		t.Fatalf("ListAlfajores fallo: %v", err)
	}
	encontrado := false
	for _, a := range lista {
		if a.ID == prueba.ID {
			encontrado = true
			break
		}
	}
	if !encontrado {
		t.Error("no se encontro el alfajor de prueba en ListAlfajores")
	}
}
func TestDeleteAlfajor(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearAlfajorDePrueba(t, q, ctx, "prueba-delete")

	if err := q.DeleteAlfajor(ctx, prueba.ID); err != nil {
		t.Fatalf("DeleteAlfajor fallo: %v", err)
	}
	if _, err := q.GetAlfajor(ctx, prueba.ID); err == nil {
		t.Error("deberia haber dado error, ya que se elimino el alfajor")
	}
}

func TestCreateUsuario(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearUsuarioDePrueba(t, q, ctx, "prueba-create@test.com")

	if prueba.Email != "prueba-create@test.com" {
		t.Errorf("Email esperado %q, obtuve %q", "prueba-create@test.com", prueba.Email)
	}
	if prueba.ID == 0 {
		t.Error("el id del usuario deberia haberse generado")
	}
}

func TestGetUsuario(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearUsuarioDePrueba(t, q, ctx, "prueba-get@test.com")

	usuario, err := q.GetUsuario(ctx, prueba.ID)
	if err != nil {
		t.Fatalf("GetUsuario fallo: %v", err)
	}
	if usuario.Name != prueba.Name {
		t.Errorf("Name esperado %q, obtuve %q", prueba.Name, usuario.Name)
	}
}

func TestListUsuarios(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearUsuarioDePrueba(t, q, ctx, "prueba-list@test.com")

	lista, err := q.ListUsuarios(ctx)
	if err != nil {
		t.Fatalf("ListUsuarios fallo: %v", err)
	}
	encontrado := false
	for _, u := range lista {
		if u.ID == prueba.ID {
			encontrado = true
			break
		}
	}
	if !encontrado {
		t.Error("no se encontro el usuario de prueba en ListUsuarios")
	}
}

func TestDeleteUsuario(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	prueba := crearUsuarioDePrueba(t, q, ctx, "prueba-delete@test.com")

	if err := q.DeleteUsuario(ctx, prueba.ID); err != nil {
		t.Fatalf("DeleteUsuario fallo: %v", err)
	}
	if _, err := q.GetUsuario(ctx, prueba.ID); err == nil {
		t.Error("deberia haber dado error, ya que se elimino el usuario")
	}
}

func TestCreateReview(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-review-create")
	usu := crearUsuarioDePrueba(t, q, ctx, "prueba-review-create@test.com")

	review, err := q.CreateReview(ctx, CreateReviewParams{
		IDAlfajor:    alf.ID,
		IDUser:       usu.ID,
		Calificacion: 7,
		Comentario:   sql.NullString{String: "Muy rico", Valid: true},
	})
	if err != nil {
		t.Fatalf("CreateReview fallo: %v", err)
	}
	if review.Calificacion != 7 {
		t.Errorf("Calificacion esperada %d, obtuve %d", 7, review.Calificacion)
	}
	if review.Comentario.String != "Muy rico" {
		t.Errorf("Comentario esperado %q, obtuve %q", "Muy rico", review.Comentario.String)
	}
}

func TestCreateReviewCalificacionInvalida(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-review-invalida")
	usu := crearUsuarioDePrueba(t, q, ctx, "prueba-review-invalida@test.com")

	_, err := q.CreateReview(ctx, CreateReviewParams{
		IDAlfajor:    alf.ID,
		IDUser:       usu.ID,
		Calificacion: 11,
	})
	if err == nil {
		t.Error("deberia haber dado error, ya que la calificacion maxima es 10")
	}
}

func TestGetReview(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-review-get")
	usu := crearUsuarioDePrueba(t, q, ctx, "prueba-review-get@test.com")

	if _, err := q.CreateReview(ctx, CreateReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID, Calificacion: 4}); err != nil {
		t.Fatalf("CreateReview fallo: %v", err)
	}

	review, err := q.GetReview(ctx, GetReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID})
	if err != nil {
		t.Fatalf("GetReview fallo: %v", err)
	}
	if review.Calificacion != 4 {
		t.Errorf("Calificacion esperada %d, obtuve %d", 4, review.Calificacion)
	}
}

func TestListReviewsByAlfajor(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-review-list")
	usu := crearUsuarioDePrueba(t, q, ctx, "prueba-review-list@test.com")

	if _, err := q.CreateReview(ctx, CreateReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID, Calificacion: 8}); err != nil {
		t.Fatalf("CreateReview fallo: %v", err)
	}

	lista, err := q.ListReviewsByAlfajor(ctx, alf.ID)
	if err != nil {
		t.Fatalf("ListReviewsByAlfajor fallo: %v", err)
	}
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 review, obtuve %d", len(lista))
	}
	if lista[0].Usuario != usu.Name {
		t.Errorf("Usuario esperado %q, obtuve %q", usu.Name, lista[0].Usuario)
	}
}

func TestListUltimasReviews(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-review-ultimas")
	usu := crearUsuarioDePrueba(t, q, ctx, "prueba-review-ultimas@test.com")

	if _, err := q.CreateReview(ctx, CreateReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID, Calificacion: 5}); err != nil {
		t.Fatalf("CreateReview fallo: %v", err)
	}

	lista, err := q.ListUltimasReviews(ctx, 100)
	if err != nil {
		t.Fatalf("ListUltimasReviews fallo: %v", err)
	}
	encontrado := false
	for _, r := range lista {
		if r.IDAlfajor == alf.ID && r.IDUser == usu.ID {
			encontrado = true
			if r.Alfajor != alf.Name {
				t.Errorf("Alfajor esperado %q, obtuve %q", alf.Name, r.Alfajor)
			}
			break
		}
	}
	if !encontrado {
		t.Error("no se encontro la review de prueba en ListUltimasReviews")
	}
}

func TestGetPromedioAlfajor(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-promedio")
	usu1 := crearUsuarioDePrueba(t, q, ctx, "prueba-promedio1@test.com")
	usu2 := crearUsuarioDePrueba(t, q, ctx, "prueba-promedio2@test.com")
	usu3 := crearUsuarioDePrueba(t, q, ctx, "prueba-promedio3@test.com")

	//10, 9 y 9 dan 9.33, que en estrellas es 4.67
	for i, usu := range []Usuario{usu1, usu2, usu3} {
		calificacion := int32(9)
		if i == 0 {
			calificacion = 10
		}
		if _, err := q.CreateReview(ctx, CreateReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID, Calificacion: calificacion}); err != nil {
			t.Fatalf("CreateReview fallo: %v", err)
		}
	}

	promedio, err := q.GetPromedioAlfajor(ctx, alf.ID)
	if err != nil {
		t.Fatalf("GetPromedioAlfajor fallo: %v", err)
	}
	if promedio != 4.67 {
		t.Errorf("Promedio esperado %v, obtuve %v", 4.67, promedio)
	}
}

func TestListAlfajoresPopulares(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-populares")

	lista, err := q.ListAlfajoresPopulares(ctx)
	if err != nil {
		t.Fatalf("ListAlfajoresPopulares fallo: %v", err)
	}
	encontrado := false
	for _, a := range lista {
		if a.ID == alf.ID {
			encontrado = true
			if a.Cantidad != 0 || a.Promedio != 0 {
				t.Errorf("un alfajor sin reviews deberia tener promedio 0, obtuve %v con %d reviews", a.Promedio, a.Cantidad)
			}
			break
		}
	}
	if !encontrado {
		t.Error("no se encontro el alfajor de prueba en ListAlfajoresPopulares")
	}
}

func TestDeleteReview(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()
	alf := crearAlfajorDePrueba(t, q, ctx, "prueba-review-delete")
	usu := crearUsuarioDePrueba(t, q, ctx, "prueba-review-delete@test.com")

	if _, err := q.CreateReview(ctx, CreateReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID, Calificacion: 6}); err != nil {
		t.Fatalf("CreateReview fallo: %v", err)
	}
	if err := q.DeleteReview(ctx, DeleteReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID}); err != nil {
		t.Fatalf("DeleteReview fallo: %v", err)
	}
	if _, err := q.GetReview(ctx, GetReviewParams{IDAlfajor: alf.ID, IDUser: usu.ID}); err == nil {
		t.Error("deberia haber dado error, ya que se elimino la review")
	}
}
