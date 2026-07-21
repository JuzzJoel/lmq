<script lang="ts">
</script>

<svelte:head>
  <title>Help &amp; Features - LMQ</title>
  <meta name="description" content="Learn how to use LMQ link shortener — expiration, custom aliases, A/B routing, password locks, and burn-after-reading." />
</svelte:head>

<div class="space-y-8 animate-slide-up">
  <div class="border-b-4 border-black pb-4">
    <h1 class="text-4xl font-bold uppercase tracking-tighter">Help &amp; Features</h1>
    <p class="text-gray-600 font-mono mt-2 uppercase text-sm font-bold">Everything you need to know about LMQ's link shortening features.</p>
  </div>

  <div class="bg-white border-4 border-black shadow-[8px_8px_0px_0px_#000] p-6 lg:p-10 space-y-10">
    <section class="bg-warning border-4 border-black p-6 lg:p-8 shadow-hard-lg -mx-2 lg:-mx-4 space-y-10">
      <h2 class="text-3xl md:text-4xl font-black uppercase tracking-tighter border-b-4 border-black pb-4">Core Features</h2>

      <div class="space-y-10">
        <div class="bg-white border-4 border-black shadow-hard p-6">
          <div class="flex items-start gap-4">
            <span class="text-3xl flex-shrink-0 bg-black text-white w-12 h-12 flex items-center justify-center font-bold text-xl border-2 border-black leading-none">01</span>
            <div class="space-y-3 flex-1 min-w-0">
              <h3 class="text-2xl font-black uppercase tracking-tight">Expiration Time</h3>
              <p class="font-mono text-sm leading-relaxed">
                Set a time limit on any short link. Once the expiration passes, the link stops working and returns a 404. Useful for time-sensitive promotions, temporary access, or event-based campaigns.
              </p>
              <div class="bg-black text-green-400 p-4 border-2 border-black font-mono text-xs overflow-x-auto">
<pre>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;
    "url": "https://example.com/promo",
    "expires_in": 24
  &rbrace;'

// The link will self-destruct 24 hours after creation.
// Response includes "expires_at" timestamp.</pre>
              </div>
              <p class="text-xs font-mono uppercase font-bold">Set <span class="bg-gray-200 px-1 border border-black">expires_in</span> to the number of hours. Omit or set to <span class="bg-gray-200 px-1 border border-black">0</span> for a permanent link.</p>
            </div>
          </div>
        </div>

        <div class="bg-white border-4 border-black shadow-hard p-6">
          <div class="flex items-start gap-4">
            <span class="text-3xl flex-shrink-0 bg-black text-white w-12 h-12 flex items-center justify-center font-bold text-xl border-2 border-black leading-none">02</span>
            <div class="space-y-3 flex-1 min-w-0">
              <h3 class="text-2xl font-black uppercase tracking-tight">Custom Alias / Token</h3>
              <p class="font-mono text-sm leading-relaxed">
                Replace the random short code with a memorable, branded slug. Your alias must be 3&ndash;20 characters long and can only contain letters, numbers, and hyphens.
              </p>
              <div class="bg-black text-green-400 p-4 border-2 border-black font-mono text-xs overflow-x-auto">
<pre>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;
    "url": "https://example.com/whitepaper",
    "custom_token": "whitepaper-q2"
  &rbrace;'

// Your short link becomes: https://lmq.name.ng/whitepaper-q2
// Custom tokens cannot be used with bulk shortening.</pre>
              </div>
              <p class="text-xs font-mono uppercase font-bold">Use <span class="bg-gray-200 px-1 border border-black">custom_token</span> to claim a specific slug. A conflict returns HTTP 409.</p>
            </div>
          </div>
        </div>

        <div class="bg-white border-4 border-black shadow-hard p-6">
          <div class="flex items-start gap-4">
            <span class="text-3xl flex-shrink-0 bg-black text-white w-12 h-12 flex items-center justify-center font-bold text-xl border-2 border-black leading-none">03</span>
            <div class="space-y-3 flex-1 min-w-0">
              <h3 class="text-2xl font-black uppercase tracking-tight">A/B Routing</h3>
              <p class="font-mono text-sm leading-relaxed">
                Split traffic across multiple destination URLs with weighted routing. Each visit randomly selects a destination based on the weights you assign. Great for A/B testing landing pages or gradually rolling out new content.
              </p>
              <div class="bg-black text-green-400 p-4 border-2 border-black font-mono text-xs overflow-x-auto">
<pre>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;
    "url": "https://example.com/control",
    "routes": [
      &lbrace; "url": "https://example.com/variant-a", "weight": 70 &rbrace;,
      &lbrace; "url": "https://example.com/variant-b", "weight": 30 &rbrace;
    ]
  &rbrace;'

// 70% of clicks go to variant-a, 30% to variant-b.
// The base "url" is used as fallback if routes are empty.</pre>
              </div>
              <p class="text-xs font-mono uppercase font-bold">Weights must be positive integers. Routes are optional &mdash; omit them for a standard single-destination link.</p>
            </div>
          </div>
        </div>

        <div class="bg-white border-4 border-black shadow-hard p-6">
          <div class="flex items-start gap-4">
            <span class="text-3xl flex-shrink-0 bg-black text-white w-12 h-12 flex items-center justify-center font-bold text-xl border-2 border-black leading-none">04</span>
            <div class="space-y-3 flex-1 min-w-0">
              <h3 class="text-2xl font-black uppercase tracking-tight">Password Lock</h3>
              <p class="font-mono text-sm leading-relaxed">
                Protect a short link with a password. Visitors must enter the password on the unlock page before they are redirected to the destination. The password is hashed with bcrypt before storage &mdash; LMQ never stores plaintext passwords.
              </p>
              <div class="bg-black text-green-400 p-4 border-2 border-black font-mono text-xs overflow-x-auto">
<pre>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;
    "url": "https://example.com/private",
    "password": "s3cret!"
  &rbrace;'

// Visiting the short link shows a password prompt.
// POST /api/v1/verify-password with the correct
// password returns a redirect to the destination.</pre>
              </div>
              <p class="text-xs font-mono uppercase font-bold">Set <span class="bg-gray-200 px-1 border border-black">password</span> to any string. There is no way to recover a lost password.</p>
            </div>
          </div>
        </div>

        <div class="bg-white border-4 border-black shadow-hard p-6">
          <div class="flex items-start gap-4">
            <span class="text-3xl flex-shrink-0 bg-black text-white w-12 h-12 flex items-center justify-center font-bold text-xl border-2 border-black leading-none">05</span>
            <div class="space-y-3 flex-1 min-w-0">
              <h3 class="text-2xl font-black uppercase tracking-tight">Burn After Reading</h3>
              <p class="font-mono text-sm leading-relaxed">
                Create a self-destructing link that can only be viewed once. After the first visit, the link is atomically deleted from the database. Any subsequent visit returns a 404. Ideal for sharing sensitive, one-time information.
              </p>
              <div class="bg-black text-green-400 p-4 border-2 border-black font-mono text-xs overflow-x-auto">
<pre>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;
    "url": "https://example.com/sensitive",
    "burn_after_reading": true
  &rbrace;'

// The first click redirects and the link vanishes.
// No analytics are recorded for BAR links.</pre>
              </div>
              <p class="text-xs font-mono uppercase font-bold">Combine with <span class="bg-gray-200 px-1 border border-black">password</span> for a double-locked one-time link. The link is consumed on successful password verification.</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="space-y-4">
      <h2 class="text-xl font-bold uppercase bg-warning inline-block px-2 border-2 border-black">Quick Start</h2>
      <p class="font-mono text-sm leading-relaxed">
        The simplest way to shorten a link is a single POST with just a <span class="bg-gray-200 px-1 border border-black font-mono">url</span> field. All other features are optional and can be mixed and matched freely.
      </p>
      <div class="bg-black text-green-400 p-4 border-2 border-black font-mono text-xs overflow-x-auto">
<pre>curl -X POST https://api.lmq.name.ng/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '&lbrace;"url": "https://example.com"&rbrace;'

// Response (201 Created):
&lbrace;
  "data": &lbrace;
    "results": [&lbrace;
      "token": "aB3xY7",
      "short_url": "https://lmq.name.ng/aB3xY7",
      "long_url": "https://example.com",
      "created_at": "2026-07-21T00:00:00Z",
      "has_password": false,
      "burn_after_reading": false
    &rbrace;]
  &rbrace;,
  "error": null
&rbrace;</pre>
      </div>
      <p class="font-mono text-xs uppercase font-bold">All features can be combined &mdash; for example, a password-protected, burn-after-reading link with a custom alias and A/B routes.</p>
    </section>

    <section class="space-y-4">
      <h2 class="text-xl font-bold uppercase bg-warning inline-block px-2 border-2 border-black">Need More Help?</h2>
      <p class="font-mono text-sm leading-relaxed">
        Visit the <a href="/dashboard/api-docs" class="font-bold underline hover:bg-black hover:text-white transition-none">API documentation</a> for the full developer reference, including analytics, bulk CSV upload, and the analytics export endpoint.
      </p>
    </section>
  </div>
</div>