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
import {
  crearDialogoAviso,
  enlazarDialogo,
  htmlAccionesDialogo,
  htmlBotonDialogo,
  htmlShellDialogo,
} from './dialogo';

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
      <div id="montaje-dialogos-integrado"></div>
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
  const slotDialogos = contenedor.querySelector('#montaje-dialogos-integrado') as HTMLElement;
  const aviso = crearDialogoAviso(slotDialogos, 'dlg-aviso', 'compacto');

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
