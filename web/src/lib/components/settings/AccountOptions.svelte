<script>
    import { PREFS_PROFILE_KEY } from '$lib/util.svelte';

    let prefs = JSON.parse(localStorage.getItem(PREFS_PROFILE_KEY) ?? '')?.prefs;
</script>

<h2>Account</h2>

<article>
    <form>
        <label>
            Default note visibility
            <select name="default-visibility">
                <option value="public" selected={prefs?.default_visibility === 'public'}>Public</option>
                <option value="protected" selected={prefs?.default_visibility === 'protected'}>Protected</option>
                <option value="private" selected={prefs?.default_visibility === 'private'}>Private</option>
            </select>
        </label>

        <button type="submit" class="secondary" disabled>Save</button>
    </form>
</article>

<style>
    button {
        width: auto;
        padding: 0.25rem 0.5rem;
    }
</style>
