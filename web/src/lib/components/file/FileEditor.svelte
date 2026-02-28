<script>
    import { getContext } from 'svelte';

    /** @type {{ markdown: string, html: string }} */
    const content = getContext('content');
</script>

<section>{@html content.html}</section>

<style>
    section {
        flex-grow: 1;
        padding: 0.5rem;
        margin: 0;
        overflow: auto;
    }
</style>
