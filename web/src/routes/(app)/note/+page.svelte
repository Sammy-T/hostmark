<script>
    import NoteNav from '$lib/components/note/NoteNav.svelte';
    import NoteBrowser from '$lib/components/note/NoteBrowser.svelte';
    import { SvelteSet } from 'svelte/reactivity';
    import { onMount, setContext } from 'svelte';

    let notes = $state({ value: [] });
    setContext('notes', notes);

    /** @type {SvelteSet<string>} */
    let tags = new SvelteSet();
    setContext('tags', tags);

    /** @type {SvelteSet<string>} */
    let selectedTags = new SvelteSet();
    setContext('selectedTags', selectedTags);

    $effect(() => {
        console.log($state.snapshot(selectedTags));
        if(selectedTags) loadNotes();
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
        await fetch('/api/auth/refresh');

        const url = new URL('/api/note/list', location.origin);

        selectedTags.forEach((tag) => {
            url.searchParams.append('tags', tag);
        });

        let resp = await fetch(url);
        if(!resp.ok) return;

        const data = await resp.json();

        notes.value = data;
    }

    onMount(() => {
        loadTags(); //// TODO: Reload on note submit
    });
</script>

<div>
    <NoteNav />
    <NoteBrowser />
</div>

<style>
    div {
        display: flex;
    }
</style>
