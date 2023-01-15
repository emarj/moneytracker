<script lang="ts">
    import { isExpense, isIncome, isInternal } from "../transactions";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getOperationsByEntity } from "../api";
    import { entityID } from "../store";
    import CircularProgress from "@smui/circular-progress";
    import OperationPreview from "../lib/Operation/OperationPreview.svelte";

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
                    <OperationPreview {op} />
                </li>{/each}
        </ol>
    {/if}
</div>

<style lang="scss">
    ol {
        list-style: none;
        width: 100%;
        padding: 0;

        & > li {
            --top-color: rgb(185, 185, 185);

            width: 100%;
            padding: 1rem 1rem 0.5rem;
            border-radius: 10px;
            margin-bottom: 1em;
            background: rgb(238, 238, 238);
            background: linear-gradient(
                180deg,
                var(--top-color) 5px,
                rgb(238, 238, 238) 5px
            );
            box-shadow: 10px 10px 10px rgba(192, 192, 192, 0.4);
            position: relative;

            font-size: 1.1rem;

            &.expense {
                --top-color: red;
            }

            &.income {
                --top-color: green;
            }

            :global(.amount) {
                font-weight: bold;
            }

            .transactions {
                display: none;
            }

            &.expanded .transactions {
                display: block;
            }
        }
    }
</style>
