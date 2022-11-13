<script>
  import { getAccountsForEntity } from "../data";

  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";

  const eID = 1;

  const queryResult = useQuery(["accounts", eID], () =>
    getAccountsForEntity(eID)
  );
</script>

<div>
  <h2>My Accounts</h2>

  {#if $queryResult.isLoading}
    <span>Loading...</span>
  {:else if $queryResult.error}
    <span>An error has occurred: {$queryResult.error.message}</span>
  {:else}
    <ul>
      {#each $queryResult.data as account}
        <li><AccountCard {account} /></li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  ul {
    list-style: none;
    display: flex;
    gap: 1em;
  }
</style>
