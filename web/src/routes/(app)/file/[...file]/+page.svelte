<script>
    import FileBrowser from '$lib/components/FileBrowser.svelte';
    import { setContext } from 'svelte';

    /** @type {import('./$types').PageProps} */
    let { params } = $props();

    // svelte-ignore state_referenced_locally
    // setContext('content', data.content);

    let file = $state({ value: params.file });
    setContext('file', file);

    $effect(() => {
        file.value = params.file;
    });
</script>

<FileBrowser />
