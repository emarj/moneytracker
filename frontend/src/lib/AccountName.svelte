<script>
  import { userShares } from "../store";
  export let account;

  $: defaultEntity = $userShares.find((s) => s?.priority === 0);
</script>

{#if account.owner_id === 0}
  🌐
{:else}
  {#if account.owner_id !== defaultEntity?.user_id}{account.owner.name +
      ":"}{/if}{account.name}
{/if}
