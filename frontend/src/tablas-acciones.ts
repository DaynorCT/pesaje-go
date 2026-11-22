const SVG_TABLA =
  'xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"';

/** Documento / boleta PDF */
export const ICONO_BOLETA_TABLA = `<svg ${SVG_TABLA}><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6"/><path d="M8 13h8M8 17h5"/></svg>`;

/** Cerrar / finalizar pesaje */
export const ICONO_FINALIZAR_TABLA = `<svg ${SVG_TABLA}><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m22 4-10 10-3.5-3.5"/></svg>`;

/** Editar fila (lápiz) */
export const ICONO_EDITAR_TABLA = `<svg ${SVG_TABLA}><path d="M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/></svg>`;

export function etiquetaSwitchActivo(activo: boolean): string {
  return activo ? 'ACTIVO' : 'INACTIVO';
}

export function sincronizarSwitchEstado(input: HTMLInputElement, activo: boolean): void {
  const label = input.closest('.abm-switch') as HTMLElement | null;
  const txt = etiquetaSwitchActivo(activo);
  input.checked = activo;
  input.setAttribute('aria-label', txt);
  input.setAttribute('aria-checked', activo ? 'true' : 'false');
  if (label) label.dataset.tip = txt;
}

export function htmlSwitchFila(activo: boolean): string {
  const estadoTxt = etiquetaSwitchActivo(activo);
  return `<label class="abm-switch" data-tip="${estadoTxt}">
    <input type="checkbox" class="abm-switch-input" role="switch" ${activo ? 'checked' : ''} aria-label="${estadoTxt}" aria-checked="${activo ? 'true' : 'false'}" />
    <span class="abm-switch-pista"></span>
  </label>`;
}