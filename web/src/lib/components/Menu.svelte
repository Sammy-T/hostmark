<script>
    import icMenu from '$lib/assets/menu-2.svg?raw';
    import icClose from '$lib/assets/square-x.svg?raw';
    import icFiles from '$lib/assets/book-2.svg?raw';
    import icProfile from '$lib/assets/user-circle.svg?raw';
    import icSettings from '$lib/assets/settings.svg?raw';
    import icSignout from '$lib/assets/logout.svg?raw';
    import Sidebar from './Sidebar.svelte';

    /**
     * @typedef {Object} Props
     * @property {import('svelte').Snippet} [children]
     * @property {String} [current]
     */

    /** @type {Props} */
    let { current } = $props();
    // TODO: Set current based on location

    let showMobile = $state(false);

    function toggleMobile() {
        showMobile = !showMobile;
    }

    /**
     * @param {HTMLElement} node
     * @param {{ title: string; placement?: string; }} params
     */
    function delayedTip(node, params) {
        if(!params) return;

        /** @type {Number} */
        let tid;

        function handleMouseEnter() {
            tid = setTimeout(() => {
                node.dataset.tooltip = params?.['title'];
                node.dataset.placement = params?.['placement'] ?? 'right';
            }, 1000);
        }

        function handleMouseLeave() {
            clearTimeout(tid);
            node.removeAttribute('data-tooltip');
            node.removeAttribute('data-placement');
        }

        $effect(() => {
            node.addEventListener('mouseenter', handleMouseEnter);
            node.addEventListener('mouseleave', handleMouseLeave);

            return () => {
                node.removeEventListener('mouseenter', handleMouseEnter);
                node.removeEventListener('mouseleave', handleMouseLeave);
            };
        });
    }
</script>

<nav class="toggle">
    <button onclick={toggleMobile}>{@html icMenu}</button>
</nav>

{#if showMobile}
<Sidebar mobileOnly>
    <ul>
        <li><button onclick={toggleMobile}>{@html icClose}</button></li>
        <li>hm</li>
        <li><a href="#">{@html icFiles} Files</a></li>
        <li><a href="#">{@html icProfile} Profile</a></li>
    </ul>
    
    <ul>
        <li><a href="#">{@html icSettings} Settings</a></li>
        <li><a href="#">{@html icSignout} Sign out</a></li>
    </ul>
</Sidebar>
{/if}

<Sidebar desktopOnly>
    <ul>
        <li>hm</li>
        <li><a href="#" class:highlight={current === 'files'} use:delayedTip={{ title: 'Files' }}>{@html icFiles}</a></li>
        <li><a href="#" class:highlight={current === 'profile'} use:delayedTip={{ title: 'Profile' }}>{@html icProfile}</a></li>
    </ul>
    
    <ul>
        <li><a href="#" class:highlight={current === 'settings'} use:delayedTip={{ title: 'Settings' }}>{@html icSettings}</a></li>
        <li><a href="#" class:highlight={current === 'signout'} use:delayedTip={{ title: 'Sign out' }}>{@html icSignout}</a></li>
    </ul>
</Sidebar>

<style>
    .toggle {
        & button {
            padding: 0;
            margin: 0;
            background-color: transparent;
            border: none;
        }

        & button:hover {
            color: var(--pico-primary);
        }
    }

    .toggle {
        outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
    }

    .highlight {
        background-color: oklch(from var(--pico-contrast) l c h / 0.15);
    }

    @media (min-width: 768px) {
        .toggle {
            display: none;
        }
    }
</style>
