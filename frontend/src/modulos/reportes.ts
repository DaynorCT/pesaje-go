import { AbrirArchivo, GenerarBoletaPdf, GenerarReportePdf, ReportesListar } from '../../wailsjs/go/main/App';
import { enlazarMayusculasVivoEnContenedor, normalizarTextoGuardado } from '../normalizar-campo';

const ICONO_CALENDARIO = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>`;

type Fila = {
  id: number;
  guia: string;
  placa: string;
  nombreChofer: string;
  nombreProveedor: string;
  pesoBruto: number;
  tara: number;
  pesoNeto: number;
  estadoBalanza: string;
  fecha: string;
  hora: string;
};

/** Texto dd/mm/aaaa (o yyyy-mm-dd) → yyyy-mm-dd para la API. */
function aIsoFecha(raw: string): string {
  const t = raw.trim();
  if (!t) return '';
  if (/^\d{4}-\d{2}-\d{2}$/.test(t)) return t;
  const m = t.match(/^(\d{1,2})[/.-](\d{1,2})[/.-](\d{4})$/);
  if (m) {
    const dd = m[1].padStart(2, '0');
    const mm = m[2].padStart(2, '0');
    return `${m[3]}-${mm}-${dd}`;
  }
  return '';
}

function esFechaIsoValida(iso: string): boolean {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(iso)) return false;
  const [y, m, d] = iso.split('-').map((x) => Number(x));
  const dt = new Date(y, m - 1, d);
  return dt.getFullYear() === y && dt.getMonth() === m - 1 && dt.getDate() === d;
}

function isoAVisual(iso: string): string {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(iso)) return '';
  const [y, m, d] = iso.split('-');
  return `${d}/${m}/${y}`;
}

function htmlCampoFecha(id: string, etiqueta: string): string {
  const pickerId = `${id}-cal`;
  return `
    <label class="rep-filtro rep-filtro--fecha">
      <span class="etiqueta-mini">${etiqueta}</span>
      <span class="rep-fecha-wrap">
        <input
          type="text"
          class="text-field rep-fecha"
          id="${id}"
          data-casing="preserve"
          inputmode="numeric"
          placeholder="dd/mm/aaaa"
          autocomplete="off"
          spellcheck="false"
          maxlength="10"
        />
        <span class="rep-fecha-cal" title="Calendario">
          <span class="rep-fecha-icono-visual">${ICONO_CALENDARIO}</span>
          <input type="date" id="${pickerId}" class="rep-fecha-picker-native" data-text="${id}" title="Elegir fecha" />
        </span>
      </span>
    </label>`;
}

function enlazarCampoFecha(host: HTMLElement, textId: string, pickerId: string): void {
  const text = host.querySelector(`#${textId}`) as HTMLInputElement;
  const picker = host.querySelector(`#${pickerId}`) as HTMLInputElement;
  if (!text || !picker) return;

  const sincronizarTextoAPicker = () => {
    const iso = aIsoFecha(text.value);
    if (esFechaIsoValida(iso)) {
      text.value = isoAVisual(iso);
      picker.value = iso;
    }
  };

  text.addEventListener('blur', sincronizarTextoAPicker);

  picker.addEventListener('change', () => {
    if (picker.value) text.value = isoAVisual(picker.value);
  });

  picker.addEventListener('click', () => {
    const iso = aIsoFecha(text.value);
    if (esFechaIsoValida(iso)) picker.value = iso;
  });
}

export function montarReportes(host: HTMLElement): () => void {
  let seleccion = 0;

  host.innerHTML = `
    <div class="reportes-panel">
      <div class="reportes-barra">
        ${htmlCampoFecha('rep-desde', 'Desde')}
        ${htmlCampoFecha('rep-hasta', 'Hasta')}
        <label class="rep-filtro rep-filtro--placa">
          <span class="etiqueta-mini">Placa</span>
          <input class="text-field" id="rep-placa" placeholder="Todas" />
        </label>
        <div class="rep-filtro-acciones">
          <button type="button" class="boton-secundario" id="rep-actualizar">Actualizar</button>
          <button type="button" class="boton-primario" id="rep-pdf">Reporte PDF</button>
          <button type="button" class="boton-primario" id="rep-boleta">Boleta del pesaje</button>
        </div>
      </div>
      <p id="rep-mensaje" class="subtitulo"></p>
      <div class="tarjeta tarjeta-sin-padding tabla-modulo-wrap">
        <table class="tabla-modulo" id="tabla-reportes">
          <thead><tr>
            <th>Guía</th><th>Vehículo</th><th>Chofer</th><th>Proveedor</th>
            <th>Bruto</th><th>Tara</th><th>Neto</th><th class="tabla-th-estado">Estado</th><th>Fecha</th>
          </tr></thead>
          <tbody id="rep-body"></tbody>
        </table>
      </div>
    </div>
  `;

  enlazarMayusculasVivoEnContenedor(host);
  enlazarCampoFecha(host, 'rep-desde', 'rep-desde-cal');
  enlazarCampoFecha(host, 'rep-hasta', 'rep-hasta-cal');

  const msg = host.querySelector('#rep-mensaje') as HTMLElement;

  const isoDesdeCampo = (id: string): string => {
    const raw = (host.querySelector(`#${id}`) as HTMLInputElement).value.trim();
    return aIsoFecha(raw);
  };

  const validarFiltros = (): boolean => {
    for (const [id, etiqueta] of [
      ['rep-desde', 'Desde'],
      ['rep-hasta', 'Hasta'],
    ] as const) {
      const raw = (host.querySelector(`#${id}`) as HTMLInputElement).value.trim();
      if (!raw) continue;
      const iso = aIsoFecha(raw);
      if (!esFechaIsoValida(iso)) {
        msg.textContent = `Fecha inválida en «${etiqueta}». Use dd/mm/aaaa.`;
        msg.className = 'mensaje-error';
        return false;
      }
    }
    const desde = isoDesdeCampo('rep-desde');
    const hasta = isoDesdeCampo('rep-hasta');
    if (desde && hasta && desde > hasta) {
      msg.textContent = '«Desde» no puede ser posterior a «Hasta».';
      msg.className = 'mensaje-error';
      return false;
    }
    msg.className = 'subtitulo';
    return true;
  };

  const filtros = () => ({
    desde: isoDesdeCampo('rep-desde'),
    hasta: isoDesdeCampo('rep-hasta'),
    placa: normalizarTextoGuardado('placa', (host.querySelector('#rep-placa') as HTMLInputElement).value),
  });

  const cargar = async () => {
    if (!validarFiltros()) return;
    const filas = (await ReportesListar(filtros(), 500)) as Fila[];
    const body = host.querySelector('#rep-body') as HTMLTableSectionElement;
    body.innerHTML = '';
    seleccion = 0;
    if (filas.length === 0) {
      msg.textContent = 'No hay pesajes con esos filtros.';
    } else {
      msg.textContent = `${filas.length} registro(s).`;
    }
    msg.className = 'subtitulo';
    for (const f of filas) {
      const tr = document.createElement('tr');
      const enProceso = f.estadoBalanza.toUpperCase() === 'EN_PROCESO';
      const estadoTxt = f.estadoBalanza.replace(/_/g, ' ');
      const estadoClase = enProceso ? 'proceso' : 'finalizado';
      tr.innerHTML = `
        <td>${f.guia || f.id}</td>
        <td>${f.placa}</td>
        <td>${f.nombreChofer}</td>
        <td>${f.nombreProveedor}</td>
        <td>${f.pesoBruto}</td>
        <td>${f.tara}</td>
        <td>${f.pesoNeto}</td>
        <td class="tabla-celda-estado"><span class="estado-etiqueta estado-etiqueta--${estadoClase}">${estadoTxt}</span></td>
        <td>${f.fecha} ${f.hora}</td>
      `;
      tr.addEventListener('click', () => {
        seleccion = f.id;
        body.querySelectorAll('tr').forEach((r) => r.classList.remove('seleccionado'));
        tr.classList.add('seleccionado');
      });
      body.appendChild(tr);
    }
  };

  const abrirPdf = async (res: { ok: boolean; ruta?: string; mensaje: string; esError?: boolean }) => {
    msg.textContent = res.mensaje;
    msg.className = res.esError || !res.ok ? 'mensaje-error' : 'subtitulo mensaje-exito';
    if (res.ok && res.ruta) await AbrirArchivo(res.ruta);
  };

  host.querySelector('#rep-actualizar')?.addEventListener('click', () => void cargar());
  host.querySelector('#rep-pdf')?.addEventListener('click', () => {
    if (!validarFiltros()) return;
    void GenerarReportePdf(filtros()).then((r) =>
      abrirPdf(r as { ok: boolean; ruta?: string; mensaje: string; esError?: boolean }),
    );
  });
  host.querySelector('#rep-boleta')?.addEventListener('click', () => {
    if (!seleccion) {
      msg.textContent = 'Seleccione un pesaje de la lista.';
      msg.className = 'mensaje-error';
      return;
    }
    void GenerarBoletaPdf(seleccion).then((r) =>
      abrirPdf(r as { ok: boolean; ruta?: string; mensaje: string; esError?: boolean }),
    );
  });

  void cargar();
  return () => {};
}