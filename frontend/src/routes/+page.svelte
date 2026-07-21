<script lang="ts">
  import Seo from '$lib/components/Seo.svelte';
  import UrlForm from '$lib/components/UrlForm.svelte';
</script>

<Seo title="Lightning-Fast URL Shortener" description="Transform long URLs into powerful, trackable short links with enterprise-grade analytics." />

<div class="py-12 md:py-20 flex flex-col items-center text-center">
  <div class="border-4 border-black bg-white p-8 mb-12 shadow-hard-lg max-w-4xl mx-auto transform -rotate-1">
    <h1 class="text-6xl md:text-8xl font-bold tracking-tighter mb-4 text-black uppercase">
      SHORTEN.<br/>TRACK.<br/><span class="text-accent">OPTIMIZE.</span>
    </h1>
    <p class="text-xl md:text-2xl text-black font-bold uppercase tracking-widest mt-6 border-t-4 border-black pt-6">
      Transform long URLs into powerful, trackable short links.
    </p>
  </div>

  <div class="w-full mb-24 z-10">
    <UrlForm />
  </div>

  <div class="grid grid-cols-1 md:grid-cols-3 gap-0 w-full max-w-5xl border-4 border-black shadow-hard bg-white">
    <div class="p-8 border-b-4 md:border-b-0 md:border-r-4 border-black hover:bg-warning transition-none group">
      <div class="w-16 h-16 bg-black flex items-center justify-center text-white mb-6 text-3xl font-bold border-2 border-black group-hover:bg-white group-hover:text-black">⚡</div>
      <h3 class="text-2xl font-bold mb-4 uppercase tracking-wider text-black">LIGHTNING FAST</h3>
      <p class="text-black font-mono font-bold text-sm">Built on Go and SvelteKit for sub-millisecond redirects and instantaneous dashboard updates.</p>
    </div>
    <div class="p-8 border-b-4 md:border-b-0 md:border-r-4 border-black hover:bg-accent hover:text-white transition-none group text-right md:text-center">
      <div class="w-16 h-16 bg-black flex items-center justify-center text-white mb-6 text-3xl font-bold border-2 border-black group-hover:bg-white group-hover:text-black ml-auto md:mx-auto">📊</div>
      <h3 class="text-2xl font-bold mb-4 uppercase tracking-wider">ADVANCED ANALYTICS</h3>
      <p class="font-mono font-bold text-sm">Track clicks, geographic data, browsers, and devices in real-time with comprehensive charts.</p>
    </div>
    <div class="p-8 hover:bg-danger hover:text-white transition-none group">
      <div class="w-16 h-16 bg-black flex items-center justify-center text-white mb-6 text-3xl font-bold border-2 border-black group-hover:bg-white group-hover:text-black">🛡️</div>
      <h3 class="text-2xl font-bold mb-4 uppercase tracking-wider">SECURE & RELIABLE</h3>
      <p class="font-mono font-bold text-sm">Enterprise-grade security, rate limiting, and redundant infrastructure for your links.</p>
    </div>
  </div>
</div>
