<script>
    import AlertMessage from '$lib/components/AlertMessage.svelte';
    import { goto } from '$app/navigation';
    import { PREFS_PROFILE_KEY } from '$lib/util.svelte';

    /** @type {HTMLInputElement} */
    let pwd;

    /** @type {HTMLInputElement} */
    let confirmPwd;

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    function validateConfirmPwd() {
        if(pwd.value === confirmPwd.value) {
            confirmPwd.setCustomValidity('');
        } else {
            confirmPwd.setCustomValidity('Must match password.');
        }

        if(confirmPwd.value.length < pwd.value.length) return;

        confirmPwd.reportValidity();
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
        
        const resp = await fetch('/api/auth/signup', opts);
        if(!resp.ok) {
            errText = await resp.text();
            console.error(errText, resp.status);

            alertMsg.show();
            return
        }

        const userInfo = await resp.text();
        localStorage.setItem(PREFS_PROFILE_KEY, userInfo);

        goto('/');
    }
</script>

<h3>hostmark</h3>
<h4>Create your account</h4>

<form onsubmit={onSubmit}>
    <label for="username">Username</label>
    <input id="username" name="username" type="text" placeholder="Username" autocomplete="username" minlength="3" maxlength="32" required />

    <label for="password">Password</label>
    <input id="password" name="password" type="password" placeholder="Password (at least 15 characters)" minlength="15" maxlength="64" 
        required bind:this={pwd} />

    <label for="confirm">Confirm Password</label>
    <input id="confirm" name="confirm" type="password" placeholder="Password" required bind:this={confirmPwd} oninput={validateConfirmPwd} />

    <button type="submit">Sign up</button>
    <small>Already have an account? <a href="/login">Sign in</a></small>
</form>

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

<style>
    h3 {
        text-align: center;
    }
</style>
