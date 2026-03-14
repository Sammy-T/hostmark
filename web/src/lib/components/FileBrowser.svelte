<script>
    import FileNav from './file/FileNav.svelte';
    import FileHeader from './file/FileHeader.svelte';
    import FileEditor from './file/FileEditor.svelte';
    import { setContext } from 'svelte';
    import { page } from '$app/state';
    import { refreshAll } from '$app/navigation';

    let file = $derived(page.params.file);
    let content = $derived(page.data.content);

    let directory = $state({ value: '' });
    setContext('directory', directory);

    let naming = $state({ value: false });
    setContext('naming', naming);

    let editing = $state({ value: false });
    setContext('editing', editing);

    let edited = $derived(content?.markdown);

    $effect(() => {
        if(!editing.value && content?.markdown !== edited) submitChanges();
    });

    async function submitChanges() {
        /** @type {RequestInit} */
        const opts = {
            method: 'POST',
            body: edited,
        };

        const resp = await fetch(`/api/file/${page.params.file}`, opts);
        if(!resp.ok) {
            console.error('unable to post changes');
            return;
        }

        refreshAll();
    }

    function finish() {
        editing.value = false;
    }
</script>

<div class="browser">
    <FileNav />

    <div class="file-view">
        {#if file}
            <FileHeader />
        {/if}

        <FileEditor html={content?.html} bind:edited />

        {#if editing.value}
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
