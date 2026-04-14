<script>
    import icProfile from '$lib/assets/user-circle.svg?raw';
    import Note from '$lib/components/note/Note.svelte';
    import NoteNone from '$lib/components/note/NoteNone.svelte';
    import LoadNext from '$lib/components/LoadNext.svelte';
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { goto } from '$app/navigation';
    import { onMount, setContext } from 'svelte';
    import { PAGE_SIZE } from '$lib/util.svelte';

    let info = $state();

    /** @type {{ value: { id: number }[] }}*/
    let notes = $state({ value: [] });

    let showLoadNext = $state({ value: false });

    let pagesLoaded = $state({ value: 0 });

    let wasAuthed = false;

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    let loading = false;

    setContext('loadNotes', loadNotes);

    async function loadNext() {
        if(loading) return;

        loading = true;

        const nextPage = pagesLoaded.value + 1;
        await loadNotes(nextPage);

        loading = false;
    }

    async function loadNotes(page = 0) {
        const url = new URL('/api/note/list', location.origin);

        url.searchParams.set('username', info.username);
        url.searchParams.append('page_size', PAGE_SIZE.toString());
        if(page > 0) url.searchParams.append('page', page.toString());

        let resp = await fetch(url);
        if(!resp.ok) return;

        const { notes: respNotes, authed } = await resp.json();

        if(!authed && wasAuthed) {
            wasAuthed = authed;

            const refResp = await fetch('/api/auth/refresh');

            switch(refResp.status) {
                case 200:
                    loadNotes();
                    break;

                case 400:
                case 401:
                    goto('/login');
                    break;

                default:
                    errText = await refResp.text();
                    console.error(errText, refResp.status);

                    alertMsg.show();
                    break;
            }
            return;
        }

        notes.value = (page > 0) ? [...notes.value, ...respNotes] : respNotes;
        wasAuthed = authed;

        if(respNotes.length > 0) {
            pagesLoaded.value = (page > 0) ? page : 1;
        }

        showLoadNext.value = respNotes.length === PAGE_SIZE;
    }

    async function loadUser() {
        let resp = await fetch(`/api/account/me`);

        switch(resp.status) {
            case 200:
                break;

            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        loadUser();
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
                errText = await resp.text();
                console.error(errText, resp.status);

                alertMsg.show();
                return;
        }

        info = await resp.json();
        loadNotes();
    }

    onMount(() => {
        loadUser();
    });
</script>

<div class="page">
    <section class="profile">
        <div id="icon">
            {@html icProfile}
        </div>
        
        {#if info}
            <div id="info">
                <h4>{info.username}</h4>
                <p>{info.role}</p>
                <p class="muted">created {new Intl.DateTimeFormat().format(new Date(info.created_at))}</p>
            </div>
        {/if}
    </section>

    <section class="notes">
        {#each notes.value as note (note.id)}
            <Note {note} />
        {:else}
            <NoteNone />
        {/each}

        {#if showLoadNext.value}
            <LoadNext onenteredview={loadNext} />
        {/if}
    </section>
</div>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    .muted {
        color: var(--pico-muted-color);
    }

    .page {
        height: 100%;
        padding: 1rem;
        display: grid;
        grid-template-rows: auto 1fr;
        gap: 1rem;

        & section:not(.profile) {
            margin: 0;
        }
    }

    .profile {
        width: fit-content;
        margin: auto;
        display: grid;
        grid-template-columns: auto 1fr;
        gap: 0.5rem;
    }

    .notes {
        width: min(40rem, 100%);
        padding: 0 0.5rem;
        justify-self: center;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        overflow: auto;
    }

    #icon > :global(*) {
        width: 100%;
        height: 100%;
    }

    #info > * {
        margin: 0;
    }
</style>
