<script lang="ts">
  // No data fetching required for static docs layout
</script>

<svelte:head>
  <title>API Documentation - LMQ</title>
</svelte:head>

<div class="space-y-8 animate-slide-up">
  <div class="border-b-4 border-black pb-4">
    <h1 class="text-4xl font-bold uppercase tracking-tighter">Developer API</h1>
    <p class="text-gray-600 font-mono mt-2 uppercase text-sm font-bold">Integrate LMQ directly into your own infrastructure.</p>
  </div>

  <div class="bg-white border-4 border-black shadow-[8px_8px_0px_0px_#000] p-6 lg:p-10 space-y-8">
    <section>
      <h2 class="text-2xl font-bold uppercase bg-yellow-300 inline-block px-2 border-2 border-black mb-4">Link Shortening API</h2>
      <p class="font-mono text-sm mb-4">You can programmatically shorten URLs using our core REST API. The endpoint handles bulk shortening, custom aliases, passwords, and expiration rules automatically.</p>

      <div class="bg-gray-100 border-4 border-black p-4 font-mono text-sm">
        <span class="bg-black text-white px-2 py-1 font-bold mr-2">POST</span>
        <span class="font-bold text-blue-700">https://api.lmq.name.ng/api/v1/shorten</span>
      </div>
    </section>

    <section class="space-y-4">
      <h3 class="text-xl font-bold uppercase underline">Request Payload (JSON)</h3>
      <div class="bg-black text-green-400 p-4 border-4 border-black font-mono text-sm overflow-x-auto">
<pre><code>&lbrace;
  "url": "https://google.com https://example.com", // String. Separate multiple URLs with spaces or commas (max 50)
  "expires_in": 24, // Integer (Optional). Hours until link dies. 0 or omit for never.
  "password": "secret_password", // String (Optional). Lock the destination.
  "alias": "my-custom-slug" // String (Optional). Custom token (only valid if shortening a single URL).
&rbrace;</code></pre>
      </div>
    </section>

    <section class="space-y-4">
      <h3 class="text-xl font-bold uppercase underline">Response Payload (JSON)</h3>
      <div class="bg-black text-green-400 p-4 border-4 border-black font-mono text-sm overflow-x-auto">
<pre><code>&lbrace;
  "results": [
    &lbrace;
      "token": "my-custom-slug",
      "short_url": "https://lmq.name.ng/my-custom-slug",
      "long_url": "https://google.com",
      "created_at": "2026-07-18T14:30:00Z",
      "expires_at": "2026-07-19T14:30:00Z",
      "has_password": true
    &rbrace;
  ]
&rbrace;</code></pre>
      </div>
    </section>

    <section class="space-y-4">
      <h3 class="text-xl font-bold uppercase bg-blue-300 inline-block px-2 border-2 border-black">Rate Limiting</h3>
      <p class="font-mono text-sm">Our pure IP rate limiter restricts incoming shortening requests to <span class="bg-red-200 px-1 font-bold border border-black">100 requests per hour</span> per IP address to ensure network stability. Violations will return <code class="bg-gray-200 px-1 border border-black">HTTP 429 Too Many Requests</code>.</p>
    </section>

    <section class="space-y-4">
      <h3 class="text-xl font-bold uppercase underline">Usage Examples</h3>
      
      <div class="space-y-2">
        <h4 class="font-bold font-mono text-sm">cURL Example</h4>
        <div class="bg-black text-green-400 p-4 border-4 border-black font-mono text-sm overflow-x-auto">
<pre><code>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;"url": "https://google.com"&rbrace;'</code></pre>
        </div>
      </div>

      <div class="space-y-2">
        <h4 class="font-bold font-mono text-sm">JavaScript (Fetch) Example</h4>
        <div class="bg-black text-green-400 p-4 border-4 border-black font-mono text-sm overflow-x-auto">
<pre><code>const response = await fetch('https://api.lmq.name.ng/api/v1/shorten', &lbrace;
  method: 'POST',
  headers: &lbrace;
    'Content-Type': 'application/json'
  &rbrace;,
  body: JSON.stringify(&lbrace;
    url: 'https://google.com'
  &rbrace;)
&rbrace;);

const data = await response.json();
console.log(data.results[0].short_url);</code></pre>
        </div>
      </div>
    </section>

  </div>
</div>
