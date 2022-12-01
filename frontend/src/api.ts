import type { Operation } from "./model"


const baseURL = 'http://localhost:3245/api'

const fetcher = async (url: string, init?) => {
    const response = await fetch(url, init)
    if (!response.ok) {
        throw new Error(`error: ${response.statusText}`)
    }
    return response.json()
}

export const getEntities = () =>
    fetcher(`${baseURL}/entities`)

export const getCategories = () =>
    fetcher(`${baseURL}/categories`)

export const getAccounts = () =>
    fetcher(`${baseURL}/accounts`)

export const getAccountsByEntity = (eID: Number) =>
    fetcher(`${baseURL}/accounts/${eID}`)

export const getAccountBalances = (aID: Number) =>
    fetcher(`${baseURL}/balances/${aID}`)

export const getAccountBalance = (aID: Number) =>
    fetcher(`${baseURL}/balance/${aID}`)

export const getTransactionsByAccount = (aID: Number) =>
    fetcher(`${baseURL}/transactions/account/${aID}`)

export const getOperationsByEntity = (eID: Number) =>
    fetcher(`${baseURL}/operations/entity/${eID}`) as Promise<Operation[]>

export const addOperation = (op) =>
    fetcher(`${baseURL}/operation`, { method: "POST", body: JSON.stringify(op) })