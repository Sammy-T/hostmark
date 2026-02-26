<script>
    import FileBrowser from '$lib/components/FileBrowser.svelte';
    import { setContext } from 'svelte';

    /** @type {import('./$types').PageProps} */
    let { data } = $props();

    // svelte-ignore state_referenced_locally
    setContext('content', data.content);
</script>
 
<FileBrowser />
