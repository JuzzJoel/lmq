<script lang="ts">
  import { onMount } from 'svelte';
  import Seo from '$lib/components/Seo.svelte';
  import StatCard from '$lib/components/StatCard.svelte';
  import Chart from '$lib/components/Chart.svelte';
  import LinkTable from '$lib/components/LinkTable.svelte';
  import { getLinks, getOverview } from '$lib/api';
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

  let overviewTotalClicks = $state(0);
  let overviewActiveToday = $state(0);
  let overviewTopCountry = $state('XX');
  let chartLabels: string[] = $state([]);
  let chartData: number[] = $state([]);

  async function loadData(page: number = 1, search: string = '') {
    loading = true;
    error = null;

    const [linksRes, overviewRes] = await Promise.all([
      getLinks(page, 10, search),
      getOverview()
    ]);

    if (linksRes.error) {
      error = linksRes.error;
      if (linksRes.error.toLowerCase().includes("unauthorized")) {
         logout();
      }
    } else {
      links = linksRes.data.links || [];
      totalLinks = linksRes.data.total || 0;
      isMockMode = linksRes.mock || false;
      currentPage = page;
      searchQuery = search;
    }

    if (!overviewRes.error) {
      const ov = overviewRes.data;
      overviewTotalClicks = ov.total_clicks || 0;
      overviewActiveToday = ov.active_today || 0;
      overviewTopCountry = ov.top_country || 'XX';

      // Build chart data: last 7 days with zero-fill for days with no clicks
      const dayMap: Record<string, number> = {};
      if (ov.clicks_by_day) {
        for (const d of ov.clicks_by_day) {
          dayMap[d.date] = d.count;
        }
      }
      const labels: string[] = [];
      const data: number[] = [];
      const now = new Date();
      for (let i = 6; i >= 0; i--) {
        const d = new Date(now);
        d.setDate(d.getDate() - i);
        const key = d.toISOString().slice(0, 10);
        labels.push(d.toLocaleDateString('en-US', { weekday: 'short' }).toUpperCase());
        data.push(dayMap[key] || 0);
      }
      chartLabels = labels;
      chartData = data;
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
      <div class="border-2 border-black bg-white p-6 shadow-hard h-32 skeleton-pulse"></div>
    {/each}
  </div>
  <div class="border-4 border-black bg-white p-6 shadow-hard-lg mb-12">
    <div class="h-6 w-48 bg-gray-200 border-2 border-black mb-6 skeleton-pulse"></div>
    <div class="h-64 w-full bg-gray-200 border-2 border-black skeleton-pulse"></div>
  </div>
  <div class="border-2 border-black bg-white shadow-hard">
    <div class="h-12 bg-gray-100 border-b-2 border-black skeleton-pulse"></div>
    {#each Array(3) as _}
      <div class="h-16 border-b border-black flex items-center px-4 gap-4">
        <div class="h-6 w-16 bg-gray-200 border border-black skeleton-pulse"></div>
        <div class="h-6 w-48 bg-gray-200 border border-black skeleton-pulse"></div>
        <div class="h-6 w-64 bg-gray-200 border border-black skeleton-pulse"></div>
        <div class="h-6 w-12 bg-gray-200 border border-black skeleton-pulse ml-auto"></div>
      </div>
    {/each}
  </div>
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
    <StatCard label="PAGE CLICKS" value={overviewTotalClicks} icon="🖱️" delay={100} />
    <StatCard label="TOP COUNTRY" value={overviewTopCountry} icon="🌍" delay={200} />
    <StatCard label="ACTIVE TODAY" value={overviewActiveToday} icon="🔥" delay={300} />
  </div>

  <div class="bg-white border-4 border-black p-6 mb-12 shadow-hard-lg">
    <h2 class="text-2xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">CLICK ACTIVITY (7 DAYS)</h2>
    <Chart type="line" labels={chartLabels} datasets={[{
      label: 'CLICKS',
      data: chartData,
      borderColor: '#EAB308',
      backgroundColor: 'rgba(234, 179, 8, 0.15)',
      fill: true,
    }]} height="300px" />
    {#if chartData.length > 0 && chartData.reduce((a, b) => a + b, 0) === 0}
      <p class="text-center text-black font-mono font-bold text-sm mt-4 border-t-4 border-black pt-4 uppercase tracking-wider">No click activity in the last 7 days — shorten a link and visit it to see data here.</p>
    {/if}
  </div>

  <div class="mt-12">
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
