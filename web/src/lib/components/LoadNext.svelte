<script>
    import { onMount } from 'svelte';

    /**
     * @callback onEnteredView
     */

    /**
     * @typedef {Object} Props
     * @property {onEnteredView} [onenteredview] 
     */

    /** @type {Props} */
    let { onenteredview } = $props();

    /** @type {HTMLElement} */
    let loadNextEl;

    /** @type {IntersectionObserverCallback} */
    function onIntersect(entries, observer) {
        entries.forEach((entry) => {
            if(entry.target !== loadNextEl || !entry.isIntersecting) return;
            
            onenteredview?.();
        });
    }

    onMount(() => {
        const observer = new IntersectionObserver(onIntersect, { threshold: 1 });
        observer.observe(loadNextEl);

        return () => {
            observer.unobserve(loadNextEl);
            observer.disconnect();
        };
    });
</script>

<article aria-busy="true" bind:this={loadNextEl}>Loading...</article>

<style>
    article {
        text-align: center;
    }
</style>
