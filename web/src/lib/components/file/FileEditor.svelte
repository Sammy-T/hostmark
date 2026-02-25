<script>
    import { getContext } from 'svelte';

    /** @type {{ markdown: string, html: string }} */
    const content = getContext('content');
</script>

<div>
    <header>todo/breadcrumb/path</header>
    <section>{@html content.html}</section>
</div>

<style>
    div {
        height: 100dvh;
        flex-grow: 1;
        overflow: auto;

        & > * {
            padding: 0.5rem;
        }

        & header {
            outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        }

        & section {
            margin: 0;
        }
    }
</style>
