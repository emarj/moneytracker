<script>
    import { entityID } from "../store";
    import EntitySelect from "./EntitySelect.svelte";
</script>

<EntitySelect bind:value={$entityID} />
