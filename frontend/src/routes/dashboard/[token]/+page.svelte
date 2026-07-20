<script lang="ts">
  import { onMount } from 'svelte';
  import Seo from '$lib/components/Seo.svelte';
  import StatCard from '$lib/components/StatCard.svelte';
  import Chart from '$lib/components/Chart.svelte';
  import { getAnalytics } from '$lib/api';
  import type { LinkAnalytics } from '$lib/types';

  let { data } = $props();
  let token = $derived(data.token);
  let analytics: LinkAnalytics | null = $state(null);
  let loading = $state(true);
  let error: string | null = $state(null);
  let isMockMode = $state(false);

  async function loadData() {
    loading = true;
    error = null;
    const res = await getAnalytics(token);
    if (res.error) {
      error = res.error;
    } else {
      analytics = res.data;
      isMockMode = res.mock || false;
    }
    loading = false;
  }

  onMount(() => {
    loadData();
  });

  let lineChartData = $derived(analytics ? {
    labels: analytics.clicks_by_day.map(d => new Date(d.date).toLocaleDateString(undefined, {month:'short', day:'numeric'})),
    datasets: [{
      label: 'CLICKS',
      data: analytics.clicks_by_day.map(d => d.count),
      borderColor: '#0055FF',
      backgroundColor: 'rgba(0, 85, 255, 0.1)',
      fill: true,
      tension: 0
    }]
  } : null);

  let countryGroupsChartData = $derived(analytics && analytics.country_groups ? {
    labels: analytics.country_groups.slice(0, 5).map(c => c.country),
    datasets: [{
      label: 'CLICKS BY COUNTRY GROUP',
      data: analytics.country_groups.slice(0, 5).map(c => c.count),
      backgroundColor: '#FF2A2A',
    }]
  } : null);

  let regionsChartData = $derived(analytics && analytics.regions ? {
    labels: analytics.regions.slice(0, 5).map(r => r.region),
    datasets: [{
      label: 'CLICKS BY REGION',
      data: analytics.regions.slice(0, 5).map(r => r.count),
      backgroundColor: '#FFD700',
    }]
  } : null);

  let citiesChartData = $derived(analytics && analytics.cities ? {
    labels: analytics.cities.slice(0, 5).map(c => c.city),
    datasets: [{
      label: 'CLICKS BY CITY',
      data: analytics.cities.slice(0, 5).map(c => c.count),
      backgroundColor: '#0055FF',
    }]
  } : null);

  let doughnutChartData = $derived(analytics ? {
    labels: analytics.browsers.map(b => b.browser),
    datasets: [{
      label: 'BROWSERS',
      data: analytics.browsers.map(b => b.count),
      backgroundColor: ['#0055FF', '#FF2A2A', '#FFD700', '#000000', '#FFFFFF'],
    }]
  } : null);
</script>

<Seo title={`Analytics for /${token}`} description="Detailed analytics for your shortened link" />

<div class="mb-8">
  <a href="/dashboard" class="text-black font-bold uppercase border-2 border-black px-4 py-2 hover:bg-black hover:text-white transition-none inline-flex items-center gap-2 shadow-hard">
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="square" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
    BACK TO DASHBOARD
  </a>
</div>

{#if isMockMode}
  <div class="bg-red-500 border-4 border-black text-white p-4 mb-6 shadow-hard font-bold uppercase text-center">
    ⚠ DEMO / MOCK MODE — Analytics shown are fake. No database connection available.
  </div>
{/if}

{#if loading}
  <div class="border-2 border-black bg-white p-6 shadow-hard h-24 mb-8 animate-pulse"></div>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
    {#each Array(3) as _}
      <div class="border-2 border-black bg-white p-6 shadow-hard h-32 animate-pulse"></div>
    {/each}
  </div>
  <div class="border-2 border-black bg-white p-6 shadow-hard h-80 mb-8 animate-pulse"></div>
{:else if error || !analytics}
  <div class="bg-danger border-4 border-black text-white p-8 text-center shadow-hard">
    <p class="font-bold uppercase tracking-wider mb-4 text-xl">{error || 'DATA NOT FOUND'}</p>
    <button onclick={loadData} class="px-6 py-3 bg-black hover:bg-white hover:text-black border-2 border-white hover:border-black text-white font-bold uppercase transition-none">
      RETRY
    </button>
  </div>
{:else}
  <div class="bg-white border-4 border-black p-8 mb-8 shadow-hard flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
    <div>
      <h1 class="text-4xl font-bold flex items-center gap-2 uppercase tracking-tighter">
        /<span class="text-accent">{analytics.token}</span>
      </h1>
      <a href={analytics.long_url} target="_blank" rel="noopener noreferrer" class="text-black font-mono font-bold hover:bg-warning px-1 text-lg flex items-center gap-2 mt-2 truncate max-w-xl transition-none border-b-2 border-black pb-1">
        {analytics.long_url}
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="square" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
      </a>
    </div>
    <div class="text-right border-l-4 border-black pl-6">
      <p class="text-6xl font-bold tracking-tighter">{analytics.total_clicks.toLocaleString()}</p>
      <p class="text-sm text-black font-bold uppercase tracking-widest mt-1">TOTAL CLICKS</p>
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
    <StatCard label="TOTAL CLICKS" value={analytics.total_clicks} icon="🖱️" delay={0} />
    <StatCard label="UNIQUE COUNTRIES" value={analytics.countries ? analytics.countries.length : 0} icon="🌍" delay={100} />
    <StatCard label="TOP BROWSER" value={analytics.browsers && analytics.browsers.length > 0 ? analytics.browsers[0].count : 0} icon="🌐" delay={200} />
  </div>

  {#if lineChartData}
    <div class="bg-white border-4 border-black p-8 mb-12 shadow-hard-lg">
      <h2 class="text-2xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">CLICKS OVER TIME (LAST 30 DAYS)</h2>
      <Chart type="line" labels={lineChartData.labels} datasets={lineChartData.datasets} height="300px" />
    </div>
  {/if}

  <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-8 mb-12">
    {#if countryGroupsChartData}
      <div class="bg-white border-4 border-black p-6 shadow-hard-lg">
        <h2 class="text-xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">COUNTRY GROUPS</h2>
        <Chart type="bar" labels={countryGroupsChartData.labels} datasets={countryGroupsChartData.datasets} height="250px" />
      </div>
    {/if}
    {#if regionsChartData}
      <div class="bg-white border-4 border-black p-6 shadow-hard-lg">
        <h2 class="text-xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">REGIONS</h2>
        <Chart type="bar" labels={regionsChartData.labels} datasets={regionsChartData.datasets} height="250px" />
      </div>
    {/if}
    {#if citiesChartData}
      <div class="bg-white border-4 border-black p-6 shadow-hard-lg">
        <h2 class="text-xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">CITIES</h2>
        <Chart type="bar" labels={citiesChartData.labels} datasets={citiesChartData.datasets} height="250px" />
      </div>
    {/if}
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
    {#if doughnutChartData}
      <div class="bg-white border-4 border-black p-6 shadow-hard-lg">
        <h2 class="text-xl font-bold mb-6 uppercase tracking-wider border-b-4 border-black pb-2">BROWSERS</h2>
        <Chart type="doughnut" labels={doughnutChartData.labels} datasets={doughnutChartData.datasets} height="250px" />
      </div>
    {/if}
  </div>

  <div class="bg-white border-4 border-black shadow-hard-lg mb-12">
    <div class="px-6 py-4 border-b-4 border-black bg-warning">
      <h2 class="text-2xl font-bold uppercase tracking-wider">RECENT CLICKS</h2>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b-4 border-black bg-bg-body">
            <th class="text-left px-6 py-4 text-black font-bold uppercase tracking-wider border-r-2 border-black">TIME</th>
            <th class="text-left px-6 py-4 text-black font-bold uppercase tracking-wider border-r-2 border-black">LOCATION</th>
            <th class="text-left px-6 py-4 text-black font-bold uppercase tracking-wider border-r-2 border-black">BROWSER</th>
            <th class="text-left px-6 py-4 text-black font-bold uppercase tracking-wider border-r-2 border-black">OS</th>
            <th class="text-center px-6 py-4 text-black font-bold uppercase tracking-wider">DEVICE</th>
          </tr>
        </thead>
        <tbody>
          {#if !analytics.recent_clicks || analytics.recent_clicks.length === 0}
            <tr>
              <td colspan="5" class="px-6 py-12 text-center text-black font-bold uppercase text-lg">NO RECENT CLICKS</td>
            </tr>
          {:else}
            {#each analytics.recent_clicks as click}
              <tr class="border-b-2 border-black hover:bg-warning transition-none">
                <td class="px-6 py-4 text-black font-mono font-bold whitespace-nowrap border-r-2 border-black">
                  {new Date(click.clicked_at).toLocaleString()}
                </td>
                <td class="px-6 py-4 text-black font-bold border-r-2 border-black">
                  {click.city ? click.city + ', ' : ''}{click.region ? click.region + ', ' : ''}{click.country || 'Unknown'}
                </td>
                <td class="px-6 py-4 text-black font-bold border-r-2 border-black">
                  {click.browser || 'Unknown'}
                </td>
                <td class="px-6 py-4 text-black font-bold border-r-2 border-black">
                  {click.os || 'Unknown'}
                </td>
                <td class="px-6 py-4 text-center">
                  <span class="inline-block px-3 py-1 border-2 border-black text-xs font-bold uppercase bg-white text-black">
                    {click.is_mobile ? '📱 MOBILE' : '💻 DESKTOP'}
                  </span>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
{/if}
