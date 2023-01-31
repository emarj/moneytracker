<script>
    import { getAccountBalance } from "../../api";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import Amount from "../Amount.svelte";

    export let id;
    export let label = "Balance";

    const queryClient = useQueryClient();

    const balanceQuery = useQuery(["balance", id], () => getAccountBalance(id));

    export const refresh = () => {
        queryClient.invalidateQueries({ queryKey: ["balance", id] });
    };
</script>

<span>
    {label}: {#if $balanceQuery.isLoading}
        ...
    {:else if $balanceQuery.error}
        error
    {:else}
        <span> <Amount value={$balanceQuery.data.value} /></span>
    {/if}
</span>

<style lang="scss">
    span {
        font-size: 1.2rem;
        span {
            font-weight: bold;
        }
    }
</style>
