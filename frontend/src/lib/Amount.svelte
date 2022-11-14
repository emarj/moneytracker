<script>
    export let value;
    export let negative = false;
    export let hide_plus = true;

    export let currency = "euro";

    let symbols = new Map([
        ["euro", "€"],
        ["dollar", "$"],
    ]);

    const format = () => {
        let sign = "";
        if (negative) sign = "-";
        else if (!hide_plus) sign = "+";

        return `${sign}${value}${symbol}`;
    };

    $: symbol = symbols.has(currency) ? symbols.get(currency) : currency;
</script>

{format()}
