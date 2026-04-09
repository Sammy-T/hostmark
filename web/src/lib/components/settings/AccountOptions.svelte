<script>
    import AlertMessage from '../AlertMessage.svelte';
    import { goto } from '$app/navigation';
    import { PREFS_PROFILE_KEY } from '$lib/util.svelte';

    let prefs = $state(JSON.parse(localStorage.getItem(PREFS_PROFILE_KEY) ?? '')?.prefs);

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    /** @type {Promise<Response>?} */
    let updateReq = $state(null);

    async function loadUser() {
        let resp = await fetch(`/api/account/me`);
        if(!resp.ok) return;

        const info = await resp.text();
        localStorage.setItem(PREFS_PROFILE_KEY, info);

        prefs = JSON.parse(info)?.prefs;
    }

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

        updateReq = fetch(`/api/account/${prefs?.user}`, opts);

        let resp = await updateReq;

        switch(resp.status) {
            case 200:
                loadUser();
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
                errText = await resp.text();
                console.error(errText, resp.status);
                
                alertMsg.show();
        }
    }
</script>

<h2>Account</h2>

<article>
    <form onsubmit={onSubmit}>
        <label>
            Default note visibility
            <select name="default-visibility">
                <option value="public" selected={prefs?.default_visibility === 'public'}>Public</option>
                <option value="protected" selected={prefs?.default_visibility === 'protected'}>Protected</option>
                <option value="private" selected={prefs?.default_visibility === 'private'}>Private</option>
            </select>
        </label>

        {#await updateReq}
            <button type="submit" class="secondary" aria-busy="true" disabled>Save</button>
        {:then}
            <button type="submit" class="secondary">Save</button>
        {/await}
    </form>
</article>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    button {
        width: auto;
        padding: 0.25rem 0.5rem;
    }
</style>
