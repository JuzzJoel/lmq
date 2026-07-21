<script lang="ts">
  import { onMount } from 'svelte';
  import { authState, checkAuth, login, logout } from '$lib/auth.svelte';

  const { children } = $props();

  let adminTokenInput = $state('');
  let showPassword = $state(false);

  onMount(() => {
    checkAuth();
  });

  function handleLogin(e: Event) {
    e.preventDefault();
    if (adminTokenInput.trim().length > 0) {
      login(adminTokenInput);
    }
  }
</script>

{#if authState.isChecking}
  <div class="min-h-[60vh] flex items-center justify-center font-bold font-mono">LOADING...</div>
{:else if !authState.isAuthenticated}
  <div class="min-h-[60vh] flex items-center justify-center animate-slide-up w-full">
    <form onsubmit={handleLogin} class="bg-white border-4 border-black shadow-[4px_4px_0px_0px_#000] p-8 w-full max-w-lg rounded-none">
      <h2 class="text-2xl font-bold uppercase tracking-widest mb-6 border-b-4 border-black pb-4 text-center">SYSTEM LOGIN</h2>
      <div class="mb-6">
        <label for="admin_token" class="block font-bold uppercase tracking-wider mb-2">ENTER ADMIN PASSWORD</label>
        <div class="flex shadow-[4px_4px_0px_0px_#000]">
          <input
            id="admin_token"
            type={showPassword ? "text" : "password"}
            bind:value={adminTokenInput}
            class="w-full px-4 py-3 bg-white border-4 border-black border-r-0 focus:outline-none focus:bg-warning rounded-none font-mono font-bold"
            placeholder="Plaintext Password"
            required
          />
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
  <div class="space-y-6 animate-slide-up w-full">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between pb-4 border-b-4 border-black gap-4">
      <h1 class="text-3xl font-bold uppercase tracking-tighter">Dashboard</h1>
      <nav class="flex gap-4 font-mono font-bold text-sm">
        <a href="/dashboard" class="border-2 border-black px-4 py-2 hover:bg-black hover:text-white transition-colors shadow-[2px_2px_0px_0px_#000]">Overview</a>
        <a href="/dashboard/api-docs" class="border-2 border-black px-4 py-2 hover:bg-yellow-400 hover:text-black transition-colors shadow-[2px_2px_0px_0px_#000]">API Docs</a>
        <button onclick={logout} class="border-2 border-black bg-white px-4 py-2 hover:bg-red-500 hover:text-white transition-colors shadow-[2px_2px_0px_0px_#000]">Logout</button>
      </nav>
    </div>
    {@render children()}
  </div>
{/if}
