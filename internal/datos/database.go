package datos

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"sistema-pesaje/internal/balanza"
	"sistema-pesaje/internal/clave"

	_ "modernc.org/sqlite"
)

const (
	longitudClaveGenerada = 12
	alfabetoClave         = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"
)

var (
	instanciaMu sync.Mutex
	instancia   *Database
)

// Database abre pesaje.db en SistemaPesaje (paridad Java Database).
type Database struct {
	path              string
	db                *sql.DB
	mu                sync.Mutex
	claveAdministrador string
}

// Conn expone la conexión SQLite (uso interno).
func (d *Database) Conn() *sql.DB {
	return d.db
}

// RutaBase devuelve la ruta del archivo en uso.
func (d *Database) RutaBase() string {
	return d.path
}

// ClaveAdministrador solo tiene valor tras el primer arranque con semilla.
func (d *Database) ClaveAdministrador() string {
	return d.claveAdministrador
}

// ObtenerInstancia abre o reutiliza la base en la ruta resuelta.
func ObtenerInstancia() (*Database, error) {
	instanciaMu.Lock()
	defer instanciaMu.Unlock()
	if instancia != nil {
		return instancia, nil
	}
	db, err := AbrirEn(RutaResuelta())
	if err != nil {
		return nil, err
	}
	instancia = db
	return db, nil
}

// CerrarInstancia libera la instancia global.
func CerrarInstancia() {
	instanciaMu.Lock()
	defer instanciaMu.Unlock()
	if instancia != nil {
		_ = instancia.Close()
		instancia = nil
	}
}

// AbrirEn abre (o crea) la base en la ruta indicada.
func AbrirEn(path string) (*Database, error) {
	path = filepath.Clean(path)
	existia := fileExists(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	d := &Database{path: path, db: sqlDB}
	if _, err := sqlDB.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	if err := d.inicializar(existia); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return d, nil
}

func (d *Database) inicializar(baseExistente bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if baseExistente {
		_ = d.generarBackup()
	}
	if err := d.aplicarEsquemaYMigraciones(); err != nil {
		return err
	}
	if !baseExistente {
		_ = d.generarBackup()
	}
	return nil
}

func (d *Database) aplicarEsquemaYMigraciones() error {
	if err := d.aplicarEsquema(); err != nil {
		return err
	}
	if err := aplicarMigraciones(d.db); err != nil {
		return err
	}
	if err := d.sembrarAdministrador(); err != nil {
		return err
	}
	return d.sembrarConfigBalanzaPruebas()
}

func (d *Database) aplicarEsquema() error {
	for _, sentencia := range DividirSentencias(string(esquemaSQL)) {
		if _, err := d.db.Exec(sentencia); err != nil {
			return fmt.Errorf("esquema: %w", err)
		}
	}
	return nil
}

func (d *Database) sembrarAdministrador() error {
	var n int
	if err := d.db.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE rol = ?`, "administrador").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	pwd := generarClaveAleatoria()
	ahora := time.Now().Format("2006-01-02 15:04:05")
	res, err := d.db.Exec(
		`INSERT INTO personas (nombre, apellido_paterno, correo, estado, fecha_creacion) VALUES (?, ?, ?, ?, ?)`,
		"Administrador", "Sistema", EmailAdministrador, "ACTIVO", ahora,
	)
	if err != nil {
		return err
	}
	personaID, _ := res.LastInsertId()
	_, err = d.db.Exec(
		`INSERT INTO usuarios (persona_id, email, password_hash, rol, estado, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)`,
		personaID, EmailAdministrador, clave.Hash(pwd), "administrador", "ACTIVO", ahora,
	)
	if err != nil {
		return err
	}
	d.claveAdministrador = pwd
	fmt.Fprintf(os.Stderr, "[SistemaPesaje] Usuario administrador creado (%s).\n", EmailAdministrador)
	fmt.Fprintf(os.Stderr, "[SistemaPesaje] Contraseña inicial: %s (cámbiela en el primer inicio).\n", pwd)
	return nil
}

// sembrarConfigBalanzaPruebas deja puerto SIM en bases nuevas sin fila (solo si PESAJE_BALANZA_SIM=1).
func (d *Database) sembrarConfigBalanzaPruebas() error {
	if os.Getenv("PESAJE_BALANZA_SIM") != "1" {
		return nil
	}
	var n int
	if err := d.db.QueryRow(`SELECT COUNT(*) FROM config_balanza`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	cfg := balanza.ConfigPorDefecto()
	cfg.PuertoCOM = "SIM"
	if err := GuardarConfigBalanza(context.Background(), d.db, cfg); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "[SistemaPesaje] Báscula de pruebas: puerto SIM (sin hardware).\n")
	return nil
}

func generarClaveAleatoria() string {
	b := make([]byte, longitudClaveGenerada)
	for i := range b {
		var rb [1]byte
		_, _ = rand.Read(rb[:])
		b[i] = alfabetoClave[int(rb[0])%len(alfabetoClave)]
	}
	return string(b)
}

func (d *Database) generarBackup() error {
	destino := filepath.Join(filepath.Dir(d.path), ArchivoBackup)
	_ = os.Remove(destino)
	if _, err := d.db.Exec(`VACUUM INTO ?`, destino); err != nil {
		return copiarArchivo(d.path, destino)
	}
	return nil
}

func copiarArchivo(origen, destino string) error {
	in, err := os.Open(origen)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.db.Close()
}

// Contar ejecuta SELECT COUNT(*) con un argumento opcional (helper de tests).
func (d *Database) Contar(ctx context.Context, query string, args ...any) (int, error) {
	var n int
	err := d.db.QueryRowContext(ctx, query, args...).Scan(&n)
	return n, err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
