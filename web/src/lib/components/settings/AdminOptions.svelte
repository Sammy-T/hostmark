<script>
    import icPlus from '$lib/assets/plus.svg?raw';
    import icMore from '$lib/assets/dots-vertical.svg?raw';
    import icEdit from '$lib/assets/edit.svg?raw';
    import icDelete from '$lib/assets/trash.svg?raw';
</script>

<h2>Admin</h2>

{#snippet memberRow(/** @type {String} */ username, /** @type {String} */ role, isSelf = false)}
<tr>
    <td>{username}</td>
    <td>{role}</td>

    {#if isSelf}
        <td>self</td>
    {:else}
        {@const menuId = `${username}-menu`}

        <td>
            <button popovertarget={menuId}>{@html icMore}</button>

            <div id={menuId} popover>
                <div class="menu-container">
                    <button popovertarget={menuId}>{@html icEdit} Edit</button>
                    <button popovertarget={menuId}>{@html icDelete} Delete</button>
                </div>
            </div>
        </td>
    {/if}
</tr>
{/snippet}

<article>
    <header>
        <h3>Members</h3>
        <button class="secondary">{@html icPlus} Create</button>
    </header>

    <table>
        <thead>
            <tr>
                <th>username</th>
                <th>role</th>
                <th></th>
            </tr>
        </thead>
        <tbody>
            {#each { length: 5 } as _, i}
                {@const role = (i === 0) ? 'admin' : 'user'}
                {@const isAdmin = (i === 0)}

                {@render memberRow(`user-${i}`, role, isAdmin)}
            {/each}
        </tbody>
    </table>
</article>

<style>
    [popover] {
        margin: 0;
        inset: auto;
        border: none;
        top: calc(anchor(bottom) + 0.25rem);
        right: anchor(right);
        border: 1px solid oklch(from var(--pico-contrast) l c h / 0.15);
        border-radius: 0.25rem;
        background-color: var(--pico-dropdown-background-color);
        box-shadow: var(--pico-dropdown-box-shadow);
    }

    .menu-container {
        display: flex;
        flex-direction: column;

        & button {
            padding: 0 0.5rem;
            display: flex;
            justify-content: start;
            align-items: center;
            gap: 0.25rem;
            color: var(--pico-dropdown-color);

            & :global(svg) {
                width: 1rem;
            }
        }

        & button:hover {
            background-color: var(--pico-dropdown-hover-background-color);
        }
    }

    article > header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        & * {
            margin: 0;
        }
    }

    button {
        width: auto;
        padding: 0.25rem 0.5rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 0.25rem;

        & :global(svg) {
            width: 1.25rem;
        }
    }

    table { 
        & button {
            padding: 0;
            background: transparent;
            border: none;
        }

        & button:hover {
            background-color: var(--pico-dropdown-hover-background-color);
        }
    }
</style>
