export const DateFMT = (str: string) => new Date(str).toLocaleString("en-GB")
// en-GB better than it-IT since it always uses two digits per day and month