import logo from './assets/images/logo-pesaje.jpg';
import type { LoginRespuesta } from './login';
import {
  AlternarBalanza,
  CapturarPeso,
  CerrarSesion,
  EnviarPeso,
  EstadoIntegrado,
  SincronizarAhora,
  SolicitarSistemaCompleto,
} from '../wailsjs/go/main/App';

export type EstadoIntegrado = {
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
};

type Accion = {
  ok: boolean;
  aviso?: { titulo: string; mensaje: string; esError?: boolean };
  estado?: EstadoIntegrado;
};

export function montarModoIntegrado(contenedor: HTMLElement, sesion: NonNullable<LoginRespuesta['usuario']>): void {
  contenedor.innerHTML = `
    <div class="modo-integrado root">
      <header class="barra-superior">
        <div class="barra-fila-principal">
          <div class="barra-logo">
            <img alt="Sistema de Pesaje" width="160" />
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
          <span id="pill-sync" class="badge badge-activo">Al día</span>
        </div>
      </header>
      <main class="integrado-centro">
        <section class="tarjeta tarjeta-sin-padding integrado-tarjeta">
          <div class="integrado-encabezado">
            <h2 class="titulo-tarjeta">Registro de peso</h2>
          </div>
          <div class="separador-horizontal"></div>
          <div class="integrado-cuerpo">
            <div class="display-peso">
              <span class="valor-peso" id="valor-peso">0 kg</span>
            </div>
            <button type="button" class="boton-secundario integrado-ancho" id="btn-balanza">Conectar balanza</button>
            <button type="button" class="boton-secundario integrado-ancho" id="btn-capturar" disabled>Capturar peso</button>
            <p class="subtitulo integrado-captura" id="peso-capturado">Peso capturado: 0 kg</p>
            <button type="button" class="boton-primario integrado-ancho" id="btn-enviar" disabled>Enviar a la web</button>
          </div>
        </section>
      </main>
      <dialog id="dlg-aviso" class="dlg-aviso">
        <form method="dialog" class="tarjeta dlg-contenido">
          <h3 id="dlg-titulo" class="titulo-tarjeta"></h3>
          <p id="dlg-mensaje" class="subtitulo"></p>
          <button type="submit" class="boton-primario">Aceptar</button>
        </form>
      </dialog>
      <dialog id="dlg-password" class="dlg-aviso">
        <form method="dialog" class="tarjeta dlg-contenido" id="form-sistema">
          <h3 class="titulo-tarjeta">Acceso al sistema completo</h3>
          <input class="password-field" type="password" id="pwd-sistema" placeholder="Contraseña" required />
          <div class="dlg-acciones">
            <button type="button" class="boton-secundario" id="pwd-cancelar">Cancelar</button>
            <button type="submit" class="boton-primario">Continuar</button>
          </div>
        </form>
      </dialog>
    </div>
  `;

  const img = contenedor.querySelector('.barra-logo img') as HTMLImageElement;
  img.src = logo;

  const estadoBalanza = contenedor.querySelector('#estado-balanza') as HTMLElement;
  const pillSync = contenedor.querySelector('#pill-sync') as HTMLElement;
  const valorPeso = contenedor.querySelector('#valor-peso') as HTMLElement;
  const pesoCapturado = contenedor.querySelector('#peso-capturado') as HTMLElement;
  const btnBalanza = contenedor.querySelector('#btn-balanza') as HTMLButtonElement;
  const btnCapturar = contenedor.querySelector('#btn-capturar') as HTMLButtonElement;
  const btnEnviar = contenedor.querySelector('#btn-enviar') as HTMLButtonElement;
  const btnSync = contenedor.querySelector('#btn-sync') as HTMLButtonElement;
  const btnMenu = contenedor.querySelector('#btn-menu') as HTMLButtonElement;
  const menu = contenedor.querySelector('#menu-flotante') as HTMLDivElement;
  const dlgAviso = contenedor.querySelector('#dlg-aviso') as HTMLDialogElement;
  const dlgPwd = contenedor.querySelector('#dlg-password') as HTMLDialogElement;
  const formSistema = contenedor.querySelector('#form-sistema') as HTMLFormElement;

  const mostrarAviso = (titulo: string, mensaje: string, esError = false) => {
    const t = contenedor.querySelector('#dlg-titulo') as HTMLElement;
    const m = contenedor.querySelector('#dlg-mensaje') as HTMLElement;
    t.textContent = titulo;
    m.textContent = mensaje;
    m.className = esError ? 'mensaje-error' : 'subtitulo';
    dlgAviso.showModal();
  };

  const aplicarEstado = (e: EstadoIntegrado) => {
    estadoBalanza.textContent = e.balanzaTexto;
    estadoBalanza.className = e.balanzaConectada ? 'mensaje-exito' : 'subtitulo';
    pillSync.textContent = e.syncTexto;
    pillSync.className = `badge badge-${e.syncVariante === 'activo' ? 'activo' : e.syncVariante === 'inactivo' ? 'inactivo' : 'proceso'}`;
    valorPeso.textContent = formatearPeso(e.pesoVivo);
    pesoCapturado.textContent = `Peso capturado: ${formatearPeso(e.pesoCapturado)}`;
    btnBalanza.textContent = e.textoConectar;
    btnCapturar.disabled = !e.puedeCapturar;
    btnEnviar.disabled = !e.puedeEnviar;
  };

  const manejarAccion = (acc: Accion) => {
    if (acc.estado) aplicarEstado(acc.estado);
    if (acc.aviso) mostrarAviso(acc.aviso.titulo, acc.aviso.mensaje, !!acc.aviso.esError);
  };

  void EstadoIntegrado().then((e) => aplicarEstado(e as EstadoIntegrado));

  btnBalanza.addEventListener('click', () => {
    void AlternarBalanza().then((a) => manejarAccion(a as Accion));
  });
  btnCapturar.addEventListener('click', () => {
    void CapturarPeso().then((a) => manejarAccion(a as Accion));
  });
  btnEnviar.addEventListener('click', () => {
    void EnviarPeso().then((a) => manejarAccion(a as Accion));
  });
  btnSync.addEventListener('click', () => {
    void SincronizarAhora().then((a) => manejarAccion(a as Accion));
  });

  btnMenu.addEventListener('click', () => {
    menu.hidden = !menu.hidden;
  });
  document.addEventListener('click', (ev) => {
    if (!menu.hidden && !menu.contains(ev.target as Node) && ev.target !== btnMenu) {
      menu.hidden = true;
    }
  });

  menu.querySelectorAll('button[data-accion]').forEach((btn) => {
    btn.addEventListener('click', () => {
      menu.hidden = true;
      const accion = btn.getAttribute('data-accion');
      if (accion === 'sistema') {
        dlgPwd.showModal();
        (contenedor.querySelector('#pwd-sistema') as HTMLInputElement).value = '';
        return;
      }
      if (accion === 'balanza') {
        mostrarAviso('Configuración balanza RS232', 'La pantalla de configuración llega en T6/T8.');
        return;
      }
      if (accion === 'cerrar') {
        if (window.confirm('¿Seguro que desea cerrar sesión?')) {
          CerrarSesion();
          import('./login').then(({ montarLogin }) => montarLogin(contenedor));
        }
      }
    });
  });

  contenedor.querySelector('#pwd-cancelar')?.addEventListener('click', () => dlgPwd.close());
  formSistema.addEventListener('submit', (ev) => {
    ev.preventDefault();
    const pwd = (contenedor.querySelector('#pwd-sistema') as HTMLInputElement).value;
    void SolicitarSistemaCompleto(pwd).then((a) => {
      dlgPwd.close();
      const acc = a as Accion;
      if (acc.ok) {
        import('./shell').then(({ montarShell }) => montarShell(contenedor, sesion));
        return;
      }
      manejarAccion(acc);
    });
  });

  window.setInterval(() => {
    void EstadoIntegrado().then((e) => aplicarEstado(e as EstadoIntegrado));
  }, 2000);

  void sesion;
}

function formatearPeso(peso: number): string {
  return `${peso.toLocaleString('es-BO')} kg`;
}
