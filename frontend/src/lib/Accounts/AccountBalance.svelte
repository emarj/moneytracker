<script>
    import { getAccountBalance } from "../../data";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import Amount from "../Amount.svelte";

    export let id;

    const queryClient = useQueryClient();

    const balanceQuery = useQuery(["balance", id], () => getAccountBalance(id));

    export const refresh = () => {
        queryClient.invalidateQueries({ queryKey: ["balance", id] });
    };
</script>

<p>
    Balance: {#if $balanceQuery.isLoading}
        ...
    {:else if $balanceQuery.error}
        error
    {:else}
        <span> <Amount value={$balanceQuery.data.value} /></span>
    {/if}
</p>

<style lang="scss">
    p {
        font-size: 1.2rem;
        span {
            font-weight: bold;
        }
    }
</style>
