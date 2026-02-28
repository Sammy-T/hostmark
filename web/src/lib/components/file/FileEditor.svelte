<script>
    import { page } from '$app/state';

    let content = $derived(page.data.content);
</script>

<section>{@html content?.html}</section>

<style>
    section {
        flex-grow: 1;
        padding: 0.5rem;
        margin: 0;
        overflow: auto;
    }
</style>
