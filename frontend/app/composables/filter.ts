import type { Filter } from "~/types/core";

export const useFilter = () => {
    const route = useRoute()

    const kinds = ref<string[]>([])
    const clusterKinds = ref<string[]>([])

    const source = computed<string>(() => route.params.source as string)
    const category = computed<string>(() => route.params.category as string)
    const policy = computed<string>(() => route.params.policy as string)

    const filter = computed((): Filter => ({
        sources: source.value ? [source.value] : undefined,
        categories: category.value ? [category.value] : undefined,
        policies: policy.value ? [policy.value] : undefined,
        kinds: kinds.value,
        clusterKinds: clusterKinds.value
    }))

    return { kinds, clusterKinds, filter, source, category, policy }
}