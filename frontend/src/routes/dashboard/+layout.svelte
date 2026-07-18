<script lang="ts">
  const { children } = $props();
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between pb-4 border-b border-border-subtle">
    <h1 class="text-2xl font-bold">Dashboard</h1>
  </div>
  {@render children()}
</div>
