export type ModuloId =
  | 'inicio'
  | 'usuarios'
  | 'empresas'
  | 'choferes'
  | 'proveedores'
  | 'pesaje'
  | 'catalogo'
  | 'vehiculos'
  | 'productos'
  | 'origenes'
  | 'destinos'
  | 'reportes'
  | 'config'
  | 'boleta'
  | 'balanza';

export type ModuloNav = {
  id: ModuloId;
  titulo: string;
  grupo?: string;
  pendienteEn?: string;
};

export const MODULOS: ModuloNav[] = [
  { id: 'inicio', titulo: 'Inicio', pendienteEn: 'T9' },
  { id: 'usuarios', titulo: 'Usuarios', pendienteEn: 'T8' },
  { id: 'empresas', titulo: 'Empresas', pendienteEn: 'T8' },
  { id: 'choferes', titulo: 'Choferes', pendienteEn: 'T8' },
  { id: 'proveedores', titulo: 'Proveedores', pendienteEn: 'T8' },
  { id: 'pesaje', titulo: 'Pesaje', pendienteEn: 'T7' },
  { id: 'catalogo', titulo: 'Catálogo' },
  { id: 'vehiculos', titulo: 'Vehículos', grupo: 'Catálogo', pendienteEn: 'T8' },
  { id: 'productos', titulo: 'Productos', grupo: 'Catálogo', pendienteEn: 'T8' },
  { id: 'origenes', titulo: 'Orígenes', grupo: 'Catálogo', pendienteEn: 'T8' },
  { id: 'destinos', titulo: 'Destinos', grupo: 'Catálogo', pendienteEn: 'T8' },
  { id: 'reportes', titulo: 'Reportes', pendienteEn: 'T9' },
  { id: 'config', titulo: 'Configuración' },
  { id: 'boleta', titulo: 'Boleta', grupo: 'Configuración', pendienteEn: 'T9' },
  { id: 'balanza', titulo: 'Balanza', grupo: 'Configuración', pendienteEn: 'T8' },
];

export function principales(): ModuloNav[] {
  return MODULOS.filter((m) => !m.grupo);
}

export function deGrupo(nombre: string): ModuloNav[] {
  return MODULOS.filter((m) => m.grupo === nombre);
}

export function porId(id: ModuloId): ModuloNav | undefined {
  return MODULOS.find((m) => m.id === id);
}
