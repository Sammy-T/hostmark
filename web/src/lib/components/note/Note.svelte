<script>
    import icPublic from '$lib/assets/world-map.svg?raw';
    import icProtected from '$lib/assets/users.svg?raw';
    import icPrivate from '$lib/assets/lock.svg?raw';
    import icMore from '$lib/assets/dots-vertical.svg?raw';
    import icEdit from '$lib/assets/edit.svg?raw';
    import icTrash from '$lib/assets/trash.svg?raw';
    import icTag from '$lib/assets/tag.svg?raw';
    import icTagFilled from '$lib/assets/tag-filled.svg?raw';
    import { marked } from 'marked';
    import { getContext } from 'svelte';

    let { note } = $props();

    /** @type {import('svelte/reactivity').SvelteSet<string>} */
    let selectedTags = getContext('selectedTags');

    let menuId = $derived(`note-menu-${note.id}`);

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

<article>
    <header>
        <small>{new Intl.DateTimeFormat().format(new Date(note.created_at))}</small>
        
        <div class="side">
            <small>{note.owner}</small>

            {#if note.visibility === 'public'}
                {@html icPublic}
            {:else if note.visibility === 'protected'}
                {@html icProtected}
            {:else if note.visibility === 'private'}
                {@html icPrivate}
            {/if}
            
            <button popovertarget={menuId}>{@html icMore}</button>

            <div id={menuId} popover>
                <div class="menu-container">
                    <button popovertarget={menuId}>{@html icEdit} Edit</button>
                    <button popovertarget={menuId}>{@html icTrash} Delete</button>
                </div>
            </div>
        </div>
    </header>

    {#await marked.parse(note.content) then parsed}
        {@html parsed}
    {/await}

    {#if note.tags.length > 0}
        <hr class="separator">
    {/if}
    
    <div class="tags">
        {#each note.tags as tag (tag.name)}
            <a href={`#${tag}`} data-tag={tag.name} onclick={tagClicked}>
                {#if selectedTags.has(tag.name)}
                    {@html icTagFilled}
                {:else}
                    {@html icTag}
                {/if}

                {tag.name}
            </a>
        {/each}
    </div>
</article>

<style>
    [popover] {
        margin: 0;
        inset: auto;
        border: none;
        top: calc(anchor(bottom) + 0.25rem);
        right: anchor(right);
        border: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        border-radius: 0.25rem;
        background-color: var(--pico-dropdown-background-color);
        box-shadow: var(--pico-dropdown-box-shadow);
    }

    .menu-container {
        display: flex;
        flex-direction: column;

        & button {
            padding: 0 0.5rem;
            display: flex;
            align-items: center;
            gap: 0.25rem;
            color: var(--pico-dropdown-color);

            & :global(svg) {
                width: 1rem;
            }
        }

        & button:hover {
            background-color: var(--pico-dropdown-hover-background-color);
        }
    }

    article {
        margin: 0;

        & header {
            display: flex;
            justify-content: space-between;
            align-items: center;

            & > small {
                color: var(--pico-muted-color);
            }

            & button {
                padding: 0;
                background-color: transparent;
                border: none;
                color: oklch(from var(--pico-contrast) l c h / 0.5);
            }

            .side > :global(svg) {
                color: oklch(from var(--pico-contrast) l c h / 0.5);
            }

            .side > small {
                margin-right: 0.5rem;
            }
        }
    }

    .separator {
        margin: 0.25rem 0;
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
    }
</style>
