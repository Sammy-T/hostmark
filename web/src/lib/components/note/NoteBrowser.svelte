<script>
    import NoteHeader from './NoteHeader.svelte';
    import Note from './Note.svelte';
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
        <div class="empty-view">
            <h3>No notes found.</h3>
        </div>
    {/each}

    {#if showLoadNext.value}
        <LoadNext onenteredview={loadNext} />
    {/if}
</section>

<style>
    section {
        height: 100dvh;
        padding: 0.5rem;
        flex-grow: 1;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        overflow: auto;

        & .empty-view {
            flex-grow: 1;
            align-content: center;
            text-align: center;
        }
    }
</style>
