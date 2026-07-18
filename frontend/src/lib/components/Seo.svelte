<script lang="ts">
  interface Props {
    title: string;
    description: string;
    canonical?: string;
    ogImage?: string;
  }

  const { title, description, canonical = '', ogImage = '' }: Props = $props();
  const siteName = 'LMQ';
  let fullTitle = $derived(`${title} | ${siteName}`);
</script>

<svelte:head>
  <title>{fullTitle}</title>
  <meta name="description" content={description} />
  {#if canonical}
    <link rel="canonical" href={canonical} />
  {/if}
  <meta property="og:title" content={fullTitle} />
  <meta property="og:description" content={description} />
  <meta property="og:type" content="website" />
  {#if ogImage}
    <meta property="og:image" content={ogImage} />
  {/if}
  <meta property="og:site_name" content={siteName} />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content={fullTitle} />
  <meta name="twitter:description" content={description} />
</svelte:head>
