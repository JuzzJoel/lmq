<script lang="ts">
  import '../app.css';

  const { children } = $props();
  let mobileMenuOpen = $state(false);
</script>

<div class="min-h-screen flex flex-col bg-bg-body text-black font-sans">
  <header class="fixed top-0 left-0 right-0 z-50 bg-white border-b-4 border-black">
    <div class="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex-shrink-0 border-r-4 border-black h-full flex items-center pr-6">
          <a href="/" class="text-3xl font-bold uppercase tracking-tighter text-black">
            LMQ
          </a>
        </div>
        <nav class="hidden md:flex h-full">
          <div class="flex items-stretch h-full">
            <a href="/" class="flex items-center text-black hover:bg-black hover:text-white px-6 font-bold uppercase tracking-wider border-l-4 border-black transition-none">
              HOME
            </a>
            <a href="/dashboard" class="flex items-center text-black hover:bg-black hover:text-white px-6 font-bold uppercase tracking-wider border-l-4 border-black transition-none">
              DASHBOARD
            </a>
            <a href="/help" class="flex items-center text-black hover:bg-black hover:text-white px-4 font-bold text-lg border-l-4 border-black transition-none leading-none" aria-label="Help">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="square" stroke-linejoin="miter">
                <circle cx="12" cy="12" r="9" />
                <path d="M9.5 9.5c0-1.4 1.1-2.5 2.5-2.5s2.5 1.1 2.5 2.5c0 1.4-1.1 2.5-2.5 2.5" />
                <path d="M12 16v.01" />
              </svg>
            </a>
          </div>
        </nav>
        <div class="md:hidden">
          <button type="button" onclick={() => mobileMenuOpen = !mobileMenuOpen} class="text-black hover:bg-black hover:text-white p-2 border-2 border-black transition-none">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              {#if mobileMenuOpen}
                <path stroke-linecap="square" stroke-linejoin="miter" stroke-width="3" d="M6 18L18 6M6 6l12 12" />
              {:else}
                <path stroke-linecap="square" stroke-linejoin="miter" stroke-width="3" d="M4 6h16M4 12h16M4 18h16" />
              {/if}
            </svg>
          </button>
        </div>
      </div>
    </div>
    {#if mobileMenuOpen}
      <div class="md:hidden bg-white border-t-4 border-black border-b-4">
        <div class="flex flex-col">
          <a href="/" class="text-black hover:bg-black hover:text-white px-4 py-3 font-bold uppercase border-b-2 border-black transition-none">HOME</a>
          <a href="/dashboard" class="text-black hover:bg-black hover:text-white px-4 py-3 font-bold uppercase border-b-2 border-black transition-none">DASHBOARD</a>
          <a href="/help" class="text-black hover:bg-black hover:text-white px-4 py-3 font-bold uppercase transition-none">HELP</a>
        </div>
      </div>
    {/if}
  </header>

  <main class="flex-grow pt-20 w-full">
    <div class="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8">
      {@render children()}
    </div>
  </main>

  <footer class="bg-white border-t-4 border-black py-8 mt-12">
    <div class="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 flex flex-col items-center justify-center text-black font-bold uppercase text-sm space-y-2">
      <p>&copy; {new Date().getFullYear()} LMQ. ALL RIGHTS RESERVED.</p>
      <a href="https://github.com/JuzzJoel" target="_blank" rel="noopener noreferrer" class="inline-block px-3 py-1 bg-warning border-2 border-black hover:bg-black hover:text-white transition-none shadow-[2px_2px_0px_0px_#000]">
        AMANI
      </a>
    </div>
  </footer>
</div>
