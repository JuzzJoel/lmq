<script lang="ts">
  const { children } = $props();
</script>

<div class="space-y-6 animate-slide-up">
  <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between pb-4 border-b-4 border-black gap-4">
    <h1 class="text-3xl font-bold uppercase tracking-tighter">Dashboard</h1>
    <nav class="flex gap-4 font-mono font-bold text-sm">
      <a href="/dashboard" class="border-2 border-black px-4 py-2 hover:bg-black hover:text-white transition-colors shadow-[2px_2px_0px_0px_#000]">Overview</a>
      <a href="/dashboard/api-docs" class="border-2 border-black px-4 py-2 hover:bg-yellow-400 hover:text-black transition-colors shadow-[2px_2px_0px_0px_#000]">API Docs</a>
    </nav>
  </div>
  {@render children()}
</div>
