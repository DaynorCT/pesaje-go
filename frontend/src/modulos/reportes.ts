import { AbrirArchivo, GenerarBoletaPdf, GenerarReportePdf, ReportesListar } from '../../wailsjs/go/main/App';

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

export function montarReportes(host: HTMLElement): () => void {
  let seleccion = 0;

  host.innerHTML = `
    <div class="reportes-panel">
      <motionwrap class="reportes-barra">
        <input class="text-field" id="rep-desde" placeholder="Desde (yyyy-MM-dd)" />
        <input class="text-field" id="rep-hasta" placeholder="Hasta (yyyy-MM-dd)" />
        <input class="text-field" id="rep-placa" placeholder="Placa" />
        <button type="button" class="boton-secundario" id="rep-actualizar">Actualizar</button>
        <button type="button" class="boton-primario" id="rep-pdf">Reporte PDF</button>
        <button type="button" class="boton-primario" id="rep-boleta">Boleta del pesaje</button>
      </motionwrap>
      <p id="rep-mensaje" class="subtitulo"></p>
      <motionwrap class="tarjeta tarjeta-sin-padding tabla-modulo-scroll">
        <table class="tabla-modulo">
          <thead><tr>
            <th>Guía</th><th>Vehículo</th><th>Chofer</th><th>Proveedor</th>
            <th>Bruto</th><th>Tara</th><th>Neto</th><th>Estado</th><th>Fecha</th>
          </tr></thead>
          <tbody id="rep-body"></tbody>
        </table>
      </motionwrap>
    </div>
  `.replace(/motionwrap/g, 'div');

  const msg = host.querySelector('#rep-mensaje') as HTMLElement;

  const cargar = async () => {
    const filas = (await ReportesListar(200)) as Fila[];
    const body = host.querySelector('#rep-body') as HTMLTableSectionElement;
    body.innerHTML = '';
    seleccion = 0;
    for (const f of filas) {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td>${f.guia || f.id}</td>
        <td>${f.placa}</td>
        <td>${f.nombreChofer}</td>
        <td>${f.nombreProveedor}</td>
        <td>${f.pesoBruto}</td>
        <td>${f.tara}</td>
        <td>${f.pesoNeto}</td>
        <td>${f.estadoBalanza}</td>
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

  const filtros = () => ({
    desde: (host.querySelector('#rep-desde') as HTMLInputElement).value.trim(),
    hasta: (host.querySelector('#rep-hasta') as HTMLInputElement).value.trim(),
    placa: (host.querySelector('#rep-placa') as HTMLInputElement).value.trim(),
  });

  const abrirPdf = async (res: { ok: boolean; ruta?: string; mensaje: string }) => {
    msg.textContent = res.mensaje;
    if (res.ok && res.ruta) await AbrirArchivo(res.ruta);
  };

  host.querySelector('#rep-actualizar')?.addEventListener('click', () => void cargar());
  host.querySelector('#rep-pdf')?.addEventListener('click', () => {
    void GenerarReportePdf(filtros()).then((r) => abrirPdf(r as { ok: boolean; ruta?: string; mensaje: string }));
  });
  host.querySelector('#rep-boleta')?.addEventListener('click', () => {
    if (!seleccion) {
      msg.textContent = 'Seleccione un pesaje de la lista.';
      return;
    }
    void GenerarBoletaPdf(seleccion).then((r) => abrirPdf(r as { ok: boolean; ruta?: string; mensaje: string }));
  });

  void cargar();
  return () => {};
}