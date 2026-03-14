<script>
    import { getContext } from 'svelte';

    /**
     * @typedef {Object} Props
     * @property {String} html The value of the html content
     * @property {String} [edited] The working value of the markdown content
     */

    /** @type {Props} */
    let { html, edited = $bindable() } = $props();

    /** @type {{ value: boolean }}*/
    const editing = getContext('editing');
</script>

{#if editing.value}
    <textarea name="editor" bind:value={edited}></textarea>
{:else}
    <section>{@html html}</section>
{/if}

<style>
    section {
        flex-grow: 1;
        padding: 0.5rem;
        margin: 0;
        overflow: auto;
    }

    textarea {
        flex-grow: 1;
        padding: 0.5rem;
        margin: 0;
        border: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        box-shadow: none;
    }
</style>
