

const baseURL = 'http://10.99.1.25:3245/api'

export const getAccountsByEntity = (eID: Number) =>
    fetch(`${baseURL}/accounts/${eID}`)
        .then(res => res.json())

export const getAccountBalances = (aID: Number) =>
    fetch(`${baseURL}/balances/${aID}`)
        .then(res => res.json())

export const getAccountBalance = (aID: Number) =>
    fetch(`${baseURL}/balance/${aID}`)
        .then(res => res.json())

export const getTransactionsByAccount = (aID: Number) =>
    fetch(`${baseURL}/transactions/account/${aID}`)
        .then(res => res.json())

export const getOperationsByEntity = (eID: Number) =>
    fetch(`${baseURL}/operations/entity/${eID}`)
        .then(res => res.json())
