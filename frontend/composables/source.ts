import type {Ref} from "vue";
import {SourceContext} from "~/modules/core/provider/dashboard";

type Store = {
    kinds: { namespaced: string[]; cluster: string[] },
    namespaces: string[]
    categories: string[]
    key: string
}

const sources: { [source: string]: Store } = reactive({})

const getStore = (key?: string) => {
    if (!key) { key = 'global' }

    if (!sources[key]) {
        sources[key] = {
            kinds: { namespaced: [], cluster: [] },
            namespaces: [],
            categories: [],
            key,
        }
    }

    return sources[key]
} 

export const useSourceStore = (key?: string) => {
    const store = getStore(key)

    const loading = ref(false)
    const error = ref<Error | null>(null)

    const load = (source?: string[] | string) => {
        loading.value = true
        error.value = null

        let loadStore = store

        if (typeof source === 'string') {
            loadStore = getStore(source)
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
            console.log(nsKinds, clusterKinds)
            loadStore.kinds.namespaced = nsKinds || []
            loadStore.kinds.cluster = clusterKinds || []
            loadStore.namespaces = namespaces || []
            loadStore.categories = (categoryTrees || []).reduce<string[]>((categories, source) => {
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

export const useSourceContext = (source: Ref<string | undefined>) => {
    provide(SourceContext, source)
}

export const injectSourceContext = () => inject(SourceContext, ref(undefined))