<script>
    import Sidebar from '../Sidebar.svelte';
</script>

<!-- Mobile note nav -->
<Sidebar popId="mobile-note-nav" isMenu={false} mobileOnly>
    <h6>Tags</h6>
    <div class="tags">
        {#each { length: 5 } as _, i }
            <a href="#">tag-{i}</a>
        {/each}
    </div>
</Sidebar>

<!-- Desktop note nav -->
<aside>
    <h6>Tags</h6>
    <div class="tags">
        {#each { length: 5 } as _, i}
            <a href="#">tag-{i}</a>
        {/each}
    </div>
</aside>

<style>
    aside {
        display: none;
        height: 100dvh;
        max-width: 25dvw;
        padding: 0.5rem;
        overflow: auto;
        outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;

        & a {
            text-decoration: underline;
        }
    }

    @media (min-width: 768px) {
        aside {
            display: revert-layer;
        }
    }
</style>
