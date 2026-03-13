<script>
    import { fade } from 'svelte/transition';

    /**
     * @typedef {Object} Props
     * @property {import('svelte').Snippet} children
     * @property {Boolean} [mobileOnly]
     * @property {Boolean} [desktopOnly]
     */

    /** @type {Props} */
    let { children, mobileOnly = false, desktopOnly = false } = $props();
</script>

<aside class="sidebar" class:mobile={mobileOnly} class:desktop={desktopOnly} transition:fade>
    <nav>
        {@render children()}
    </nav>
</aside>

<style>
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
