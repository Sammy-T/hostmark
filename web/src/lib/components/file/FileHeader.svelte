<script>
    import icDots from '$lib/assets/dots-vertical.svg?raw';
    import icPencil from '$lib/assets/pencil.svg?raw';
    import icEdit from '$lib/assets/edit.svg?raw';
    import { page } from '$app/state';
    import { getContext } from 'svelte';

    let file = $derived(page.params.file);

    /** @type {{ value: boolean }} */
    const editing = getContext('editing');

    /**
     * @param {MouseEvent} ev
     */
    function onItemClick(ev) {
        // @ts-ignore
        const { name } = ev.target;
        
        switch(name) {
            case 'edit':
                editing.value = !editing.value;
                break;

            default:
                console.log(name);
        }
    }
</script>

<header>
    <div>{file}</div>

    <ul class="menu">
        <li><button popovertarget="pop-more">{@html icDots}</button></li>
    </ul>
</header>

<div id="pop-more" popover>
    <div class="container">
        <button popovertarget="pop-more" name="edit" onclick={onItemClick}>{@html icPencil} edit</button>
        <button popovertarget="pop-more" name="rename" onclick={onItemClick}>{@html icEdit} rename</button>
    </div>
</div>

<style>
    button {
        margin: 0;
        padding: 0;
        background-color: transparent;
        border: none;
    }

    button:hover {
        background-color: oklch(from var(--pico-contrast) l c h / 0.05);
    }

    header {
        padding: 0.5rem;
        outline: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        display: grid;
        grid-template-columns: 1fr auto;

        & > * {
            margin: 0;
        }
    }

    .menu {
        padding: 0;
        display: flex;
        flex-direction: row;
        gap: 0.25rem;

        & > li {
            list-style: none;
            margin: 0;
        }
    }

    [popover] {
        margin: 0;
        inset: auto;
        top: calc(anchor(bottom) + 0.25rem);
        right: anchor(right);
        border-width: 1px;
        border-radius: 0.25rem;
        border-color: var(--pico-dropdown-border-color);
        background-color: var(--pico-dropdown-background-color);
        box-shadow: var(--pico-dropdown-box-shadow);
        color: var(--pico-dropdown-color);

        & .container {
            display: grid;
            grid-template-columns: auto;

            & > * {
                padding: 0.15rem 0.5rem;
                display: flex;
                align-items: center;
                gap: 0.15rem;
            }
        }
    }
</style>
