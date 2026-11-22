package datos

import (
	"database/sql"
	"fmt"
	"regexp"
)

type columna struct {
	tabla      string
	nombre     string
	definicion string
}

type migracion struct {
	version  int
	columnas []columna
}

var identificadorSeguro = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func migraciones() []migracion {
	catalogos := []string{"vehiculos", "choferes", "proveedores", "productos", "origenes", "destinos"}
	var colsCatalogo []columna
	for _, c := range catalogos {
		colsCatalogo = append(colsCatalogo,
			columna{c, "version", "INTEGER NOT NULL DEFAULT 0"},
			columna{c, "pendiente_push", "INTEGER NOT NULL DEFAULT 0"},
		)
	}
	m1 := []columna{
		{"pesajes", "sincronizado", "INTEGER NOT NULL DEFAULT 0"},
		{"pesajes", "sync_event_id", "TEXT"},
		{"pesajes", "sync_error", "TEXT"},
		{"choferes", "telefono", "TEXT"},
	}
	m2 := append(append([]columna{}, colsCatalogo...),
		columna{"pesajes", "origen", "TEXT NOT NULL DEFAULT 'ESCRITORIO'"},
		columna{"pesajes", "web_id", "INTEGER"},
	)
	return []migracion{
		{version: 1, columnas: m1},
		{version: 2, columnas: m2},
	}
}

func aplicarMigraciones(db *sql.DB) error {
	versionActual, err := leerVersionEsquema(db)
	if err != nil {
		return err
	}
	for _, m := range migraciones() {
		if err := agregarColumnasFaltantes(db, m.columnas); err != nil {
			return err
		}
		if m.version > versionActual {
			if err := escribirVersionEsquema(db, m.version); err != nil {
				return err
			}
			versionActual = m.version
		}
	}
	if versionActual < VersionEsquema {
		return escribirVersionEsquema(db, VersionEsquema)
	}
	return nil
}

func leerVersionEsquema(db *sql.DB) (int, error) {
	var v int
	err := db.QueryRow(`PRAGMA user_version`).Scan(&v)
	return v, err
}

func escribirVersionEsquema(db *sql.DB, version int) error {
	_, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", version))
	return err
}

func agregarColumnasFaltantes(db *sql.DB, columnas []columna) error {
	for _, c := range columnas {
		existeTabla, err := existeTabla(db, c.tabla)
		if err != nil || !existeTabla {
			continue
		}
		existeCol, err := existeColumna(db, c.tabla, c.nombre)
		if err != nil || existeCol {
			continue
		}
		tabla := validarIdentificador(c.tabla)
		nombre := validarIdentificador(c.nombre)
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tabla, nombre, c.definicion))
		if err != nil {
			return err
		}
	}
	return nil
}

func existeColumna(db *sql.DB, tabla, columna string) (bool, error) {
	validarIdentificador(tabla)
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tabla))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return false, err
		}
		if equalsFold(name, columna) {
			return true, nil
		}
	}
	return false, rows.Err()
}

func existeTabla(db *sql.DB, tabla string) (bool, error) {
	var n int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, tabla,
	).Scan(&n)
	return n > 0, err
}

func validarIdentificador(id string) string {
	if id == "" || !identificadorSeguro.MatchString(id) {
		panic("identificador SQL no válido: " + id)
	}
	return id
}

func equalsFold(a, b string) bool {
	if len(a) != len(b) {
		// fall through to manual
	}
	return stringFold(a) == stringFold(b)
}

func stringFold(s string) string {
	bs := []byte(s)
	for i, c := range bs {
		if c >= 'A' && c <= 'Z' {
			bs[i] = c + ('a' - 'A')
		}
	}
	return string(bs)
}
