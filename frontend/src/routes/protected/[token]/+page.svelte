<script lang="ts">
  import { page } from '$app/stores';
  import Seo from '$lib/components/Seo.svelte';

  let password = $state('');
  let loading = $state(false);
  let errorMessage = $state('');
  let showPassword = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!password) {
      errorMessage = "PASSWORD REQUIRED";
      return;
    }

    loading = true;
    errorMessage = "";
    
    try {
      const token = $page.params.token;
      const res = await fetch('/api/v1/verify-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, password })
      });

      const data = await res.json();
      if (res.ok && data.long_url) {
        window.location.href = data.long_url;
      } else {
        errorMessage = data.error || "ACCESS DENIED";
      }
    } catch (err) {
      errorMessage = "VERIFICATION FAILED";
    } finally {
      loading = false;
    }
  }
</script>

<Seo title="Protected Link" description="Enter password to access this secure link." />

<div class="min-h-[70vh] flex items-center justify-center px-4">
  <div class="border-4 border-black bg-white p-8 w-full max-w-lg shadow-[8px_8px_0px_0px_#000]">
    <div class="mb-8 text-center border-b-4 border-black pb-4">
      <h1 class="text-4xl font-black uppercase tracking-tighter">SECURE LINK</h1>
      <p class="font-bold text-sm mt-2">AUTHORIZATION REQUIRED</p>
    </div>

    <form onsubmit={handleSubmit} class="space-y-6">
      <div>
        <label class="block text-xl font-bold uppercase mb-2">Password</label>
        <div class="relative flex border-4 border-black">
          <input
            type={showPassword ? "text" : "password"}
            bind:value={password}
            placeholder="ENTER VAULT KEY"
            class="w-full px-4 py-3 bg-white text-black placeholder:text-gray-400 focus:outline-none focus:bg-warning transition-none font-mono text-xl font-bold rounded-none"
            disabled={loading}
          />
          <button
            type="button"
            onclick={() => showPassword = !showPassword}
            class="px-4 bg-black text-white font-bold uppercase hover:bg-warning hover:text-black transition-none"
          >
            {showPassword ? 'HIDE' : 'SHOW'}
          </button>
        </div>
      </div>

      {#if errorMessage}
        <div class="p-4 bg-white border-4 border-black shadow-[4px_4px_0px_0px_#FF2A2A]">
          <p class="text-red-500 font-bold uppercase text-center">{errorMessage}</p>
        </div>
      {/if}

      <button
        type="submit"
        disabled={loading}
        class="w-full py-4 bg-accent text-white border-4 border-black text-2xl font-black uppercase tracking-widest hover:bg-black transition-none shadow-[4px_4px_0px_0px_#000] active:translate-x-1 active:translate-y-1 active:shadow-none disabled:opacity-50"
      >
        {loading ? 'VERIFYING...' : 'DECRYPT'}
      </button>
    </form>
  </div>
</div>
