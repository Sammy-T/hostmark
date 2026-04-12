<script>
    /**
     * @callback onSubmitted
     * @param {number} status
     * @param {string} respText
     */

    /**
     * @typedef {Object} Props
     * @property {string} username
     * @property {onSubmitted} [onsubmitted]
     */

    /** @type {Props} */
    let { username, onsubmitted } = $props();

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

        const resp = await fetch(`/api/account/${username}`, opts);
        const respText = await resp.text();

        onsubmitted?.(resp.status, respText);
    }

    function onClose() {
        form.reset();
    }
</script>

<dialog id="edit-user" closedby="any" onclose={onClose} bind:this={dialog}>
    <div class="dialog-container">
        <h3>Update password</h3>

        <form onsubmit={onSubmit} bind:this={form}>
            <label for="current-password">Current Password</label>
            <input id="current-password" type="password" name="current-password" placeholder="Current password" minlength="15" maxlength="64" 
                required />

            <label for="password">New Password</label>
            <input id="password" type="password" name="password" placeholder="Password (at least 15 characters)" minlength="15" maxlength="64"
                required bind:this={pwd} />
            
            <label for="confirm">Confirm New Password</label>
            <input id="confirm" type="password" name="confirm" placeholder="Password" required bind:this={confirmPwd} oninput={validateConfirmPwd} />

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
