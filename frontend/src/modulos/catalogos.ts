import type { ModuloId } from '../navegacion';
import { porId } from '../navegacion';
import { montarAbm, type ConfigAbm } from './abm';

type ConfigAbmBase = Omit<ConfigAbm, 'tituloModulo'>;

const catalogoSimple = (entidad: string, titulo: string, tituloEntidad: string): ConfigAbmBase => ({
  entidad,
  tituloEntidad,
  columnas: [
    { clave: 'nombre', titulo: titulo },
    { clave: 'descripcion', titulo: 'Descripción' },
  ],
  campos: [
    { clave: 'nombre', etiqueta: 'Nombre', requerido: true },
    { clave: 'descripcion', etiqueta: 'Descripción', anchoCompleto: true },
    { clave: 'estado', etiqueta: 'Estado', tipo: 'estado', soloEdicion: true },
  ],
});

const CONFIGS: Partial<Record<ModuloId, ConfigAbmBase>> = {
  productos: catalogoSimple('productos', 'Nombre', 'producto'),
  origenes: catalogoSimple('origenes', 'Nombre', 'origen'),
  destinos: catalogoSimple('destinos', 'Nombre', 'destino'),
  vehiculos: {
    entidad: 'vehiculos',
    tituloEntidad: 'vehículo',
    columnas: [
      { clave: 'placa', titulo: 'Placa' },
      { clave: 'tipo', titulo: 'Tipo' },
      { clave: 'tara', titulo: 'Tara' },
    ],
    campos: [
      { clave: 'placa', etiqueta: 'Placa', requerido: true },
      { clave: 'tipo', etiqueta: 'Tipo' },
      { clave: 'tara', etiqueta: 'Tara (kg)' },
      { clave: 'estado', etiqueta: 'Estado', tipo: 'estado', soloEdicion: true },
    ],
  },
  choferes: {
    entidad: 'choferes',
    tituloEntidad: 'chofer',
    columnas: [
      { clave: 'nombre_completo', titulo: 'Nombre' },
      { clave: 'ci', titulo: 'CI' },
      { clave: 'telefono', titulo: 'Teléfono' },
      { clave: 'licencia', titulo: 'Licencia' },
    ],
    campos: [
      { clave: 'nombre', etiqueta: 'Nombre', requerido: true },
      { clave: 'apellido_paterno', etiqueta: 'Ap. paterno' },
      { clave: 'apellido_materno', etiqueta: 'Ap. materno' },
      { clave: 'ci', etiqueta: 'CI' },
      { clave: 'telefono', etiqueta: 'Teléfono' },
      { clave: 'licencia', etiqueta: 'Licencia', requerido: true },
      { clave: 'estado', etiqueta: 'Estado', tipo: 'estado', soloEdicion: true },
    ],
  },
  proveedores: {
    entidad: 'proveedores',
    tituloEntidad: 'proveedor',
    columnas: [
      { clave: 'nombre_completo', titulo: 'Nombre' },
      { clave: 'ci', titulo: 'CI' },
      { clave: 'telefono', titulo: 'Teléfono' },
      { clave: 'empresa', titulo: 'Empresa' },
      { clave: 'descripcion', titulo: 'Descripción' },
    ],
    campos: [
      { clave: 'nombre', etiqueta: 'Nombre' },
      { clave: 'apellido_paterno', etiqueta: 'Ap. paterno' },
      { clave: 'apellido_materno', etiqueta: 'Ap. materno' },
      { clave: 'ci', etiqueta: 'CI' },
      { clave: 'telefono', etiqueta: 'Teléfono' },
      { clave: 'empresa', etiqueta: 'Empresa', requerido: true },
      { clave: 'descripcion', etiqueta: 'Descripción', anchoCompleto: true },
      { clave: 'estado', etiqueta: 'Estado', tipo: 'estado', soloEdicion: true },
    ],
  },
  empresas: {
    entidad: 'empresas',
    tituloEntidad: 'empresa',
    columnas: [
      { clave: 'nombre', titulo: 'Empresa' },
      { clave: 'actividad', titulo: 'Actividad' },
      { clave: 'correo', titulo: 'Correo' },
      { clave: 'telefono', titulo: 'Teléfono' },
    ],
    campos: [
      { clave: 'nombre', etiqueta: 'Nombre', requerido: true },
      { clave: 'actividad', etiqueta: 'Actividad' },
      { clave: 'correo', etiqueta: 'Correo' },
      { clave: 'telefono', etiqueta: 'Teléfono' },
      { clave: 'estado', etiqueta: 'Estado', tipo: 'estado', soloEdicion: true },
    ],
  },
  usuarios: {
    entidad: 'usuarios',
    tituloEntidad: 'usuario',
    columnas: [
      { clave: 'ci', titulo: 'Nro. documento' },
      { clave: 'nombre_completo', titulo: 'Nombres' },
      { clave: 'correo', titulo: 'Correo electrónico' },
      { clave: 'telefono', titulo: 'Teléfono' },
      { clave: 'rol', titulo: 'Rol' },
    ],
    campos: [
      { clave: 'nombre', etiqueta: 'Nombre', requerido: true },
      { clave: 'apellido_paterno', etiqueta: 'Ap. paterno' },
      { clave: 'apellido_materno', etiqueta: 'Ap. materno' },
      { clave: 'ci', etiqueta: 'CI' },
      { clave: 'telefono', etiqueta: 'Teléfono' },
      { clave: 'correo', etiqueta: 'Correo', requerido: true },
      { clave: 'clave', etiqueta: 'Contraseña', tipo: 'password', requerido: true },
      { clave: 'rol', etiqueta: 'Rol', tipo: 'select', opciones: ['usuario', 'administrador'] },
      { clave: 'estado', etiqueta: 'Estado', tipo: 'estado', soloEdicion: true },
    ],
  },
};

export function montarCatalogo(host: HTMLElement, id: ModuloId): (() => void) | undefined {
  const cfg = CONFIGS[id];
  if (!cfg) return undefined;
  const tituloModulo = porId(id)?.titulo ?? cfg.entidad;
  return montarAbm(host, { ...cfg, tituloModulo });
}