<script>
    import mdHello from '$lib/hello.md?raw';
    import { marked } from 'marked';
    import { onMount } from 'svelte';

    /** @type {String} */
    let content = $state('');

    async function loadContent() {
        console.log('loading content...');
        content = await marked.parse(mdHello);
    }

    onMount(() => {
        loadContent();
    });
</script>

{@html content}
