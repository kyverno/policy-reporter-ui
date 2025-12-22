import {Severity} from "~/modules/core/types";
import {ShowedSeverities} from "~/modules/core/provider/dashboard";
import type {Ref} from "vue";

export const useSeveritiesProvider = (data?: Ref<{ severities: Severity[] } | null>) => {
    provide(ShowedSeverities, computed(() => {
        const severities = data?.value?.severities
        if (severities && severities.length) {
            return [Severity.UNKNOWN,Severity.INFO,Severity.LOW,Severity.MEDIUM,Severity.HIGH,Severity.CRITICAL].reduce<Severity[]>((acc, s) => {
                if (severities.includes(s)) { return [...acc, s] }

                return acc
            }, [])
        }

        return [Severity.UNKNOWN,Severity.INFO,Severity.LOW,Severity.MEDIUM,Severity.HIGH,Severity.CRITICAL]
    }))
}
export const useSeveritiesInjection = () => {
    return inject(ShowedSeverities, ref([Severity.UNKNOWN,Severity.INFO,Severity.LOW,Severity.MEDIUM,Severity.HIGH,Severity.CRITICAL]))
}

export const useSeverityFilter = () => {
    const severities = useSeveritiesInjection()

    if (severities.value.length === 6) {
        return ref(undefined)
    }

    return severities
}