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
	})
	if err != nil {
		t.Fatalf("UpdateAlfajor fallo: %v", err)
	}
	if actualizado.Marca != "fantoche" {
		t.Errorf("Marca esperada %q, obtuve %q", "fantoche", actualizado.Marca)
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
