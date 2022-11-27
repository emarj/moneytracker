<script lang="ts">
  import CircularProgress from "@smui/circular-progress";

  import { getAccountsByEntity } from "../../api";

  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { entityID } from "../../store";

  const queryResult = useQuery(["accounts", $entityID], () =>
    getAccountsByEntity($entityID)
  );
</script>

<div>
  <h2>My Accounts</h2>

  {#if $queryResult.isLoading}
    <span
      ><CircularProgress
        style="height: 32px; width: 32px;"
        indeterminate
      /></span
    >
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

<style lang="scss">
  ul {
    padding: 0;
    list-style: none;
    display: flex;
    gap: 1rem;
    justify-content: center;
    flex-wrap: wrap;
  }
</style>
