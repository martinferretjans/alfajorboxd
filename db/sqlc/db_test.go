package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *Queries {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("el DATABASE_URL no esta seteado")
	}

	conexion, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("no se puedo abrir la conexion: %v", err)
	}

	t.Cleanup(func() {
		_ = conexion.Close()
	})

	return New(conexion)
}

func TestAlfajorCRUD(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	const testID = "prueba1"
	if err := queries.DeleteAlfajor(ctx, testID); err != nil {
		t.Fatalf("no se pudo limpiar el alfajor de una prueba anterior: %v", err)
	}

	nuevoAlf, err := queries.CreateAlfajor(ctx, CreateAlfajorParams{
		ID:          testID,
		Name:        "Alfajor de Prueba",
		Relleno:     "Dulce de leche",
		Cobertura:   sql.NullString{String: "negro", Valid: true},
		Precio:      1500,
		Tapas:       2,
		Marca:       "guaymallen",
		Descripcion: sql.NullString{String: "Alfajor usado para probar el CRUD", Valid: true},
	})
	if err != nil {
		t.Fatalf("CreateAlfajor fallo: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteAlfajor(ctx, nuevoAlf.ID)
	})

	alfajor, err := queries.GetAlfajor(ctx, nuevoAlf.ID)
	if err != nil {
		t.Fatalf("GetAlfajor fallo: %v", err)
	}
	if alfajor.ID != nuevoAlf.ID {
		t.Errorf("ID esperado %q, obtuve %q", nuevoAlf.ID, alfajor.ID)
	}
	if alfajor.Name != "Alfajor de Prueba" {
		t.Errorf("esperabamos 'Alfajor de Prueba', obtuve %q", alfajor.Name)
	}

	update, err := queries.UpdateAlfajor(ctx, UpdateAlfajorParams{
		ID:          alfajor.ID,
		Name:        alfajor.Name,
		Relleno:     alfajor.Relleno,
		Cobertura:   alfajor.Cobertura,
		Precio:      alfajor.Precio,
		Tapas:       alfajor.Tapas,
		Marca:       "fantoche",
		Descripcion: alfajor.Descripcion,
	})
	if err != nil {
		t.Fatalf("UpdateAlfajor fallo: %v", err)
	}
	if update.Marca != "fantoche" {
		t.Errorf("Marca esperada %q, obtuve %q", "fantoche", update.Marca)
	}

	listaAlf, err := queries.ListAlfajores(ctx)
	if err != nil {
		t.Fatalf("ListAlfajores fallo: %v", err)
	}

	encontro := false
	for _, alf := range listaAlf {
		if alf.ID == nuevoAlf.ID {
			encontro = true
			break
		}
	}
	if !encontro {
		t.Error("no se encontro el alfajor de prueba")
	}

	if err := queries.DeleteAlfajor(ctx, nuevoAlf.ID); err != nil {
		t.Fatalf("DeleteAlfajor fallo: %v", err)
	}

	_, err = queries.GetAlfajor(ctx, nuevoAlf.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("GetAlfajor devolvio %v; se esperaba sql.ErrNoRows", err)
	}
}
