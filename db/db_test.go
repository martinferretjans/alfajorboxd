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

func TestAlfajorCRUD(t *testing.T){
	queries := setupTestDB(t)
	ctx :=context.Background()
	
	nuevoAlf, err:= queries.CreateAlfajor(ctx, CreateAlfajorParams{
		ID: "prueba1",
		Name: "Alfajor de Prueba",
		Relleno: "Dulce de leche",
		Cobertura: "negro",
		Precio: 1500,
		Tapas: sql.NullInt32{Int32: 2, Valid: true},
		Marca: sql.NullString{String: "guaymallen", Valid: true},
	})
	if err != nil{
		t.Fatalf("CreateAlfajor fallo: %v", err)
	}
	t.Cleanup(func(){
		_= queries.DeleteAlfajor(ctx, nuevoAlf.ID)
	})
	//-----------------------------------------------------------------

	alfajor, err := queries.GetAlfajor(ctx, nuevoAlf.ID)
	if err!= nil{
		t.Fatalf("GetAlfajor fallo: %v", err)		
	}
	if alfajor.ID != nuevoAlf.ID {
		t.Errorf("ID esperado %q, obtuve %q", nuevoAlf.ID, alfajor.ID)
	}
	if alfajor.Name != "Alfajor de Prueba" {
		t.Errorf("esperabamos 'Alfajor de Prueba', obtuve %q", alfajor.Name)
	}
	//-----------------------------------------------------------------

	update, err := queries.UpdateAlfajor(ctx, UpdateAlfajorParams{
		ID: alfajor.ID,
		Name: alfajor.Name,
		Relleno: alfajor.Relleno,
		Cobertura: alfajor.Cobertura,
		Precio: alfajor.Precio,
		Tapas: alfajor.Tapas,
		Marca: sql.NullString{String: "fantoche", Valid: true},
		Descripcion: alfajor.Descripcion,})
	if err!= nil{
		t.Fatalf("UpdateAlfajor fallo %v", err)	
	}
	if update.Marca.String != "fantoche"{
		t.Fatalf("updateAlfajor fallo %v", err)
	}
	//-------------------------------------------------------------------
	listaAlf, err := queries.ListAlfajores(ctx)
	if err!= nil{
		t.Fatalf("ListAlfajores fallo %v", err)
	}
	encontro := false
	for _, alf := range listaAlf{
		if alf.ID == nuevoAlf.ID{
			encontro = true
			break
		}
	}
	if !encontro{
		t.Error("no se encontro el alfajor de prueba")
	}
	//------------------------------------------------------------------
	if err := queries.DeleteAlfajor(ctx, nuevoAlf.ID); err != nil {
		t.Fatalf("DeleteAlfajor falló: %v", err)
	}
	_, err = queries.GetAlfajor(ctx, nuevoAlf.ID)
	if  err == nil {
		t.Error("deberia haber dado error, ya que se elimino el alfajor")
	}

}
	

