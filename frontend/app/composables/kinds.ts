import type { Ref } from "vue";
import { ShowedNamespacedKinds, ShowedClusterKinds } from "~/provider/dashboard";

export const useNamespacedKindProvider = (data: Ref<{ filter: { namespaceKinds: string[] } } | null>) => {
    provide(ShowedNamespacedKinds, computed(() => data.value?.filter.namespaceKinds || []))
}

export const useNamespacedKindsInjection = (): Ref<string[]> => {
    return inject(ShowedNamespacedKinds, ref([]))
}

export const useNamespacedKindsFilter = () => {
    const kinds = useNamespacedKindsInjection()

    if (!kinds.value.length) {
        return ref([])
    }

    return kinds
}


export const useClusterKindProvider = (data: Ref<{ filter: { clusterKinds: string[] } } | null>) => {
    provide(ShowedClusterKinds, computed(() => data.value?.filter.clusterKinds || []))
}

export const useClusterKindsInjection = (): Ref<string[]> => {
    return inject(ShowedClusterKinds, ref([]))
}

export const useClusterKindsFilter = () => {
    const kinds = useClusterKindsInjection()

    if (!kinds.value.length) {
        return ref([])
    }

    return kinds
}