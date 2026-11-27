import logo from './assets/images/logo-pesaje.jpg';
import type { LoginRespuesta } from './login';
import { montarPesaje } from './modulos/pesaje';
import {
  CerrarSesion,
  EstadoIntegrado,
  SincronizarAhora,
  SolicitarSistemaCompleto,
} from '../wailsjs/go/main/App';
import {
  crearDialogoAviso,
  enlazarDialogo,
  htmlAccionesDialogo,
  htmlBotonDialogo,
  htmlShellDialogo,
} from './dialogo';

export type EstadoIntegradoView = {
  balanzaConectada: boolean;
  balanzaTexto: string;
  syncTexto: string;
  syncVariante: string;
  pesoVivo: number;
  pesoCapturado: number;
  textoConectar: string;
  puedeCapturar: boolean;
  puedeEnviar: boolean;
  syncDisponible: boolean;
  centralizadorOfflineFirst?: boolean;
};

type Accion = {
  ok: boolean;
  aviso?: { titulo: string; mensaje: string; esError?: boolean };
  estado?: EstadoIntegradoView;
};

export function montarModoIntegrado(
  contenedor: HTMLElement,
  sesion: NonNullable<LoginRespuesta['usuario']>,
): () => void {
  contenedor.innerHTML = `
    <div class="modo-integrado root">
      <header class="barra-superior">
        <div class="barra-fila-principal">
          <div class="barra-logo">
            <img alt="Sistema de Pesaje" width="132" />
          </div>
          <div class="barra-acciones">
            <button type="button" class="boton-secundario" id="btn-sync">Sincronizar ahora</button>
            <div class="menu-contenedor">
              <button type="button" class="boton-secundario" id="btn-menu">Menú</button>
              <div class="menu-flotante" id="menu-flotante" hidden>
                <button type="button" data-accion="sistema">Sistema completo</button>
                <hr />
                <button type="button" data-accion="balanza">Configurar balanza</button>
                <hr />
                <button type="button" data-accion="cerrar">Cerrar sesión</button>
              </div>
            </div>
          </div>
        </div>
        <div class="barra-fila-estado subtitulo">
          <span>Balanza</span>
          <strong id="estado-balanza">desconectada</strong>
          <span id="pill-offline" class="badge badge-proceso" hidden>Offline-first</span>
          <span id="pill-sync" class="badge badge-activo">Al día</span>
        </div>
      </header>
      <main class="integrado-centro integrado-centro--pesaje shell-contenido-abm" id="integrado-pesaje-host"></main>
      <div id="montaje-dialogos-integrado"></div>
    </div>
  `;

  const img = contenedor.querySelector('.barra-logo img') as HTMLImageElement;
  img.src = logo;

  const estadoBalanza = contenedor.querySelector('#estado-balanza') as HTMLElement;
  const pillOffline = contenedor.querySelector('#pill-offline') as HTMLElement;
  const pillSync = contenedor.querySelector('#pill-sync') as HTMLElement;
  const btnSync = contenedor.querySelector('#btn-sync') as HTMLButtonElement;
  const btnMenu = contenedor.querySelector('#btn-menu') as HTMLButtonElement;
  const menu = contenedor.querySelector('#menu-flotante') as HTMLDivElement;
  const slotDialogos = contenedor.querySelector('#montaje-dialogos-integrado') as HTMLElement;
  const hostPesaje = contenedor.querySelector('#integrado-pesaje-host') as HTMLElement;
  const aviso = crearDialogoAviso(slotDialogos, 'dlg-aviso', 'compacto');
  const limpiarPesaje = montarPesaje(hostPesaje);

  const wrapPwd = document.createElement('div');
  wrapPwd.innerHTML = htmlShellDialogo({
    id: 'dlg-password',
    ancho: 'compacto',
    etiqueta: 'form',
    idCuerpo: 'form-sistema',
    attrsCuerpo: 'autocomplete="off"',
    interior: `
      <h3 class="titulo-tarjeta">Acceso al sistema completo</h3>
      <p class="subtitulo dlg-intro">
        Ingrese la contraseña de acceso al escritorio completo (catálogos, pesaje y reportes).
      </p>
      <label class="pesaje-campo" for="pwd-sistema">
        <span class="etiqueta-mini">Contraseña</span>
        <input
          class="password-field"
          type="password"
          id="pwd-sistema"
          name="contrasena-sistema"
          autocomplete="current-password"
          placeholder="Contraseña"
          required
        />
      </label>
      ${htmlAccionesDialogo(
        `${htmlBotonDialogo('button', 'boton-secundario', 'pwd-cancelar', 'CANCELAR')}${htmlBotonDialogo('submit', 'boton-primario', 'pwd-continuar', 'CONTINUAR')}`,
      )}
    `,
  });
  slotDialogos.appendChild(wrapPwd.firstElementChild!);

  const dlgPwd = slotDialogos.querySelector('#dlg-password') as HTMLDialogElement;
  const formSistema = slotDialogos.querySelector('#form-sistema') as HTMLFormElement;
  const pwdInput = slotDialogos.querySelector('#pwd-sistema') as HTMLInputElement;
  const btnPwdCancelar = slotDialogos.querySelector('#pwd-cancelar') as HTMLButtonElement;
  enlazarDialogo(dlgPwd, [btnPwdCancelar]);

  const mostrarAviso = (titulo: string, mensaje: string, esError = false) => {
    aviso.mostrar(titulo, mensaje, esError);
  };

  const aplicarEstado = (e: EstadoIntegradoView) => {
    estadoBalanza.textContent = e.balanzaTexto;
    estadoBalanza.className = e.balanzaConectada ? 'mensaje-exito' : 'subtitulo';
    pillSync.textContent = e.syncTexto;
    pillSync.className = `badge badge-${e.syncVariante === 'activo' ? 'activo' : e.syncVariante === 'inactivo' ? 'inactivo' : 'proceso'}`;
    btnSync.disabled = !e.syncDisponible;
    const showOffline = e.centralizadorOfflineFirst ?? e.syncDisponible;
    pillOffline.hidden = !showOffline;
  };

  const manejarAccion = (acc: Accion) => {
    if (acc.estado) aplicarEstado(acc.estado);
    if (acc.aviso) mostrarAviso(acc.aviso.titulo, acc.aviso.mensaje, !!acc.aviso.esError);
  };

  void EstadoIntegrado().then((e) => aplicarEstado(e as EstadoIntegradoView));

  btnSync.addEventListener('click', () => {
    void SincronizarAhora().then((a) => manejarAccion(a as Accion));
  });

  btnMenu.addEventListener('click', () => {
    menu.hidden = !menu.hidden;
  });

  const cerrarMenuFlotante = (ev: MouseEvent) => {
    if (!menu.hidden && !menu.contains(ev.target as Node) && ev.target !== btnMenu) {
      menu.hidden = true;
    }
  };
  document.addEventListener('click', cerrarMenuFlotante);

  const alSalir = () => {
    document.removeEventListener('click', cerrarMenuFlotante);
    window.clearInterval(timerSync);
    limpiarPesaje();
  };

  menu.querySelectorAll('button[data-accion]').forEach((btn) => {
    btn.addEventListener('click', () => {
      menu.hidden = true;
      const accion = btn.getAttribute('data-accion');
      if (accion === 'sistema') {
        pwdInput.value = '';
        dlgPwd.showModal();
        requestAnimationFrame(() => pwdInput.focus());
        return;
      }
      if (accion === 'balanza') {
        mostrarAviso(
          'Configuración balanza',
          'Use Menú → Sistema completo → Configuración → Balanza para RS232.',
        );
        return;
      }
      if (accion === 'cerrar') {
        if (window.confirm('¿Seguro que desea cerrar sesión?')) {
          alSalir();
          CerrarSesion();
          import('./login').then(({ montarLogin }) => montarLogin(contenedor));
        }
      }
    });
  });

  contenedor.querySelector('#pwd-cancelar')?.addEventListener('click', () => dlgPwd.close());
  formSistema.addEventListener('submit', (ev) => {
    ev.preventDefault();
    const pwd = pwdInput.value;
    void SolicitarSistemaCompleto(pwd).then((a) => {
      dlgPwd.close();
      const acc = a as Accion;
      if (acc.ok) {
        alSalir();
        import('./shell').then(({ montarShell }) => montarShell(contenedor, sesion));
        return;
      }
      manejarAccion(acc);
    });
  });

  const timerSync = window.setInterval(() => {
    void EstadoIntegrado().then((e) => aplicarEstado(e as EstadoIntegradoView));
  }, 2000);

  void sesion;
  return alSalir;
}
