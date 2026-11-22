package datos

import (
	"context"
	"database/sql"
)

// BoletaConfig fila boleta_config.
type BoletaConfig struct {
	ID              int64
	Salida          string
	Direccion       string
	Celular1        string
	Celular2        string
	Ciudad          string
	TituloSuperior  string
	Marca           string
	TituloDocumento string
	Acreditacion    string
}

func ObtenerBoletaConfig(ctx context.Context, db *sql.DB) (*BoletaConfig, error) {
	row := db.QueryRowContext(ctx, `
SELECT id, salida, direccion, celular1, celular2, ciudad, titulo_superior, marca, titulo_documento, acreditacion
FROM boleta_config ORDER BY id DESC LIMIT 1`)
	var b BoletaConfig
	err := row.Scan(&b.ID, &b.Salida, &b.Direccion, &b.Celular1, &b.Celular2, &b.Ciudad,
		&b.TituloSuperior, &b.Marca, &b.TituloDocumento, &b.Acreditacion)
	if err == sql.ErrNoRows {
		return &BoletaConfig{
			Salida: "Salida a Potosi", Direccion: "Tarija km 2",
			TituloSuperior: "BOLETA DE PESAJE", Marca: "SISTEMA DE PESAJE",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func GuardarBoletaConfig(ctx context.Context, db *sql.DB, b BoletaConfig) error {
	fecha := Ahora()
	if b.ID <= 0 {
		act, err := ObtenerBoletaConfig(ctx, db)
		if err != nil {
			return err
		}
		if act != nil && act.ID > 0 {
			b.ID = act.ID
		}
	}
	if b.ID <= 0 {
		_, err := db.ExecContext(ctx, `
INSERT INTO boleta_config (salida, direccion, celular1, celular2, ciudad, titulo_superior, marca,
 titulo_documento, acreditacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			b.Salida, b.Direccion, b.Celular1, b.Celular2, b.Ciudad,
			b.TituloSuperior, b.Marca, b.TituloDocumento, b.Acreditacion, fecha)
		return err
	}
	_, err := db.ExecContext(ctx, `
UPDATE boleta_config SET salida=?, direccion=?, celular1=?, celular2=?, ciudad=?, titulo_superior=?,
 marca=?, titulo_documento=?, acreditacion=?, fecha_modificacion=? WHERE id=?`,
		b.Salida, b.Direccion, b.Celular1, b.Celular2, b.Ciudad,
		b.TituloSuperior, b.Marca, b.TituloDocumento, b.Acreditacion, fecha, b.ID)
	return err
}