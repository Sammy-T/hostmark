<script>
    import { setContext } from 'svelte';
    import FileEditor from './file/FileEditor.svelte';
    import FileNav from './file/FileNav.svelte';

    let directory = $state({ value: '' });
    setContext('directory', directory);
</script>

<div class="browser">
    <FileNav />

    <div class="file-view">
        <header>todo/breadcrumb/path</header>
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

        & header {
            padding: 0.5rem;
            outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        }
    }
</style>
