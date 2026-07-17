import type { Ref } from "vue";
import { ShowedNamespaces } from "~/provider/dashboard";

export const useNamespacesProvider = (data: Ref<{ namespaces: string[] } | null>) => {
    provide(ShowedNamespaces, computed(() => data.value?.namespaces || []))
}

export const useNamespacesInjection = (): Ref<string[]> => {
    return inject(ShowedNamespaces, ref([]))
}

export const useNamespacesFilter = () => {
    const namespaces = useNamespacesInjection()

    if (!namespaces.value.length) {
        return ref([])
    }

    return namespaces
}
