<script lang="ts">
    import Textfield from "@smui/textfield";

    export let value: number | null = null;
    export let decimalDigits = undefined;
    export let disabled = false;
    export let validationPattern = `\\d+(\\.\\d${
        decimalDigits != null ? `{0,${decimalDigits}}` : "+"
    })?`;

    let valueText = "";

    $: value, updateText();

    const updateText = () => {
        valueText = value === null ? "" : value.toString();
    };

    const updateValue = () => {
        // This converts string to number
        value = +valueText;
    };

    const checkAndUpdateValue = () => {
        valueText = checkDecimalPart(valueText);

        updateValue();
    };

    const checkDecimalPart = (text: string): string => {
        let res = text;
        if (!decimalDigits) return res;

        const diff = decimalPart(text).length - decimalDigits;
        if (diff > 0) {
            res = text.slice(0, text.length - diff);
        }

        return res;
    };

    const decimalPart = (str: string) => {
        let res = "";
        const dotPos = str.indexOf(".");
        if (dotPos > -1) {
            res = valueText.slice(dotPos + 1);
        }
        return res;
    };

    const filter = (e: InputEvent) => {
        if (e.data === null) return;

        const char = e.data;
        const isDigit = new RegExp("^[0-9]+$").test(char);

        if (isDigit) {
            return;
        } else if (char === ".") {
            const firstDot = valueText.indexOf(".");
            if (firstDot < 0) {
                // There are no dots
                return;
            }
        }
        // In all other cases
        e.preventDefault();
    };

    const validate = (e: InputEvent) => {
        if (e.data === null) {
            // This is not insert, see e.inputType for details
            // we don't care, we just validate the value
            const n = parseDecimal(valueText, decimalDigits);
            if (!Number.isNaN(n) || n === null) {
                value = n;
            } else {
                valueText = "";
            }
            return;
        }

        // In all other cases
        checkAndUpdateValue();
    };

    // parseDecimal will return the parsed number if it is valid or NaN otherwise
    // an integer number ending with a dot will be parsed as valid and the dot will be ignored
    const parseDecimal = (
        str: string,
        decimalPlaces?: number
    ): number | null => {
        if (str === "") {
            return null;
        }

        const parts = str.split(".");

        if (parts.length > 2) {
            return NaN;
        }
        const regex = new RegExp("^[0-9]+$");

        if (!regex.test(parts[0])) {
            return NaN;
        }

        let outStr = parts[0];

        if (parts.length == 2 && parts[1] !== "") {
            let decimalpart = parts[1];
            if (!regex.test(decimalpart)) {
                return NaN;
            }

            if (decimalPlaces && decimalpart.length > decimalPlaces) {
                decimalpart = decimalpart.slice(0, decimalPlaces);
            }

            outStr += "." + decimalpart;
        }

        const result = +outStr; //Convert to number

        return result;
    };

    /*     const testArray = {
        "": null,
        "123..": NaN,
        "123.": 123,
        "123.3434.434": NaN,
        "123": 123,
        "123.23": 123.23,
        "123.232": 123.23,
        ".0": NaN,
        ".": NaN,
        "0.": 0,
    };

    for (const [input, output] of Object.entries(testArray)) {
        const res = parseDecimal(input, 2);
        if (res !== output && !(Number.isNaN(res) && Number.isNaN(output))) {
            //console.log(input, output, "FAIL", res, output);
        }
    } */
</script>

<Textfield
    variant="outlined"
    bind:value={valueText}
    on:beforeinput={(e) => filter(e)}
    on:input={(e) => validate(e)}
    label="Amount"
    suffix="€"
    input$pattern={validationPattern}
    input$inputmode="numeric"
    {disabled}
/>
