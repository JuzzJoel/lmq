<script lang="ts">
  import { onMount } from 'svelte';
  import Seo from '$lib/components/Seo.svelte';
  import StatCard from '$lib/components/StatCard.svelte';
  import Chart from '$lib/components/Chart.svelte';
  import LinkTable from '$lib/components/LinkTable.svelte';
  import { getLinks } from '$lib/api';
  import type { Link } from '$lib/types';

  let links: Link[] = $state([]);
  let totalLinks = $state(0);
  let loading = $state(true);
  let error: string | null = $state(null);

  let isAuthenticated = $state(false);
  let adminTokenInput = $state('');
  let showPassword = $state(false);

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
         isAuthenticated = false;
         sessionStorage.removeItem('admin_token');
      }
    } else {
      links = res.data.links || [];
      totalLinks = res.data.total || 0;
      currentPage = page;
      searchQuery = search;
    }
    loading = false;
  }

  function handleLogin(e: Event) {
    e.preventDefault();
    if (adminTokenInput.trim().length > 0) {
      sessionStorage.setItem('admin_token', adminTokenInput.trim());
      isAuthenticated = true;
      loadData(1, '');
    }
  }

  onMount(() => {
    const token = sessionStorage.getItem('admin_token');
    if (token) {
      isAuthenticated = true;
      loadData(1, '');
    } else {
      loading = false;
    }
  });
</script>

<Seo title="Dashboard" description="Overview of your links and analytics" />

{#if !isAuthenticated}
<div class="min-h-[60vh] flex items-center justify-center">
  <form onsubmit={handleLogin} class="bg-white border-4 border-black shadow-[4px_4px_0px_0px_#000] p-8 w-full max-w-lg rounded-none">
    <h2 class="text-2xl font-bold uppercase tracking-widest mb-6 border-b-4 border-black pb-4 text-center">SYSTEM LOGIN</h2>
    <div class="mb-6">
      <label for="admin_token" class="block font-bold uppercase tracking-wider mb-2">ENTER 64-CHARACTER ADMINISTRATIVE SYSTEM KEY</label>
      <div class="flex shadow-[4px_4px_0px_0px_#000]">
        {#if showPassword}
          <input
            id="admin_token"
            type="text"
            bind:value={adminTokenInput}
            class="w-full px-4 py-3 bg-white border-4 border-black border-r-0 focus:outline-none focus:bg-warning rounded-none font-mono font-bold"
            placeholder="••••••••••••••••"
            required
          />
        {:else}
          <input
            id="admin_token"
            type="password"
            bind:value={adminTokenInput}
            class="w-full px-4 py-3 bg-white border-4 border-black border-r-0 focus:outline-none focus:bg-warning rounded-none font-mono font-bold"
            placeholder="••••••••••••••••"
            required
          />
        {/if}
        <button 
          type="button" 
          onclick={() => showPassword = !showPassword} 
          class="px-6 bg-white border-4 border-black text-black font-bold uppercase tracking-wider hover:bg-black hover:text-white transition-none whitespace-nowrap"
        >
          {showPassword ? '🙈 HIDE' : '👁 SHOW'}
        </button>
      </div>
    </div>
    <button type="submit" class="w-full py-4 bg-black text-white hover:bg-white hover:text-black border-4 border-black uppercase font-bold tracking-widest transition-none">
      AUTHENTICATE
    </button>
  </form>
</div>
{:else}
  <div class="mb-8 border-b-4 border-black pb-4">
    <h1 class="text-4xl font-bold uppercase tracking-tighter text-black">DASHBOARD</h1>
  </div>

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
{/if}
