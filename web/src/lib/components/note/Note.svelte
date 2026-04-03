<script>
    import icPublic from '$lib/assets/world-map.svg?raw';
    import icProtected from '$lib/assets/users.svg?raw';
    import icPrivate from '$lib/assets/lock.svg?raw';
    import icMore from '$lib/assets/dots-vertical.svg?raw';

    let { note } = $props();
</script>

<article>
    <header>
        <small>{new Intl.DateTimeFormat().format(new Date(note.created_at))}</small>
        
        <div class="more">
            {#if note.visibility === 'public'}
                {@html icPublic}
            {:else if note.visibility === 'protected'}
                {@html icProtected}
            {:else if note.visibility === 'private'}
                {@html icPrivate}
            {/if}
            
            <button>{@html icMore}</button> <!-- TODO: add menu -->
        </div>
    </header>

    {note.content}
    
    <div class="tags">
        {#each note.tags as tag (tag.name)}
            <a href="#">{tag.name}</a>
        {/each}
    </div>
</article>

<style>
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

            .more > :global(svg) {
                color: oklch(from var(--pico-contrast) l c h / 0.5);
            }
        }
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;

        & a {
            text-decoration: underline;
        }
    }
</style>
