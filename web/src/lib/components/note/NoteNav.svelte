<script>
    import Sidebar from '../Sidebar.svelte';
    import { getContext } from 'svelte';

    /** @type {import('svelte/reactivity').SvelteSet<string>} */
    let tags = getContext('tags');

    /** @type {import('svelte/reactivity').SvelteSet<string>} */
    let selectedTags = getContext('selectedTags');

    /**
     * @param {Event} event
     */
    function tagClicked(event) {
        event.preventDefault();

        // @ts-ignore
        const tag = event.target?.innerText;

        if(selectedTags.has(tag)) {
            selectedTags.delete(tag);
            return;
        }

        selectedTags.add(tag);
    }
</script>

<!-- Mobile note nav -->
<Sidebar popId="mobile-note-nav" isMenu={false} mobileOnly>
    <h6>Tags</h6>
    <div class="tags">
        {#each tags.values() as tag}
            <a href={`#${tag}`} onclick={tagClicked}>{tag}</a>
        {/each}
    </div>
</Sidebar>

<!-- Desktop note nav -->
<aside>
    <h6>Tags</h6>
    <div class="tags">
        {#each tags.values() as tag}
            <a href={`#${tag}`} onclick={tagClicked}>{tag}</a>
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
