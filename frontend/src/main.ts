import './css/tema.css';
import './css/jacobitus.css';
import './css/pesaje.css';
import './style.css';
import './app.css';

import logo from './assets/images/logo-pesaje.jpg';

const app = document.querySelector('#app')!;
app.innerHTML = `
  <main class="scaffold root">
    <img class="scaffold-logo" alt="Sistema de Pesaje" />
    <h1 class="titulo-modulo">Sistema de Pesaje</h1>
    <p class="subtitulo">Go + Wails — paridad con el módulo Java :pesaje</p>
    <div class="tarjeta">
      <p>T1: scaffold, tema y hojas CSS web copiadas desde <code>pesaje.css</code> / <code>jacobitus.css</code>.</p>
      <span class="version-banner">Referencia Java en ../jacobitus-escritorio</span>
    </div>
  </main>
`;

const img = app.querySelector('.scaffold-logo') as HTMLImageElement;
img.src = logo;

