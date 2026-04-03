<script>
    import NoteHeader from './NoteHeader.svelte';
    import Note from './Note.svelte';
</script>

<section>
    <NoteHeader />

    {#each { length: 50 } as _, i }
    <Note content={`note-${i+1}`} />
    {/each}

    <!-- <div class="empty-view">
        <h3>No notes found.</h3>
    </div> -->
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

        /* & .empty-view {
            flex-grow: 1;
            align-content: center;
            text-align: center;
        } */
    }
</style>
