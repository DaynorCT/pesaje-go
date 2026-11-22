import { BalanzaGuardarConfig, BalanzaObtenerConfig } from '../../wailsjs/go/main/App';

type BalanzaCfg = {
  puertoCom: string;
  baudRate: number;
  dataBits: number;
  paridad: string;
  stopBits: number;
  flowControl: string;
  timeoutMs: number;
  metodoLectura: string;
  posicionInicio: number;
  posicionLongitud: number;
};

export function montarBalanza(host: HTMLElement): () => void {
  host.innerHTML = `
    <section class="tarjeta balanza-config">
      <h2 class="titulo-tarjeta">Configuración de balanza</h2>
      <div class="pesaje-grid-datos" id="balanza-grid"></div>
      <button type="button" class="boton-primario" id="balanza-guardar">Guardar</button>
      <p id="balanza-mensaje" class="subtitulo"></p>
    </section>
  `;

  const campos: { clave: keyof BalanzaCfg; etiqueta: string }[] = [
    { clave: 'puertoCom', etiqueta: 'Puerto COM' },
    { clave: 'baudRate', etiqueta: 'Baudrate' },
    { clave: 'dataBits', etiqueta: 'Data bits' },
    { clave: 'paridad', etiqueta: 'Paridad' },
    { clave: 'stopBits', etiqueta: 'Stop bits' },
    { clave: 'flowControl', etiqueta: 'Flow control' },
    { clave: 'timeoutMs', etiqueta: 'Timeout (ms)' },
    { clave: 'metodoLectura', etiqueta: 'Método lectura' },
    { clave: 'posicionInicio', etiqueta: 'Posición inicio' },
    { clave: 'posicionLongitud', etiqueta: 'Longitud' },
  ];

  const grid = host.querySelector('#balanza-grid') as HTMLElement;
  grid.innerHTML = campos
    .map(
      (c) => `<label class="pesaje-campo"><span class="etiqueta-mini">${c.etiqueta}</span>
      <input class="text-field" id="balanza-${c.clave}" data-clave="${c.clave}" /></label>`,
    )
    .join('');

  const msg = host.querySelector('#balanza-mensaje') as HTMLElement;

  const aplicar = (cfg: BalanzaCfg) => {
    for (const c of campos) {
      const el = host.querySelector(`#balanza-${c.clave}`) as HTMLInputElement;
      const v = cfg[c.clave];
      el.value = String(v ?? '');
    }
  };

  const leer = (): BalanzaCfg => ({
    puertoCom: val('puertoCom'),
    baudRate: num('baudRate'),
    dataBits: num('dataBits'),
    paridad: val('paridad'),
    stopBits: num('stopBits'),
    flowControl: val('flowControl'),
    timeoutMs: num('timeoutMs'),
    metodoLectura: val('metodoLectura') || 'AUTO',
    posicionInicio: num('posicionInicio'),
    posicionLongitud: num('posicionLongitud'),
  });

  const val = (k: string) =>
    (host.querySelector(`#balanza-${k}`) as HTMLInputElement)?.value ?? '';
  const num = (k: string) => parseInt(val(k), 10) || 0;

  void BalanzaObtenerConfig().then((c) => aplicar(c as BalanzaCfg));

  host.querySelector('#balanza-guardar')?.addEventListener('click', async () => {
    const res = (await BalanzaGuardarConfig(leer())) as { ok: boolean; mensaje: string; esError?: boolean };
    msg.textContent = res.mensaje;
    msg.className = res.esError ? 'mensaje-error' : 'subtitulo mensaje-exito';
  });

  return () => {};
}