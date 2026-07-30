import type { Dashboard, ResultView } from "~/types/core";
import type { Ref } from "vue";

export const useResultView = (data: Ref<Dashboard | undefined>) => {
    const route = useRoute()
    const router = useRouter()

    const allowedViews = computed<ResultView[]>(() => {
        const views = data.value?.renderOptions.allowedResultViews
        if (views?.length) return views

        return [data.value?.renderOptions.defaultResultView || data.value?.renderOptions.resultView || 'resources'] as ResultView[]
    })
    const defaultView = computed<ResultView>(() => data.value?.renderOptions.defaultResultView || allowedViews.value[0])
    const view = computed<ResultView>(() => {
        const requested = route.query.view
        return typeof requested === 'string' && allowedViews.value.includes(requested as ResultView)
            ? requested as ResultView
            : defaultView.value
    })
    const canSwitchResultView = computed(() => allowedViews.value.length > 1)

    const setResultView = (next: ResultView) => {
        if (!allowedViews.value.includes(next)) return

        router.push({ query: { ...route.query, view: next } })
    }

    return { allowedViews, canSwitchResultView, setResultView, view }
}
