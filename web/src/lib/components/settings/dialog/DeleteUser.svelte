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

    export function show() {
        dialog.showModal();
    }

    async function onConfirm() {
        dialog.close();
        
        /** @type {RequestInit} */
        const opts = {
            method: 'delete',
        };

        const resp = await fetch(`/api/account/${username}`, opts);
        const respText = await resp.text();

        onsubmitted?.(resp.status, respText);
    }
</script>

<dialog id="delete-user" closedby="any" bind:this={dialog}>
    <div class="dialog-container">
        <h3>Delete {username}?</h3>

        <footer>
            <button type="button" class="secondary" popovertarget="delete-user">Cancel</button>
            <button type="button" popovertarget="delete-user" onclick={onConfirm}>Confirm</button>
        </footer>
    </div>
</dialog>

<style>
    footer {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 0.5rem;
    }
</style>
