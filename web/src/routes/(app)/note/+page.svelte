<script>
    import NoteNav from '$lib/components/note/NoteNav.svelte';
    import NoteBrowser from '$lib/components/note/NoteBrowser.svelte';
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { SvelteSet } from 'svelte/reactivity';
    import { onMount, setContext } from 'svelte';
    import { goto } from '$app/navigation';
    import { PAGE_SIZE, STORAGE_WAS_AUTHED_KEY } from '$lib/util.svelte';

    let notes = $state({ value: [] });
    setContext('notes', notes);

    /** @type {SvelteSet<string>} */
    let tags = new SvelteSet();
    setContext('tags', tags);

    /** @type {SvelteSet<string>} */
    let selectedTags = new SvelteSet();
    setContext('selectedTags', selectedTags);

    let showLoadNext = $state({ value: false });
    setContext('showLoadNext', showLoadNext);

    let pagesLoaded = $state({ value: 0 });
    setContext('pagesLoaded', pagesLoaded);

    setContext('loadTags', loadTags);
    setContext('loadNotes', loadNotes);

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    $effect(() => {
        if(selectedTags.values()) loadNotes();
    });

    async function loadTags() {
        const resp = await fetch('/api/tags');
        if(!resp.ok) return;

        const respJson = await resp.json();

        tags.clear();
        
        // @ts-ignore
        respJson.forEach((tag) => tags.add(tag.name));
    }

    async function loadNotes(page = 0) {
        const url = new URL('/api/note/list', location.origin);

        url.searchParams.append('page_size', PAGE_SIZE.toString());
        if(page > 0) url.searchParams.append('page', page.toString());

        selectedTags.forEach((tag) => {
            url.searchParams.append('tags', tag);
        });

        let resp = await fetch(url);
        if(!resp.ok) return;

        const { notes: respNotes, authed } = await resp.json();

        const wasAuthed = sessionStorage.getItem(STORAGE_WAS_AUTHED_KEY) === 'true';

        if(!authed && wasAuthed) {
            sessionStorage.setItem(STORAGE_WAS_AUTHED_KEY, authed.toString());

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
        sessionStorage.setItem(STORAGE_WAS_AUTHED_KEY, authed.toString());

        if(respNotes.length > 0) {
            pagesLoaded.value = (page > 0) ? page : 1;
        }

        showLoadNext.value = respNotes.length === PAGE_SIZE;
    }

    onMount(() => {
        loadTags();
    });
</script>

<div class="page">
    <NoteNav />
    <NoteBrowser />
</div>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>
