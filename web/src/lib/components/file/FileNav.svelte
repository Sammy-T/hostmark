<script>
    import icHome from '$lib/assets/home.svg?raw';
    import icPlus from '$lib/assets/plus.svg?raw';
    import icFolder from '$lib/assets/folder.svg?raw';
    import Sidebar from '../Sidebar.svelte';
    import AlertMessage from '../AlertMessage.svelte';
    import { getContext } from 'svelte';
    import { goto } from '$app/navigation';

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

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    $effect(() => {
        load(workingDir.value);
    });

    /**
     * @param {string} dir
     */
    async function load(dir) {
        const resp = await fetch(`/api/dir/${dir}`);

        switch(resp.status) {
            case 200:
                break;

            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        load(dir);
                        break;

                    case 400:
                    case 401:
                        goto('/login');
                        break;
                        
                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                }
                return;

            default:
                console.error('unable to load directory');
                
                errText = await resp.text();
                console.error(errText, resp.status);
            
                alertMsg.show();
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
</script>

{#snippet pathEntry(/** @type {string} */ href, /** @type {string} */ title, /** @type {string} */ icon = '')}
    <li>
        <a {href} onclick={onEntryClick}>
            {#if icon}{@html icon} {/if}{title}
        </a>
    </li>
{/snippet}

<!-- Mobile file nav -->
<Sidebar popId="mobile-file-nav" mobileOnly>
    <ul data-sveltekit-preload-data="off">
        {@render pathEntry('/#[home]', 'home', icHome)}
        {@render pathEntry('/#[back]', '..')}
        {@render pathEntry('/#[new]', 'new file', icPlus)}

        {#each entries as entry}
            {@const type = (entry.isDir) ? '' : '/file'}
            {@const href = (workingDir.value) ? [type, workingDir.value, entry.name].join('/') : [type, entry.name].join('/')}
            {@const icon = (entry.isDir) ? icFolder : ''}

            {@render pathEntry(href, entry.name, icon)}
        {/each}
    </ul>
</Sidebar>

<!-- Desktop file nav -->
<aside>
    <nav>
        <ul data-sveltekit-preload-data="off">
            {@render pathEntry('/#[home]', 'home', icHome)}
            {@render pathEntry('/#[back]', '..')}
            {@render pathEntry('/#[new]', 'new file', icPlus)}

            {#each entries as entry}
                {@const type = (entry.isDir) ? '' : '/file'}
                {@const href = (workingDir.value) ? [type, workingDir.value, entry.name].join('/') : [type, entry.name].join('/')}
                {@const icon = (entry.isDir) ? icFolder : ''}

                {@render pathEntry(href, entry.name, icon)}
            {/each}
        </ul>
    </nav>
</aside>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

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

    li a {
        display: flex;
        align-items: center;
        gap: 0.25rem;

        & :global(svg) {
            width: 1rem;
        }
    }

    @media (min-width: 768px) {
        aside {
            display: revert-layer;
        }
    }
</style>
