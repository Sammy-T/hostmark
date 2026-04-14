<script>
    import AccountOptions from '$lib/components/settings/AccountOptions.svelte';
    import AdminOptions from '$lib/components/settings/AdminOptions.svelte';
    import { STORAGE_PROFILE_KEY } from '$lib/util.svelte';
    import { onMount } from 'svelte';

    let userInfo = $state();

    onMount(() => {
        userInfo = JSON.parse(localStorage?.getItem(STORAGE_PROFILE_KEY) ?? '');
    });
</script>

<div class="page">
    <AccountOptions />

    {#if userInfo?.role === 'admin'}
        <AdminOptions />
    {/if}
</div>

<style>
    .page {
        padding: 1rem;
        flex-direction: column;
        overflow: auto;
    }
</style>
