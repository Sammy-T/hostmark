<script>
    let value = $state('');

    /**
     * @param {SubmitEvent} event
     */
    function onSubmit(event) {
        event.preventDefault();

        // @ts-ignore
        const data = new FormData(event.target);

        let msg = '';
        for(const [key, value] of data.entries()) {
            msg += `${key}: ${value}`;
        }
        console.log(msg);
    }
</script>

<header>
    <form onsubmit={onSubmit}>
        <textarea name="content" placeholder="New Note..." autocapitalize="on" spellcheck required bind:value></textarea>

        <fieldset>
            <button type="submit" class="secondary" disabled={value === ''}>Save</button>
        </fieldset>
    </form>
</header>

<style>
    header {
        position: sticky;
        top: 0;
        padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
        background-color: var(--pico-background-color);
        border: 2px solid oklch(from var(--pico-muted-color) l c h / 0.25);
        border-radius: var(--pico-border-radius);
        background-color: oklch(from var(--pico-background-color) l c h / 0.5);
        backdrop-filter: blur(10px);

        & form * {
            margin: 0;
        }

        & form {
            display: grid;
            gap: 0.5rem;
        }

        & fieldset {
            display: flex;
            justify-content: end;
        }

        & button {
            width: auto;
            padding: 0.25rem 0.5rem;
        }
    }
</style>
