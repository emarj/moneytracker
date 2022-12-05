<script>
    import { showBalances } from "../store";

    export let value;
    export let negative = false;
    export let hide_plus = true;

    let fmtString;

    $: {
        let sign = "";
        if (negative) sign = "-";
        else if (!hide_plus) sign = "+";

        fmtString = `${sign}${new Number(value).toLocaleString("it-IT", {
            style: "currency",
            currency: "EUR",
        })}`;
    }
</script>

<span
    class="amount amount-{negative ? 'neg' : 'pos'}"
    class:hidden={!$showBalances}>{fmtString}</span
>

<style>
    .amount.hidden {
        filter: blur(0.5rem);
    }
</style>
