<script>
    import icMenu from '$lib/assets/menu-2.svg?raw';
    import icClose from '$lib/assets/square-x.svg?raw';
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
            <li><a href="#">Files</a></li>
            <li><a href="#">Profile</a></li>
        </ul>

        <ul>
            <li><a href="#">Settings</a></li>
            <li><a href="#">Sign out</a></li>
        </ul>
    </nav>
</aside>
{/if}

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

    .toggle {
        outline: 1px solid rgba(128, 128, 128, 0.485);
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
</style>
