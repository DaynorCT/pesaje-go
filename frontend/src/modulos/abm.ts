import { AbmGuardar, AbmListar } from '../../wailsjs/go/main/App';

export type CampoAbm = {
  clave: string;
  etiqueta: string;
  tipo?: 'text' | 'password' | 'estado' | 'select';
  requerido?: boolean;
  opciones?: string[];
  soloEdicion?: boolean;
  soloNuevo?: boolean;
};

export type ColumnaAbm = {
  clave: string;
  titulo: string;
};

export type ConfigAbm = {
  entidad: string;
  columnas: ColumnaAbm[];
  campos: CampoAbm[];
};

type Fila = { id: number; estado: string; celdas: Record<string, string> };

export function montarAbm(host: HTMLElement, cfg: ConfigAbm): () => void {
  host.innerHTML = `
    <div class="abm-panel">
      <div class="abm-tabla-wrap tarjeta tarjeta-sin-padding">
        <table class="tabla-modulo" id="abm-tabla">
          <thead><tr id="abm-head"></tr></thead>
          <tbody id="abm-body"></tbody>
        </table>
      </div>
      <aside class="abm-form tarjeta">
        <div class="abm-campos" id="abm-campos"></div>
        <div class="abm-acciones">
          <button type="button" class="boton-secundario" id="abm-nuevo">Nuevo</button>
          <button type="button" class="boton-primario" id="abm-guardar">Guardar</button>
        </div>
        <p id="abm-mensaje" class="subtitulo abm-mensaje"></p>
      </aside>
    </div>
  `;

  const head = host.querySelector('#abm-head') as HTMLTableRowElement;
  head.innerHTML =
    cfg.columnas.map((c) => `<th>${c.titulo}</th>`).join('') + '<th>Estado</th>';

  const camposHost = host.querySelector('#abm-campos') as HTMLElement;
  camposHost.innerHTML = cfg.campos
    .map((f) => {
      if (f.tipo === 'estado' || f.tipo === 'select') {
        const opts = (f.opciones ?? ['ACTIVO', 'INACTIVO'])
          .map((o) => `<option value="${o}">${o}</option>`)
          .join('');
        return `<label class="pesaje-campo" data-campo="${f.clave}"><span class="etiqueta-mini">${f.etiqueta}</span>
          <select class="combo-web" id="campo-${f.clave}">${opts}</select></label>`;
      }
      const type = f.tipo === 'password' ? 'password' : 'text';
      return `<label class="pesaje-campo" data-campo="${f.clave}"><span class="etiqueta-mini">${f.etiqueta}</span>
        <input class="text-field" type="${type}" id="campo-${f.clave}" /></label>`;
    })
    .join('');

  let seleccionId = 0;
  let modoNuevo = true;
  let datosExtra: Record<string, string> = {};

  const mensaje = host.querySelector('#abm-mensaje') as HTMLElement;

  const mostrarMsg = (texto: string, error = false) => {
    mensaje.textContent = texto;
    mensaje.className = error ? 'mensaje-error abm-mensaje' : 'subtitulo abm-mensaje mensaje-exito';
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
      el.closest('[data-campo]')?.classList.remove('oculto');
      if (f.soloNuevo) el.closest('[data-campo]')?.classList.toggle('oculto', !modoNuevo);
      if (f.soloEdicion) el.closest('[data-campo]')?.classList.toggle('oculto', modoNuevo);
    }
    mensaje.textContent = '';
  };

  const rellenar = (fila: Fila) => {
    modoNuevo = false;
    seleccionId = fila.id;
    datosExtra = {};
    if (fila.celdas.logo !== undefined) datosExtra.logo = fila.celdas.logo;
    for (const f of cfg.campos) {
      const el = host.querySelector(`#campo-${f.clave}`) as HTMLInputElement | HTMLSelectElement;
      if (!el) continue;
      if (f.soloNuevo) {
        el.closest('[data-campo]')?.classList.add('oculto');
        continue;
      }
      el.closest('[data-campo]')?.classList.remove('oculto');
      if (f.tipo === 'estado') {
        (el as HTMLSelectElement).value = fila.estado || 'ACTIVO';
      } else {
        el.value = fila.celdas[f.clave] ?? '';
      }
    }
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
      tr.innerHTML =
        celdas +
        `<td><span class="badge badge-${activo ? 'activo' : 'inactivo'}">${f.estado}</span></td>`;
      tr.addEventListener('click', () => rellenar(f));
      body.appendChild(tr);
    }
  };

  host.querySelector('#abm-nuevo')?.addEventListener('click', () => {
    limpiar();
    (host.querySelector('#abm-body') as HTMLElement)
      .querySelectorAll('tr')
      .forEach((tr) => tr.classList.remove('seleccionado'));
  });

  host.querySelector('#abm-guardar')?.addEventListener('click', async () => {
    const campos = leerCampos();
    for (const f of cfg.campos) {
      if (!f.requerido) continue;
      if (modoNuevo && f.soloEdicion) continue;
      if (!modoNuevo && f.soloNuevo) continue;
      if (!modoNuevo && f.tipo === 'password') continue;
      if (!campos[f.clave]?.trim()) {
        mostrarMsg(`El campo ${f.etiqueta} es obligatorio.`, true);
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
      mostrarMsg(res.mensaje, true);
      return;
    }
    mostrarMsg(res.mensaje);
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