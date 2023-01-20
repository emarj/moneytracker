import { writable, derived } from 'svelte/store';

import { writable as lsWritable } from 'svelte-local-storage-store'

// Persistent stores
export const authStore = lsWritable("auth", null)
export const entityID = lsWritable("entity_id", 1);
export const showBalances = lsWritable("showBalances", true);

export const user = derived(
    authStore,
    $authStore => $authStore?.user
);

export const userShares = derived(
    user,
    $user => $user?.shares.sort(cmp)
);

const cmp = (s, t) => fix(s.priority) - fix(t.priority)
const fix = (p) => ((p === null) ? 999 : p)

export const userEntities = derived(
    userShares,
    $userShares => $userShares?.map((s) => s.entity)
);

export const userEntitiesID = derived(
    userShares,
    $userShares => $userShares?.map((s) => s.entity_id)
);

export const userDefaultEntity = derived(
    userShares,
    $userShares => $userShares.find((s) => s?.priority === 0)
)

// In-memory stores
export const messageStore = writable(null);

// TODO: Type store
//export const typesStore = writable(null);



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

