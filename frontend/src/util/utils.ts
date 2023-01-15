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
