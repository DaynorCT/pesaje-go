import logo from './assets/images/logo-pesaje.jpg';
import { Autenticar, AvisoLoginInicial } from '../wailsjs/go/main/App';

export type LoginRespuesta = {
  exito: boolean;
  estado: string;
  mensaje?: string;
  usuario?: {
    id: number;
    email: string;
    nombre: string;
    rol: string;
  };
};

export function montarLogin(contenedor: HTMLElement): void {
  contenedor.innerHTML = `
    <div class="login-shell root">
      <form class="tarjeta login" id="form-login" autocomplete="on">
        <div class="login-logo">
          <img alt="Sistema de Pesaje" width="320" />
        </div>
        <h1 class="titulo-modulo">Ingresar al sistema</h1>
        <p class="subtitulo">Ingrese con su usuario y contraseña para continuar.</p>
        <p id="login-aviso" class="subtitulo login-aviso" hidden></p>
        <input id="login-email" class="text-field" type="email" name="email"
               autocomplete="username" placeholder="Correo electrónico" required />
        <input id="login-clave" class="password-field" type="password" name="password"
               autocomplete="current-password" placeholder="Contraseña" required />
        <p id="login-mensaje" class="login-mensaje" role="alert" aria-live="polite"></p>
        <button type="submit" class="boton-primario login-submit">Ingresar</button>
      </form>
    </div>
  `;

  const img = contenedor.querySelector('.login-logo img') as HTMLImageElement;
  img.src = logo;

  const aviso = contenedor.querySelector('#login-aviso') as HTMLParagraphElement;
  const mensaje = contenedor.querySelector('#login-mensaje') as HTMLParagraphElement;
  const email = contenedor.querySelector('#login-email') as HTMLInputElement;
  const clave = contenedor.querySelector('#login-clave') as HTMLInputElement;
  const form = contenedor.querySelector('#form-login') as HTMLFormElement;
  const boton = contenedor.querySelector('.login-submit') as HTMLButtonElement;

  void AvisoLoginInicial().then((texto) => {
    if (texto?.trim()) {
      aviso.textContent = texto;
      aviso.hidden = false;
    }
  });

  const limpiarMensaje = () => {
    mensaje.textContent = '';
    mensaje.className = 'login-mensaje';
  };

  form.addEventListener('submit', (ev) => {
    ev.preventDefault();
    limpiarMensaje();
    boton.disabled = true;
    Autenticar(email.value.trim(), clave.value)
      .then((res: LoginRespuesta) => {
        if (res.exito && res.usuario) {
          montarPostLogin(contenedor, res);
          return;
        }
        mensaje.textContent = res.mensaje || 'No se pudo iniciar sesión.';
        mensaje.className = 'login-mensaje mensaje-error';
      })
      .catch(() => {
        mensaje.textContent = 'No se pudo iniciar sesión.';
        mensaje.className = 'login-mensaje mensaje-error';
      })
      .finally(() => {
        boton.disabled = false;
      });
  });

  email.focus();
}

function montarPostLogin(contenedor: HTMLElement, res: LoginRespuesta): void {
  const nombre = res.usuario?.nombre || res.usuario?.email || '';
  contenedor.innerHTML = `
    <div class="login-shell root">
      <div class="tarjeta login login-ok">
        <h1 class="titulo-modulo">Sesión iniciada</h1>
        <p class="subtitulo">Bienvenido, <strong>${escapeHtml(nombre)}</strong>.</p>
        <p class="subtitulo">T4 (modo integrado) sigue pendiente; la ventana ya quedó en 480×680.</p>
      </div>
    </div>
  `;
}

function escapeHtml(texto: string): string {
  return texto
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
