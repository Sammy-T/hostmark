<script>
    import NoteNav from '$lib/components/note/NoteNav.svelte';
    import NoteBrowser from '$lib/components/note/NoteBrowser.svelte';
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { SvelteSet } from 'svelte/reactivity';
    import { onMount, setContext } from 'svelte';
    import { goto } from '$app/navigation';

    let notes = $state({ value: [] });
    setContext('notes', notes);

    /** @type {SvelteSet<string>} */
    let tags = new SvelteSet();
    setContext('tags', tags);

    /** @type {SvelteSet<string>} */
    let selectedTags = new SvelteSet();
    setContext('selectedTags', selectedTags);

    let wasAuthed = false;

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    $effect(() => {
        if(selectedTags.values()) loadNotes();
    });

    async function loadTags() {
        const resp = await fetch('/api/tags');
        if(!resp.ok) return;

        const respJson = await resp.json();`/api/note/list`

        tags.clear();
        
        // @ts-ignore
        respJson.forEach((tag) => tags.add(tag.name));
    }

    async function loadNotes() {
        const url = new URL('/api/note/list', location.origin);

        selectedTags.forEach((tag) => {
            url.searchParams.append('tags', tag);
        });

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

        wasAuthed = authed;
        notes.value = respNotes;
    }

    onMount(() => {
        loadTags(); //// TODO: Reload on note submit
    });
</script>

<div>
    <NoteNav />
    <NoteBrowser />
</div>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    div {
        display: flex;
    }
</style>
