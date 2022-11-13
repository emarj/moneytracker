

const baseURL = 'http://localhost:3245/api'

export const getAccountsByEntity = (eID: Number) =>
    fetch(`${baseURL}/accounts/${eID}`).then(res => res.json())

export const getAccountBalance = (aID: Number) =>
    fetch(`${baseURL}/balance/${aID}`).then(res => res.json())

export const getTransactionsByAccount = (aID: Number) =>
    fetch(`${baseURL}/transactions/account/${aID}`).then(res => res.json())

export const getTransactionsByEntity = (eID: Number) =>
    fetch(`${baseURL}/transactions/entity/${eID}`).then(res => res.json())
