import { InicioResumen, SincronizarAhora } from '../../wailsjs/go/main/App';

type Resumen = {
  pesajesTotal: number;
  pesajesHoy: number;
  vehiculosActivos: number;
  productosActivos: number;
  choferesActivos: number;
  proveedoresActivos: number;
  syncConfigurada: boolean;
  syncEstado: string;
  syncPendientes: number;
  syncErrores: number;
};

export function montarInicio(host: HTMLElement): () => void {
  host.innerHTML = `
    <section class="tarjeta inicio-panel">
      <h2 class="titulo-modulo">Resumen del sistema</h2>
      <dl class="inicio-grid" id="inicio-stats"></dl>
      <h3 class="titulo-tarjeta">Sincronización con la web</h3>
      <p id="sync-estado" class="subtitulo"></p>
      <p id="sync-pend" class="subtitulo"></p>
      <p id="sync-err" class="subtitulo"></p>
      <button type="button" class="boton-secundario" id="btn-sync">Sincronizar ahora</button>
    </section>
  `;

  const stats = [
    ['Pesajes totales', 'pesajesTotal'],
    ['Pesajes hoy', 'pesajesHoy'],
    ['Vehículos activos', 'vehiculosActivos'],
    ['Productos activos', 'productosActivos'],
    ['Choferes activos', 'choferesActivos'],
    ['Proveedores activos', 'proveedoresActivos'],
  ] as const;

  const aplicar = (r: Resumen) => {
    const grid = host.querySelector('#inicio-stats') as HTMLElement;
    grid.innerHTML = stats
      .map(
        ([label, key]) =>
          `<div><dt class="subtitulo">${label}</dt><dd class="titulo-tarjeta inicio-valor">${r[key]}</dd></div>`,
      )
      .join('');
    (host.querySelector('#sync-estado') as HTMLElement).textContent = r.syncEstado;
    (host.querySelector('#sync-pend') as HTMLElement).textContent =
      `Pendientes de enviar: ${r.syncPendientes}`;
    (host.querySelector('#sync-err') as HTMLElement).textContent = `Con error: ${r.syncErrores}`;
  };

  void InicioResumen().then((r) => aplicar(r as Resumen));

  host.querySelector('#btn-sync')?.addEventListener('click', () => {
    void SincronizarAhora().then((a) => {
      const aviso = (a as { aviso?: { mensaje: string } }).aviso;
      if (aviso) {
        (host.querySelector('#sync-estado') as HTMLElement).textContent = aviso.mensaje;
      }
    });
  });

  return () => {};
}