<script>
  import { entityID } from "../store";
  export let account;
</script>

{account.owner.id !== $entityID ? account.owner.name + ":" : ""}{account.name}
