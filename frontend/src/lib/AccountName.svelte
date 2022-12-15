<script>
  import { entityID } from "../store";
  export let account;
</script>

{#if account.owner.id === 0}
  🌐
{:else}
  {account.owner.id !== 0 && account.owner.id !== $entityID
    ? account.owner.name + ":"
    : ""}{account.name}
{/if}
