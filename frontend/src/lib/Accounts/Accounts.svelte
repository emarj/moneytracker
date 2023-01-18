<script lang="ts">
  import CircularProgress from "@smui/circular-progress";
  import MdAddCircleOutline from "svelte-icons/md/MdAddCircleOutline.svelte";

  import { getAccounts } from "../../api";

  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { authStore, entityID } from "../../store";
  import { push } from "svelte-spa-router";

  const queryResult = useQuery(["accounts", $authStore.user.id], () =>
    getAccounts()
  );
</script>

<div class="container">
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
    {#each [...Object.entries($queryResult.data)] as [entity, list]}
      <h2>{entity}</h2>
      <ul>
        {#each list as account}
          <li><AccountCard {account} /></li>
        {/each}
      </ul>
    {/each}
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
