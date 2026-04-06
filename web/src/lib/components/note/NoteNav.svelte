<script>
    import icTag from '$lib/assets/tag.svg?raw';
    import icTagFilled from '$lib/assets/tag-filled.svg?raw';
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
        const tag = event.target?.dataset.tag;

        if(selectedTags.has(tag)) {
            selectedTags.delete(tag);
            return;
        }

        selectedTags.add(tag);
    }
</script>

{#snippet tagItem(/** @type {string} */ tag)}
    <a href={`#${tag}`} data-tag={tag} onclick={tagClicked}>
        {#if selectedTags.has(tag)}
            {@html icTagFilled}
        {:else}
            {@html icTag}
        {/if}

        {tag}
    </a>
{/snippet}

<!-- Mobile note nav -->
<Sidebar popId="mobile-note-nav" isMenu={false} mobileOnly>
    <h6>Tags</h6>
    <div class="tags">
        {#each tags.values() as tag}
            {@render tagItem(tag)}
        {/each}
    </div>
</Sidebar>

<!-- Desktop note nav -->
<aside>
    <h6>Tags</h6>
    <div class="tags">
        {#each tags.values() as tag}
            {@render tagItem(tag)}
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

        & :global(svg) {
            width: 1rem;
        }
    }

    @media (min-width: 768px) {
        aside {
            display: revert-layer;
        }
    }
</style>
