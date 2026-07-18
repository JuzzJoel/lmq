<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    label: string;
    value: number;
    icon?: string;
    trend?: string;
    delay?: number;
  }

  const { label, value, icon = '📊', trend = '', delay = 0 }: Props = $props();

  let displayValue = $state(0);
  let visible = $state(false);

  onMount(() => {
    setTimeout(() => {
      visible = true;
      // Animate count up
      const duration = 1000;
      const start = performance.now();
      function tick(now: number) {
        const elapsed = now - start;
        const progress = Math.min(elapsed / duration, 1);
        const eased = 1 - Math.pow(1 - progress, 3); // ease-out cubic
        displayValue = Math.floor(eased * value);
        if (progress < 1) requestAnimationFrame(tick);
      }
      requestAnimationFrame(tick);
    }, delay);
  });
</script>

<div
  class="relative border-2 border-black bg-white p-6 shadow-hard transition-none hover:bg-warning"
  class:opacity-0={!visible}
  class:animate-slide-up={visible}
>
  <div class="relative z-10">
    <div class="flex items-center justify-between mb-3 border-b-2 border-black pb-2">
      <span class="text-2xl">{icon}</span>
      {#if trend}
        <span class="text-xs font-bold px-2 py-1 border-2 border-black bg-success text-white uppercase">{trend}</span>
      {/if}
    </div>
    <p class="text-4xl font-bold tracking-tight text-black font-mono">
      {displayValue.toLocaleString()}
    </p>
    <p class="text-sm text-black mt-2 uppercase tracking-widest font-bold">{label}</p>
  </div>
</div>
