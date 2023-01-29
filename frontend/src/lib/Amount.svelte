<script>
    import { showBalances } from "../store";

    export let value;
    export let invert = false;
    export let hide_plus = true;

    let fmtString;

    $: {
        if (invert) {
            value = -value;
        }

        fmtString = new Number(value).toLocaleString("it-IT", {
            style: "currency",
            currency: "EUR",
            signDisplay: hide_plus ? "auto" : "exceptZero",
        });
    }
</script>

<span
    class="amount"
    class:positive={value > 0}
    class:negative={value < 0}
    class:hidden={!$showBalances}>{fmtString}</span
>

<style>
    .amount.hidden {
        filter: blur(0.5rem);
    }
</style>
