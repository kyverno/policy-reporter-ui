type Store = {
    kinds: { namespaced: string[]; cluster: string[] },
    namespaces: string[]
    categories: string[]
}

const sources: { [source: string]: Store } = reactive({})

export const useSourceStore = (key?: string) => {
    if (!key) { key = 'global' }

    if (!sources[key]) {
        sources[key] = {
            kinds: { namespaced: [], cluster: [] },
            namespaces: [],
            categories: [],
        }
    }

    const store = sources[key]

    const loading = ref(false)
    const error = ref<Error | null>(null)

    const load = (source?: string[] | string) => {
        loading.value = true
        error.value = null

        if (typeof source === 'string') {
            source = [source]
        } else if (Array.isArray(source) && !source.length) {
            source = undefined
        }
    
        return Promise.all([
            callAPI(api => api.namespaces({ sources: source as string[] })),
            callAPI(api => api.namespacedKinds(source as string[])),
            callAPI(api => api.clusterKinds(source as string[])),
            callAPI(api => api.categoryTree(undefined, { sources: source as string[] })),
        ]).then(([namespaces, nsKinds, clusterKinds, categoryTrees]) => {
            store.kinds.namespaced = nsKinds || []
            store.kinds.cluster = clusterKinds || []
            store.namespaces = namespaces || []
            store.categories = (categoryTrees || []).reduce<string[]>((categories, source) => {
                return [...categories, ...source.categories.map(c => c.name)]
            }, [])
        }).catch((err) => {
            error.value = err
        }).finally(() => {
            loading.value = false
        })
    }

    return {
        store,
        loading,
        error,
        load
    }
}