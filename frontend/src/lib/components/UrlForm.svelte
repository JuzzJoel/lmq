<script lang="ts">
  import { invalidateAll } from '$app/navigation';
  import { API_BASE } from '$lib/api';

  let inputUrlValue = $state('');
  let expiresIn = $state('never');
  let passwordLock = $state('');
  let showPassword = $state(false);
  let showAliasField = $state(false);
  let alias = $state('');
  
  let shortenedLinks = $state([]);
  let errorMessage = $state('');
  let isLoading = $state(false);
  let copyTargetIndex = $state(-1);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!inputUrlValue.trim()) return;

    isLoading = true;
    errorMessage = '';
    shortenedLinks = [];

    try {
      const res = await fetch(`${API_BASE}/shorten`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          url: inputUrlValue,
          expires_in: expiresIn === 'never' ? 0 : parseInt(expiresIn),
          password: passwordLock || null,
          custom_token: showAliasField && alias ? alias : null
        })
      });

      const data = await res.json();

      if (res.ok) {
        if (data && data.results && Array.isArray(data.results)) {
          shortenedLinks = data.results;
        } else if (data && data.short_url) {
          shortenedLinks = [data];
        } else {
          shortenedLinks = [data];
        }
      } else {
        errorMessage = data.error || 'Failed to shorten URLs.';
      }
    } catch (err: any) {
      errorMessage = 'Network or connection error occurred.';
    } finally {
      isLoading = false;
    }
  }

  function handleCopy(text: string, index: number) {
    navigator.clipboard.writeText(text);
    copyTargetIndex = index;
    setTimeout(() => {
      copyTargetIndex = -1;
    }, 300);
  }

  import QRCode from 'qrcode';

  function downloadSVG(svgString: string, filename: string) {
    const blob = new Blob([svgString], { type: 'image/svg+xml' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  async function getQRSvg(url: string): Promise<string> {
    return await QRCode.toString(url, {
      type: 'svg',
      margin: 1,
      width: 250,
      color: { dark: '#000000', light: '#ffffff' }
    });
  }
</script>

<form onsubmit={handleSubmit} class="w-full max-w-2xl mx-auto space-y-4 group">
  <div class="flex flex-col md:flex-row gap-2 transition-transform duration-300">
    <input
      type="text"
      bind:value={inputUrlValue}
      placeholder="Enter URL(s) separated by spaces or commas..."
      class="flex-1 bg-white border-4 border-black rounded-none p-4 text-lg font-mono focus:outline-none shadow-[4px_4px_0px_0px_#000] hover:shadow-[6px_6px_0px_0px_#FF2A2A] hover:-translate-y-1 transition-all duration-200"
      required
    />
    <button
      type="submit"
      disabled={isLoading}
      class="bg-black text-white rounded-none px-8 py-4 font-mono font-bold text-lg border-4 border-black hover:bg-yellow-400 hover:text-black hover:-translate-y-1 shadow-[4px_4px_0px_0px_#ccc] hover:shadow-[6px_6px_0px_0px_#000] active:translate-y-0 active:shadow-none transition-all duration-200 disabled:opacity-50 cursor-pointer"
    >
      {isLoading ? 'WORKING...' : 'SHORTEN'}
    </button>
  </div>

  <div class="bg-gray-100 border-4 border-black p-4 rounded-none space-y-3 font-mono">
    <div class="flex flex-wrap gap-4 text-sm font-bold">
      <label class="flex items-center gap-2">
        <span>EXPIRATION:</span>
        <select bind:value={expiresIn} class="bg-white border-2 border-black p-1 rounded-none outline-none">
          <option value="never">NEVER</option>
          <option value="1">1 HOUR</option>
          <option value="24">24 HOURS</option>
        </select>
      </label>
      
      <button type="button" onclick={() => showAliasField = !showAliasField} class="underline hover:text-red-500">
        {showAliasField ? '- REMOVE ALIAS' : '+ ADD ALIAS'}
      </button>
    </div>

    {#if showAliasField}
      <div class="mt-2">
        <input 
          type="text" 
          bind:value={alias} 
          placeholder="Custom vanity slug (single URL requests only)" 
          class="w-full bg-white border-2 border-black p-2 rounded-none outline-none text-xs"
        />
      </div>
    {/if}

    <div class="flex items-center gap-2 text-sm font-bold mt-2">
      <span>PASSWORD LOCK:</span>
      <div class="flex border-2 border-black bg-white items-center">
        <input
          type={showPassword ? 'text' : 'password'}
          bind:value={passwordLock}
          placeholder="Optional password protection"
          class="p-1 outline-none text-xs w-48"
        />
        <button type="button" onclick={() => showPassword = !showPassword} class="px-2 border-l-2 border-black text-xs hover:bg-gray-200">
          {showPassword ? 'HIDE' : 'SHOW'}
        </button>
      </div>
    </div>
  </div>
</form>

{#if errorMessage}
  <div class="max-w-2xl mx-auto mt-4 bg-red-100 border-4 border-black p-4 rounded-none font-mono font-bold text-red-600">
    {errorMessage}
  </div>
{/if}

{#if shortenedLinks.length === 0 && !errorMessage}
  <div class="max-w-2xl mx-auto mt-6 border-4 border-dashed border-gray-400 bg-gray-50 p-6 rounded-none font-mono text-center shadow-[4px_4px_0px_0px_#ccc]">
    <div class="flex flex-col items-center justify-center py-6 text-gray-400">
      <span class="text-3xl mb-2">📊</span>
      <p class="text-sm font-bold uppercase tracking-wider">Link & High-Res Vector QR Code Preview Canvas</p>
      <p class="text-xs mt-1">Your generated asset blocks will populate right here upon shortening your source destinations.</p>
    </div>
  </div>
{/if}

{#if shortenedLinks.length > 0}
  <div class="max-w-2xl mx-auto mt-6 space-y-4">
    {#each shortenedLinks as link, i}
      {#await getQRSvg(link.short_url) then svg}
        <div class="bg-white border-4 border-black rounded-none p-4 shadow-[4px_4px_0px_0px_#000] flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 hover:shadow-[6px_6px_0px_0px_#FFD700] hover:-translate-y-1 transition-all duration-300">
          <div class="w-28 h-28 shrink-0 border-4 border-black shadow-[2px_2px_0px_0px_#ccc] p-1 bg-white flex items-center justify-center">
              {@html svg}
          </div>
          <div class="truncate w-full sm:max-w-xs xl:max-w-md ml-0 sm:ml-4 flex-1">
            <p class="text-xs font-bold text-gray-400 truncate">{link.long_url}</p>
            <a href={link.short_url} target="_blank" class="text-lg font-bold underline text-black hover:text-red-600 font-mono">{link.short_url}</a>
          </div>
          <div class="flex flex-col gap-2 w-full sm:w-auto shrink-0 font-mono">
            <button
              type="button"
              onclick={() => handleCopy(link.short_url, i)}
              class="flex-1 sm:flex-none border-2 border-black rounded-none px-3 py-1 text-xs font-bold shadow-[2px_2px_0px_0px_#000] transition-all uppercase"
              class:bg-yellow-400={copyTargetIndex === i}
              class:bg-white={copyTargetIndex !== i}
            >
              {copyTargetIndex === i ? 'COPIED!' : 'COPY SHORT LINK'}
            </button>
            <button
              type="button"
              onclick={() => downloadSVG(svg, `lmq-qr-${link.short_url.split('/').pop()}.svg`)}
              class="flex-1 sm:flex-none bg-black text-white border-2 border-black rounded-none px-3 py-1 text-xs font-bold shadow-[2px_2px_0px_0px_#000] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none transition-all uppercase"
            >
              DOWNLOAD VECTOR QR
            </button>
          </div>
        </div>
      {/await}
    {/each}
  </div>
{/if}
