<script>
    import FileNav from './file/FileNav.svelte';
    import FileHeader from './file/FileHeader.svelte';
    import FileEditor from './file/FileEditor.svelte';
    import { setContext } from 'svelte';
    import { page } from '$app/state';

    let file = $derived(page.params.file);

    let directory = $state({ value: '' });
    setContext('directory', directory);

    let editing = $state({ value: false });
    setContext('editing', editing);
</script>

<div class="browser">
    <FileNav />

    <div class="file-view">
        {#if file}
            <FileHeader />
        {/if}

        <FileEditor />
    </div>
</div>

<style>
    .browser {
        display: flex;
    }

    .file-view {
        height: 100dvh;
        flex-grow: 1;
        display: flex;
        flex-direction: column;
    }
</style>
