<script>
    import Note from '$lib/components/note/Note.svelte';
    import { page } from '$app/state';
</script>

<Note note={page.data.note} />
