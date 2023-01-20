<script lang="ts">
  import CircularProgress from "@smui/circular-progress";
  import { getAccountsByEntity } from "../../api";
  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { push } from "svelte-spa-router";
  import IconButton from "@smui/icon-button/src/IconButton.svelte";

  export let entity;

  const queryResult = useQuery(["accounts", "entity", entity.id], () =>
    getAccountsByEntity(entity.id)
  );
</script>

<section class="container">
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
    <h2>{entity.display_name}</h2>
    <ul>
      {#each $queryResult.data as account}
        <li><AccountCard {account} /></li>
      {/each}
    </ul>
  {/if}
</section>

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
