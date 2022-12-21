<script>
    import CircularProgress from "@smui/circular-progress";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getOperation } from "../api";
    import Operation from "../lib/Operation/Operation.svelte";

    export let params = {};

    const opQuery = useQuery(["operation", params.id], () =>
        getOperation(params.id)
    );
</script>

{#if $opQuery.isLoading}
    <span
        ><CircularProgress
            style="height: 32px; width: 32px;"
            indeterminate
        /></span
    >
{:else if $opQuery.error}
    <span>An error has occurred: {$opQuery.error.message}</span>
{:else if $opQuery.data}
    <Operation op={$opQuery.data} />
{/if}
