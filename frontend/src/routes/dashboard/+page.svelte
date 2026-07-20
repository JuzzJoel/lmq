<script lang="ts">
  import { onMount } from 'svelte';
  import Seo from '$lib/components/Seo.svelte';
  import StatCard from '$lib/components/StatCard.svelte';
  import Chart from '$lib/components/Chart.svelte';
  import LinkTable from '$lib/components/LinkTable.svelte';
  import { getLinks } from '$lib/api';
  import type { Link } from '$lib/types';
  import { logout } from '$lib/auth.svelte';

  let links: Link[] = $state([]);
  let totalLinks = $state(0);
  let loading = $state(true);
  let error: string | null = $state(null);
  let isMockMode = $state(false);

  let totalClicks = $derived(links.reduce((sum, link) => sum + link.click_count, 0));

  let currentPage = $state(1);
  let searchQuery = $state('');

  const chartLabels = ['MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT', 'SUN'];
  const chartDatasets = [
    {
      label: 'CLICKS (7 DAYS)',
      data: [120, 190, 150, 220, 180, 250, 210],
      borderColor: '#0055FF',
      backgroundColor: 'rgba(0, 85, 255, 0.1)',
      fill: true,
      tension: 0
    }
  ];

  async function loadData(page: number = 1, search: string = '') {
    loading = true;
    error = null;
    const res = await getLinks(page, 10, search);
    if (res.error) {
      error = res.error;
      if (res.error.toLowerCase().includes("unauthorized")) {
         logout();
      }
    } else {
      links = res.data.links || [];
      totalLinks = res.data.total || 0;
      isMockMode = res.mock || false;
      currentPage = page;
      searchQuery = search;
    }
    loading = false;
  }

  onMount(() => {
    loadData(1, '');
  });
</script>

<Seo title="Dashboard" description="Overview of your links and analytics" />

<div class="mb-8 border-b-4 border-black pb-4">
  <h1 class="text-4xl font-bold uppercase tracking-tighter text-black">OVERVIEW</h1>
</div>

{#if isMockMode}
  <div class="bg-red-500 border-4 border-black text-white p-4 mb-6 shadow-hard font-bold uppercase text-center">
    ⚠ DEMO / MOCK MODE — Data shown is not real. No database connection available.
  </div>
{/if}

{#if loading && links.length === 0}
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
    {#each Array(4) as _}
      <div class="border-2 border-black bg-white p-6 shadow-hard h-32 animate-pulse"></div>
    {/each}
  </div>
  <div class="border-2 border-black bg-white p-6 shadow-hard h-80 mb-8 animate-pulse"></div>
  <div class="border-2 border-black bg-white p-6 shadow-hard h-96 animate-pulse"></div>
{:else if error && links.length === 0}
  <div class="bg-danger border-4 border-black text-white p-8 text-center shadow-hard">
    <p class="font-bold uppercase tracking-wider mb-4 text-xl">{error}</p>
    <button onclick={() => loadData(1, '')} class="px-6 py-3 bg-black hover:bg-white hover:text-black border-2 border-white hover:border-black text-white font-bold uppercase transition-none">
      RETRY
    </button>
  </div>
{:else}
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
    <StatCard label="TOTAL LINKS" value={totalLinks} icon="🔗" delay={0} />
    <StatCard label="PAGE CLICKS" value={totalClicks} icon="🖱️" delay={100} />
    <StatCard label="TOP COUNTRY" value={0} icon="🌍" delay={200} />
    <StatCard label="ACTIVE TODAY" value={0} icon="🔥" delay={300} />
  </div>

  <div class="bg-white border-4 border-black p-6 mb-12 shadow-hard-lg">
    <h2 class="text-2xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">CLICK ACTIVITY</h2>
    <Chart type="line" labels={chartLabels} datasets={chartDatasets} height="300px" />
  </div>

  <div>
    <h2 class="text-2xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2 inline-block">YOUR LINKS</h2>
    <LinkTable 
      {links} 
      {totalLinks} 
      {currentPage} 
      {searchQuery} 
      onLoadData={(p, s) => loadData(p, s)} 
    />
  </div>
{/if}
