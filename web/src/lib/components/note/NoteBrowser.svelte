<script>
    import NoteHeader from './NoteHeader.svelte';
    import Note from './Note.svelte';
</script>

<section>
    <NoteHeader />

    {#each { length: 50 } as _, i }
    <Note content={`note-${i+1}`} />
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
    }
</style>
