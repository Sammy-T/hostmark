<script>
    import LoadingView from '$lib/components/LoadingView.svelte';
    import FileView from '$lib/components/FileView.svelte';
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { marked } from 'marked';
    import markedAlert from 'marked-alert';

    marked.use({ gfm: true }, markedAlert());

    let markdown = $state('');
    let html = $state('');

    let loading = $state(true);

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    async function loadReadme() {
        loading = true;

        let resp = await fetch(`/api/file/readme.md`);

        switch(resp.status) {
            case 200:
                break;
            
            case 401:
                const refResp = await fetch('/api/auth/refresh');
            
                switch(refResp.status) {
                    case 200:
                        loadReadme();
                        break;
                    
                    case 400:
                    case 401:
                        goto('/login');
                        break;
                    
                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                }
                return;
            
            default:
                errText = await resp.text();
                console.error(errText, resp.status);

                alertMsg.show();
                return;
        }

        markdown = await resp.text();
        html = await marked.parse(markdown);

        loading = false;
    }

    onMount(() => {
        loadReadme();
    });
</script>

{#if loading}
    <LoadingView />
{:else}
    <FileView file="readme.md" {markdown} {html} />
{/if}

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>
