<script>
    import Sidebar from '../Sidebar.svelte';
    import { cbLibPlus, showFolderSidebar } from '$lib/util.svelte';
    import { getContext, onMount } from 'svelte';

    /**
     * @callback onAddFile
     */

    /**
     * @typedef {Object} Props
     * @property {onAddFile} [onaddfile]
     */

    /** @type {Props} */
    let { onaddfile } = $props();

    /** @type {{ value: string }} */
    let workingDir = getContext('workingDir');

    // @ts-ignore
    let entries = $state([]);

    $effect(() => {
        load(workingDir.value);
    });

    /**
     * @param {string} dir
     */
    async function load(dir) {
        const resp = await fetch(`/api/dir/${dir}`);
        if(!resp.ok) {
            console.error('unable to load directory');
            return;
        }

        entries = await resp.json();
        console.log($state.snapshot(entries));
    }

    /**
     * @param {Event} ev
     */
    function onEntryClick(ev) {
        /** @type {HTMLAnchorElement} */
        // @ts-ignore
        const anchor = ev.target;

        const hrefPath = anchor.href.replace(/https?:\/\//, '').replace(`${location.host}/`, '')
        console.log(hrefPath);

        if(hrefPath.startsWith('file/')) return;

        ev.preventDefault();

        let newPath;

        switch(hrefPath) {
            case '#[home]':
                workingDir.value = '';
                break;

            case '#[back]':
                newPath = workingDir.value.split('/').slice(0, -1).filter(v => v !== '').join('/');
                workingDir.value = newPath;
                break;

            case '#[new]':
                onaddfile?.();
                break;
            default:
                newPath = [...workingDir.value.split('/').filter(v => v !== ''), hrefPath.split('/').at(-1)].join('/');
                workingDir.value = newPath;
        }
    }

    onMount(() => {
        if(onaddfile) cbLibPlus.cb = onaddfile;
    });
</script>

{#snippet pathEntry(/** @type {string} */ href, /** @type {string} */ title)}
    <li><a {href} onclick={onEntryClick}>{title}</a></li>
{/snippet}

<!-- Mobile file nav -->
{#if showFolderSidebar.value}
<Sidebar mobileOnly>
    <ul data-sveltekit-preload-data="off">
        {@render pathEntry('/#[home]', '[home]')}
        {@render pathEntry('/#[back]', '..')}
        {@render pathEntry('/#[new]', '+new file')}

        {#each entries as entry}
            {@const type = (entry.isDir) ? '' : '/file'}
            {@const href = (workingDir.value) ? [type, workingDir.value, entry.name].join('/') : [type, entry.name].join('/')}

            {@render pathEntry(href, `${entry.name}${entry.isDir ? '/' : ''}`)}
        {/each}
</ul>
</Sidebar>
{/if}

<aside>
    <nav>
        <ul data-sveltekit-preload-data="off">
            {@render pathEntry('/#[home]', '[home]')}
            {@render pathEntry('/#[back]', '..')}
            {@render pathEntry('/#[new]', '+new file')}

            {#each entries as entry}
                {@const type = (entry.isDir) ? '' : '/file'}
                {@const href = (workingDir.value) ? [type, workingDir.value, entry.name].join('/') : [type, entry.name].join('/')}

                {@render pathEntry(href, `${entry.name}${entry.isDir ? '/' : ''}`)}
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
