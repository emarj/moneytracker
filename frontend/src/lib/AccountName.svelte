<script>
  import { entityID } from "../store";
  export let account;
</script>

{#if account.owner_id === 0}
  🌐
{:else}
  {account.owner_id !== 0 && account.owner_id !== $entityID
    ? account.owner.name + ":"
    : ""}{account.name}
{/if}
