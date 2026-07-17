import { DashboardType } from "~/provider/dashboard";
import type {Ref} from "vue";
import type { Dashboard, ViewType } from "~/types/core";

export const injectDashboardType = () => {
    return inject(DashboardType, ref('status'))
}

export const useDashboardType = (data: Ref<Dashboard | null>) => {
    provide(DashboardType, computed(() => data.value?.renderOptions.dataType || 'status'))
}

export const useDashboardHelper = (data: Ref<Dashboard | undefined>) => {
    const { mode, isCompact } = useMode(data)

    const showResults = computed<boolean>(() => data.value?.renderOptions.resultView === 'results')
    const dataType = computed<ViewType>(() => data.value?.renderOptions.dataType || 'status')
    const isSeverity = computed<boolean>(() => data.value?.renderOptions.dataType === 'severity')

    return { showResults, isSeverity, dataType, mode, isCompact }
}

export const useDashboardProvider = (data: Ref<Dashboard | null>) => {
    useStatusProvider(data)
    useSeveritiesProvider(data)
    useDashboardType(data)
    useNamespacedKindProvider(data)
    useClusterKindProvider(data)
    useNamespacesProvider(data)
    useCategoriesProvider(data)
}