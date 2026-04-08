<script>
    import NoteHeader from './NoteHeader.svelte';
    import Note from './Note.svelte';
    import NoteNone from './NoteNone.svelte';
    import LoadNext from '../LoadNext.svelte';
    import { getContext } from 'svelte';

    let notes = getContext('notes');

    /** @type {{ value: boolean }} */
    let showLoadNext = getContext('showLoadNext');

    /** @type {{ value: number }} */
    let pagesLoaded = getContext('pagesLoaded');

    /** @type {function} */
    let loadNotes = getContext('loadNotes');

    let loading = false;

    async function loadNext() {
        if(loading) return;

        loading = true;

        const nextPage = pagesLoaded.value + 1;
        await loadNotes(nextPage);

        loading = false;
    }
</script>

<section>
    <NoteHeader />

    {#each notes.value as note (note.id)}
        <Note {note} />
    {:else}
        <NoteNone />
    {/each}

    {#if showLoadNext.value}
        <LoadNext onenteredview={loadNext} />
    {/if}
</section>

<style>
    section {
        padding: 0.5rem;
        margin: 0;
        flex-grow: 1;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        overflow: auto;
    }
</style>
