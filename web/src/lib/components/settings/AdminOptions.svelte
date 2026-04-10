<script>
    import icPlus from '$lib/assets/plus.svg?raw';
    import icMore from '$lib/assets/dots-vertical.svg?raw';
    import icEdit from '$lib/assets/edit.svg?raw';
    import icDelete from '$lib/assets/trash.svg?raw';
    import EditUser from './dialog/EditUser.svelte';
    import DeleteUser from './dialog/DeleteUser.svelte';
    import AlertMessage from '../AlertMessage.svelte';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { STORAGE_PROFILE_KEY } from '$lib/util.svelte';

    let userInfo = JSON.parse(localStorage.getItem(STORAGE_PROFILE_KEY) ?? '');

    // @ts-ignore
    let users = $state([]);

    let updateMode = $state('edit');

    let editingUser = $state('');

    /** @type {EditUser} */
    let editUserDialog;

    /** @type {AlertMessage} */
    let alertMsg;
    
    let errText = $state('');

    function setCreating() {
        updateMode = 'create';
        editUserDialog.show();
    }

    /**
     * @param {string} username
     */
    function setEditing(username) {
        updateMode = 'edit';
        editingUser = username;
        editUserDialog.show();
    }

    /**
     * @param {string} username
     */
    function setDeleting(username) {
        editingUser = username;
    }

    /**
     * @param {number} status
     * @param {string} respText
     */
    function onEditSubmitted(status, respText) {
        editingUser = '';

        switch(status) {
            case 200:
                loadUsers();
                break;

            default:
                errText = respText;
                alertMsg.show();
        }
    }

    async function loadUsers() {
        const resp = await fetch('/api/account/list');

        switch(resp.status) {
            case 200:
                break;
            
            case 401:
                const refResp = await fetch('/api/auth/refresh');

                switch(refResp.status) {
                    case 200:
                        loadUsers();
                        return;

                    case 400:
                    case 401:
                        goto('/login');
                        return;

                    default:
                        errText = await refResp.text();
                        console.error(errText, refResp.status);

                        alertMsg.show();
                        return;
                }
            
            default:
                errText = await resp.text();
                console.error(errText, resp.status);
                
                alertMsg.show();
        }

        const data = await resp.json();
        users = data;
    }

    onMount(() => {
        loadUsers();
    });
</script>

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
                    <button popovertarget={menuId} popovertargetaction="hide" onclick={() => setEditing(username)}>{@html icEdit} Edit</button>
                    <button popovertarget={menuId} popovertargetaction="hide" onclick={() => setDeleting(username)}>{@html icDelete} Delete</button>
                </div>
            </div>
        </td>
    {/if}
</tr>
{/snippet}

<h2>Admin</h2>

<article>
    <header>
        <h3>Members</h3>
        <button class="secondary" onclick={setCreating}>{@html icPlus} Create</button>
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
            {#each users as user (user.username)}
                {@const isAdmin = user.username === userInfo?.username}

                {@render memberRow(user.username, user.role, isAdmin)}                
            {/each}
        </tbody>
    </table>
</article>

<EditUser mode={updateMode} username={editingUser} onsubmitted={onEditSubmitted} bind:this={editUserDialog} />
<DeleteUser />

<AlertMessage type="warning" heading="Error" bind:this={alertMsg}>
    {errText}
</AlertMessage>

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
