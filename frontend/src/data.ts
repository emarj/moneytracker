

const baseURL = 'http://localhost:3245/api'

export const getEntities = () =>
    fetch(`${baseURL}/entities`)
        .then(res => res.json())

export const getAccounts = () =>
    fetch(`${baseURL}/accounts`)
        .then(res => res.json())

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


