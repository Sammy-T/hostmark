<script>
    import Note from '$lib/components/note/Note.svelte';
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    let noteId = $state('');

    let note = $state();

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    async function loadNote() {
        let resp = await fetch(`/api/note/${noteId}`);

        switch(resp.status) {
            case 200:
                break;

            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        resp = await fetch(`/api/note/${noteId}`);
                        if(!resp.ok) {
                            errText = await resp.text();
                            console.error(errText, resp.status);

                            alertMsg.show();
                            return;
                        }
                        break;

                    case 400:
                    case 401:
                        goto('/login');
                        return;

                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                        return
                }
                break;

            default:
                errText = await resp.text();
                console.error(errText, resp.status);

                alertMsg.show();
                return;
        }

        note = await resp.json();
    }

    onMount(() => {
        noteId = location.href.split('/note/').at(-1) ?? '';
        loadNote();
    });
</script>

<div class="page">
    <div class="note-container">
        {#if note}
            <Note {note} />
        {/if}
    </div>
</div>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

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
