import {type Dashboard, Status} from "~/modules/core/types";
import {ShowedStatus} from "~/modules/core/provider/dashboard";
import type {Ref} from "vue";

export const useStatusProvider = (data?: Ref<{ status: Status[] } | null>) => {
    provide(ShowedStatus, computed(() => {
        const status = data?.value?.status
        if (status && status.length) {
            return [Status.SKIP, Status.PASS, Status.WARN, Status.FAIL, Status.ERROR, Status.SUMMARY].reduce<Status[]>((acc, s) => {
                if (status.includes(s)) { return [...acc, s] }

                return acc
            }, [])
        }

        return [Status.SKIP, Status.PASS, Status.WARN, Status.FAIL, Status.ERROR]
    }))
}
export const useStatusInjection = () => {
    return inject(ShowedStatus, ref([Status.SKIP, Status.PASS, Status.WARN, Status.FAIL, Status.ERROR]))
}

export const useStatusFilter = () => {
    const status = useStatusInjection()

    if (status.value.length === 5) {
        return ref(undefined)
    }

    return status
}