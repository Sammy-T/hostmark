<script>
    import { page } from '$app/state';

    let content = $derived(page.data.content);
    let editing = $state(false);

    let lastClick = 0;

    function onDoubleClick() {
        editing = !editing;
    }

    function onClick() {
        const clickTime = new Date().getTime();

        if(clickTime - lastClick < 250) onDoubleClick();

        lastClick = clickTime;
    }
</script>

{#if editing}
    <textarea onclick={onClick}>{content?.markdown}</textarea>
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
