export type OpcionCombo = { id: number; label: string };

const opcionesPorCombo = new WeakMap<HTMLElement, OpcionCombo[]>();
let documentoCerrarEnlazado = false;

function normalizarBusqueda(s: string): string {
  return s
    .normalize('NFD')
    .replace(/\p{M}/gu, '')
    .toLowerCase();
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

/** HTML del campo buscable (id va en el input hidden). */
export function htmlComboBuscador(id: string): string {
  return `
    <span class="combo-buscador" data-combo-id="${id}">
      <input type="hidden" id="${id}" value="" />
      <input
        type="text"
        class="combo-buscador-input text-field"
        autocomplete="off"
        spellcheck="false"
        placeholder="Buscar…"
        aria-autocomplete="list"
        aria-expanded="false"
      />
      <div class="combo-buscador-lista" role="listbox" hidden></div>
    </span>`;
}

export function actualizarOpcionesCombo(root: ParentNode, id: string, items: OpcionCombo[]): void {
  const wrap = root.querySelector(`[data-combo-id="${id}"]`) as HTMLElement | null;
  if (!wrap) return;
  opcionesPorCombo.set(wrap, items);
  const hidden = wrap.querySelector(`#${CSS.escape(id)}`) as HTMLInputElement;
  const inp = wrap.querySelector('.combo-buscador-input') as HTMLInputElement;
  const actual = Number(hidden.value);
  if (actual > 0 && !items.some((o) => o.id === actual)) {
    establecerValorCombo(root, id, undefined, false);
  } else if (actual > 0) {
    const lab = items.find((o) => o.id === actual)?.label ?? '';
    if (document.activeElement !== inp) inp.value = lab;
  }
}

export function leerIdCombo(root: ParentNode, id: string): number {
  const el = root.querySelector(`#${CSS.escape(id)}`) as HTMLInputElement | null;
  const v = Number(el?.value);
  return v > 0 ? v : 0;
}

export function establecerValorCombo(
  root: ParentNode,
  id: string,
  valueId?: number,
  dispararChange = true,
): void {
  const wrap = root.querySelector(`[data-combo-id="${id}"]`) as HTMLElement | null;
  if (!wrap) return;
  const hidden = wrap.querySelector(`#${CSS.escape(id)}`) as HTMLInputElement;
  const inp = wrap.querySelector('.combo-buscador-input') as HTMLInputElement;
  const items = opcionesPorCombo.get(wrap) ?? [];
  if (valueId && valueId > 0) {
    const o = items.find((x) => x.id === valueId);
    hidden.value = String(valueId);
    inp.value = o?.label ?? '';
  } else {
    hidden.value = '';
    inp.value = '';
  }
  if (dispararChange) hidden.dispatchEvent(new Event('change', { bubbles: true }));
}

export function seleccionarPrimeraOpcion(root: ParentNode, id: string): void {
  const wrap = root.querySelector(`[data-combo-id="${id}"]`) as HTMLElement | null;
  const items = wrap ? (opcionesPorCombo.get(wrap) ?? []) : [];
  if (items.length > 0) establecerValorCombo(root, id, items[0].id, true);
  else establecerValorCombo(root, id, undefined, true);
}

export function enlazarCombosBuscadores(root: ParentNode): void {
  root.querySelectorAll<HTMLElement>('.combo-buscador').forEach((wrap) => {
    if (wrap.dataset.comboEnlazado === '1') return;
    wrap.dataset.comboEnlazado = '1';
    const id = wrap.dataset.comboId ?? '';
    const hidden = wrap.querySelector(`#${CSS.escape(id)}`) as HTMLInputElement;
    const inp = wrap.querySelector('.combo-buscador-input') as HTMLInputElement;
    const lista = wrap.querySelector('.combo-buscador-lista') as HTMLDivElement;

    const cerrar = () => {
      lista.hidden = true;
      inp.setAttribute('aria-expanded', 'false');
    };

    const render = (filtro: string) => {
      const items = opcionesPorCombo.get(wrap) ?? [];
      const q = normalizarBusqueda(filtro.trim());
      const visibles = q ? items.filter((o) => normalizarBusqueda(o.label).includes(q)) : items;
      const corto = visibles.slice(0, 250);
      if (corto.length === 0) {
        lista.innerHTML = q
          ? '<p class="combo-buscador-vacio subtitulo">Sin coincidencias</p>'
          : '';
        lista.hidden = !q;
        inp.setAttribute('aria-expanded', q ? 'true' : 'false');
        return;
      }
      lista.innerHTML = corto
        .map(
          (o) =>
            `<button type="button" class="combo-buscador-opcion" role="option" data-id="${o.id}">${escapeHtml(o.label)}</button>`,
        )
        .join('');
      lista.hidden = false;
      inp.setAttribute('aria-expanded', 'true');
    };

    const elegir = (optId: number) => {
      establecerValorCombo(root, id, optId, true);
      cerrar();
      inp.blur();
    };

    inp.addEventListener('focus', () => render(inp.value));
    inp.addEventListener('input', () => {
      hidden.value = '';
      render(inp.value);
    });
    lista.addEventListener('mousedown', (ev) => {
      ev.preventDefault();
      const btn = (ev.target as HTMLElement).closest('.combo-buscador-opcion') as HTMLButtonElement | null;
      if (!btn) return;
      elegir(Number(btn.dataset.id));
    });
    inp.addEventListener('keydown', (ev) => {
      if (ev.key === 'Escape') cerrar();
    });
    inp.addEventListener('blur', () => {
      window.setTimeout(() => {
        if (wrap.contains(document.activeElement)) return;
        cerrar();
        const actual = Number(hidden.value);
        const items = opcionesPorCombo.get(wrap) ?? [];
        if (actual > 0) {
          const o = items.find((x) => x.id === actual);
          if (o) inp.value = o.label;
        } else {
          inp.value = '';
        }
      }, 120);
    });
  });

  if (!documentoCerrarEnlazado) {
    documentoCerrarEnlazado = true;
    document.addEventListener('mousedown', (ev) => {
      const t = ev.target as Node;
      document.querySelectorAll<HTMLElement>('.combo-buscador').forEach((wrap) => {
        if (wrap.contains(t)) return;
        const lista = wrap.querySelector('.combo-buscador-lista') as HTMLElement;
        if (!lista.hidden) {
          lista.hidden = true;
          wrap.querySelector('.combo-buscador-input')?.setAttribute('aria-expanded', 'false');
        }
      });
    });
  }
}
