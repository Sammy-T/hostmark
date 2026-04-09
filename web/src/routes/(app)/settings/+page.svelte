<script>
    import AccountOptions from '$lib/components/settings/AccountOptions.svelte';
    import AdminOptions from '$lib/components/settings/AdminOptions.svelte';
</script>

<div class="page">
    <AccountOptions />
    <AdminOptions />
</div>

<style>
    .page {
        padding: 1rem;
        flex-direction: column;
        overflow: auto;
    }
</style>
