import { useConfigStore } from "~/store/config"

export const doAPIRefresh = (refresh: () => Promise<void>, onLoading?: (loading: boolean) => void) => {
    const store = useConfigStore()

    let interval: any = null

    watch(() => store.refreshInterval, (current) => {
        if (interval) {
            clearInterval(interval)
            interval = null
        }

        if (current <= 0) return

        interval = setInterval(() => refresh(), current)
    }, { immediate: true })

    watch(() => store.currentCluster, async () => {
        onLoading && onLoading(true)
        await refresh()
        onLoading && onLoading(false)
    })

    onUnmounted(() => interval && clearInterval(interval))
}