<script>
    export let value;
    export let negative = false;
    export let hide_plus = true;

    export let currency = "EUR";

    let fmtString;

    $: {
        let sign = "";
        if (negative) sign = "-";
        else if (!hide_plus) sign = "+";

        fmtString = `${sign}${value.toLocaleString("it-IT", {
            style: "currency",
            currency: "EUR",
            minimimunFractionDigits: 2,
        })}`;
    }
</script>

<span class="amount amount-{negative ? 'neg' : 'pos'}">{fmtString}</span>
