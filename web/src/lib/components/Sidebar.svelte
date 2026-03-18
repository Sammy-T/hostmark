<script>
    /**
     * @typedef {Object} Props
     * @property {Boolean} [mobileOnly]
     * @property {Boolean} [desktopOnly]
     * @property {String} [popId]
     * @property {import('svelte').Snippet} children
     */

    /** @type {Props} */
    let { mobileOnly = false, desktopOnly = false, popId, children } = $props();
</script>

<div id={popId} class:mobile={mobileOnly} class:desktop={desktopOnly} popover={(popId) ? 'auto' : null}>
    <aside class="sidebar">
        <nav>
            {@render children()}
        </nav>
    </aside>
</div>

<style>
    [popover] {
        margin: 0;
        inset: auto;
        border: none;
    }

    .sidebar {
        :global(button) {
            padding: 0;
            margin: 0;
            background-color: transparent;
            border: none;
        }

        :global(button:hover) {
            color: var(--pico-primary);
        }

        :global(a:hover) {
            text-decoration: none;
        }
    }

    aside {
        height: 100dvh;
        background-color: oklch(from var(--pico-background-color) calc(l * 0.85) c h);

        & nav {
            min-height: 100%;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
        }
    }

    nav {
        padding: 0.5rem;
        margin: 0;

        & :global(*), & :global(ul:first-of-type), & :global(ul:last-of-type) {
            padding: 0;
            margin: 0;
        }

        & :global(li) {
            padding: 0.5rem 0;
        }
    }

    .mobile {
        min-width: 10rem;
        position: absolute;
        top: 0;
        left: 0;
    }

    .desktop {
        display: none;
        outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);

        & :global(a) {
            padding: 0.15rem;
            display: flex;
            justify-content: center;
            align-items: center;
        }
    }

    @media (min-width: 768px) {
        .mobile {
            display: none;
        }

        .desktop {
            display: revert-layer;
        }
    }
</style>
