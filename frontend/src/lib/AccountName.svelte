<script>
      import { entityID } from "../entity";
    export let account;
    </script>

{account.owner.id !== $entityID
    ? account.owner.name + ":"
    : ""}{account.name}