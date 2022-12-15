<script lang="ts">
  import CircularProgress from "@smui/circular-progress";
  import MdAddCircleOutline from "svelte-icons/md/MdAddCircleOutline.svelte";

  import { getAccountsByEntity } from "../../api";

  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { entityID } from "../../store";
  import { push } from "svelte-spa-router";

  const queryResult = useQuery(["accounts", $entityID], () =>
    getAccountsByEntity($entityID)
  );
</script>

<div class="container">
  <h2>My Accounts</h2>

  <button
    class="add-account"
    title="New Account"
    on:click={() => push("/newaccount")}><MdAddCircleOutline /></button
  >

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
  .container {
    position: relative;

    button.add-account {
      position: absolute;
      top: 0;
      right: 0;
      width: 32px;
      height: 32px;
    }

    ul {
      padding: 0;
      list-style: none;
      display: flex;
      gap: 1rem;
      justify-content: center;
      flex-wrap: wrap;
    }
  }
</style>
