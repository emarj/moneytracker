<script lang="ts">
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getOperations } from "../../api";
    import { authStore, entityID } from "../../store";
    import CircularProgress from "@smui/circular-progress";
    import { link } from "svelte-spa-router";
    import OperationPreview from "../Operation/OperationPreview.svelte";

    const operationsQuery = useQuery(
        ["operations", "user", $authStore.user.id, "latest"],
        () => getOperations()
    );

    let expanded = false;
</script>

<div>
    <h2>Lastest Operations</h2>
    <!--<button
        on:click={() => {
            console.log(`Invalidating Transactions of Entity ${eID}`);
            queryClient.invalidateQueries(["transactions", "entity", eID]);
        }}>Refresh</button
    >-->
    {#if $operationsQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $operationsQuery.error}
        Error: {$operationsQuery.error?.message}
    {:else if $operationsQuery.data}
        <ul>
            {#each $operationsQuery.data as op (op.id)}
                <li class:expanded>
                    <OperationPreview {op} />
                </li>
            {/each}
        </ul>
    {/if}
    <a href="/operations" use:link>More</a>
</div>

<style lang="scss">
    ul {
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
