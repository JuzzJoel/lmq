<script lang="ts">
  import Seo from '$lib/components/Seo.svelte';
</script>

<Seo title="404: Link Terminated" description="This link has expired or been deleted." />

<div class="min-h-[70vh] flex items-center justify-center">
  <div class="border-8 border-black bg-white p-12 text-center shadow-[16px_16px_0px_0px_#FF2A2A] max-w-3xl transform rotate-1">
    <h1 class="text-6xl md:text-8xl font-black uppercase tracking-tighter mb-6">
      LINK<br/>TERMINATED
    </h1>
    <div class="h-2 bg-black w-full my-8"></div>
    <p class="text-2xl font-bold uppercase tracking-widest mb-10">
      EXPIRED OR DELETED
    </p>
    <a href="/" class="inline-block px-10 py-4 bg-black text-white text-xl font-bold uppercase tracking-widest border-4 border-black hover:bg-white hover:text-black transition-none shadow-[8px_8px_0px_0px_#FFD700] active:translate-x-2 active:translate-y-2 active:shadow-none">
      RETURN TO BASE
    </a>
  </div>
</div>
