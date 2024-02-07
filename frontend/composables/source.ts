type Store = {
    kinds: { namespaced: string[]; cluster: string[] },
    namespaces: string[]
    categories: string[]
}

const sources: { [source: string]: Store } = {}

export const useSourceStore = (source?: string) => {
    const key = source || 'global'

    if (!sources[key]) {
        sources[key] = reactive({
            kinds: { namespaced: [], cluster: [] },
            namespaces: [],
            categories: [],
        })
    }

    const store = sources[key]

    const loading = ref(false)
    const error = ref<Error | null>(null)

    const load = () => {
        loading.value = true
        error.value = null

        return Promise.all([
            callAPI(api => api.namespaces({ sources: source ? [source] : undefined })),
            callAPI(api => api.namespacedKinds(source)),
            callAPI(api => api.clusterKinds(source)),
            callAPI(api => api.categoryTree(undefined, { sources: source ? [source] : undefined })),
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