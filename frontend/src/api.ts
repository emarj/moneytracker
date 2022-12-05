import type { Operation } from "./model"


const baseURL = 'http://localhost:3245/api'

const fetcher = async (url: string, init?) => {
    const response = await fetch(url, init)
    if (!response.ok) {
        throw new Error(`error: ${response.statusText}`)
    }
    return response.json()
}

export const login = (l) =>
    fetcher(`${baseURL}/login`, { method: "POST", body: JSON.stringify(l) })

export const logout = () =>
    fetcher(`${baseURL}/logout`, { method: "POST" })

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

export const adjustBalance = (bal) =>
    fetcher(`${baseURL}/balance`, { method: "POST", body: JSON.stringify(bal) })

export const getTransactionsByAccount = (aID: number) =>
    fetcher(`${baseURL}/transactions/account/${aID}`)

export const getOperationsByEntity = (eID: number, limit?: number) => {
    const limitStr = (limit) ? `?limit=${limit}` : ""
    return fetcher(`${baseURL}/operations/entity/${eID}${limitStr}`) as Promise<Operation[]>
}

export const addOperation = (op) =>
    fetcher(`${baseURL}/operation`, { method: "POST", body: JSON.stringify(op) })

export const deleteOperation = (opID: number) =>
    fetcher(`${baseURL}/operation/${opID}`, { method: "DELETE" })