<script>
    import OperationForm from "./lib/Operation/OperationForm.svelte";
    import { querystring } from "svelte-spa-router";
    import { entityID } from "./entity";

    /*let params = new URLSearchParams($querystring);

    const from = params.get("from");
    const to = params.get("to");*/

    let op = {
        description: null,
        owner: { id: $entityID },
        transactions: [{}],
    };
</script>

<OperationForm {op} />
