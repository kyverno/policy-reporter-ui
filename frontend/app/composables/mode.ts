import type { Dashboard, Mode } from "~/types/core";
import type { Ref } from "vue";

export const useMode = (data: Ref<Dashboard | undefined>) => {
    const route = useRoute()
    const mode = ref<Mode>((route.query.mode as Mode))

    watch(data, (d) => {
        if (mode.value) return;

        mode.value = d?.renderOptions.dashboardMode ?? 'detailed'
    })

    const isCompact = computed(() => mode.value === 'compact')
    const isDetailed = computed(() => mode.value === 'detailed')

    return { mode, isCompact, isDetailed }
}