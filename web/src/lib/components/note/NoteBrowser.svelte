<script>
    import NoteHeader from './NoteHeader.svelte';
    import Note from './Note.svelte';
    import { page } from '$app/state';

    let notes = $derived(page.data.notes);

    $effect(() => {
        console.log($state.snapshot(notes));
    });
</script>

<section>
    <NoteHeader />

    {#each notes as note (note.id)}
        <Note {note} />
    {:else}
        <div class="empty-view">
            <h3>No notes found.</h3>
        </div>
    {/each}
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
