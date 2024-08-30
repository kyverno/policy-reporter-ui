import { DashboardType } from "~/modules/core/provider/dashboard";
import type {Ref} from "vue";
import {da} from "vuetify/locale";

export const injectDashboardType = () => {
    return inject(DashboardType, ref('status'))
}

export const useDashboardType = (data: Ref<{ type: string } | null>) => {
    provide(DashboardType, computed(() => data.value?.type || 'status'))
}