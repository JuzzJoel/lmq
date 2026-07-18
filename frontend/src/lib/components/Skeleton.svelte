<script lang="ts">
  interface Props {
    width?: string;
    height?: string;
    rounded?: string;
    count?: number;
  }

  const { width = '100%', height = '1rem', rounded = 'rounded-none', count = 1 }: Props = $props();
</script>

{#each Array(count) as _}
  <div 
    class="bg-white border-2 border-black animate-pulse shadow-hard {rounded}"
    style="width: {width}; height: {height};"
  ></div>
{/each}
