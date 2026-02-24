<script>
    import icMenu from '$lib/assets/menu-2.svg?raw';
    import icClose from '$lib/assets/square-x.svg?raw';
    import icFiles from '$lib/assets/book-2.svg?raw';
    import icProfile from '$lib/assets/user-circle.svg?raw';
    import icSettings from '$lib/assets/settings.svg?raw';
    import icSignout from '$lib/assets/logout.svg?raw';
    import { fade } from 'svelte/transition';

    let showMobile = $state(false);

    function toggleMobile() {
        showMobile = !showMobile;
    }
</script>

<nav class="toggle">
    <button onclick={toggleMobile}>{@html icMenu}</button>
</nav>

{#if showMobile}
<aside class="mobile" transition:fade>
    <nav>
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
    </nav>
</aside>
{/if}

<aside class="sidebar" transition:fade>
    <nav>
        <ul>
            <li>hm</li>
            <li><a href="#" data-tooltip="Files" data-placement="right">{@html icFiles}</a></li>
            <li><a href="#" data-tooltip="Profile" data-placement="right">{@html icProfile}</a></li>
        </ul>

        <ul>
            <li><a href="#" data-tooltip="Settings" data-placement="right">{@html icSettings}</a></li>
            <li><a href="#" data-tooltip="Sign out" data-placement="right">{@html icSignout}</a></li>
        </ul>
    </nav>
</aside>

<style>
    button {
        padding: 0;
        margin: 0;
        background-color: transparent;
        border: none;
    }

    button:hover {
        color: var(--pico-primary);
    }

    a:hover {
        text-decoration: none;
    }

    nav {
        padding: 0.5rem;
        margin: 0;

        & *, & ul:first-of-type, & ul:last-of-type {
            padding: 0;
            margin: 0;
        }

        & li {
            padding: 0.5rem 0;
        }
    }

    .toggle, .sidebar {
        outline: 1px solid rgba(128, 128, 128, 0.485);
    }

    .sidebar {
        display: none;
    }

    aside {
        height: 100dvh;

        & nav {
            min-height: 100%;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
        }
    }

    aside.mobile {
        min-width: 10rem;
        position: absolute;
        top: 0;
        left: 0;
        background-color: oklch(from var(--pico-background-color) calc(l * 0.85) c h);
    }

    @media (min-width: 768px) {
        .toggle, .mobile {
            display: none;
        }

        .sidebar {
            display: revert-layer;
        }
    }
</style>
