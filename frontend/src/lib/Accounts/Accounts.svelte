<script lang="ts">
  import CircularProgress from "@smui/circular-progress";

  import { getAccounts } from "../../api";

  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { user, entityID } from "../../store";
  import { push } from "svelte-spa-router";
  import IconButton from "@smui/icon-button/src/IconButton.svelte";

  const queryResult = useQuery(["accounts", $user?.id], () => getAccounts());
</script>

<div class="container">
  <IconButton
    class="add-account material-icons"
    on:click={() => push("/newaccount")}>add_circle_outline</IconButton
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

    :global(button.add-account) {
      position: absolute;
      top: 0;
      right: 0;
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
