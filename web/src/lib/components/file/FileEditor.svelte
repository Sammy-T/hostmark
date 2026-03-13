<script>
    import { refreshAll } from '$app/navigation';
    import { page } from '$app/state';

    let content = $derived(page.data.content);

    let editing = $state(false);
    let edited = $state();

    let lastClick = 0;

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

    function onDoubleClick() {
        editing = !editing;

        if(editing) return;

        console.log(content?.markdown === edited, edited);

        if(content?.markdown !== edited) submitChanges();
    }

    function onClick() {
        const clickTime = new Date().getTime();

        if(clickTime - lastClick < 250) onDoubleClick();

        lastClick = clickTime;
    }
</script>

{#if editing}
    <textarea onclick={onClick} bind:value={edited}>{content?.markdown}</textarea>
{:else}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <section onclick={onClick}>{@html content?.html}</section>
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
