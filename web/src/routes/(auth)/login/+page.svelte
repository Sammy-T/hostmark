<script>
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { goto } from '$app/navigation';

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
        
        const resp = await fetch('/api/auth/login', opts);
        if(!resp.ok) {
            errText = await resp.text();
            console.error(errText, resp.status);

            alertMsg.show();
            return
        }

        goto('/');
    }
</script>

<h3>hostmark</h3>

<form onsubmit={onSubmit}>
    <label for="username">Username</label>
    <input id="username" name="username" type="text" placeholder="Username" required />

    <label for="password">Password</label>
    <input id="password" name="password" type="password" placeholder="Password" required />

    <button type="submit">Sign in</button>
    <small>Don't have an account? <a href="/signup">Sign up</a></small>
</form>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    h3 {
        text-align: center;
    }
</style>
