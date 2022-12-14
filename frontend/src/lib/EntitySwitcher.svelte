<script>
    import { entityID } from "../store";
    import EntitySelect from "./EntitySelect.svelte";
    export let style = "material";
</script>

<EntitySelect bind:value={$entityID} {style} />
