<script>
    /**
     * @callback onSubmitted
     * @param {number} status
     * @param {string} respText
     */

    /**
     * @typedef {Object} Props
     * @property {string} [username]
     * @property {string} [mode]
     * @property {onSubmitted} [onsubmitted]
     */

    /** @type {Props} */
    let { username, mode = 'edit', onsubmitted } = $props();

    /** @type {HTMLDialogElement} */
    let dialog;

    /** @type {HTMLFormElement} */
    let form;

    /** @type {HTMLInputElement} */
    let pwd;

    /** @type {HTMLInputElement} */
    let confirmPwd;

    export function show() {
        dialog.showModal();
    }

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

        dialog.close();

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

        const endpoint = (mode === 'create') ? '/api/account' : `/api/account/${username}`;

        const resp = await fetch(endpoint, opts);
        const respText = await resp.text();

        onsubmitted?.(resp.status, respText);
    }

    function onClose() {
        form.reset();
    }
</script>

<dialog id="edit-user" closedby="any" onclose={onClose} bind:this={dialog}>
    <div class="dialog-container">
        {#if mode === 'create'}
            <h3>Create user</h3>
        {:else}
            <h3>Edit {username}</h3>
        {/if}

        <form onsubmit={onSubmit} bind:this={form}>
            {#if mode === 'create'}
                <label for="username">Username</label>
                <input id="username" name="username" type="text" placeholder="Username" autocomplete="username" minlength="3" maxlength="32" required />
            {/if}

            <label for="password">Password</label>
            <input id="password" type="password" name="password" placeholder="Password (at least 15 characters)" minlength="15" maxlength="64"
                required={mode === 'create'} bind:this={pwd} />
            
            <label for="confirm">Confirm Password</label>
            <input id="confirm" name="confirm" type="password" placeholder="Password" required={mode === 'create'} 
                bind:this={confirmPwd} oninput={validateConfirmPwd} />

            <label for="role">Role</label>
            <select id="role" name="role" required={mode === 'create'}>
                <option value="" selected={mode === 'edit'}></option>
                <option value="user" selected={mode === 'create'}>User</option>
                <option value="admin">Admin</option>
            </select>

            <footer>
                <button type="button" class="secondary" commandfor="edit-user" command="close">Cancel</button>
                <button type="submit">Submit</button>
            </footer>
        </form>
    </div>
</dialog>

<style>
    footer {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 0.5rem;
    }
</style>
