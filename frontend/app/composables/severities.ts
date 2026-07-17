import {Severity} from "~/types/core";
import {ShowedSeverities} from "~/provider/dashboard";
import type {Ref} from "vue";

export const useSeveritiesProvider = (data?: Ref<{ filter: { severities: Severity[] } } | null>) => {
    provide(ShowedSeverities, computed(() => {
        const severities = data?.value?.filter.severities
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