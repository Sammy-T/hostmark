<script>
    import Sidebar from '../Sidebar.svelte';
    import { cbLibPlus, showFolderSidebar } from '$lib/util.svelte';
    import { onMount } from 'svelte';

    function addFile() {
        console.log('add file');
    }

    onMount(() => {
        cbLibPlus.cb = addFile;
    });
</script>

<!-- Mobile file nav -->
{#if showFolderSidebar.value}
<Sidebar mobileOnly>
    <ul>
        <li><a href="#">..</a></li>
        
        {#each { length: 5 } as _, i}
            <li><a href="#">file-{i + 1}</a></li>
        {/each}
    </ul>
</Sidebar>
{/if}

<aside>
    <nav>
        <ul>
            <li><a href="#">..</a></li>

            {#each { length: 5 } as _, i}
                <li><a href="#">file-{i + 1}</a></li>
            {/each}
        </ul>
    </nav>
</aside>

<style>
    aside {
        display: none;
        height: 100dvh;
        width: 15dvw;
        overflow: auto;
        outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);

        & *, & nav ul:first-of-type, & nav ul:last-of-type {
            padding: 0;
            margin: 0;
        }

        & a {
            padding: 0.25rem 0.5rem;
        }

        & a:hover {
            background-color: oklch(from var(--pico-contrast) l c h / 0.05);
        }
    }

    @media (min-width: 768px) {
        aside {
            display: revert-layer;
        }
    }
</style>
