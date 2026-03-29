<script>
    /**
     * @typedef {Object} Props
     * @property {'info'|'warning'} [type]
     * @property {String} [heading]
     * @property {import('svelte').Snippet} children
     */

    /** @type {Props} */
    let { type, heading, children } = $props();

    /** @type {HTMLDialogElement} */
    let dialog;

    export function show() {
        dialog.showModal();
    }
</script>

<dialog closedby="any" bind:this={dialog}>
    <div class="dialog-container">
        {#if heading}
        <h3 class={type}>{heading}</h3>
        {/if}
        <p>{@render children()}</p>
    </div>
</dialog>

<style>
    .dialog-container {
        padding: 1rem;
        border: 2px solid var(--pico-muted-color);
        border-radius: 1rem;
    }

    .warning {
        color: var(--pico-form-element-invalid-focus-color);
    }

    p {
        margin: 0;
    }
</style>
