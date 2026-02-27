<script>
    import { setContext } from 'svelte';
    import FileEditor from './file/FileEditor.svelte';
    import FileNav from './file/FileNav.svelte';

    let directory = $state({ value: '' });
    setContext('directory', directory);
</script>

<div>
    <FileNav />
    <FileEditor />
</div>

<style>
    div {
        display: flex;
    }
</style>
