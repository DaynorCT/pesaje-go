/** Shell HTML/CSS unificado para `<dialog>` en toda la app. */

export type AnchoDialogo = 'compacto' | 'medio' | 'amplio';

export function claseAnchoDialogo(ancho: AnchoDialogo): string {
  return `dlg-ancho-${ancho}`;
}

/** Misma lógica que ABM (usuarios, choferes, etc.). */
export function columnasDialogoPorCampos(visibles: number): 2 | 3 {
  return visibles > 7 ? 3 : 2;
}

export function anchoDialogoPorCampos(visibles: number): AnchoDialogo {
  const cols = columnasDialogoPorCampos(visibles);
  if (cols === 3) return 'amplio';
  if (visibles <= 3) return 'compacto';
  return 'medio';
}

export function aplicarLayoutDialogoCampos(
  contenedor: HTMLElement,
  cantidadCampos: number,
  dialog: HTMLDialogElement,
): void {
  const cols = columnasDialogoPorCampos(cantidadCampos);
  contenedor.className = `abm-dialog-grid abm-dialog-grid--${cols}`;
  fijarAnchoDialogo(dialog, anchoDialogoPorCampos(cantidadCampos));
}

export function htmlAccionesDialogo(contenido: string): string {
  return `<div class="dlg-acciones">${contenido}</div>`;
}

export function htmlBotonDialogo(
  tipo: 'button' | 'submit',
  clases: string,
  id: string,
  texto: string,
): string {
  return `<button type="${tipo}" class="${clases}" id="${id}">${texto}</button>`;
}

export type ShellDialogoOpts = {
  id: string;
  ancho?: AnchoDialogo;
  etiqueta?: 'div' | 'form';
  idCuerpo?: string;
  clasesCuerpo?: string;
  attrsCuerpo?: string;
  interior: string;
};

/** Marco `<dialog>` + `.tarjeta.dlg-contenido` con contenido interior libre. */
export function htmlShellDialogo(opts: ShellDialogoOpts): string {
  const tag = opts.etiqueta ?? 'div';
  const ancho = opts.ancho ?? 'medio';
  const clases = ['tarjeta', 'dlg-contenido', opts.clasesCuerpo ?? ''].filter(Boolean).join(' ');
  const idAttr = opts.idCuerpo ? ` id="${opts.idCuerpo}"` : '';
  const extra = opts.attrsCuerpo ? ` ${opts.attrsCuerpo}` : '';
  return `
    <dialog id="${opts.id}" class="dlg-aviso ${claseAnchoDialogo(ancho)}">
      <${tag} class="${clases}"${idAttr}${extra}>
        ${opts.interior}
      </${tag}>
    </dialog>
  `.trim();
}

export function enlazarDialogo(
  dialog: HTMLDialogElement,
  botonesCerrar?: (HTMLElement | null)[],
): void {
  dialog.addEventListener('cancel', (ev) => {
    ev.preventDefault();
    dialog.close();
  });
  for (const btn of botonesCerrar ?? []) {
    btn?.addEventListener('click', () => dialog.close());
  }
}

export type DialogoAviso = {
  dialog: HTMLDialogElement;
  mostrar: (titulo: string, mensaje: string, esError?: boolean) => void;
};

/** Aviso simple: título, mensaje, Aceptar. */
export function crearDialogoAviso(
  root: ParentNode,
  id = 'dlg-aviso',
  ancho: AnchoDialogo = 'compacto',
): DialogoAviso {
  const idTitulo = `${id}-titulo`;
  const idMensaje = `${id}-mensaje`;
  const idOk = `${id}-ok`;
  const wrap = document.createElement('div');
  wrap.innerHTML = htmlShellDialogo({
    id,
    ancho,
    interior: `
      <h3 class="titulo-tarjeta" id="${idTitulo}"></h3>
      <p class="subtitulo dlg-intro" id="${idMensaje}"></p>
      ${htmlAccionesDialogo(htmlBotonDialogo('button', 'boton-primario', idOk, 'ACEPTAR'))}
    `,
  });
  root.appendChild(wrap.firstElementChild!);
  const dialog = root.querySelector(`#${id}`) as HTMLDialogElement;
  const titulo = root.querySelector(`#${idTitulo}`) as HTMLElement;
  const mensaje = root.querySelector(`#${idMensaje}`) as HTMLElement;
  const ok = root.querySelector(`#${idOk}`) as HTMLButtonElement;
  enlazarDialogo(dialog, [ok]);
  return {
    dialog,
    mostrar: (t, m, esError = false) => {
      titulo.textContent = t;
      mensaje.textContent = m;
      mensaje.className = esError ? 'mensaje-error dlg-intro' : 'subtitulo dlg-intro';
      dialog.showModal();
    },
  };
}

export type DialogoFormGuardarOpts = {
  dialogId: string;
  formId: string;
  tituloId: string;
  cuerpoId: string;
  mensajeId: string;
  cancelarId: string;
  ancho?: AnchoDialogo;
  scroll?: boolean;
  textoGuardar?: string;
  tituloHtml?: string;
  introHtml?: string;
  clasesCuerpoId?: string;
  attrsForm?: string;
};

/** Formulario modal: título, zona de campos, mensaje, Cancelar + Guardar. */
export function htmlDialogoFormGuardar(opts: DialogoFormGuardarOpts): string {
  const titulo =
    opts.tituloHtml ?? `<h3 class="titulo-tarjeta" id="${opts.tituloId}"></h3>`;
  const intro = opts.introHtml ?? '';
  const clasesCampos = opts.clasesCuerpoId ? ` class="${opts.clasesCuerpoId}"` : '';
  const scroll = opts.scroll ? ' dlg-contenido--scroll' : '';
  const guardar = opts.textoGuardar ?? 'GUARDAR';
  return htmlShellDialogo({
    id: opts.dialogId,
    ancho: opts.ancho ?? 'medio',
    etiqueta: 'form',
    idCuerpo: opts.formId,
    clasesCuerpo: scroll.trim() || undefined,
    attrsCuerpo: opts.attrsForm,
    interior: `
      ${titulo}
      ${intro}
      <div id="${opts.cuerpoId}"${clasesCampos}></div>
      <p id="${opts.mensajeId}" class="subtitulo dlg-dialog-msg"></p>
      ${htmlAccionesDialogo(
        `${htmlBotonDialogo('button', 'boton-secundario', opts.cancelarId, 'CANCELAR')}${htmlBotonDialogo('submit', 'boton-primario', `${opts.formId}-submit`, guardar)}`,
      )}
    `,
  });
}

export function fijarAnchoDialogo(dialog: HTMLDialogElement, ancho: AnchoDialogo): void {
  dialog.classList.remove('dlg-ancho-compacto', 'dlg-ancho-medio', 'dlg-ancho-amplio');
  dialog.classList.add(claseAnchoDialogo(ancho));
}