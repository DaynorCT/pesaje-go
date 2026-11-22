import { AbmGuardar, AbmListar } from '../../wailsjs/go/main/App';

export type CampoAbm = {
  clave: string;
  etiqueta: string;
  tipo?: 'text' | 'password' | 'estado' | 'select';
  requerido?: boolean;
  opciones?: string[];
  soloEdicion?: boolean;
  soloNuevo?: boolean;
  /** Ocupa todo el ancho del diálogo en la rejilla. */
  anchoCompleto?: boolean;
};

export type ColumnaAbm = {
  clave: string;
  titulo: string;
};

export type ConfigAbm = {
  entidad: string;
  /** Título del módulo en pantalla, p. ej. «Usuarios». */
  tituloModulo: string;
  /** Título singular del formulario, p. ej. «usuario». */
  tituloEntidad: string;
  columnas: ColumnaAbm[];
  campos: CampoAbm[];
};

const ICONO_LAPIZ = `<svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/></svg>`;

type Fila = { id: number; estado: string; celdas: Record<string, string> };

function tituloCapitalizado(tituloEntidad: string): string {
  return tituloEntidad.charAt(0).toUpperCase() + tituloEntidad.slice(1);
}

function columnasDialogo(camposVisibles: number): 2 | 3 {
  return camposVisibles > 7 ? 3 : 2;
}

function etiquetaEstadoSwitch(activo: boolean): string {
  return activo ? 'Activo' : 'Inactivo';
}

function sincronizarSwitchUI(input: HTMLInputElement, activo: boolean): void {
  const label = input.closest('.abm-switch') as HTMLElement | null;
  const txt = etiquetaEstadoSwitch(activo);
  input.checked = activo;
  input.setAttribute('aria-label', txt);
  input.setAttribute('aria-checked', activo ? 'true' : 'false');
  if (label) label.dataset.estado = txt;
}

export function montarAbm(host: HTMLElement, cfg: ConfigAbm): () => void {
  const entLabel = tituloCapitalizado(cfg.tituloEntidad);

  host.innerHTML = `
    <div class="abm-panel">
      <header class="abm-cabecera">
        <h2 class="titulo-tarjeta abm-titulo-modulo">${escapeHtml(cfg.tituloModulo)}</h2>
        <div class="abm-cabecera-der">
          <p id="abm-mensaje" class="subtitulo abm-mensaje"></p>
          <button type="button" class="boton-primario" id="abm-crear">Crear ${entLabel}</button>
        </div>
      </header>
      <div class="abm-tabla-wrap tarjeta tarjeta-sin-padding">
        <table class="tabla-modulo" id="abm-tabla">
          <thead><tr id="abm-head"></tr></thead>
          <tbody id="abm-body"></tbody>
        </table>
      </div>
      <dialog id="abm-dialog" class="dlg-aviso abm-dialog">
        <form class="tarjeta dlg-contenido abm-dialog-form" id="abm-form">
          <h3 class="titulo-tarjeta" id="abm-dialog-titulo"></h3>
          <div class="abm-dialog-grid" id="abm-campos"></div>
          <p id="abm-dialog-msg" class="subtitulo abm-dialog-msg"></p>
          <div class="dlg-acciones">
            <button type="button" class="boton-secundario" id="abm-cancelar">Cancelar</button>
            <button type="submit" class="boton-primario" id="abm-guardar">Guardar</button>
          </div>
        </form>
      </dialog>
    </div>
  `;

  const head = host.querySelector('#abm-head') as HTMLTableRowElement;
  head.innerHTML =
    cfg.columnas.map((c) => `<th>${c.titulo}</th>`).join('') + '<th class="abm-th-acciones">Acciones</th>';

  const camposHost = host.querySelector('#abm-campos') as HTMLElement;
  camposHost.innerHTML = cfg.campos
    .map((f) => {
      const ancho = f.anchoCompleto || f.clave === 'descripcion' ? 'data-ancho-completo="1"' : '';
      if (f.tipo === 'estado' || f.tipo === 'select') {
        const opts = (f.opciones ?? ['ACTIVO', 'INACTIVO'])
          .map((o) => `<option value="${o}">${o}</option>`)
          .join('');
        return `<label class="pesaje-campo" data-campo="${f.clave}" ${ancho}><span class="etiqueta-mini">${f.etiqueta}</span>
          <select class="combo-web" id="campo-${f.clave}">${opts}</select></label>`;
      }
      const type = f.tipo === 'password' ? 'password' : 'text';
      return `<label class="pesaje-campo" data-campo="${f.clave}" ${ancho}><span class="etiqueta-mini">${f.etiqueta}</span>
        <input class="text-field" type="${type}" id="campo-${f.clave}" /></label>`;
    })
    .join('');

  const dialog = host.querySelector('#abm-dialog') as HTMLDialogElement;
  const tituloDialog = host.querySelector('#abm-dialog-titulo') as HTMLElement;
  const mensaje = host.querySelector('#abm-mensaje') as HTMLElement;
  const mensajeDialog = host.querySelector('#abm-dialog-msg') as HTMLElement;
  const form = host.querySelector('#abm-form') as HTMLFormElement;

  let seleccionId = 0;
  let modoNuevo = true;
  let datosExtra: Record<string, string> = {};

  const mostrarMsg = (texto: string, error = false) => {
    mensaje.textContent = texto;
    mensaje.className = error ? 'mensaje-error abm-mensaje' : 'subtitulo abm-mensaje mensaje-exito';
  };

  const mostrarMsgDialog = (texto: string, error = false) => {
    mensajeDialog.textContent = texto;
    mensajeDialog.className = error ? 'mensaje-error abm-dialog-msg' : 'subtitulo abm-dialog-msg mensaje-exito';
  };

  const actualizarVisibilidadCampos = () => {
    let visibles = 0;
    for (const f of cfg.campos) {
      const wrap = camposHost.querySelector(`[data-campo="${f.clave}"]`) as HTMLElement | null;
      if (!wrap) continue;
      const oculto =
        (modoNuevo && f.soloEdicion) || (!modoNuevo && f.soloNuevo);
      wrap.classList.toggle('oculto', oculto);
      if (!oculto) visibles += 1;
    }
    const cols = columnasDialogo(visibles);
    camposHost.className = `abm-dialog-grid abm-dialog-grid--${cols}`;
    dialog.classList.toggle('abm-dialog--3', cols === 3);
    dialog.classList.toggle('abm-dialog--2', cols === 2);
    return visibles;
  };

  const actualizarTituloDialog = () => {
    tituloDialog.textContent = modoNuevo ? `Crear ${entLabel}` : `Editar ${entLabel}`;
  };

  const marcarFila = (id: number) => {
    const body = host.querySelector('#abm-body') as HTMLElement;
    body.querySelectorAll('tr').forEach((tr) => {
      tr.classList.toggle('seleccionado', tr.dataset.id === String(id) && id > 0);
    });
  };

  const camposDesdeFila = (fila: Fila): Record<string, string> => {
    const out: Record<string, string> = { estado: fila.estado || 'ACTIVO' };
    for (const f of cfg.campos) {
      if (f.clave === 'estado') continue;
      if (fila.celdas[f.clave] !== undefined) out[f.clave] = fila.celdas[f.clave];
    }
    return out;
  };

  const leerCampos = (): Record<string, string> => {
    const out: Record<string, string> = { ...datosExtra };
    for (const f of cfg.campos) {
      const el = host.querySelector(`#campo-${f.clave}`) as HTMLInputElement | HTMLSelectElement;
      if (!el) continue;
      if (modoNuevo && f.soloEdicion) continue;
      if (!modoNuevo && f.soloNuevo) continue;
      out[f.clave] = el.value;
    }
    return out;
  };

  const limpiar = () => {
    seleccionId = 0;
    modoNuevo = true;
    datosExtra = {};
    for (const f of cfg.campos) {
      const el = host.querySelector(`#campo-${f.clave}`) as HTMLInputElement | HTMLSelectElement;
      if (!el) continue;
      if (f.tipo === 'estado' || f.tipo === 'select') {
        (el as HTMLSelectElement).value = f.clave === 'estado' ? 'ACTIVO' : (f.opciones?.[0] ?? '');
      } else {
        el.value = f.clave === 'tara' ? '0' : '';
      }
    }
    mensajeDialog.textContent = '';
    actualizarVisibilidadCampos();
    actualizarTituloDialog();
    marcarFila(0);
  };

  const rellenar = (fila: Fila) => {
    modoNuevo = false;
    seleccionId = fila.id;
    datosExtra = {};
    if (fila.celdas.logo !== undefined) datosExtra.logo = fila.celdas.logo;
    for (const f of cfg.campos) {
      const el = host.querySelector(`#campo-${f.clave}`) as HTMLInputElement | HTMLSelectElement;
      if (!el) continue;
      if (f.soloNuevo) continue;
      if (f.tipo === 'estado') {
        (el as HTMLSelectElement).value = fila.estado || 'ACTIVO';
      } else {
        el.value = fila.celdas[f.clave] ?? '';
      }
    }
    mensajeDialog.textContent = '';
    actualizarVisibilidadCampos();
    actualizarTituloDialog();
    marcarFila(fila.id);
  };

  const abrirDialogo = () => {
    actualizarVisibilidadCampos();
    actualizarTituloDialog();
    mensajeDialog.textContent = '';
    if (!dialog.open) dialog.showModal();
  };

  const cerrarDialogo = () => {
    if (dialog.open) dialog.close();
  };

  const cambiarEstadoFila = async (fila: Fila, activo: boolean, input: HTMLInputElement) => {
    const nuevo = activo ? 'ACTIVO' : 'INACTIVO';
    const campos = { ...camposDesdeFila(fila), estado: nuevo };
    const res = (await AbmGuardar({
      entidad: cfg.entidad,
      id: fila.id,
      campos,
    })) as { ok: boolean; mensaje: string; esError?: boolean };
    if (!res.ok) {
      sincronizarSwitchUI(input, !activo);
      mostrarMsg(res.mensaje, true);
      return;
    }
    fila.estado = nuevo;
    sincronizarSwitchUI(input, activo);
    if (seleccionId === fila.id && dialog.open) {
      const estEl = host.querySelector('#campo-estado') as HTMLSelectElement | null;
      if (estEl) estEl.value = nuevo;
    }
    mostrarMsg(res.mensaje);
  };

  const cargar = async () => {
    const filas = (await AbmListar(cfg.entidad)) as Fila[];
    const body = host.querySelector('#abm-body') as HTMLTableSectionElement;
    body.innerHTML = '';
    for (const f of filas) {
      const tr = document.createElement('tr');
      tr.dataset.id = String(f.id);
      const celdas = cfg.columnas
        .map((c) => `<td>${escapeHtml(f.celdas[c.clave] ?? '')}</td>`)
        .join('');
      const activo = f.estado.toUpperCase() === 'ACTIVO';
      const estadoTxt = etiquetaEstadoSwitch(activo);
      tr.innerHTML =
        celdas +
        `<td class="abm-celda-acciones">
          <button type="button" class="abm-btn-editar" data-tip="Editar ${entLabel}" aria-label="Editar ${entLabel}">${ICONO_LAPIZ}</button>
          <label class="abm-switch" data-estado="${estadoTxt}">
            <input type="checkbox" class="abm-switch-input" role="switch" ${activo ? 'checked' : ''} aria-label="${estadoTxt}" aria-checked="${activo ? 'true' : 'false'}" />
            <span class="abm-switch-pista"></span>
          </label>
        </td>`;
      tr.querySelector('.abm-btn-editar')?.addEventListener('click', (ev) => {
        ev.stopPropagation();
        rellenar(f);
        abrirDialogo();
      });
      const sw = tr.querySelector('.abm-switch-input') as HTMLInputElement;
      sw.addEventListener('click', (ev) => ev.stopPropagation());
      sw.addEventListener('change', () => {
        void cambiarEstadoFila(f, sw.checked, sw);
      });
      body.appendChild(tr);
    }
    marcarFila(seleccionId);
  };

  host.querySelector('#abm-crear')?.addEventListener('click', () => {
    limpiar();
    abrirDialogo();
  });

  host.querySelector('#abm-cancelar')?.addEventListener('click', () => {
    cerrarDialogo();
  });

  dialog.addEventListener('cancel', (ev) => {
    ev.preventDefault();
    cerrarDialogo();
  });

  form.addEventListener('submit', async (ev) => {
    ev.preventDefault();
    const campos = leerCampos();
    for (const f of cfg.campos) {
      if (!f.requerido) continue;
      if (modoNuevo && f.soloEdicion) continue;
      if (!modoNuevo && f.soloNuevo) continue;
      if (!modoNuevo && f.tipo === 'password') continue;
      if (!campos[f.clave]?.trim()) {
        mostrarMsgDialog(`El campo ${f.etiqueta} es obligatorio.`, true);
        return;
      }
    }
    if (
      !modoNuevo &&
      campos.estado === 'INACTIVO' &&
      !window.confirm('Se cambiará el estado a INACTIVO. ¿Continuar?')
    ) {
      return;
    }
    const res = (await AbmGuardar({
      entidad: cfg.entidad,
      id: seleccionId,
      campos,
    })) as { ok: boolean; mensaje: string; esError?: boolean };
    if (!res.ok) {
      mostrarMsgDialog(res.mensaje, true);
      return;
    }
    mostrarMsg(res.mensaje);
    cerrarDialogo();
    await cargar();
    limpiar();
  });

  void cargar();
  limpiar();
  return () => {};
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}