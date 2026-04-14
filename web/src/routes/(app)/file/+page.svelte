<script>
    import FileView from '$lib/components/FileView.svelte';
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { onMount } from 'svelte';
    import { goto, onNavigate } from '$app/navigation';
    import { marked } from 'marked';
    import markedAlert from 'marked-alert';

    marked.use({ gfm: true }, markedAlert());

    let file = $state('');
    let markdown = $state('');
    let html = $state('');

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    $effect(() => {
        if(file) loadFile();
    });

    async function loadFile() {
        let resp = await fetch(`/api/file/${file}`);

        switch(resp.status) {
            case 200:
                break;
            
            // case 400:
            //     return;
            
            case 401:
                const refResp = await fetch('/api/auth/refresh');
            
                switch(refResp.status) {
                    case 200:
                        resp = await fetch(`/api/file/${file}`);
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
                        return;
                }
                break;
            
            default:
                errText = await resp.text();
                console.error(errText, resp.status);

                alertMsg.show();
                return;
        }

        markdown = await resp.text();
        html = await marked.parse(markdown);
    }

    onNavigate(() => {
        file = location.href.split('/file/').at(-1) ?? '';
    });

    onMount(() => {
        file = location.href.split('/file/').at(-1) ?? '';
    });
</script>

<FileView {file} {markdown} {html} />

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>
