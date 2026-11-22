import logo from './assets/images/logo-pesaje.jpg';
import type { LoginRespuesta } from './login';
import { montarPesaje } from './modulos/pesaje';
import { montarCatalogo } from './modulos/catalogos';
import { montarBalanza } from './modulos/balanza';
import { montarInicio } from './modulos/inicio';
import { montarReportes } from './modulos/reportes';
import { montarBoletaConfig } from './modulos/boleta-config';
import { deGrupo, porId, principales, type ModuloId } from './navegacion';
import {
  AbrirSistemaWeb,
  CerrarSesion,
  SesionActual,
  VolverAModoIntegrado,
} from '../wailsjs/go/main/App';

export function montarShell(
  contenedor: HTMLElement,
  sesionInicial?: NonNullable<LoginRespuesta['usuario']>,
): void {
  let modulo: ModuloId = 'inicio';
  let sesion = sesionInicial;
  let limpiarModulo: (() => void) | undefined;

  const render = () => {
    limpiarModulo?.();
    limpiarModulo = undefined;
    contenedor.replaceChildren();
    const mod = porId(modulo)!;
    contenedor.innerHTML = `
      <div class="shell root">
        <header class="barra-superior shell-cabecera">
          <div class="shell-fila-principal">
            <img class="shell-logo" alt="Sistema de Pesaje" width="160" />
            <nav class="barra-navegacion shell-nav" id="shell-nav"></nav>
            <div class="shell-usuario">
              <button type="button" class="boton-secundario boton-usuario-cabecera" id="btn-usuario">
                Usuario
              </button>
              <div class="menu-flotante shell-menu-usuario" id="menu-usuario" hidden></div>
            </div>
          </div>
          <h1 class="titulo-modulo shell-titulo">${mod.titulo}</h1>
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

    limpiarModulo = pintarContenido(contenedor.querySelector('#shell-contenido') as HTMLElement, mod);
    enlazarUsuario(contenedor, sesion);
  };

  if (!sesion) {
    void SesionActual().then((s) => {
      sesion = s ?? undefined;
      render();
    });
  } else {
    render();
  }
}

function pintarNav(
  root: HTMLElement,
  activo: ModuloId,
  alElegir: (id: ModuloId) => void,
): void {
  const nav = root.querySelector('#shell-nav') as HTMLElement;
  nav.innerHTML = '';
  for (const m of principales()) {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'item-menu shell-item-menu';
    btn.textContent = m.titulo;
    const hijos = m.id === 'catalogo' ? deGrupo('Catálogo') : m.id === 'config' ? deGrupo('Configuración') : [];
    const estaActivo =
      activo === m.id || (!!porId(activo)?.grupo && porId(activo)?.grupo === m.titulo);
    if (estaActivo) btn.classList.add('activo');

    if (hijos.length > 0) {
      btn.addEventListener('click', (ev) => {
        ev.stopPropagation();
        alternarSubmenu(nav, btn, m.titulo, hijos, alElegir);
      });
    } else {
      btn.addEventListener('click', () => {
        cerrarSubmenus(nav);
        if (m.id !== 'catalogo' && m.id !== 'config') alElegir(m.id);
      });
    }
    nav.appendChild(btn);
  }
  document.addEventListener(
    'click',
    () => cerrarSubmenus(nav),
    { once: true },
  );
}

function alternarSubmenu(
  nav: HTMLElement,
  ancla: HTMLElement,
  titulo: string,
  hijos: ReturnType<typeof deGrupo>,
  alElegir: (id: ModuloId) => void,
): void {
  cerrarSubmenus(nav);
  const pop = document.createElement('div');
  pop.className = 'menu-flotante shell-submenu';
  pop.innerHTML = `<p class="grupo-menu">${titulo}</p>`;
  for (const h of hijos) {
    const b = document.createElement('button');
    b.type = 'button';
    b.textContent = h.titulo;
    b.addEventListener('click', (ev) => {
      ev.stopPropagation();
      pop.remove();
      alElegir(h.id);
    });
    pop.appendChild(b);
  }
  ancla.style.position = 'relative';
  ancla.appendChild(pop);
}

function cerrarSubmenus(nav: HTMLElement): void {
  nav.querySelectorAll('.shell-submenu').forEach((el) => el.remove());
}

function pintarContenido(host: HTMLElement, mod: ReturnType<typeof porId>): (() => void) | undefined {
  if (!mod) return undefined;
  if (mod.id === 'inicio') {
    host.className = 'shell-contenido';
    host.replaceChildren();
    return montarInicio(host);
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
    host.className = 'shell-contenido shell-contenido-pesaje';
    host.replaceChildren();
    return montarPesaje(host);
  }
  if (mod.id === 'balanza') {
    host.className = 'shell-contenido';
    host.replaceChildren();
    return montarBalanza(host);
  }
  const abmIds: ModuloId[] = [
    'usuarios', 'empresas', 'choferes', 'proveedores',
    'vehiculos', 'productos', 'origenes', 'destinos',
  ];
  if (abmIds.includes(mod.id)) {
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
): void {
  const btn = root.querySelector('#btn-usuario') as HTMLButtonElement;
  const menu = root.querySelector('#menu-usuario') as HTMLDivElement;
  const nombre = sesion?.nombre || sesion?.email || 'Usuario';
  btn.textContent = nombre;

  menu.innerHTML = `
    <p class="subtitulo">${escapeHtml(nombre)}</p>
    <p class="subtitulo">${escapeHtml(sesion?.email ?? '')}</p>
    <hr />
    <button type="button" data-u="web">Ir al sistema web</button>
    <button type="button" data-u="integrado">Volver al modo báscula</button>
    <hr />
    <button type="button" data-u="salir">Cerrar sesión</button>
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
