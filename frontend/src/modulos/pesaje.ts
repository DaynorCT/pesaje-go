import {
  PesajeActualizarTara,
  PesajeAlternarBalanza,
  PesajeBoleta,
  PesajeCapturarBruto,
  PesajeCapturarTara,
  PesajeCambiarEstado,
  PesajeCerrar,
  PesajeCrearChofer,
  PesajeCrearProveedor,
  PesajeCrearVehiculo,
  PesajeDetalleChofer,
  PesajeDetalleVehiculo,
  PesajeEstado,
  PesajeFinalizar,
  PesajeListarGrilla,
  PesajeOpciones,
  PesajeRegistrar,
} from '../../wailsjs/go/main/App';
import { crearDialogoAviso, enlazarDialogo, htmlDialogoFormGuardar } from '../dialogo';

type Opcion = { id: number; label: string };

type EstadoPane = {
  balanzaConectada: boolean;
  textoConectar: string;
  pesoVivo: number;
  pesoBruto: number;
  pesoTara: number;
  pesoNeto: number;
  puedeCapturar: boolean;
  mensaje: string;
  mensajeEsError: boolean;
};

type Fila = {
  id: number;
  nombreChofer: string;
  placa: string;
  licencia: string;
  tipoVehiculo: string;
  nombreProveedor: string;
  producto: string;
  origen: string;
  destino: string;
  costoBs: number;
  fleteBs: number;
  fecha: string;
  hora: string;
  pesoBruto: number;
  tara: number;
  pesoNeto: number;
  estado: string;
  estadoBalanza: string;
};

type AccionPesaje = {
  ok: boolean;
  estadoPesaje?: EstadoPane;
  aviso?: { titulo: string; mensaje: string; esError?: boolean };
};

export function montarPesaje(host: HTMLElement): () => void {
  let timer = 0;
  let taraManual = false;
  let opciones: {
    vehiculos: Opcion[];
    choferes: Opcion[];
    proveedores: Opcion[];
    productos: Opcion[];
    origenes: Opcion[];
    destinos: Opcion[];
  } = { vehiculos: [], choferes: [], proveedores: [], productos: [], origenes: [], destinos: [] };

  host.innerHTML = `
    <div class="panel-pesaje">
      <div class="panel-pesaje-medio">
        <section class="tarjeta pesaje-registro">
          <h2 class="titulo-tarjeta">Registro de peso</h2>
          <div class="separador-horizontal"></div>
          <div class="display-peso"><span class="valor-peso" id="pesaje-vivo">0 kg</span></div>
          <div class="separador-horizontal"></div>
          <div class="pesaje-acciones-balanza">
            <div class="pesaje-col-conexion">
              <button type="button" class="boton-secundario" id="btn-conectar">Conectar</button>
              <span class="etiqueta-mini">Conexión</span>
            </div>
            <button type="button" class="boton-primario" id="btn-cap-bruto" disabled>Cap. peso</button>
            <button type="button" class="boton-info" id="btn-cap-tara" disabled>Cap. tara</button>
          </div>
          <div class="pesaje-resumen">
            <div class="caja-peso-resumen"><span class="etiqueta-mini">Peso Bruto</span><strong id="res-bruto" class="valor-caja-peso">0</strong></div>
            <div class="caja-peso-resumen"><span class="etiqueta-mini">Peso tara</span><strong id="res-tara" class="valor-caja-peso">0</strong></div>
            <div class="caja-peso-resumen"><span class="etiqueta-mini">Peso Neto</span><strong id="res-neto" class="valor-caja-peso">0</strong></div>
          </div>
          <p id="pesaje-mensaje" class="subtitulo pesaje-mensaje"></p>
        </section>
        <section class="tarjeta pesaje-datos">
          <h2 class="titulo-tarjeta">Datos del Pesaje</h2>
          <div class="separador-horizontal"></div>
          <div class="pesaje-grid-datos" id="grid-datos"></div>
          <div class="separador-horizontal"></div>
          <div class="pesaje-guardar">
            <button type="button" class="boton-primario" id="btn-guardar">Registrar Pesaje</button>
          </div>
        </section>
      </div>
      <section class="tarjeta tarjeta-sin-padding pesaje-grilla-wrap tabla-modulo-wrap pesaje-tabla-ajuste">
        <table class="tabla-modulo tabla-modulo--ajuste-ventana" id="tabla-pesajes">
          <colgroup>
            <col class="pesaje-col-id" />
            <col class="pesaje-col-texto" />
            <col class="pesaje-col-placa" />
            <col class="pesaje-col-corto" />
            <col class="pesaje-col-corto" />
            <col class="pesaje-col-texto" />
            <col class="pesaje-col-texto" />
            <col class="pesaje-col-texto" />
            <col class="pesaje-col-texto" />
            <col class="pesaje-col-num" />
            <col class="pesaje-col-num" />
            <col class="pesaje-col-fecha" />
            <col class="pesaje-col-num" />
            <col class="pesaje-col-num" />
            <col class="pesaje-col-num" />
            <col class="pesaje-col-estado" />
            <col class="pesaje-col-estado" />
            <col class="pesaje-col-acciones" />
          </colgroup>
          <thead><tr id="tabla-head"></tr></thead>
          <tbody id="tabla-body"></tbody>
        </table>
      </section>
      <div id="pesaje-dialogos"></div>
    </div>
  `;

  const qs = <T extends HTMLElement>(sel: string) => host.querySelector(sel) as T;
  const slotDialogos = qs<HTMLElement>('#pesaje-dialogos');
  const aviso = crearDialogoAviso(slotDialogos, 'dlg-aviso', 'medio');
  const wrapRapido = document.createElement('div');
  wrapRapido.innerHTML = htmlDialogoFormGuardar({
    dialogId: 'dlg-rapido',
    formId: 'form-rapido',
    tituloId: 'rapido-titulo',
    cuerpoId: 'rapido-campos',
    mensajeId: 'rapido-msg',
    cancelarId: 'rapido-cancelar',
    ancho: 'medio',
  });
  slotDialogos.appendChild(wrapRapido.firstElementChild!);
  const dlgRapido = qs<HTMLDialogElement>('#dlg-rapido');
  const formRapido = qs<HTMLFormElement>('#form-rapido');
  const gridDatos = qs<HTMLElement>('#grid-datos');
  enlazarDialogo(dlgRapido, [qs<HTMLButtonElement>('#rapido-cancelar')]);

  const mostrarAviso = (titulo: string, mensaje: string) => {
    aviso.mostrar(titulo, mensaje);
  };

  const aplicarEstado = (e: EstadoPane) => {
    qs('#pesaje-vivo').textContent = `${e.pesoVivo.toLocaleString('es-BO')} kg`;
    qs('#btn-conectar').textContent = e.textoConectar;
    qs<HTMLButtonElement>('#btn-cap-bruto').disabled = !e.puedeCapturar;
    qs<HTMLButtonElement>('#btn-cap-tara').disabled = !e.puedeCapturar;
    qs('#res-bruto').textContent = String(e.pesoBruto);
    qs('#res-tara').textContent = String(e.pesoTara);
    qs('#res-neto').textContent = String(e.pesoNeto);
    const msg = qs('#pesaje-mensaje');
    msg.textContent = e.mensaje;
    msg.className = e.mensajeEsError ? 'mensaje-error pesaje-mensaje' : 'subtitulo pesaje-mensaje mensaje-exito';
    if (!e.mensaje) msg.className = 'subtitulo pesaje-mensaje';
  };

  const manejarAccion = (a: AccionPesaje) => {
    if (a.estadoPesaje) aplicarEstado(a.estadoPesaje);
    if (a.aviso) mostrarAviso(a.aviso.titulo, a.aviso.mensaje);
  };

  const selectHTML = (id: string, label: string, plus?: string) => `
    <label class="pesaje-campo">
      <span class="etiqueta-mini">${label}</span>
      <span class="pesaje-fila-combo">
        <select class="combo-web" id="${id}"></select>
        ${plus ? `<button type="button" class="boton-mas" data-plus="${plus}">+</button>` : ''}
      </span>
    </label>
  `;

  gridDatos.innerHTML = `
    ${selectHTML('sel-chofer', 'Chofer', 'chofer')}
    ${selectHTML('sel-vehiculo', 'Placa *', 'vehiculo')}
    <label class="pesaje-campo"><span class="etiqueta-mini">Licencia</span><input class="text-field" id="licencia" readonly placeholder="—" /></label>
    <label class="pesaje-campo"><span class="etiqueta-mini">Tipo vehículo</span><input class="text-field" id="tipo-vehiculo" readonly placeholder="—" /></label>
    ${selectHTML('sel-proveedor', 'Proveedor', 'proveedor')}
    ${selectHTML('sel-producto', 'Producto')}
    ${selectHTML('sel-origen', 'Origen')}
    ${selectHTML('sel-destino', 'Destino')}
    <label class="pesaje-campo"><span class="etiqueta-mini">Costo (Bs)</span><input class="text-field" id="costo" value="0" inputmode="numeric" /></label>
    <label class="pesaje-campo"><span class="etiqueta-mini">Flete pend. (Bs)</span><input class="text-field" id="flete" value="0" inputmode="numeric" /></label>
    <label class="pesaje-campo pesaje-tara-campo">
      <label class="pesaje-check"><input type="checkbox" id="tara-manual" /> Tara manual</label>
      <span class="etiqueta-mini">Tara (kg)</span>
      <input class="text-field" id="tara" value="0" inputmode="numeric" readonly />
    </label>
  `;

  const llenarSelect = (id: string, items: Opcion[]) => {
    const sel = qs<HTMLSelectElement>(`#${id}`);
    sel.innerHTML = items.map((o) => `<option value="${o.id}">${escapeHtml(o.label)}</option>`).join('');
  };

  const recargarOpciones = async () => {
    const op = (await PesajeOpciones()) as typeof opciones & {
      vehiculos: Opcion[];
    };
    opciones = {
      vehiculos: mapOp(op.vehiculos),
      choferes: mapOp(op.choferes),
      proveedores: mapOp(op.proveedores),
      productos: mapOp(op.productos),
      origenes: mapOp(op.origenes),
      destinos: mapOp(op.destinos),
    };
    llenarSelect('sel-vehiculo', opciones.vehiculos);
    llenarSelect('sel-chofer', opciones.choferes);
    llenarSelect('sel-proveedor', opciones.proveedores);
    llenarSelect('sel-producto', opciones.productos);
    llenarSelect('sel-origen', opciones.origenes);
    llenarSelect('sel-destino', opciones.destinos);
    await alCambiarVehiculo();
    await alCambiarChofer();
  };

  const mapOp = (xs: { id: number; label: string }[]): Opcion[] =>
    (xs || []).map((x) => ({ id: x.id, label: x.label }));

  const alCambiarChofer = async () => {
    const id = Number(qs<HTMLSelectElement>('#sel-chofer').value);
    if (!id) {
      qs<HTMLInputElement>('#licencia').value = '';
      return;
    }
    const det = (await PesajeDetalleChofer(id)) as { licencia: string };
    qs<HTMLInputElement>('#licencia').value = det?.licencia ?? '';
  };

  const alCambiarVehiculo = async () => {
    const id = Number(qs<HTMLSelectElement>('#sel-vehiculo').value);
    if (!id) {
      qs<HTMLInputElement>('#tipo-vehiculo').value = '';
      return;
    }
    const det = (await PesajeDetalleVehiculo(id)) as { tipo: string; tara: number };
    qs<HTMLInputElement>('#tipo-vehiculo').value = det?.tipo ?? '';
    if (!taraManual && det) {
      qs<HTMLInputElement>('#tara').value = String(det.tara ?? 0);
      void PesajeActualizarTara(det.tara ?? 0).then((est) => aplicarEstado(est as EstadoPane));
    }
  };

  const columnas = [
    'ID', 'Chofer', 'Placa', 'Lic.', 'Tipo', 'Prov.', 'Prod.', 'Orig.', 'Dest.',
    'Costo', 'Flete', 'Fecha', 'Bruto', 'Tara', 'Neto', 'Est.', 'Balanza', 'Acc.',
  ];
  const claseTh = (c: string) => {
    if (c === 'Est.' || c === 'Balanza') return 'tabla-th-estado';
    if (c === 'Acc.') return 'tabla-th-acciones';
    return '';
  };
  qs('#tabla-head').innerHTML = columnas
    .map((c) => {
      const cl = claseTh(c);
      return cl ? `<th class="${cl}">${c}</th>` : `<th>${c}</th>`;
    })
    .join('');

  const recargarGrilla = async () => {
    const filas = (await PesajeListarGrilla(50)) as Fila[];
    const body = qs<HTMLTableSectionElement>('#tabla-body');
    body.innerHTML = '';
    for (const f of filas) {
      const tr = document.createElement('tr');
      const activo = f.estado.toUpperCase() === 'ACTIVO';
      const enProceso = f.estadoBalanza.toUpperCase() === 'EN_PROCESO';
      let acciones = '';
      if (activo) {
        acciones += `<button type="button" class="boton-secundario boton-tabla-accion" data-tip="Desactivar" data-act="estado" data-id="${f.id}">Des.</button> `;
        if (enProceso) {
          acciones += `<button type="button" class="boton-secundario boton-tabla-accion" data-tip="Finalizar pesaje" data-act="finalizar" data-id="${f.id}">Fin.</button>`;
        } else {
          acciones += `<button type="button" class="boton-secundario boton-tabla-accion" data-tip="Boleta PDF" data-act="boleta" data-id="${f.id}">Bol.</button>`;
        }
      } else {
        acciones += `<button type="button" class="boton-secundario boton-tabla-accion" data-tip="Reactivar" data-act="activar" data-id="${f.id}">Act.</button>`;
      }
      const balanzaTxt = etiquetaBalanzaCorta(f.estadoBalanza);
      const fechaHora = `${f.fecha} ${f.hora}`.trim();
      tr.innerHTML = `
        ${tdCelda(String(f.id))}
        ${tdCelda(mayus(f.nombreChofer))}
        ${tdCelda(mayus(f.placa))}
        ${tdCelda(mayus(f.licencia))}
        ${tdCelda(mayus(f.tipoVehiculo))}
        ${tdCelda(mayus(f.nombreProveedor))}
        ${tdCelda(mayus(f.producto))}
        ${tdCelda(mayus(f.origen))}
        ${tdCelda(mayus(f.destino))}
        ${tdCelda(String(f.costoBs))}
        ${tdCelda(String(f.fleteBs))}
        ${tdCelda(fechaHora)}
        ${tdCelda(String(f.pesoBruto))}
        ${tdCelda(String(f.tara))}
        ${tdCelda(String(f.pesoNeto))}
        <td class="tabla-celda-estado" title="${escapeAttr(activo ? 'ACTIVO' : 'INACTIVO')}"><span class="estado-etiqueta estado-etiqueta--${activo ? 'activo' : 'inactivo'}">${activo ? 'ACT.' : 'INA.'}</span></td>
        <td class="tabla-celda-estado" title="${escapeAttr(f.estadoBalanza.replace(/_/g, ' '))}"><span class="estado-etiqueta estado-etiqueta--${enProceso ? 'proceso' : 'finalizado'}">${balanzaTxt}</span></td>
        <td class="tabla-celda-acciones">${acciones}</td>
      `;
      body.appendChild(tr);
    }
  };

  qs('#btn-conectar').addEventListener('click', () => {
    void PesajeAlternarBalanza().then((r) => manejarAccion(r as AccionPesaje));
  });
  qs('#btn-cap-bruto').addEventListener('click', () => {
    void PesajeCapturarBruto().then((r) => manejarAccion(r as AccionPesaje));
  });
  qs('#btn-cap-tara').addEventListener('click', () => {
    void PesajeCapturarTara().then((r) => {
      manejarAccion(r as AccionPesaje);
      taraManual = true;
      qs<HTMLInputElement>('#tara-manual').checked = true;
      qs<HTMLInputElement>('#tara').readOnly = false;
    });
  });

  qs('#sel-chofer').addEventListener('change', () => void alCambiarChofer());
  qs('#sel-vehiculo').addEventListener('change', () => void alCambiarVehiculo());

  qs<HTMLInputElement>('#tara-manual').addEventListener('change', (ev) => {
    taraManual = (ev.target as HTMLInputElement).checked;
    qs<HTMLInputElement>('#tara').readOnly = !taraManual;
    if (!taraManual) void alCambiarVehiculo();
  });

  qs<HTMLInputElement>('#tara').addEventListener('input', () => {
    if (!taraManual) return;
    const t = parseInt(qs<HTMLInputElement>('#tara').value, 10) || 0;
    void PesajeActualizarTara(t).then((est) => aplicarEstado(est as EstadoPane));
  });

  qs('#btn-guardar').addEventListener('click', () => {
    const vehiculoId = Number(qs<HTMLSelectElement>('#sel-vehiculo').value);
    const payload = {
      vehiculoId,
      choferId: opcionId(host, 'sel-chofer'),
      proveedorId: opcionId(host, 'sel-proveedor'),
      productoId: opcionId(host, 'sel-producto'),
      origenId: opcionId(host, 'sel-origen'),
      destinoId: opcionId(host, 'sel-destino'),
      pesoBruto: parseInt(qs('#res-bruto').textContent || '0', 10),
      tara: taraManual ? parseInt(qs<HTMLInputElement>('#tara').value, 10) || 0 : parseInt(qs<HTMLInputElement>('#tara').value, 10) || 0,
      costoBs: parseInt(qs<HTMLInputElement>('#costo').value, 10) || 0,
      fleteBsPendiente: parseInt(qs<HTMLInputElement>('#flete').value, 10) || 0,
      taraManual,
    };
    void PesajeRegistrar(payload).then(async (r) => {
      manejarAccion(r as AccionPesaje);
      if ((r as AccionPesaje).ok) await recargarGrilla();
    });
  });

  host.querySelector('#tabla-body')?.addEventListener('click', (ev) => {
    const btn = (ev.target as HTMLElement).closest('button[data-act]') as HTMLButtonElement | null;
    if (!btn) return;
    const id = Number(btn.dataset.id);
    const act = btn.dataset.act;
    if (act === 'estado' && window.confirm('Se cambiará el estado a INACTIVO. ¿Continuar?')) {
      void PesajeCambiarEstado(id, false).then(async (r) => {
        if ((r as AccionPesaje).ok) await recargarGrilla();
      });
    }
    if (act === 'activar' && window.confirm('¿Reactivar este pesaje?')) {
      void PesajeCambiarEstado(id, true).then(async (r) => {
        if ((r as AccionPesaje).ok) await recargarGrilla();
      });
    }
    if (act === 'finalizar') {
      if (!window.confirm(`¿Cerrar el pesaje #${id} como FINALIZADO?`)) return;
      void PesajeFinalizar(id).then(async (r) => {
        manejarAccion(r as AccionPesaje);
        if ((r as AccionPesaje).ok) await recargarGrilla();
      });
    }
    if (act === 'boleta') {
      void PesajeBoleta(id).then((r) => manejarAccion(r as AccionPesaje));
    }
  });

  type RapidoDef = { tipo: string; titulo: string; campos: string[]; defaults: string[] };
  const rapidos: Record<string, RapidoDef> = {
    vehiculo: { tipo: 'vehiculo', titulo: 'Nuevo vehículo', campos: ['Placa *', 'Tipo', 'Tara (kg)'], defaults: ['', '', '0'] },
    chofer: {
      tipo: 'chofer',
      titulo: 'Nuevo chofer',
      campos: ['Nombre *', 'Apellido paterno', 'Apellido materno', 'Nro. documento', 'Teléfono', 'Licencia *'],
      defaults: ['', '', '', '', '', ''],
    },
    proveedor: {
      tipo: 'proveedor',
      titulo: 'Nuevo proveedor',
      campos: ['Nombre *', 'Apellido paterno', 'Apellido materno', 'CI', 'Teléfono', 'Empresa *', 'Descripción'],
      defaults: ['', '', '', '', '', '', ''],
    },
  };

  let rapidoActual: RapidoDef | null = null;

  gridDatos.querySelectorAll('button[data-plus]').forEach((btn) => {
    btn.addEventListener('click', () => {
      const t = btn.getAttribute('data-plus')!;
      rapidoActual = rapidos[t];
      if (!rapidoActual) return;
      qs('#rapido-titulo').textContent = rapidoActual.titulo;
      const cont = qs('#rapido-campos');
      cont.innerHTML = rapidoActual.campos
        .map(
          (c, i) =>
            `<label class="pesaje-campo"><span class="etiqueta-mini">${c}</span><input class="text-field" data-idx="${i}" value="${escapeAttr(rapidoActual!.defaults[i])}" /></label>`,
        )
        .join('');
      dlgRapido.showModal();
    });
  });

  formRapido.addEventListener('submit', (ev) => {
    ev.preventDefault();
    if (!rapidoActual) return;
    const vals = Array.from(formRapido.querySelectorAll('input[data-idx]')).map(
      (el) => (el as HTMLInputElement).value,
    );
    const fin = async () => {
      dlgRapido.close();
      await recargarOpciones();
      const e = (await PesajeEstado()) as EstadoPane;
      aplicarEstado(e);
    };
    if (rapidoActual.tipo === 'vehiculo') {
      void PesajeCrearVehiculo(vals[0], vals[1], parseInt(vals[2], 10) || 0).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (rapidoActual.tipo === 'chofer') {
      void PesajeCrearChofer(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (rapidoActual.tipo === 'proveedor') {
      void PesajeCrearProveedor(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    }
  });

  void PesajeEstado().then((e) => aplicarEstado(e as EstadoPane));
  void recargarOpciones();
  void recargarGrilla();

  timer = window.setInterval(() => {
    void PesajeEstado().then((e) => aplicarEstado(e as EstadoPane));
  }, 1500);

  return () => {
    window.clearInterval(timer);
    PesajeCerrar();
  };
}

function opcionId(root: HTMLElement, selId: string): number | undefined {
  const v = (root.querySelector(`#${selId}`) as HTMLSelectElement)?.value;
  return v ? Number(v) : undefined;
}

function mayus(s: string): string {
  return (s || '').toUpperCase();
}

function etiquetaBalanzaCorta(estado: string): string {
  const u = (estado || '').toUpperCase();
  if (u === 'EN_PROCESO') return 'PROC.';
  if (u === 'FINALIZADO') return 'FIN.';
  return mayus(estado).replace(/_/g, ' ').slice(0, 8);
}

function tdCelda(texto: string, clases = ''): string {
  const t = texto ?? '';
  const cls = clases ? ` class="${clases}"` : '';
  const tip = t ? ` title="${escapeAttr(t)}"` : '';
  return `<td${cls}${tip}>${escapeHtml(t)}</td>`;
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

function escapeAttr(s: string): string {
  return escapeHtml(s).replace(/"/g, '&quot;');
}