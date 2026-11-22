package auth_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"sistema-pesaje/internal/auth"
	"sistema-pesaje/internal/clave"
	"sistema-pesaje/internal/datos"
)

func preparar(t *testing.T) (*datos.Database, *auth.Service, func()) {
	t.Helper()
	dir := t.TempDir()
	ruta := filepath.Join(dir, "pesaje.db")
	t.Setenv(datos.PropiedadSistemaRuta, ruta)
	datos.CerrarInstancia()
	base, err := datos.ObtenerInstancia()
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	return base, auth.NewService(base), func() {
		datos.CerrarInstancia()
		_ = os.RemoveAll(dir)
	}
}

func TestElAdministradorSembradoPuedeIniciarSesion(t *testing.T) {
	base, svc, limpiar := preparar(t)
	defer limpiar()

	pwd := base.ClaveAdministrador()
	if pwd == "" {
		t.Fatal("el primer arranque debe sembrar el administrador")
	}

	res := svc.Autenticar(context.Background(), datos.EmailAdministrador, pwd)
	if !res.EsExito() {
		t.Fatalf("estado=%s", res.Estado)
	}
	if res.Usuario == nil || res.Usuario.Email != datos.EmailAdministrador {
		t.Fatal("usuario inesperado")
	}
	if !res.Usuario.EsAdministrador() {
		t.Fatal("debe ser administrador")
	}
}

func TestUnaContrasenaIncorrectaSeRechaza(t *testing.T) {
	base, svc, limpiar := preparar(t)
	defer limpiar()

	pwd := base.ClaveAdministrador()
	res := svc.Autenticar(context.Background(), datos.EmailAdministrador, pwd+"x")
	if res.EsExito() {
		t.Fatal("debería fallar")
	}
	if res.Estado != auth.EstadoContrasenaInvalida {
		t.Fatalf("estado=%s", res.Estado)
	}
}

func TestUnCorreoDesconocidoSeRechaza(t *testing.T) {
	_, svc, limpiar := preparar(t)
	defer limpiar()

	res := svc.Autenticar(context.Background(), "nadie@sistema.com", "lo-que-sea")
	if res.EsExito() {
		t.Fatal("debería fallar")
	}
	if res.Estado != auth.EstadoEmailInvalido {
		t.Fatalf("estado=%s", res.Estado)
	}
}

func TestUnUsuarioInactivoNoPuedeIniciarSesion(t *testing.T) {
	base, svc, limpiar := preparar(t)
	defer limpiar()

	crearUsuario(t, base, "inactivo@sistema.com", "clave123", "INACTIVO")
	res := svc.Autenticar(context.Background(), "inactivo@sistema.com", "clave123")
	if res.EsExito() {
		t.Fatal("debería fallar")
	}
	if res.Estado != auth.EstadoUsuarioInactivo {
		t.Fatalf("estado=%s", res.Estado)
	}
}

func TestLaVerificacionDeContrasenaSirveParaCambiarDeModo(t *testing.T) {
	base, svc, limpiar := preparar(t)
	defer limpiar()

	pwd := base.ClaveAdministrador()
	ctx := context.Background()
	if !svc.VerificarContrasena(ctx, datos.EmailAdministrador, pwd) {
		t.Fatal("contraseña válida")
	}
	if svc.VerificarContrasena(ctx, datos.EmailAdministrador, "otra-clave") {
		t.Fatal("contraseña inválida")
	}
	if svc.VerificarContrasena(ctx, "nadie@sistema.com", pwd) {
		t.Fatal("email desconocido")
	}
}

func TestMigraHashSha1APbkdf2AlIniciarSesion(t *testing.T) {
	base, svc, limpiar := preparar(t)
	defer limpiar()

	legado := "sha1:NTU3YzU1NWNlNzg1ZGViMDZjNjA0NjhhZjViYzk3YTUxMDMyYjFiOA==:Zlg1akJleTdqbmRWaU44VA=="
	res, err := base.Conn().Exec(
		`INSERT INTO personas (nombre, correo, estado, fecha_creacion) VALUES (?, ?, ?, ?)`,
		"Legacy", "legacy@sistema.com", "ACTIVO", "2026-11-20 10:00:00",
	)
	if err != nil {
		t.Fatal(err)
	}
	personaID, _ := res.LastInsertId()
	_, err = base.Conn().Exec(
		`INSERT INTO usuarios (persona_id, email, password_hash, rol, estado, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)`,
		personaID, "legacy@sistema.com", legado, "usuario", "ACTIVO", "2026-11-20 10:00:00",
	)
	if err != nil {
		t.Fatal(err)
	}

	authRes := svc.Autenticar(context.Background(), "legacy@sistema.com", "admin123")
	if !authRes.EsExito() {
		t.Fatalf("estado=%s", authRes.Estado)
	}

	var hash string
	err = base.Conn().QueryRow(`SELECT password_hash FROM usuarios WHERE email = ?`, "legacy@sistema.com").Scan(&hash)
	if err != nil {
		t.Fatal(err)
	}
	if !clave.EsHashModerno(hash) {
		t.Fatal("debe migrar a pbkdf2")
	}
	if !clave.Verificar("admin123", hash) {
		t.Fatal("hash migrado debe verificar")
	}
}

func TestElUltimoLoginQuedaRegistrado(t *testing.T) {
	base, svc, limpiar := preparar(t)
	defer limpiar()

	pwd := base.ClaveAdministrador()
	_ = svc.Autenticar(context.Background(), datos.EmailAdministrador, pwd)

	var ultimo sql.NullInt64
	err := base.Conn().QueryRow(`SELECT ultimo_login FROM usuarios WHERE email = ?`, datos.EmailAdministrador).Scan(&ultimo)
	if err != nil {
		t.Fatal(err)
	}
	if !ultimo.Valid || ultimo.Int64 <= 0 {
		t.Fatal("ultimo login debe quedar registrado")
	}
}

func crearUsuario(t *testing.T, base *datos.Database, email, contrasena, estado string) {
	t.Helper()
	res, err := base.Conn().Exec(
		`INSERT INTO personas (nombre, correo, estado, fecha_creacion) VALUES (?, ?, ?, ?)`,
		"Usuario", email, "ACTIVO", "2026-11-20 10:00:00",
	)
	if err != nil {
		t.Fatal(err)
	}
	personaID, _ := res.LastInsertId()
	_, err = base.Conn().Exec(
		`INSERT INTO usuarios (persona_id, email, password_hash, rol, estado, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)`,
		personaID, email, clave.Hash(contrasena), "usuario", estado, "2026-11-20 10:00:00",
	)
	if err != nil {
		t.Fatal(err)
	}
}
