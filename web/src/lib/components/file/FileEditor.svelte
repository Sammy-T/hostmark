<script>
    import { page } from '$app/state';
    import { refreshAll } from '$app/navigation';
    import { getContext } from 'svelte';

    let content = $derived(page.data.content);

    /** @type {{ value: boolean }}*/
    const editing = getContext('editing');

    let edited = $derived(content?.markdown);

    $effect(() => {
        if(!editing.value && content?.markdown !== edited) submitChanges();
    });

    async function submitChanges() {
        /** @type {RequestInit} */
        const opts = {
            method: 'POST',
            body: edited,
        };

        const resp = await fetch(`/api/file/${page.params.file}`, opts);
        if(!resp.ok) {
            console.error('unable to post changes');
            return;
        }

        refreshAll();
    }
</script>

{#if editing.value}
    <textarea bind:value={edited}></textarea>
{:else}
    <section>{@html content?.html}</section>
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
    }
</style>
