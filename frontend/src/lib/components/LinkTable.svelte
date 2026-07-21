<script lang="ts">
  import type { Link } from '$lib/types';
  import QRCode from 'qrcode';

  interface Props {
    links: Link[];
    totalLinks: number;
    currentPage: number;
    searchQuery: string;
    onLoadData: (page: number, search: string) => void;
  }

  const { links, totalLinks, currentPage, searchQuery, onLoadData }: Props = $props();

  let searchInput = $state(searchQuery);
  let tagFilterInput = $state('');
  let activeQR = $state<number | null>(null);

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short', day: 'numeric', year: 'numeric'
    });
  }

  function formatNumber(n: number): string {
    if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M';
    if (n >= 1000) return (n / 1000).toFixed(1) + 'K';
    return n.toString();
  }

  function handleSearch(e: Event) {
    e.preventDefault();
    onLoadData(1, searchInput);
  }

  async function downloadQR(url: string) {
    try {
      const svg = await QRCode.toString(url, {
        type: 'svg',
        margin: 1,
        color: { dark: '#000000', light: '#ffffff' }
      });
      const blob = new Blob([svg], { type: 'image/svg+xml' });
      const blobUrl = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = 'lmq-qr-code.svg';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(blobUrl);
    } catch (err) {
      console.error("Failed to generate QR Code SVG", err);
    }
  }

  let copiedToken = $state('');
  async function copyToClipboard(url: string, token: string) {
    try {
      await navigator.clipboard.writeText(url);
      copiedToken = token;
      setTimeout(() => { if (copiedToken === token) copiedToken = ''; }, 300);
    } catch (err) {}
  }

  let limit = 10;
  let maxPages = $derived(Math.ceil(totalLinks / limit) || 1);
</script>

<div class="border-2 border-black bg-white shadow-hard w-full">
  <!-- Top Controls -->
  <div class="flex flex-col sm:flex-row justify-between items-center p-4 border-b-2 border-black bg-gray-50">
      <form onsubmit={handleSearch} class="flex w-full sm:w-1/2 mb-4 sm:mb-0 shadow-[2px_2px_0px_0px_#000]">
          <input 
              type="text" 
              bind:value={searchInput} 
              placeholder="SEARCH ALIAS OR URL..." 
              class="w-full px-3 py-2 border-2 border-black border-r-0 font-mono font-bold text-sm outline-none focus:bg-warning rounded-none"
          />
          <button type="submit" class="bg-black text-white px-4 py-2 font-bold uppercase border-2 border-black hover:bg-warning hover:text-black transition-none">
              GO
          </button>
      </form>

      <div class="flex gap-4">
          <button 
              disabled={currentPage <= 1}
              onclick={() => onLoadData(currentPage - 1, searchInput)}
              class="px-3 py-2 border-2 border-black font-bold uppercase shadow-[2px_2px_0px_0px_#000] hover:bg-warning active:translate-x-1 active:translate-y-1 active:shadow-none disabled:opacity-50 transition-none"
          >
              &lt; PREV
          </button>
          <div class="px-4 py-2 border-2 border-black font-mono font-bold text-sm bg-white text-center min-w-[80px]">
              {currentPage} / {maxPages}
          </div>
          <button 
              disabled={currentPage >= maxPages}
              onclick={() => onLoadData(currentPage + 1, searchInput)}
              class="px-3 py-2 border-2 border-black font-bold uppercase shadow-[2px_2px_0px_0px_#000] hover:bg-warning active:translate-x-1 active:translate-y-1 active:shadow-none disabled:opacity-50 transition-none"
          >
              NEXT &gt;
          </button>
      </div>
  </div>

  <div class="overflow-x-auto">
  <table class="w-full text-sm table-fixed">
    <thead>
      <tr class="bg-bg-body border-b-2 border-black uppercase">
        <th class="text-left px-4 py-3 text-black font-bold tracking-wider w-[12%]">Token</th>
        <th class="text-left px-4 py-3 text-black font-bold tracking-wider w-[22%]">Short Link</th>
        <th class="text-left px-4 py-3 text-black font-bold tracking-wider w-[30%]">Destination</th>
        <th class="text-center px-4 py-3 text-black font-bold tracking-wider w-[10%]">Clicks</th>
        <th class="text-right px-4 py-3 text-black font-bold tracking-wider w-[12%]">Created</th>
        <th class="text-center px-4 py-3 text-black font-bold tracking-wider w-[14%]">Actions</th>
      </tr>
    </thead>
    <tbody>
      {#if links.length === 0}
        <tr>
          <td colspan="6" class="px-4 py-12 text-center text-black font-bold border-b border-black">
            <p class="text-lg mb-1 uppercase">No links found</p>
          </td>
        </tr>
      {:else}
        {#each links as link (link.id)}
          <tr class="border-b border-black hover:bg-warning transition-none">
            <td class="px-4 py-3 border-r border-black relative overflow-hidden">
              <a
                href="/dashboard/{link.token}"
                class="font-mono text-black hover:bg-black hover:text-white px-1 font-bold uppercase transition-none text-xs"
              >
                /{link.token}
              </a>
              {#if link.has_password}
                 <span title="Password Protected" class="ml-1 text-[10px]">🔒</span>
              {/if}
              {#if link.routes && link.routes.length > 0}
                 <span title="A/B Testing: {link.routes.length} routes" class="ml-1 text-[10px]">🔀</span>
              {/if}
              {#if link.burn_after_reading}
                 <span title="Burn after reading" class="ml-1 text-[10px]">☠</span>
              {/if}
              {#if link.tags && link.tags.length > 0}
                <div class="flex gap-1 mt-1 flex-wrap">
                  {#each link.tags as tag}
                    <span class="text-[8px] bg-blue-100 border border-black px-1 font-bold uppercase leading-tight">{tag}</span>
                  {/each}
                </div>
              {/if}
            </td>
            <td class="px-4 py-3 truncate border-r border-black">
              <a href={link.short_url || (window.location.origin + '/' + link.token)} target="_blank" class="font-mono text-[11px] text-black hover:text-accent underline font-bold block truncate">{link.short_url || (window.location.origin + '/' + link.token)}</a>
              <button type="button" onclick={() => copyToClipboard(link.short_url || (window.location.origin + '/' + link.token), link.token)} class="mt-1 text-[9px] border border-black px-1 hover:bg-warning uppercase font-bold" title="Copy short link">📋 COPY</button>
            </td>
            <td class="px-4 py-3 truncate text-black font-mono text-[11px] border-r border-black font-bold" title={link.long_url}>
              {link.long_url}
            </td>
            <td class="px-4 py-3 text-center border-r border-black">
              <span class="inline-block px-2 py-1 border-2 border-black text-[11px] font-bold bg-accent text-white uppercase leading-tight">
                {formatNumber(link.click_count)}
              </span>
            </td>
            <td class="px-4 py-3 text-right text-black font-bold uppercase text-[11px] border-r border-black truncate">
              {formatDate(link.created_at)}
            </td>
            <td class="px-4 py-3 text-center">
                <div class="flex flex-col gap-1 justify-center items-center">
                    <button type="button" onclick={() => copyToClipboard(link.short_url || (window.location.origin + '/' + link.token), link.token)} 
                            class="w-full border-2 border-black rounded-none px-2 py-1 font-mono text-[10px] font-bold shadow-[2px_2px_0px_0px_#000] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none transition-all uppercase text-center leading-tight"
                            class:bg-warning={copiedToken === link.token}
                            class:bg-white={copiedToken !== link.token}
                            class:text-black={true}>
                        {copiedToken === link.token ? 'COPIED!' : 'COPY'}
                    </button>
                    <button type="button" onclick={() => downloadQR(link.short_url || (window.location.origin + '/' + link.token))} 
                            class="w-full bg-white text-black border-2 border-black rounded-none px-2 py-1 font-mono text-[10px] font-bold shadow-[2px_2px_0px_0px_#000] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none transition-all hover:bg-black hover:text-white uppercase text-center leading-tight">
                        QR
                    </button>
                </div>
            </td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
  </div>
</div>
