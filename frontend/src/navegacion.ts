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
  { id: 'inicio', titulo: 'Inicio' },
  { id: 'usuarios', titulo: 'Usuarios' },
  { id: 'empresas', titulo: 'Empresas' },
  { id: 'choferes', titulo: 'Choferes' },
  { id: 'proveedores', titulo: 'Proveedores' },
  { id: 'pesaje', titulo: 'Pesaje' },
  { id: 'catalogo', titulo: 'Catálogo' },
  { id: 'vehiculos', titulo: 'Vehículos', grupo: 'Catálogo' },
  { id: 'productos', titulo: 'Productos', grupo: 'Catálogo' },
  { id: 'origenes', titulo: 'Orígenes', grupo: 'Catálogo' },
  { id: 'destinos', titulo: 'Destinos', grupo: 'Catálogo' },
  { id: 'reportes', titulo: 'Reportes' },
  { id: 'config', titulo: 'Configuración' },
  { id: 'boleta', titulo: 'Boleta', grupo: 'Configuración' },
  { id: 'balanza', titulo: 'Balanza', grupo: 'Configuración' },
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
