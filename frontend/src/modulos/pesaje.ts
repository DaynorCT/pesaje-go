import {
  PesajeActualizarRegistro,
  PesajeActualizarTara,
  PesajeCancelarEdicion,
  PesajeCargarEdicion,
  PesajeAlternarBalanza,
  PesajeBoleta,
  PesajeCapturarBruto,
  PesajeCapturarTara,
  PesajeCambiarEstado,
  PesajeCerrar,
  PesajeCrearChofer,
  PesajeCrearDestino,
  PesajeCrearOrigen,
  PesajeCrearProducto,
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
import {
  aplicarLayoutDialogoCampos,
  crearDialogoAviso,
  crearDialogoConfirmar,
  enlazarDialogo,
  htmlDialogoFormGuardar,
} from '../dialogo';
import {
  aMayusculas,
  enlazarMayusculasVivo,
  enlazarMayusculasVivoEnContenedor,
  normalizarTextoGuardado,
} from '../normalizar-campo';
import {
  ICONO_BOLETA_TABLA,
  ICONO_EDITAR_TABLA,
  ICONO_FINALIZAR_TABLA,
  htmlSwitchFila,
  sincronizarSwitchEstado,
} from '../tablas-acciones';
import { iconoModuloHTML } from '../iconos-modulo';

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

/** null = se puede finalizar; string = motivo por el que no. */
function motivoNoFinalizar(f: Fila): string | null {
  if ((f.estadoBalanza || '').toUpperCase().trim() !== 'EN_PROCESO') {
    return 'El pesaje no está EN PROCESO en balanza.';
  }
  if (f.pesoBruto <= 0) return 'Falta el peso bruto (use Capt. peso).';
  if (f.tara <= 0) return 'Falta la tara (vehículo o Capt. tara / tara manual).';
  if (f.pesoNeto <= 0 || f.pesoBruto <= f.tara) {
    return 'El peso neto debe ser mayor que cero: bruto tiene que superar a la tara.';
  }
  return null;
}

export type MontarPesajeOpciones = {
  /** tabla: solo grilla (sistema completo). completo: operación + grilla (báscula). */
  alcance?: 'completo' | 'tabla';
};

export function montarPesaje(host: HTMLElement, montaje?: MontarPesajeOpciones): () => void {
  const soloTabla = montaje?.alcance === 'tabla';
  let timer = 0;
  let taraManual = false;
  let pesajeEdicionId = 0;
  const filasGrilla = new Map<number, Fila>();
  let opciones: {
    vehiculos: Opcion[];
    choferes: Opcion[];
    proveedores: Opcion[];
    productos: Opcion[];
    origenes: Opcion[];
    destinos: Opcion[];
  } = { vehiculos: [], choferes: [], proveedores: [], productos: [], origenes: [], destinos: [] };

  host.innerHTML = `
    <div class="panel-pesaje abm-panel${soloTabla ? ' panel-pesaje--solo-tabla' : ''}">
      <header class="abm-cabecera">
        <h2 class="titulo-tarjeta abm-titulo-modulo">Pesaje</h2>
      </header>
      ${
        soloTabla
          ? ''
          : `
      <div class="panel-pesaje-medio">
        <section class="tarjeta pesaje-registro">
          <h2 class="titulo-tarjeta">Registro de peso</h2>
          <div class="separador-horizontal"></div>
          <div class="display-peso"><span class="valor-peso" id="pesaje-vivo">0 kg</span></div>
          <div class="separador-horizontal"></div>
          <div class="pesaje-acciones-balanza">
            <div class="pesaje-col-conexion" id="pesaje-col-conexion">
              <div class="pesaje-conexion-fila">
                <span class="pesaje-icono-balanza" id="pesaje-icono-balanza" aria-hidden="true"></span>
                <label class="abm-switch pesaje-switch-balanza" id="pesaje-switch-balanza" data-tip="Conectar balanza">
                  <input
                    type="checkbox"
                    class="abm-switch-input"
                    id="sw-balanza"
                    role="switch"
                    aria-label="Conexión balanza"
                    aria-checked="false"
                  />
                  <span class="abm-switch-pista"></span>
                </label>
              </div>
              <span class="etiqueta-mini pesaje-leyenda-conexion">Conexión</span>
              <span class="estado-etiqueta estado-etiqueta--inactivo" id="estado-conexion">DESCONECTADO</span>
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
            <button type="button" class="boton-secundario oculto" id="btn-cancelar-edicion">CANCELAR</button>
            <button type="button" class="boton-primario" id="btn-guardar">REGISTRAR PESAJE</button>
          </div>
        </section>
      </div>
      `}
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
  const aviso = crearDialogoAviso(slotDialogos, 'dlg-aviso', 'compacto');
  const confirmar = crearDialogoConfirmar(slotDialogos, 'dlg-confirmar-pesaje');

  let aplicarEstado: (e: EstadoPane) => void = () => {};

  const mostrarAviso = (titulo: string, mensaje: string) => {
    aviso.mostrar(titulo, mensaje);
  };

  const manejarAccion = (a: AccionPesaje) => {
    if (!soloTabla && a.estadoPesaje) aplicarEstado(a.estadoPesaje);
    if (a.aviso) mostrarAviso(a.aviso.titulo, a.aviso.mensaje);
  };

  let recargarOpciones: () => Promise<void> = async () => {};
  let alCambiarChofer: () => Promise<void> = async () => {};
  let alCambiarVehiculo: () => Promise<void> = async () => {};
  let sincronizarConexionBalanza: (conectado: boolean) => void = () => {};
  let gridDatos: HTMLElement | undefined;
  let dlgRapido: HTMLDialogElement | undefined;
  let formRapido: HTMLFormElement | undefined;

  if (!soloTabla) {
  const wrapRapido = document.createElement('div');
  wrapRapido.innerHTML = htmlDialogoFormGuardar({
    dialogId: 'dlg-rapido',
    formId: 'form-rapido',
    tituloId: 'rapido-titulo',
    cuerpoId: 'rapido-campos',
    mensajeId: 'rapido-msg',
    cancelarId: 'rapido-cancelar',
    clasesCuerpoId: 'abm-dialog-grid',
    scroll: true,
    ancho: 'medio',
  });
  slotDialogos.appendChild(wrapRapido.firstElementChild!);
  dlgRapido = qs<HTMLDialogElement>('#dlg-rapido');
  formRapido = qs<HTMLFormElement>('#form-rapido');
  gridDatos = qs<HTMLElement>('#grid-datos');
  enlazarDialogo(dlgRapido, [qs<HTMLButtonElement>('#rapido-cancelar')]);

  sincronizarConexionBalanza = (conectado: boolean) => {
    const sw = qs<HTMLInputElement>('#sw-balanza');
    const tip = conectado ? 'Desconectar balanza' : 'Conectar balanza';
    sw.checked = conectado;
    sw.setAttribute('aria-checked', conectado ? 'true' : 'false');
    sw.setAttribute('aria-label', conectado ? 'Balanza conectada' : 'Balanza desconectada');
    const label = qs<HTMLElement>('#pesaje-switch-balanza');
    label.dataset.tip = tip;
    qs('#pesaje-col-conexion').classList.toggle('pesaje-col-conexion--conectada', conectado);
    const pill = qs('#estado-conexion');
    pill.textContent = conectado ? 'CONECTADO' : 'DESCONECTADO';
    pill.className = `estado-etiqueta ${conectado ? 'estado-etiqueta--activo' : 'estado-etiqueta--inactivo'}`;
  };

  aplicarEstado = (e: EstadoPane) => {
    qs('#pesaje-vivo').textContent = `${e.pesoVivo.toLocaleString('es-BO')} kg`;
    sincronizarConexionBalanza(e.balanzaConectada);
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
    ${selectHTML('sel-producto', 'Producto', 'producto')}
    ${selectHTML('sel-origen', 'Origen', 'origen')}
    ${selectHTML('sel-destino', 'Destino', 'destino')}
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

  recargarOpciones = async () => {
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

  alCambiarChofer = async () => {
    const id = Number(qs<HTMLSelectElement>('#sel-chofer').value);
    if (!id) {
      qs<HTMLInputElement>('#licencia').value = '';
      return;
    }
    const det = (await PesajeDetalleChofer(id)) as { licencia: string };
    qs<HTMLInputElement>('#licencia').value = det?.licencia ?? '';
  };

  alCambiarVehiculo = async () => {
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

  }

  const columnas = [
    'ID', 'Chofer', 'Placa', 'Lic.', 'Tipo', 'Prov.', 'Prod.', 'Orig.', 'Dest.',
    'Costo', 'Flete', 'Fecha', 'Bruto', 'Tara', 'Neto', 'Estado', 'Balanza', 'Acciones.',
  ];
  const claseTh = (c: string) => {
    if (c === 'Estado' || c === 'Balanza') return 'tabla-th-estado';
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
    filasGrilla.clear();
    for (const f of filas) {
      filasGrilla.set(f.id, f);
      const tr = document.createElement('tr');
      const activo = f.estado.toUpperCase() === 'ACTIVO';
      const balanza = (f.estadoBalanza || '').toUpperCase().trim();
      const enProceso = balanza === 'EN_PROCESO';
      const finalizado = balanza === 'FINALIZADO';
      let acciones = '';
      if (finalizado) {
        acciones += `<button type="button" class="tabla-btn-icono" data-tip="Boleta PDF" data-act="boleta" data-id="${f.id}" aria-label="Boleta PDF">${ICONO_BOLETA_TABLA}</button>`;
      } else if (activo && enProceso) {
        if (!soloTabla) {
          acciones += `<button type="button" class="tabla-btn-icono" data-tip="Editar pesaje" data-act="editar" data-id="${f.id}" aria-label="Editar pesaje">${ICONO_EDITAR_TABLA}</button>`;
        }
        acciones += `<button type="button" class="tabla-btn-icono" data-tip="Finalizar pesaje" data-act="finalizar" data-id="${f.id}" aria-label="Finalizar pesaje">${ICONO_FINALIZAR_TABLA}</button>`;
      }
      if (!finalizado) {
        acciones += htmlSwitchFila(activo);
      }
      const balanzaTxt = aMayusculas(f.estadoBalanza).replace(/_/g, ' ');
      const fechaHora = `${f.fecha} ${f.hora}`.trim();
      tr.dataset.id = String(f.id);
      tr.innerHTML = `
        ${tdCelda(String(f.id))}
        ${tdCelda(aMayusculas(f.nombreChofer))}
        ${tdCelda(aMayusculas(f.placa))}
        ${tdCelda(aMayusculas(f.licencia))}
        ${tdCelda(aMayusculas(f.tipoVehiculo))}
        ${tdCelda(aMayusculas(f.nombreProveedor))}
        ${tdCelda(aMayusculas(f.producto))}
        ${tdCelda(aMayusculas(f.origen))}
        ${tdCelda(aMayusculas(f.destino))}
        ${tdCelda(String(f.costoBs))}
        ${tdCelda(String(f.fleteBs))}
        ${tdCelda(fechaHora)}
        ${tdCelda(String(f.pesoBruto))}
        ${tdCelda(String(f.tara))}
        ${tdCelda(String(f.pesoNeto))}
        <td class="tabla-celda-estado"><span class="estado-etiqueta estado-etiqueta--${activo ? 'activo' : 'inactivo'}">${activo ? 'ACTIVO' : 'INACTIVO'}</span></td>
        <td class="tabla-celda-estado"><span class="estado-etiqueta estado-etiqueta--${enProceso ? 'proceso' : 'finalizado'}">${escapeHtml(balanzaTxt)}</span></td>
        <td class="tabla-celda-acciones">${acciones}</td>
      `;
      const sw = tr.querySelector('.abm-switch-input') as HTMLInputElement | null;
      if (sw) {
        sw.addEventListener('click', (ev) => ev.stopPropagation());
        sw.addEventListener('change', () => {
          void cambiarEstadoGrilla(f, sw);
        });
      }
      if (pesajeEdicionId === f.id) tr.classList.add('seleccionado');
      body.appendChild(tr);
    }
  };

  const cambiarEstadoGrilla = async (fila: Fila, sw: HTMLInputElement) => {
    const activo = sw.checked;
    if (!activo) {
      const ok = await confirmar.preguntar(
        'Inactivar pesaje',
        'Se cambiará el estado a INACTIVO. ¿Continuar?',
        { confirmar: 'INACTIVAR' },
      );
      if (!ok) {
        sincronizarSwitchEstado(sw, true);
        return;
      }
    } else if (fila.estado.toUpperCase() !== 'ACTIVO') {
      const ok = await confirmar.preguntar('Reactivar pesaje', '¿Reactivar este pesaje?', {
        confirmar: 'REACTIVAR',
      });
      if (!ok) {
        sincronizarSwitchEstado(sw, false);
        return;
      }
    }
    const r = (await PesajeCambiarEstado(fila.id, activo)) as AccionPesaje;
    if (!r.ok) {
      sincronizarSwitchEstado(sw, !activo);
      if (r.aviso) mostrarAviso(r.aviso.titulo, r.aviso.mensaje);
      return;
    }
    await recargarGrilla();
  };

  const marcarFilaGrilla = (id: number) => {
    qs('#tabla-body').querySelectorAll('tr').forEach((tr) => {
      tr.classList.toggle('seleccionado', Number((tr as HTMLTableRowElement).dataset.id) === id && id > 0);
    });
  };

  let cancelarEdicion: () => void = () => {};
  let aplicarEdicionEnFormulario: (id: number) => Promise<void> = async () => {};

  if (!soloTabla) {
  const setModoEdicion = (id: number) => {
    pesajeEdicionId = id;
    const btn = qs<HTMLButtonElement>('#btn-guardar');
    btn.textContent = id > 0 ? 'GUARDAR CAMBIOS' : 'REGISTRAR PESAJE';
    qs('#btn-cancelar-edicion')?.classList.toggle('oculto', id <= 0);
  };

  const restablecerFormularioNuevo = async () => {
    const primer = (selId: string) => {
      const el = qs<HTMLSelectElement>(`#${selId}`);
      if (el.options.length > 0) el.selectedIndex = 0;
    };
    primer('sel-vehiculo');
    primer('sel-chofer');
    primer('sel-proveedor');
    primer('sel-producto');
    primer('sel-origen');
    primer('sel-destino');
    qs<HTMLInputElement>('#costo').value = '0';
    qs<HTMLInputElement>('#flete').value = '0';
    taraManual = false;
    qs<HTMLInputElement>('#tara-manual').checked = false;
    qs<HTMLInputElement>('#tara').readOnly = true;
    await alCambiarVehiculo();
    await alCambiarChofer();
  };

  cancelarEdicion = () => {
    setModoEdicion(0);
    marcarFilaGrilla(0);
    void restablecerFormularioNuevo().then(() => {
      void PesajeCancelarEdicion().then((est) => aplicarEstado(est as EstadoPane));
    });
  };

  aplicarEdicionEnFormulario = async (id: number) => {
    let det: {
      id: number;
      vehiculoId: number;
      choferId?: number;
      proveedorId?: number;
      productoId?: number;
      origenId?: number;
      destinoId?: number;
      pesoBruto: number;
      tara: number;
      costoBs: number;
      fleteBsPendiente: number;
    };
    try {
      det = (await PesajeCargarEdicion(id)) as typeof det;
    } catch (e) {
      mostrarAviso('Pesaje', String(e));
      return;
    }
    const sel = (idSel: string, val?: number) => {
      const el = qs<HTMLSelectElement>(`#${idSel}`);
      if (val && val > 0) el.value = String(val);
    };
    sel('sel-vehiculo', det.vehiculoId);
    sel('sel-chofer', det.choferId);
    sel('sel-proveedor', det.proveedorId);
    sel('sel-producto', det.productoId);
    sel('sel-origen', det.origenId);
    sel('sel-destino', det.destinoId);
    qs<HTMLInputElement>('#costo').value = String(det.costoBs ?? 0);
    qs<HTMLInputElement>('#flete').value = String(det.fleteBsPendiente ?? 0);
    qs<HTMLInputElement>('#tara').value = String(det.tara ?? 0);
    taraManual = det.tara > 0;
    qs<HTMLInputElement>('#tara-manual').checked = taraManual;
    qs<HTMLInputElement>('#tara').readOnly = !taraManual;
    await alCambiarVehiculo();
    await alCambiarChofer();
    void PesajeEstado().then((est) => aplicarEstado(est as EstadoPane));
    setModoEdicion(id);
    marcarFilaGrilla(id);
  };

  const iconoBalanza = qs('#pesaje-icono-balanza');
  iconoBalanza.innerHTML = iconoModuloHTML('balanza').replace(/\bwidth="20"/, 'width="18"').replace(/\bheight="20"/, 'height="18"');

  const swBalanza = qs<HTMLInputElement>('#sw-balanza');
  swBalanza.addEventListener('click', (ev) => ev.stopPropagation());
  swBalanza.addEventListener('change', () => {
    const previo = !swBalanza.checked;
    void PesajeAlternarBalanza().then(async (r) => {
      const acc = r as AccionPesaje;
      if (acc.estadoPesaje) aplicarEstado(acc.estadoPesaje);
      else {
        const est = (await PesajeEstado()) as EstadoPane;
        aplicarEstado(est);
      }
      if (acc.aviso) mostrarAviso(acc.aviso.titulo, acc.aviso.mensaje);
      if (!acc.ok && !acc.estadoPesaje) sincronizarConexionBalanza(previo);
    });
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

  qs('#btn-cancelar-edicion').addEventListener('click', () => cancelarEdicion());

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
    const guardar = pesajeEdicionId > 0
      ? PesajeActualizarRegistro({ ...payload, id: pesajeEdicionId })
      : PesajeRegistrar(payload);
    void guardar.then(async (r) => {
      manejarAccion(r as AccionPesaje);
      if ((r as AccionPesaje).ok) {
        setModoEdicion(0);
        await recargarGrilla();
      }
    });
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
    producto: {
      tipo: 'producto',
      titulo: 'Nuevo producto',
      campos: ['Nombre *', 'Descripción'],
      defaults: ['', ''],
    },
    origen: {
      tipo: 'origen',
      titulo: 'Nuevo origen',
      campos: ['Nombre *', 'Descripción'],
      defaults: ['', ''],
    },
    destino: {
      tipo: 'destino',
      titulo: 'Nuevo destino',
      campos: ['Nombre *', 'Descripción'],
      defaults: ['', ''],
    },
  };

  let rapidoActual: RapidoDef | null = null;

  gridDatos!.querySelectorAll('button[data-plus]').forEach((btn) => {
    btn.addEventListener('click', () => {
      const t = btn.getAttribute('data-plus')!;
      rapidoActual = rapidos[t];
      if (!rapidoActual) return;
      qs('#rapido-titulo').textContent = rapidoActual.titulo;
      const cont = qs('#rapido-campos');
      const def = rapidoActual;
      cont.innerHTML = def.campos
        .map((c, i) => {
          const ancho = c.toLowerCase().includes('descripción') ? ' data-ancho-completo="1"' : '';
          return `<label class="pesaje-campo"${ancho}><span class="etiqueta-mini">${c}</span><input class="text-field" data-idx="${i}" value="${escapeAttr(def.defaults[i])}" /></label>`;
        })
        .join('');
      cont.querySelectorAll<HTMLInputElement>('input[data-idx]').forEach((inp) => {
        const label = inp.closest('label')?.querySelector('.etiqueta-mini')?.textContent ?? '';
        if (/\btara\b/i.test(label)) inp.dataset.casing = 'preserve';
        enlazarMayusculasVivo(inp);
      });
      aplicarLayoutDialogoCampos(cont, def.campos.length, dlgRapido!);
      dlgRapido!.showModal();
    });
  });

  formRapido!.addEventListener('submit', (ev) => {
    ev.preventDefault();
    const def = rapidoActual;
    if (!def) return;
    const vals = Array.from(formRapido!.querySelectorAll('input[data-idx]')).map((el) => {
      const inp = el as HTMLInputElement;
      if (inp.dataset.casing === 'preserve') return inp.value.trim();
      const idx = Number(inp.dataset.idx);
      const clave =
        def.tipo === 'vehiculo' && idx === 0
          ? 'placa'
          : def.tipo === 'chofer' && idx === 0
            ? 'nombre'
            : (def.tipo === 'producto' || def.tipo === 'origen' || def.tipo === 'destino') && idx === 0
              ? 'nombre'
              : 'texto';
      return normalizarTextoGuardado(clave, inp.value);
    });
    const fin = async () => {
      dlgRapido!.close();
      await recargarOpciones();
      const e = (await PesajeEstado()) as EstadoPane;
      aplicarEstado(e);
    };
    if (def.tipo === 'vehiculo') {
      void PesajeCrearVehiculo(vals[0], vals[1], parseInt(vals[2], 10) || 0).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (def.tipo === 'chofer') {
      void PesajeCrearChofer(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (def.tipo === 'proveedor') {
      void PesajeCrearProveedor(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (def.tipo === 'producto') {
      void PesajeCrearProducto(vals[0], vals[1]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (def.tipo === 'origen') {
      void PesajeCrearOrigen(vals[0], vals[1]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    } else if (def.tipo === 'destino') {
      void PesajeCrearDestino(vals[0], vals[1]).then(async (r) => {
        if ((r as AccionPesaje).ok) await fin();
        else manejarAccion(r as AccionPesaje);
      });
    }
  });

  enlazarMayusculasVivoEnContenedor(host);
  }

  host.querySelector('#tabla-body')?.addEventListener('click', (ev) => {
    const btn = (ev.target as HTMLElement).closest('button[data-act]') as HTMLButtonElement | null;
    if (!btn) return;
    const id = Number(btn.dataset.id);
    const act = btn.dataset.act;
    if (act === 'editar') {
      void aplicarEdicionEnFormulario(id);
      return;
    }
    if (act === 'finalizar') {
      void (async () => {
        const fila = filasGrilla.get(id);
        const bloqueo = fila ? motivoNoFinalizar(fila) : null;
        if (bloqueo) {
          aviso.mostrar('No se puede finalizar', bloqueo, true);
          return;
        }
        const acepta = await confirmar.preguntar(
          'Finalizar pesaje',
          `¿Cerrar el pesaje #${id} como FINALIZADO?\nBruto: ${fila!.pesoBruto} kg · Tara: ${fila!.tara} kg · Neto: ${fila!.pesoNeto} kg`,
          { confirmar: 'FINALIZAR' },
        );
        if (!acepta) return;
        const r = (await PesajeFinalizar(id)) as AccionPesaje;
        manejarAccion(r);
        if (r.ok) {
          if (!soloTabla && pesajeEdicionId === id) cancelarEdicion();
          await recargarGrilla();
        }
      })();
      return;
    }
    if (act === 'boleta') {
      void PesajeBoleta(id).then((r) => manejarAccion(r as AccionPesaje));
    }
  });

  if (!soloTabla) {
    void PesajeEstado().then((e) => aplicarEstado(e as EstadoPane));
    void recargarOpciones();
    timer = window.setInterval(() => {
      void PesajeEstado().then((e) => aplicarEstado(e as EstadoPane));
    }, 1500);
  }
  void recargarGrilla();

  return () => {
    window.clearInterval(timer);
    if (!soloTabla) PesajeCerrar();
  };
}

function opcionId(root: HTMLElement, selId: string): number | undefined {
  const v = (root.querySelector(`#${selId}`) as HTMLSelectElement)?.value;
  return v ? Number(v) : undefined;
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