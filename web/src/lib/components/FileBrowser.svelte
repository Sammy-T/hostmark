<script>
    import FileNav from './file/FileNav.svelte';
    import FileHeader from './file/FileHeader.svelte';
    import FileEditor from './file/FileEditor.svelte';
    import { setContext } from 'svelte';
    import { page } from '$app/state';
    import { goto, refreshAll } from '$app/navigation';

    let file = $derived(page.params.file);
    let content = $derived(page.data.content);

    let directory = $state({ value: '' });
    setContext('directory', directory);

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
        editedFile = (directory.value.length > 0) ? directory.value + '/' : directory.value;
        edited = '';

        addingFile = true;
        naming.value = true;
        editing.value = true;
    }

    /**
     * @param {String} path 
     * @param {'PUT' | 'PATCH' | 'DELETE'} method
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
        if(editedFile !== file) {
            const created = await requestChange(`/api/file/${editedFile}`, 'PUT', edited);
            if(!created) return;
            
            if(addingFile) {
                addingFile = false;
            } else {
                await requestChange(`/api/file/${file}`, 'DELETE');
            }

            goto(`/file/${editedFile}`);
        } else {
            await requestChange(`/api/file/${file}`, 'PATCH', edited);

            refreshAll();
        }
    }

    function finish() {
        naming.value = false;
        editing.value = false;
    }
</script>

<div class="browser">
    <FileNav onaddfile={addFile} />

    <div class="file-view">
        {#if file}
            <FileHeader bind:editedFile />
        {/if}

        <FileEditor html={content?.html} bind:edited />

        {#if naming.value || editing.value}
            <button class="secondary" onclick={finish}>finish</button>
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
