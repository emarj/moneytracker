<script>
    import { entityID } from "../entity";
    import EntitySelect from "./EntitySelect.svelte";
</script>

<EntitySelect bind:value={$entityID} />
