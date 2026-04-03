<script>
    import icMore from '$lib/assets/dots-vertical.svg?raw';
    /**
     * @typedef {Object} Props
     * @property {string} content
     */
    let { content } = $props();
</script>

<article>
    <header>
        <small>time</small>
        <button>{@html icMore}</button>
    </header>

    {content}
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

            & > button {
                padding: 0;
                background-color: transparent;
                border: none;
                color: oklch(from var(--pico-contrast) l c h / 0.5);
            }
        }
    }
</style>
