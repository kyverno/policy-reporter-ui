import { DashboardType } from "~/provider/dashboard";
import type {Ref} from "vue";
import type { Dashboard, ViewType } from "~/types/core";

export const injectDashboardType = () => {
    return inject(DashboardType, ref('status'))
}

export const useDashboardType = (data: Ref<Dashboard>) => {
    provide(DashboardType, computed(() => data.value?.renderOptions.dataType || 'status'))
}

export const useDashboardHelper = (data: Ref<Dashboard | undefined>) => {
    const { mode, isCompact } = useMode(data)
    const { allowedViews, canSwitchResultView, setResultView, view } = useResultView(data)

    const showResults = computed<boolean>(() => view.value === 'results')
    const dataType = computed<ViewType>(() => data.value?.renderOptions.dataType || 'status')
    const isSeverity = computed<boolean>(() => data.value?.renderOptions.dataType === 'severity')

    return { allowedViews, canSwitchResultView, setResultView, showResults, isSeverity, dataType, mode, isCompact, view }
}
