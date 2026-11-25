package boleta

import (
	"context"
	"database/sql"
	"strings"

	"sistema-pesaje/internal/datos"
)

func logoEmpresaActiva(ctx context.Context, db *sql.DB) []byte {
	if db == nil {
		return nil
	}
	var logo sql.NullString
	err := db.QueryRowContext(ctx, `
SELECT logo FROM empresas WHERE estado = 'ACTIVO' ORDER BY id DESC LIMIT 1`).Scan(&logo)
	if err != nil || !logo.Valid || strings.TrimSpace(logo.String) == "" {
		return nil
	}
	return datos.DecodificarLogoDataURL(logo.String)
}

// ResolverLogoParaPDF empresa activa o logo del sistema.
func ResolverLogoParaPDF(ctx context.Context, db *sql.DB) []byte {
	if b := logoEmpresaActiva(ctx, db); len(b) > 0 {
		return b
	}
	return datos.LogoSistemaJPEG()
}