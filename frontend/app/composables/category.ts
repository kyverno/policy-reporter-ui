import type { Ref } from "vue";
import { ShowedCategories } from "~/provider/dashboard";

export const useCategoriesProvider = (data: Ref<{ filter: { categories: string[] } } | null>) => {
    provide(ShowedCategories, computed(() => data.value?.filter.categories || []))
}

export const useCategoriesInjection = (): Ref<string[]> => {
    return inject(ShowedCategories, ref([]))
}

export const useCategoriesFilter = () => {
    const Categories = useCategoriesInjection()

    if (!Categories.value.length) {
        return ref([])
    }

    return Categories
}
