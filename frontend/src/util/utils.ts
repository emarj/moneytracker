export const TimestampFMT = (dt: Date) =>
    // en-GB better than it-IT since it always uses two digits per day and month
    new Date(dt).toLocaleString("en-GB")


const formatter = new Intl.RelativeTimeFormat('en', {
    numeric: 'auto'
});

export const DateFMT = (dt: Date) => {
    const diff = new Date(dt).getTime() - Date.now()
    return capitalize(formatter.format(Math.round(diff / 86400000), 'day'))
}

export const capitalize = (str: string) =>
    str.charAt(0).toUpperCase() + str.slice(1)

export const JSONPretty = (obj: any) => JSONString(obj, 4)


export const JSONString = (obj: any, space?: number | string) => JSON.stringify(obj, null, space)

export const removeUnderscore = (obj: any) => {
    var dup = {};
    for (var key in obj) {
        let newKey = key
        if (key.startsWith("_")) {
            newKey = key.slice(1)
        }

        dup[newKey] = obj[key];
    }
    return dup;
}