<script>
    import icPublic from '$lib/assets/world-map.svg?raw';
    import icProtected from '$lib/assets/users.svg?raw';
    import icPrivate from '$lib/assets/lock.svg?raw';
    import icMore from '$lib/assets/dots-vertical.svg?raw';
    import icEdit from '$lib/assets/edit.svg?raw';
    import icTrash from '$lib/assets/trash.svg?raw';
    import icTag from '$lib/assets/tag.svg?raw';
    import icTagFilled from '$lib/assets/tag-filled.svg?raw';
    import AlertMessage from '../AlertMessage.svelte';
    import { marked } from 'marked';
    import { getContext } from 'svelte';
    import { goto } from '$app/navigation';

    let { note } = $props();

    /** @type {import('svelte/reactivity').SvelteSet<string>} */
    let selectedTags = getContext('selectedTags');

    /** @type {function} */
    let loadNotes = getContext('loadNotes');

    let menuId = $derived(`note-menu-${note.id}`);

    let editing = $state(false);

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    /**
     * @param {SubmitEvent} event
     */
    async function onSubmitEdit(event) {
        event.preventDefault();

        // @ts-ignore
        const data = new FormData(event.target);
        const content = data.get('content');
        
        if(content === note.content) {
            editing = false;
            return;
        }

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

        const resp = await fetch(`/api/note/${note.id}`, opts);

        switch(resp.status) {
            case 200:
                loadNotes();
                break;

            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        onSubmitEdit(event);
                        return;

                    case 400:
                    case 401:
                        goto('/login');
                        return;

                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                        editing = false;
                        return;
                }
            
            default:
                errText = await resp.text();
                console.error(errText, resp.status);

                alertMsg.show();
        }

        editing = false;
    }

    async function deleteNote() {
        /** @type {RequestInit} */
        const opts = {
            method: 'delete',
        };

        const resp = await fetch(`/api/note/${note.id}`, opts);

        switch(resp.status) {
            case 200:
                loadNotes();
                break;

            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        deleteNote();
                        return;

                    case 400:
                    case 401:
                        goto('/login');
                        return;

                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                        editing = false;
                        return;
                }
            
            default:
                errText = await resp.text();
                console.error(errText, resp.status);

                alertMsg.show();
        }

    }

    function toggleEditing() {
        editing = !editing;
    }

    /**
     * @param {Event} event
     */
    function tagClicked(event) {
        event.preventDefault();

        // @ts-ignore
        const tag = event.target?.dataset.tag;

        if(selectedTags.has(tag)) {
            selectedTags.delete(tag);
            return;
        }

        selectedTags.add(tag);
    }
</script>

<article>
    <header>
        <small>{new Intl.DateTimeFormat().format(new Date(note.created_at))}</small>
        
        <div class="side">
            <small>{note.owner}</small>

            {#if note.visibility === 'public'}
                {@html icPublic}
            {:else if note.visibility === 'protected'}
                {@html icProtected}
            {:else if note.visibility === 'private'}
                {@html icPrivate}
            {/if}
            
            <button popovertarget={menuId}>{@html icMore}</button>

            <div id={menuId} popover>
                <div class="menu-container">
                    <button popovertarget={menuId} onclick={toggleEditing}>{@html icEdit} Edit</button>
                    <button popovertarget={menuId} onclick={deleteNote}>{@html icTrash} Delete</button>
                </div>
            </div>
        </div>
    </header>

    {#if editing}
        <form onsubmit={onSubmitEdit}>
            <textarea name="content">{note.content}</textarea>
            <button type="submit">Save</button>
        </form>
    {:else}
        {#await marked.parse(note.content) then parsed}
            {@html parsed}
        {/await}
    {/if}

    {#if note.tags.length > 0}
        <hr class="separator">
    {/if}
    
    <div class="tags">
        {#each note.tags as tag (tag.name)}
            {#if selectedTags}
                <a href={`#${tag}`} data-tag={tag.name} onclick={tagClicked}>
                    {#if selectedTags.has(tag.name)}
                        {@html icTagFilled}
                    {:else}
                        {@html icTag}
                    {/if}

                    {tag.name}
                </a>
            {:else}
                <span>{@html icTag} {tag.name}</span>
            {/if}
        {/each}
    </div>
</article>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    [popover] {
        margin: 0;
        inset: auto;
        border: none;
        top: calc(anchor(bottom) + 0.25rem);
        right: anchor(right);
        border: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        border-radius: 0.25rem;
        background-color: var(--pico-dropdown-background-color);
        box-shadow: var(--pico-dropdown-box-shadow);
    }

    .menu-container {
        display: flex;
        flex-direction: column;

        & button {
            padding: 0 0.5rem;
            display: flex;
            align-items: center;
            gap: 0.25rem;
            color: var(--pico-dropdown-color);

            & :global(svg) {
                width: 1rem;
            }
        }

        & button:hover {
            background-color: var(--pico-dropdown-hover-background-color);
        }
    }

    article {
        margin: 0;

        & header {
            display: flex;
            justify-content: space-between;
            align-items: center;

            & > small {
                color: var(--pico-muted-color);
            }

            & button {
                padding: 0;
                background-color: transparent;
                border: none;
                color: oklch(from var(--pico-contrast) l c h / 0.5);
            }

            .side > :global(svg) {
                color: oklch(from var(--pico-contrast) l c h / 0.5);
            }

            .side > small {
                margin-right: 0.5rem;
            }
        }
    }

    form {
        padding-bottom: 0.25rem;
        display: flex;
        flex-direction: column;
        align-items: end;
        gap: 0.5rem;

        & > * {
            margin: 0;
        }

        & button {
            width: auto;
            padding: 0.25rem 0.5rem;
        }
    }

    .separator {
        margin: 0.25rem 0;
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
    }
</style>
