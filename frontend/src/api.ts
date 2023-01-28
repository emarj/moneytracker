import type { Operation } from "./model"
import { localLogout } from "./store"
import { JSONString } from "./util/utils"


const baseURL = '/api'

const fetcher = async (url: string, init?) => {
    const response = await fetch(url, init)
    if (!response.ok) {
        if (response.status === 401) {
            await logout()
            localLogout()
        }
        throw new Error(`error: ${response.statusText}`)
    }
    return response.json()
}

export const login = (l) =>
    fetcher(`${baseURL}/login`, { method: "POST", body: JSONString(l) })

export const logout = () =>
    fetcher(`${baseURL}/logout`, { method: "POST" })

/////

export const getTypes = () =>
    fetcher(`${baseURL}/types`)
/////

export const getUserEntities = () =>
    fetcher(`${baseURL}/entities`)

export const getAllEntities = () =>
    fetcher(`${baseURL}/entities/all`)

export const getCategories = () =>
    fetcher(`${baseURL}/categories`)
export const addCategory = (cat) =>
    fetcher(`${baseURL}/category`, { method: "POST", body: JSONString(cat) })

export const getAccounts = () =>
    fetcher(`${baseURL}/accounts`)

export const getAccountsByEntity = (eID: number) =>
    fetcher(`${baseURL}/accounts/${eID}`)

export const getAccount = (aID: number) =>
    fetcher(`${baseURL}/account/${aID}`)

export const deleteAccount = (aID: number) =>
    fetcher(`${baseURL}/account/${aID}`, { method: "DELETE" })

export const addAccount = (a) =>
    fetcher(`${baseURL}/account`, { method: "POST", body: JSONString(a) })


export const getAccountBalance = (aID: number) =>
    fetcher(`${baseURL}/balance/${aID}`)

export const getAccountBalances = (aID: number) =>
    fetcher(`${baseURL}/balance/history/${aID}`)

export const adjustBalance = (bal) =>
    fetcher(`${baseURL}/balance`, { method: "POST", body: JSONString(bal) })

export const getTransactionsByAccount = (aID: number) =>
    fetcher(`${baseURL}/transactions/account/${aID}`)

export const getOperations = (limit?: number) => {
    const limitStr = (limit) ? `?limit=${limit}` : ""
    return fetcher(`${baseURL}/operations${limitStr}`) as Promise<Operation[]>
}

export const getOperationsByEntity = (eID: number, limit?: number) => {
    const limitStr = (limit) ? `?limit=${limit}` : ""
    return fetcher(`${baseURL}/operations/entity/${eID}${limitStr}`) as Promise<Operation[]>
}

export const getOperation = (oID: number) =>
    fetcher(`${baseURL}/operation/${oID}`)

export const addOperation = (op) =>
    fetcher(`${baseURL}/operation`, { method: "POST", body: JSONString(op) })

export const deleteOperation = (opID: number) =>
    fetcher(`${baseURL}/operation/${opID}`, { method: "DELETE" })

export const addExpense = (e) =>
    fetcher(`${baseURL}/expense`, { method: "POST", body: JSONString(e) })
