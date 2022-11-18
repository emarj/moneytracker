<script>
    import OperationForm from "./lib/Operation/OperationForm.svelte";
    import { querystring } from "svelte-spa-router";

    let params = new URLSearchParams($querystring);

    const from = params.get("from");
    let op = {};
    if (from) {
        op.from = from;
    }
</script>

<OperationForm {op} />
