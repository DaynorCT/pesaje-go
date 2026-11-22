import logo from './assets/images/logo-pesaje.jpg';
import type { LoginRespuesta } from './login';
import { montarPesaje } from './modulos/pesaje';
import { montarCatalogo } from './modulos/catalogos';
import { montarBalanza } from './modulos/balanza';
import { montarInicio } from './modulos/inicio';
import { montarReportes } from './modulos/reportes';
import { montarBoletaConfig } from './modulos/boleta-config';
import {
  ICONO_MENU_INTEGRADO,
  ICONO_MENU_SALIR,
  ICONO_MENU_WEB,
  iconoUsuarioCabeceraHTML,
  insertarIconoModulo,
} from './iconos-modulo';
import { deGrupo, porId, principales, type ModuloId } from './navegacion';

const IDS_ABM: ModuloId[] = [
  'usuarios', 'empresas', 'choferes', 'proveedores',
  'vehiculos', 'productos', 'origenes', 'destinos',
];
import {
  AbrirSistemaWeb,
  CerrarSesion,
  SesionActual,
  VolverAModoIntegrado,
} from '../wailsjs/go/main/App';

export function montarShell(
  contenedor: HTMLElement,
  sesionInicial?: NonNullable<LoginRespuesta['usuario']>,
): () => void {
  let modulo: ModuloId = 'inicio';
  let sesion = sesionInicial;
  let limpiarModulo: (() => void) | undefined;

  const cerrarMenusDocumento = (ev: MouseEvent) => {
    const nav = contenedor.querySelector('#shell-nav');
    if (!nav) return;
    const t = ev.target as HTMLElement | null;
    if (t?.closest('.shell-submenu, .shell-menu-ancla')) return;
    cerrarSubmenus(nav);
  };
  document.addEventListener('click', cerrarMenusDocumento);

  const render = () => {
    limpiarModulo?.();
    limpiarModulo = undefined;
    contenedor.replaceChildren();
    const mod = porId(modulo)!;
    contenedor.innerHTML = `
      <div class="shell root">
        <header class="barra-superior shell-cabecera">
          <div class="shell-fila-principal">
            <img class="shell-logo" alt="Sistema de Pesaje" width="132" />
            <nav class="barra-navegacion shell-nav" id="shell-nav"></nav>
            <div class="shell-usuario">
              <button
                type="button"
                class="boton-usuario-cabecera"
                id="btn-usuario"
                aria-label="Menú de usuario"
                title="Usuario"
              >
                <span class="shell-item-icon shell-usuario-icon" aria-hidden="true"></span>
              </button>
              <div class="menu-flotante shell-menu-usuario" id="menu-usuario" hidden></div>
            </div>
          </div>
          ${
            IDS_ABM.includes(mod.id) || mod.id === 'pesaje' || mod.id === 'inicio'
              ? ''
              : `<h1 class="titulo-modulo shell-titulo">${mod.titulo}</h1>`
          }
        </header>
        <main class="shell-contenido" id="shell-contenido"></main>
      </div>
    `;

    const img = contenedor.querySelector('.shell-logo') as HTMLImageElement;
    img.src = logo;

    pintarNav(contenedor, modulo, (id) => {
      modulo = id;
      render();
    });

    limpiarModulo = pintarContenido(contenedor.querySelector('#shell-contenido') as HTMLElement, mod, (id) => {
      modulo = id;
      render();
    });
    enlazarUsuario(contenedor, sesion, destruir);
  };

  const destruir = () => {
    document.removeEventListener('click', cerrarMenusDocumento);
    limpiarModulo?.();
    cerrarSubmenus(document);
  };

  if (!sesion) {
    void SesionActual().then((s) => {
      sesion = s ?? undefined;
      render();
    });
  } else {
    render();
  }

  return destruir;
}

function pintarNav(
  root: HTMLElement,
  activo: ModuloId,
  alElegir: (id: ModuloId) => void,
): void {
  const nav = root.querySelector('#shell-nav') as HTMLElement;
  nav.innerHTML = '';
  for (const m of principales()) {
    const hijos = m.id === 'catalogo' ? deGrupo('Catálogo') : m.id === 'config' ? deGrupo('Configuración') : [];
    const estaActivo =
      activo === m.id ||
      (m.id === 'catalogo' && porId(activo)?.grupo === 'Catálogo') ||
      (m.id === 'config' && porId(activo)?.grupo === 'Configuración');

    const ancla = document.createElement('div');
    ancla.className = 'shell-menu-ancla';
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'item-menu shell-item-menu';
    if (estaActivo) btn.classList.add('activo');
    insertarIconoModulo(btn, m.id);
    const etiqueta = document.createElement('span');
    etiqueta.className = 'shell-item-label';
    etiqueta.textContent = m.titulo;
    btn.appendChild(etiqueta);
    if (hijos.length > 0) {
      const chevron = document.createElement('span');
      chevron.className = 'shell-item-chevron';
      chevron.setAttribute('aria-hidden', 'true');
      chevron.textContent = '▾';
      btn.appendChild(chevron);
    }
    ancla.appendChild(btn);

    if (hijos.length > 0) {
      btn.addEventListener('click', (ev) => {
        ev.stopPropagation();
        const tituloSub =
          m.id === 'config' ? 'Configuración' : m.id === 'catalogo' ? 'Catálogo' : m.titulo;
        alternarSubmenu(ancla, tituloSub, hijos, alElegir);
      });
    } else {
      btn.addEventListener('click', (ev) => {
        ev.stopPropagation();
        cerrarSubmenus(nav);
        if (m.id !== 'catalogo' && m.id !== 'config') alElegir(m.id);
      });
    }
    nav.appendChild(ancla);
  }
}

function alternarSubmenu(
  ancla: HTMLElement,
  titulo: string,
  hijos: ReturnType<typeof deGrupo>,
  alElegir: (id: ModuloId) => void,
): void {
  const abierto = ancla.querySelector('.shell-submenu');
  document.querySelectorAll('.shell-submenu').forEach((el) => el.remove());
  if (abierto) {
    return;
  }
  const pop = document.createElement('div');
  pop.className = 'menu-flotante shell-submenu';
  const enc = document.createElement('p');
  enc.className = 'subtitulo menu-encabezado';
  enc.textContent = titulo;
  pop.appendChild(enc);
  const sep = document.createElement('hr');
  pop.appendChild(sep);
  for (const h of hijos) {
    const b = document.createElement('button');
    b.type = 'button';
    b.className = 'item-menu shell-submenu-item';
    insertarIconoModulo(b, h.id);
    const lab = document.createElement('span');
    lab.textContent = h.titulo;
    b.appendChild(lab);
    b.addEventListener('click', (ev) => {
      ev.stopPropagation();
      cerrarSubmenus(document);
      alElegir(h.id);
    });
    pop.appendChild(b);
  }
  ancla.appendChild(pop);
}

function cerrarSubmenus(root: ParentNode): void {
  root.querySelectorAll('.shell-submenu').forEach((el) => el.remove());
}

function pintarContenido(
  host: HTMLElement,
  mod: ReturnType<typeof porId>,
  irAModulo?: (id: ModuloId) => void,
): (() => void) | undefined {
  if (!mod) return undefined;
  if (mod.id === 'inicio') {
    host.className = 'shell-contenido shell-contenido-abm';
    host.replaceChildren();
    return montarInicio(host, { irAModulo });
  }
  if (mod.id === 'reportes') {
    host.className = 'shell-contenido shell-contenido-abm';
    host.replaceChildren();
    return montarReportes(host);
  }
  if (mod.id === 'boleta') {
    host.className = 'shell-contenido';
    host.replaceChildren();
    return montarBoletaConfig(host);
  }
  if (mod.id === 'pesaje') {
    host.className = 'shell-contenido shell-contenido-abm';
    host.replaceChildren();
    return montarPesaje(host);
  }
  if (mod.id === 'balanza') {
    host.className = 'shell-contenido';
    host.replaceChildren();
    return montarBalanza(host);
  }
  if (IDS_ABM.includes(mod.id)) {
    host.className = 'shell-contenido shell-contenido-abm';
    host.replaceChildren();
  return montarCatalogo(host, mod.id) ?? undefined;
  }
  host.className = 'shell-contenido';
  if (mod.id === 'catalogo' || mod.id === 'config') {
    host.innerHTML = `
      <section class="tarjeta shell-panel">
        <h2 class="titulo-tarjeta">${mod.titulo}</h2>
        <p class="subtitulo">Elija una opción del submenú en la barra superior.</p>
      </section>
    `;
    return undefined;
  }
  host.innerHTML = `
    <section class="tarjeta shell-panel">
      <h2 class="titulo-tarjeta">${mod.titulo}</h2>
      <p class="subtitulo">Módulo no disponible.</p>
    </section>
  `;
  return undefined;
}

function enlazarUsuario(
  root: HTMLElement,
  sesion: NonNullable<LoginRespuesta['usuario']> | undefined,
  alDestruirShell?: () => void,
): void {
  const btn = root.querySelector('#btn-usuario') as HTMLButtonElement;
  const menu = root.querySelector('#menu-usuario') as HTMLDivElement;
  const nombre = sesion?.nombre || sesion?.email || 'Usuario';
  btn.title = nombre;
  btn.setAttribute('aria-label', `Menú de usuario: ${nombre}`);
  const icono = btn.querySelector('.shell-usuario-icon');
  if (icono) icono.innerHTML = iconoUsuarioCabeceraHTML();

  menu.innerHTML = `
    <div class="shell-menu-usuario-cabecera">
      <p class="shell-menu-usuario-nombre">${escapeHtml(nombre)}</p>
      <p class="shell-menu-usuario-correo">${escapeHtml(sesion?.email ?? '')}</p>
    </div>
    <hr />
    <button type="button" class="shell-menu-usuario-item" data-u="web">
      <span class="shell-item-icon shell-menu-item-icono">${ICONO_MENU_WEB}</span>
      <span>Ir al sistema web</span>
    </button>
    <button type="button" class="shell-menu-usuario-item" data-u="integrado">
      <span class="shell-item-icon shell-menu-item-icono">${ICONO_MENU_INTEGRADO}</span>
      <span>Volver al modo báscula</span>
    </button>
    <hr />
    <button type="button" class="shell-menu-usuario-item shell-menu-usuario-item--salir" data-u="salir">
      <span class="shell-item-icon shell-menu-item-icono">${ICONO_MENU_SALIR}</span>
      <span>Cerrar sesión</span>
    </button>
  `;

  btn.onclick = (ev) => {
    ev.stopPropagation();
    menu.hidden = !menu.hidden;
  };

  menu.querySelectorAll('button[data-u]').forEach((b) => {
    b.addEventListener('click', () => {
      menu.hidden = true;
      const u = b.getAttribute('data-u');
      if (u === 'web') void AbrirSistemaWeb();
      const host = root;
      if (u === 'integrado' && sesion) {
        void VolverAModoIntegrado().then(() =>
          import('./modo-integrado').then(({ montarModoIntegrado }) =>
            montarModoIntegrado(host, sesion),
          ),
        );
      }
      if (u === 'salir' && window.confirm('¿Seguro que desea cerrar sesión?')) {
        alDestruirShell?.();
        CerrarSesion();
        import('./login').then(({ montarLogin }) => montarLogin(host));
      }
    });
  });
}

function escapeHtml(texto: string): string {
  return texto
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}
