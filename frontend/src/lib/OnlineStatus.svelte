<script lang="ts">
    import { onDestroy, onMount } from "svelte";

    export let online: boolean = navigator.onLine;

    onMount(() => {
        window.addEventListener("online", updateOnlineStatus);
        window.addEventListener("offline", updateOnlineStatus);
    });

    onDestroy(() => {
        window.removeEventListener("online", updateOnlineStatus);
        window.removeEventListener("offline", updateOnlineStatus);
    });

    const updateOnlineStatus = (event) => {
        online = navigator.onLine;
    };
</script>
