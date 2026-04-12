<script>
    import Note from '$lib/components/note/Note.svelte';
    import { page } from '$app/state';
</script>

<div>
    <Note note={page.data.note} />
</div>

<style>
    div {
        height: 100dvh;
        overflow: auto;
    }
</style>
