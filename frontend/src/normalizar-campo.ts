/** Texto de negocio en MAYÚSCULAS; correo y contraseña sin tocar. */

const CLAVES_PRESERVAR = new Set(['correo', 'clave', 'email', 'contrasena', 'password']);

export type TipoCampoTexto = 'text' | 'password' | 'estado' | 'select';

export function debePreservarCasing(clave: string, tipo?: TipoCampoTexto): boolean {
  if (tipo === 'password') return true;
  const k = clave.trim().toLowerCase();
  return CLAVES_PRESERVAR.has(k);
}

export function normalizarTextoGuardado(clave: string, valor: string, tipo?: TipoCampoTexto): string {
  if (tipo === 'select' || tipo === 'estado') return valor.trim();
  if (debePreservarCasing(clave, tipo)) return valor.trim();
  return valor.trim().toUpperCase();
}

export function aMayusculas(s: string): string {
  return (s || '').toUpperCase();
}

function debeEnlazarMayus(input: HTMLInputElement): boolean {
  if (input.dataset.casing === 'preserve') return false;
  if (input.dataset.mayusVivo === '1') return false;
  if (input.classList.contains('password-field')) return false;
  const t = (input.type || 'text').toLowerCase();
  if (t === 'password' || t === 'email') return false;
  return input.classList.contains('text-field');
}

/** Convierte a mayúsculas mientras escribes (el valor en el input queda en MAYÚSCULAS). */
export function enlazarMayusculasVivo(input: HTMLInputElement): void {
  if (!debeEnlazarMayus(input)) return;
  input.dataset.mayusVivo = '1';
  input.addEventListener('input', () => {
    const start = input.selectionStart;
    const end = input.selectionEnd;
    const v = input.value.toUpperCase();
    if (input.value === v) return;
    input.value = v;
    if (start !== null && end !== null) {
      try {
        input.setSelectionRange(start, end);
      } catch {
        /* input types without selection */
      }
    }
  });
}

export function enlazarMayusculasVivoEnContenedor(root: ParentNode): void {
  root.querySelectorAll<HTMLInputElement>('input.text-field').forEach(enlazarMayusculasVivo);
}