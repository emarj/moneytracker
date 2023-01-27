<script lang="ts">
  import CircularProgress from "@smui/circular-progress";
  import { getAccountsByEntity } from "../../api";
  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { push } from "svelte-spa-router";
  import IconButton from "@smui/icon-button";
  import Dialog, { Title, Content, Actions } from "@smui/dialog";
  import Button, { Label } from "@smui/button";
  import AccountForm from "./AccountForm.svelte";

  export let entity;

  const queryResult = useQuery(["accounts", "entity", entity.id], () =>
    getAccountsByEntity(entity.id)
  );

  let openDialog = false;
  let resetAccountForm;
</script>

<section class="container">
  <IconButton
    class="add-account material-icons"
    on:click={() => (openDialog = true)}>add_circle_outline</IconButton
  >

  <Dialog
    bind:open={openDialog}
    on:SMUIDialog:closing={resetAccountForm}
    sheet
    aria-labelledby="event-title"
    aria-describedby="event-content"
  >
    <Title id="event-title">New Account</Title>
    <Content id="event-content">
      <IconButton action="close" class="material-icons">close</IconButton>
      <AccountForm
        on:submit={() => (openDialog = false)}
        bind:reset={resetAccountForm}
        defaultEntityID={entity.id}
      />
    </Content>
    <!--  <Actions>
      <Button>
        <Label>None of Them</Label>
      </Button>
      <Button default>
        <Label>All of Them</Label>
      </Button>
    </Actions> -->
  </Dialog>
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
