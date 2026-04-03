<script>
    import icTags from '$lib/assets/tags.svg?raw';
    import icPlus from '$lib/assets/plus.svg?raw';
    import AlertMessage from '../AlertMessage.svelte';
    import { goto, invalidateAll } from '$app/navigation';
    import { onMount } from 'svelte';
    import { SvelteSet } from 'svelte/reactivity';

    let value = $state('');

    /** @type {SvelteSet<string>}*/
    let tags = new SvelteSet();

    /** @type {HTMLFormElement} */
    let form;

    /** @type {HTMLFormElement} */
    let tagForm;

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    async function loadTags() {
        const resp = await fetch('/api/tags');
        if(!resp.ok) return;

        const respJson = await resp.json();

        tags.clear();
        
        // @ts-ignore
        respJson.forEach((tag) => tags.add(tag.name));
    }

    /**
     * @param {SubmitEvent} event
     */
    async function onTagSubmit(event) {
        event.preventDefault();

        // @ts-ignore
        const data = new FormData(event.target);

        const tag = data.get('tag');

        if(tag) tags.add(tag.toString());

        /** @type {HTMLElement} */
        // @ts-ignore
        const tagPopover = document.querySelector('#tag-input');
        tagPopover.hidePopover();

        tagForm.reset();
    }

    /**
     * @param {SubmitEvent} event
     */
    async function onSubmit(event) {
        event.preventDefault();

        // @ts-ignore
        const data = new FormData(event.target);

        // @ts-ignore
        const encoded = new URLSearchParams(data).toString();

        /** @type {RequestInit} */
        const opts = {
            method: 'post',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: encoded,
        };

        let resp = await fetch('/api/note', opts);

        switch(resp.status) {
            case 200:
                break;
            
            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        onSubmit(event);
                        return;

                    case 400:
                    case 401:
                        goto('/login');
                        return;

                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                        return;
                }
            
            default:
                console.error('unable to submit note');
                return;
        }

        form.reset();
        invalidateAll();
    }

    onMount(() => {
        loadTags();
    });
</script>

{#snippet tagItem(/** @type {string} */ tag)}
<label>
    <input type="checkbox" name="tags" value={tag} />
    {tag}
</label>
{/snippet}

<header>
    <form onsubmit={onSubmit} bind:this={form}>
        <textarea name="content" placeholder="New Note..." autocapitalize="on" spellcheck required bind:value></textarea>

        <fieldset class="footer" disabled={value === ''}>
            <button type="button" class="img" popovertarget="tag-select">{@html icTags}</button>

            <div class="group">
                <select name="visibility">
                    <option value="public">Public</option>
                    <option value="protected">Protected</option>
                    <option value="private" selected>Private</option>
                </select>

                <button type="submit" class="secondary">Save</button>
            </div>
        </fieldset>

        <div id="tag-select" popover>
            <div class="container">
                <fieldset>
                    <legend>Tags <button type="button" class="secondary" popovertarget="tag-input">{@html icPlus}</button></legend>

                    {#each tags.values() as tag}
                        {@render tagItem(tag)}
                    {/each}
                </fieldset>
            </div>
        </div>
    </form>

    <div id="tag-input" popover>
        <div class="container">
            <form onsubmit={onTagSubmit} bind:this={tagForm}>
                <!-- svelte-ignore a11y_no_redundant_roles -->
                <fieldset role="group">
                    <input type="text" name="tag" placeholder="Tag name" />
                    <button type="submit">Add</button>
                </fieldset>
            </form>
        </div>
    </div>
</header>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    [popover] {
        margin: 0;
        inset: auto;
        border: none;
        top: calc(anchor(bottom) + 0.25rem);
        left: anchor(left);
        border: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        border-radius: 0.25rem;
        background-color: var(--pico-dropdown-background-color);
        box-shadow: var(--pico-dropdown-box-shadow);
    }

    #tag-select .container {
        min-width: 10dvw;
        padding: 0.5rem;

        & legend {
            width: 100%;
            margin-bottom: 0.75rem;
            display: flex;
            justify-content: space-between;
            align-items: center;

            & button {
                padding: 0;
            }
        }

        & label {
            margin-bottom: 0.25rem;
        }
    }

    header {
        position: sticky;
        top: 0;
        padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
        background-color: var(--pico-background-color);
        border: 2px solid oklch(from var(--pico-muted-color) l c h / 0.25);
        border-radius: var(--pico-border-radius);
        background-color: oklch(from var(--pico-background-color) l c h / 0.5);
        backdrop-filter: blur(10px);

        & form * {
            margin: 0;
        }

        & form {
            display: grid;
            gap: 0.5rem;
        }

        & fieldset.footer {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        & .group {
            & > * {
                width: auto;
            }

            & > *:not(:last-child) {
                margin-right: 0.25rem;
            }
        }

        & select {
            padding: 0.25rem 0.5rem;
            padding-right: calc(var(--pico-form-element-spacing-horizontal) + 0.75rem);
            background-position: center right 0.25rem;
        }

        & button {
            width: auto;
            padding: 0.25rem 0.5rem;
        }

        & button.img {
            background-color: transparent;
            border: none;
        }

        & button.img:hover {
            background-color: var(--pico-dropdown-hover-background-color);
        }
    }
</style>
