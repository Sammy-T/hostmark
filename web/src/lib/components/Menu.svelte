<script>
    import icMenu from '$lib/assets/menu-2.svg?raw';
    import icClose from '$lib/assets/square-x.svg?raw';
    import icFiles from '$lib/assets/book-2.svg?raw';
    import icNote from '$lib/assets/note.svg?raw';
    import icProfile from '$lib/assets/user-circle.svg?raw';
    import icSettings from '$lib/assets/settings.svg?raw';
    import icSignout from '$lib/assets/logout.svg?raw';
    import icFolders from '$lib/assets/folders.svg?raw';
    import icLibraryPlus from '$lib/assets/library-plus.svg?raw';
    import icFilter from '$lib/assets/filter-2.svg?raw';
    import Sidebar from './Sidebar.svelte';
    import { cbLibPlus } from '$lib/util.svelte';
    import { page } from '$app/state';

    let current = $derived((page.url.pathname === '/') ? 'file' : page.url.pathname.split('/').filter((v) => v !== '').at(0));

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

<!-- Mobile menu -->
<nav class="toggle">
    <button popovertarget="mobile-menu">{@html icMenu}</button>

    {#if current === 'file'}
    <!-- File page item(s) -->
    <ul>
        <li><button popovertarget="mobile-file-nav">{@html icFolders}</button></li>
        <li><button onclick={cbLibPlus.cb?.()}>{@html icLibraryPlus}</button></li>
    </ul>
    {:else if current === 'note'}
    <!-- Note page item(s) -->
    <ul>
        <li><button popovertarget="mobile-note-nav">{@html icFilter}</button></li>
    </ul>
    {/if}
</nav>

<!-- Mobile page nav -->
<Sidebar popId="mobile-menu" mobileOnly>
    <ul>
        <li><button popovertarget="mobile-menu">{@html icClose}</button></li>
        <li>hm</li>
        <li><a href="/">{@html icFiles} Files</a></li>
        <li><a href="/note">{@html icNote} Notes</a></li>
        <li><a href="/profile">{@html icProfile} Profile</a></li>
    </ul>
    
    <ul>
        <li><a href="#">{@html icSettings} Settings</a></li>
        <li><a href="#">{@html icSignout} Sign out</a></li>
    </ul>
</Sidebar>

<!-- Desktop page nav -->
<Sidebar desktopOnly>
    <ul>
        <li>hm</li>
        <li><a href="/" class:highlight={current === 'file'} use:delayedTip={{ title: 'Files' }}>{@html icFiles}</a></li>
        <li><a href="/note" class:highlight={current === 'note'} use:delayedTip={{ title: 'Notes' }}>{@html icNote}</a></li>
        <li><a href="/profile" class:highlight={current === 'profile'} use:delayedTip={{ title: 'Profile' }}>{@html icProfile}</a></li>
    </ul>
    
    <ul>
        <li><a href="#" class:highlight={current === 'settings'} use:delayedTip={{ title: 'Settings' }}>{@html icSettings}</a></li>
        <li><a href="#" class:highlight={current === 'signout'} use:delayedTip={{ title: 'Sign out' }}>{@html icSignout}</a></li>
    </ul>
</Sidebar>

<style>
    .toggle {
        padding: 0 0.5rem;
        outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);

        & button {
            padding: 0;
            margin: 0;
            background-color: transparent;
            border: none;
        }

        & button:hover {
            color: var(--pico-primary);
        }

        & > button {
            padding: 1rem 0.25rem;
        }
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
