<script>
    import FileNav from './file/FileNav.svelte';
    import FileHeader from './file/FileHeader.svelte';
    import FileEditor from './file/FileEditor.svelte';
    import FileNone from './file/FileNone.svelte';
    import { setContext } from 'svelte';
    import { page } from '$app/state';
    import { goto, refreshAll } from '$app/navigation';

    let file = $derived(page.params.file);
    let content = $derived(page.data.content);

    // svelte-ignore state_referenced_locally
    let workingDir = $state({ value: file?.split('/').filter((p) => !p.match(/\.\w+$/)).join('/') ?? '' });
    setContext('workingDir', workingDir);

    let naming = $state({ value: false });
    setContext('naming', naming);

    let editing = $state({ value: false });
    setContext('editing', editing);

    let editedFile = $derived(file);
    let edited = $derived(content?.markdown);
    let addingFile = $state(false);

    $effect(() => {
        if(naming.value || editing.value || (editedFile === file && edited === content?.markdown)) return;
        submitChanges();
    });

    function addFile() {
        editedFile = (workingDir.value.length > 0) ? workingDir.value + '/' : workingDir.value;
        edited = '';

        addingFile = true;
        naming.value = true;
        editing.value = true;
    }

    /**
     * @param {String} path 
     * @param {'POST' | 'DELETE'} method
     * @param {String?} body
     */
    async function requestChange(path, method, body = null) {
        /** @type {RequestInit} */
        const opts = {
            method,
        };

        if(body) opts.body = body;

        const resp = await fetch(path, opts);
        if(!resp.ok) {
            console.error('unable to submit changes');
            return false;
        }

        return true;
    }

    async function submitChanges() {
        const success = await requestChange(`/api/file/${editedFile}`, 'POST', edited);
        if(!success) return;

        if(editedFile !== file) {
            if(!addingFile) await requestChange(`/api/file/${file}`, 'DELETE');

            addingFile = false;
            
            goto(`/file/${editedFile}`);
            return;
        }

        refreshAll();
    }

    function finish() {
        naming.value = false;
        editing.value = false;
    }
</script>

<div class="browser">
    {#key file}
        <FileNav onaddfile={addFile} />
    {/key}

    <div class="file-view">
        {#if content}
            {#if file}
                <FileHeader bind:editedFile />
            {/if}
            
            <FileEditor html={content?.html} bind:edited />

            {#if naming.value || editing.value}
                <button class="secondary" onclick={finish}>finish</button>
            {/if}
        {:else}
            <FileNone />
        {/if}
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
