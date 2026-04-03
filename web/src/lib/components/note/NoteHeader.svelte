<script>
    import AlertMessage from '../AlertMessage.svelte';
    import { goto, invalidateAll } from '$app/navigation';

    let value = $state('');

    /** @type {HTMLFormElement} */
    let form;

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    /**
     * @param {SubmitEvent} event
     */
    async function onSubmit(event) {
        event.preventDefault();

        // @ts-ignore
        const data = new FormData(event.target);

        // @ts-ignore
        const encoded = new URLSearchParams(data).toString();

        /** @type {RequestInit} */
        const opts = {
            method: 'post',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: encoded,
        };

        let resp = await fetch('/api/note', opts);

        switch(resp.status) {
            case 200:
                break;
            
            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        onSubmit(event);
                        return;

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
            
            default:
                console.error('unable to submit note');
                return;
        }

        form.reset();
        invalidateAll();
    }
</script>

<header>
    <form onsubmit={onSubmit} bind:this={form}>
        <textarea name="content" placeholder="New Note..." autocapitalize="on" spellcheck required bind:value></textarea>

        <fieldset disabled={value === ''}>
            <div>
                <select name="visibility">
                    <option value="public">Public</option>
                    <option value="protected">Protected</option>
                    <option value="private" selected>Private</option>
                </select>
            </div>

            <!-- TODO: tags -->

            <button type="submit" class="secondary">Save</button>
        </fieldset>
    </form>
</header>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    header {
        position: sticky;
        top: 0;
        padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
        background-color: var(--pico-background-color);
        border: 2px solid oklch(from var(--pico-muted-color) l c h / 0.25);
        border-radius: var(--pico-border-radius);
        background-color: oklch(from var(--pico-background-color) l c h / 0.5);
        backdrop-filter: blur(10px);

        & form * {
            margin: 0;
        }

        & form {
            display: grid;
            gap: 0.5rem;
        }

        & fieldset {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        & select {
            padding: 0.25rem 0.5rem;
            padding-right: calc(var(--pico-form-element-spacing-horizontal) + 0.75rem);
            background-position: center right 0.25rem;
        }

        & button {
            width: auto;
            padding: 0.25rem 0.5rem;
        }
    }
</style>
