<script lang="ts">
    import { isExpense, isIncome, isInternal } from "./transactions";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getOperationsByEntity } from "./api";
    import { entityID } from "./store";
    import CircularProgress from "@smui/circular-progress";
    import Operation from "./lib/Operation/Operation.svelte";

    const operationsQuery = useQuery(
        ["operations", "entity", $entityID, "all"],
        () => getOperationsByEntity($entityID, 1000)
    );

    let expanded = false;
</script>

<div>
    <h1>Operations</h1>

    {#if $operationsQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $operationsQuery.error}
        Error: {$operationsQuery.error?.message}
    {:else if $operationsQuery.data}
        <ol>
            {#each $operationsQuery.data as op (op.id)}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li on:click={() => (expanded = !expanded)} class:expanded>
                    <Operation {op} />
                </li>{/each}
        </ol>
    {/if}
</div>
