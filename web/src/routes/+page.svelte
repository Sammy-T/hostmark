<script>
    import FileBrowser from '$lib/components/FileBrowser.svelte';

    /** @type {import('./$types').PageProps} */
    let { data } = $props();
</script>
 
<!-- {@html data.content} -->
<FileBrowser />
