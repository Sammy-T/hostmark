<script>
    import Note from '$lib/components/note/Note.svelte';
    import { page } from '$app/state';
</script>

<div class="page">
    <div class="note-container">
        <Note note={page.data.note} />
    </div>
</div>

<style>
    .page {
        height: 100%;
        padding: 1rem;
        overflow: auto;
        justify-content: center;
    }

    .note-container {
        width: min(100%, 1000px);
    }
</style>
