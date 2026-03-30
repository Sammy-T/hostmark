<script>
    import icProfile from '$lib/assets/user-circle.svg?raw';
    import { page } from '$app/state';

    let info = $derived(page.data);
</script>

<div class="profile">
    <section>
        <div id="icon">
            {@html icProfile}
        </div>
        
        <div id="info">
            <h4>{info.username}</h4>
            <p>{info.role}</p>
            <p class="muted">created {new Intl.DateTimeFormat().format(new Date(info.created_at))}</p>
        </div>
    </section>
</div>

<style>
    .muted {
        color: var(--pico-muted-color);
    }

    .profile {
        height: 100%;
        padding: 1rem;
        align-content: center;
    }

    section {
        width: fit-content;
        margin: auto;
        display: grid;
        grid-template-columns: auto 1fr;
        gap: 0.5rem;
    }

    #icon > :global(*) {
        width: 100%;
        height: 100%;
    }

    #info > * {
        margin: 0;
    }
</style>
