

const baseURL = 'http://localhost:3245/api'

export const getAccountsForEntity = (eID: Number) =>
    fetch(`${baseURL}/accounts/${eID}`).then(res => res.json())

export const getAccountBalance = (aID: Number) =>
    fetch(`${baseURL}/balance/${aID}`).then(res => res.json())

export const getAccountsForAccount = (aID: Number) =>
    fetch(`${baseURL}/transactions/${aID}`).then(res => res.json())
