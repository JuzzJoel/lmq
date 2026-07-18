<script lang="ts">
  import { onMount } from 'svelte';
  import { Chart as ChartJS, type ChartConfiguration } from 'chart.js/auto';

  // Inject global Bauhaus design tokens
  ChartJS.defaults.color = '#000000';
  ChartJS.defaults.font.family = "'Inter', system-ui, sans-serif";
  ChartJS.defaults.font.weight = 'bold';
  ChartJS.defaults.scale.grid.color = '#000000';
  ChartJS.defaults.scale.grid.lineWidth = 2;
  ChartJS.defaults.scale.border.color = '#000000';
  ChartJS.defaults.scale.border.width = 2;
  ChartJS.defaults.plugins.tooltip.backgroundColor = '#FFFFFF';
  ChartJS.defaults.plugins.tooltip.titleColor = '#000000';
  ChartJS.defaults.plugins.tooltip.bodyColor = '#000000';
  ChartJS.defaults.plugins.tooltip.borderColor = '#000000';
  ChartJS.defaults.plugins.tooltip.borderWidth = 2;
  ChartJS.defaults.plugins.tooltip.cornerRadius = 0;
  ChartJS.defaults.elements.line.tension = 0; // Remove curves
  ChartJS.defaults.elements.line.borderWidth = 4;
  ChartJS.defaults.elements.line.borderColor = '#000000';
  ChartJS.defaults.elements.point.radius = 0;
  ChartJS.defaults.elements.point.hoverRadius = 0;
  ChartJS.defaults.elements.bar.borderWidth = 2;
  ChartJS.defaults.elements.bar.borderColor = '#000000';

  interface Props {
    type: 'line' | 'bar' | 'doughnut' | 'pie';
    labels: string[];
    datasets: Array<{
      label: string;
      data: number[];
      backgroundColor?: string | string[];
      borderColor?: string;
      fill?: boolean;
      tension?: number;
      borderWidth?: number;
    }>;
    height?: string;
  }

  const { type, labels, datasets, height = '300px' }: Props = $props();

  let canvas: HTMLCanvasElement;
  let chart: ChartJS | undefined = $state(undefined);

  onMount(() => {
    const config: ChartConfiguration = {
      type,
      data: { labels: [...labels], datasets: datasets.map(ds => ({ ...ds, tension: 0 })) },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: type === 'doughnut' || type === 'pie',
            position: 'bottom',
            labels: { color: '#000000', padding: 16, font: { weight: 'bold' } }
          }
        },
        scales: type !== 'doughnut' && type !== 'pie' ? {
          x: {
            ticks: { color: '#000000', font: { weight: 'bold' } }
          },
          y: {
            ticks: { color: '#000000', font: { weight: 'bold' } },
            beginAtZero: true
          }
        } : undefined
      }
    };

    chart = new ChartJS(canvas, config);

    return () => {
      chart?.destroy();
      chart = undefined;
    };
  });

  $effect(() => {
    if (chart) {
      chart.data.labels = [...labels];
      chart.data.datasets = datasets.map(ds => ({ ...ds, tension: 0 }));
      chart.update('none');
    }
  });
</script>

<div class="relative w-full border-2 border-black bg-white p-4 shadow-hard" style="height: {height};">
  <canvas bind:this={canvas}></canvas>
</div>
