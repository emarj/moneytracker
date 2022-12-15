import { writable, derived } from 'svelte/store';

import { writable as lsWritable } from 'svelte-local-storage-store'

// Persistent stores
export const authStore = lsWritable("auth", null)
export const entityID = lsWritable("entity_id", 1);

// In-memory stores
export const messageStore = writable(null);
export const showBalances = lsWritable("showBalances", true);


// History Store
function createHistoryStore() {
    const { subscribe, set, update } = writable({ stack: [], aboutToPop: false });

    return {
        subscribe,
        push: (route: string) => update(data => {
            if (!data.aboutToPop) {
                data.stack.push(route);
            } else {
                data.aboutToPop = false;
            }
            return data;
        }),
        pop: () => update((data) => {
            data.stack.pop();
            data.aboutToPop = true;
            return data;
        }),
    };
}

export const historyStore = createHistoryStore();

export const isFirstPage = derived(
    historyStore,
    $historyStore => $historyStore.stack.length <= 1
);

