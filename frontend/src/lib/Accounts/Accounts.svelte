<script lang="ts">
  import CircularProgress from "@smui/circular-progress";

  import { getAccounts } from "../../api";

  import { useQuery } from "@sveltestack/svelte-query";
  import AccountCard from "./AccountCard.svelte";
  import { user, userEntities } from "../../store";
  import { push } from "svelte-spa-router";
  import IconButton from "@smui/icon-button/src/IconButton.svelte";
  import AccountsForEntity from "./AccountsForEntity.svelte";
</script>

<div class="container">
  {#if $userEntities}
    {#each $userEntities as e}
      <AccountsForEntity entity={e} />
    {/each}
  {/if}
</div>
