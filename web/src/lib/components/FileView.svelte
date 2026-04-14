<script>
    import FileNav from './file/FileNav.svelte';
    import FileHeader from './file/FileHeader.svelte';
    import FileEditor from './file/FileEditor.svelte';
    import FileNone from './file/FileNone.svelte';
    import AlertMessage from './AlertMessage.svelte';
    import { setContext } from 'svelte';
    import { goto, refreshAll } from '$app/navigation';

    /**
     * @typedef {Object} Props
     * @property {String} file
     * @property {String} markdown
     * @property {String} html
     */

    /** @type {Props} */
    let { file, markdown, html } = $props();

    // svelte-ignore state_referenced_locally
    let workingDir = $state({ value: file?.split('/').filter((p) => !p.match(/\.\w+$/)).join('/') ?? '' });
    setContext('workingDir', workingDir);

    let naming = $state({ value: false });
    setContext('naming', naming);

    let editing = $state({ value: false });
    setContext('editing', editing);

    let editedFile = $derived(file);
    let edited = $derived(markdown);
    let addingFile = $state(false);

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    $effect(() => {
        if(naming.value || editing.value || !editedFile || (editedFile === file && edited === markdown)) return;
        submitChanges();
    });

    function addFile() {
        editedFile = (workingDir.value.length > 0) ? workingDir.value + '/' : workingDir.value;
        edited = '';

        addingFile = true;
        naming.value = true;
        editing.value = true;
    }

    async function deleteFile() {
        const success = await requestChange(`/api/file/${file}`, 'DELETE');
        if(!success) return;

        goto(`/file/${workingDir.value}`);
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

        switch(resp.status) {
            case 200:
                break;

            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        return await requestChange(path, method, body);

                    case 400:
                    case 401:
                        goto('/login');
                        return false;

                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                        return false;
                }

            default:
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

<div class="page">
    {#key file}
        <FileNav onaddfile={addFile} />
    {/key}

    <div class="file-view">
        {#if html}
            {#if file || naming.value}
                <FileHeader {file} bind:editedFile ondeletefile={deleteFile} />
            {/if}
            
            <FileEditor {html} bind:edited />

            {#if naming.value || editing.value}
                <button class="secondary" onclick={finish}>finish</button>
            {/if}
        {:else}
            <FileNone />
        {/if}
    </div>
</div>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    .file-view {
        flex-grow: 1;
        display: flex;
        flex-direction: column;
    }
</style>
