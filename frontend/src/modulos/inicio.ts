import { Chart, registerables } from 'chart.js';
import { InicioResumen } from '../../wailsjs/go/main/App';
import type { ModuloId } from '../navegacion';

Chart.register(...registerables);

type PuntoDia = { fecha: string; cantidad: number };

type Resumen = {
  pesajesTotal: number;
  pesajesHoy: number;
  pesajesEnProceso: number;
  pesajesFinalizadosHoy: number;
  vehiculosActivos: number;
  productosActivos: number;
  choferesActivos: number;
  proveedoresActivos: number;
  ultimos7Dias: PuntoDia[];
};

export type OpcionesInicio = {
  irAModulo?: (id: ModuloId) => void;
};

export function montarInicio(host: HTMLElement, opciones?: OpcionesInicio): () => void {
  host.innerHTML = `
    <div class="abm-panel inicio-abm">
      <header class="abm-cabecera">
        <h2 class="titulo-tarjeta abm-titulo-modulo">Inicio</h2>
      </header>
    <div class="inicio-dashboard">
      <section class="tarjeta inicio-panel">
        <h3 class="titulo-tarjeta">Resumen del día</h3>
        <div class="inicio-kpis" id="inicio-kpis"></div>
        <div class="inicio-atalhos">
          <button type="button" class="boton-primario" data-ir="pesaje">Ir a Pesaje</button>
          <button type="button" class="boton-secundario" data-ir="reportes">Reportes</button>
        </div>
      </section>
      <section class="tarjeta inicio-panel inicio-panel--grafico">
        <h3 class="titulo-tarjeta">Pesajes — últimos 7 días</h3>
        <motionwrap class="inicio-grafico-wrap">
          <canvas id="inicio-chart" aria-label="Gráfico de pesajes por día"></canvas>
        </motionwrap>
      </section>
      <section class="tarjeta inicio-panel">
        <h3 class="titulo-tarjeta">Catálogos activos</h3>
        <dl class="inicio-grid inicio-grid--catalogos" id="inicio-catalogos"></dl>
      </section>
    </div>
    </div>
  `.replace(/motionwrap/g, 'div');

  let chart: Chart | undefined;
  const ir = opciones?.irAModulo;

  host.querySelectorAll<HTMLButtonElement>('[data-ir]').forEach((btn) => {
    btn.addEventListener('click', () => {
      const id = btn.dataset.ir as ModuloId | undefined;
      if (id && ir) ir(id);
    });
  });

  const pintarKpis = (r: Resumen) => {
    const kpis = [
      { label: 'Pesajes hoy', valor: r.pesajesHoy, destacado: true },
      { label: 'En proceso', valor: r.pesajesEnProceso },
      { label: 'Finalizados hoy', valor: r.pesajesFinalizadosHoy },
      { label: 'Total histórico', valor: r.pesajesTotal },
    ];
    const el = host.querySelector('#inicio-kpis') as HTMLElement;
    el.innerHTML = kpis
      .map(
        (k) =>
          `<p class="inicio-kpi${k.destacado ? ' inicio-kpi--destacado' : ''}"><span class="etiqueta-mini">${k.label}</span><strong class="inicio-kpi-valor">${k.valor}</strong></p>`,
      )
      .join('');
  };

  const pintarCatalogos = (r: Resumen) => {
    const items = [
      ['Vehículos', r.vehiculosActivos],
      ['Productos', r.productosActivos],
      ['Choferes', r.choferesActivos],
      ['Proveedores', r.proveedoresActivos],
    ] as const;
    const grid = host.querySelector('#inicio-catalogos') as HTMLElement;
    grid.innerHTML = items
      .map(([label, n]) => `<div><dt class="subtitulo">${label}</dt><dd class="inicio-valor">${n}</dd></div>`)
      .join('');
  };

  const pintarGrafico = (puntos: PuntoDia[]) => {
    const canvas = host.querySelector('#inicio-chart') as HTMLCanvasElement;
    if (!canvas) return;
    const labels = puntos.map((p) => p.fecha);
    const data = puntos.map((p) => p.cantidad);
    const primario = getComputedStyle(document.documentElement).getPropertyValue('--pesaje-primario').trim() || '#255ea6';
    const suave = getComputedStyle(document.documentElement).getPropertyValue('--pesaje-primario-suave').trim() || '#e8f0fa';

    chart?.destroy();
    chart = new Chart(canvas, {
      type: 'bar',
      data: {
        labels,
        datasets: [
          {
            label: 'Pesajes',
            data,
            backgroundColor: suave,
            borderColor: primario,
            borderWidth: 2,
            borderRadius: 4,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: { precision: 0 },
          },
        },
      },
    });
    requestAnimationFrame(() => chart?.resize());
  };

  const cargar = async () => {
    const r = (await InicioResumen()) as Resumen;
    pintarKpis(r);
    pintarCatalogos(r);
    pintarGrafico(r.ultimos7Dias ?? []);
  };

  void cargar();
  const timer = window.setInterval(() => void cargar(), 60_000);

  return () => {
    window.clearInterval(timer);
    chart?.destroy();
    chart = undefined;
  };
}