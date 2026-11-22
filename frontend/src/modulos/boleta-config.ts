import { BoletaConfigGuardar, BoletaConfigObtener } from '../../wailsjs/go/main/App';

const CAMPOS: { clave: string; etiqueta: string }[] = [
  { clave: 'salida', etiqueta: 'Salida' },
  { clave: 'direccion', etiqueta: 'Dirección' },
  { clave: 'celular1', etiqueta: 'Celular 1' },
  { clave: 'celular2', etiqueta: 'Celular 2' },
  { clave: 'ciudad', etiqueta: 'Ciudad' },
  { clave: 'tituloSuperior', etiqueta: 'Título superior' },
  { clave: 'marca', etiqueta: 'Marca' },
  { clave: 'tituloDocumento', etiqueta: 'Título documento' },
  { clave: 'acreditacion', etiqueta: 'Acreditación' },
];

export function montarBoletaConfig(host: HTMLElement): () => void {
  let id = 0;
  host.innerHTML = `
    <section class="tarjeta boleta-config">
      <h2 class="titulo-tarjeta">Configuración de boleta</h2>
      <div class="pesaje-grid-datos" id="boleta-campos"></div>
      <button type="button" class="boton-primario" id="boleta-guardar">Guardar</button>
      <p id="boleta-msg" class="subtitulo"></p>
    </section>
  `;

  const cont = host.querySelector('#boleta-campos') as HTMLElement;
  cont.innerHTML = CAMPOS.map(
    (c) =>
      `<label class="pesaje-campo"><span class="etiqueta-mini">${c.etiqueta}</span>
      <input class="text-field" id="bol-${c.clave}" data-k="${c.clave}" /></label>`,
  ).join('');

  void BoletaConfigObtener().then((cfg) => {
    const c = cfg as {
      id?: number;
      salida?: string;
      direccion?: string;
      celular1?: string;
      celular2?: string;
      ciudad?: string;
      tituloSuperior?: string;
      marca?: string;
      tituloDocumento?: string;
      acreditacion?: string;
    };
    id = c.id ?? 0;
    for (const f of CAMPOS) {
      const el = host.querySelector(`#bol-${f.clave}`) as HTMLInputElement;
      el.value = String((c as Record<string, string | undefined>)[f.clave] ?? '');
    }
  });

  host.querySelector('#boleta-guardar')?.addEventListener('click', async () => {
    const dto = {
      id,
      salida: '',
      direccion: '',
      celular1: '',
      celular2: '',
      ciudad: '',
      tituloSuperior: '',
      marca: '',
      tituloDocumento: '',
      acreditacion: '',
    } as Record<string, string | number>;
    for (const f of CAMPOS) {
      dto[f.clave] = (host.querySelector(`#bol-${f.clave}`) as HTMLInputElement).value;
    }
    const res = (await BoletaConfigGuardar(dto as never)) as { ok: boolean; mensaje: string };
    (host.querySelector('#boleta-msg') as HTMLElement).textContent = res.mensaje;
  });

  return () => {};
}