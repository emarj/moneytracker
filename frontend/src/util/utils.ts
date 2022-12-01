export const DateFMT = (dt: Date) => new Date(dt).toLocaleString("en-GB")
// en-GB better than it-IT since it always uses two digits per day and month