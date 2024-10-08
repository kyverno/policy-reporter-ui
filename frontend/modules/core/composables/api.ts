import { CoreAPI } from "../api";
import { FetchError } from 'ofetch'

type Callback<T> = (api: CoreAPI) => Promise<T>

export type APIResult<T> = {
    data: Ref<T | null>,
    error: Ref<FetchError | null>
    pending: Ref<boolean>
    refresh: () => Promise<void>
}

export const useAPI = <T>(callback: Callback<T>, options?: { default?: () => T, finally?: () => void }): APIResult<T> => {
    const { $coreAPI } = useNuxtApp()

    const pending = ref<boolean>(false)
    const error = ref<FetchError | null>(null)
    const data = ref<T | null>(null) as Ref<T | null>

    if (options?.default) {
        data.value = options.default()
    }

    const refresh = async () => {
        pending.value = true

        try {
            await  callback($coreAPI as CoreAPI).then((content) => {
                error.value = null
                data.value = content
            })
        } catch (err) {
            error.value = err as FetchError

            if (options?.default) {
                data.value = options.default()
            } else {
                data.value = null
            }
        } finally {
            pending.value = false

            options?.finally && options.finally()
        }
    }

    refresh()

    return { data, error, pending, refresh }
}

export const callAPI = <T>(callback: Callback<T>): Promise<T> => {
    const { $coreAPI } = useNuxtApp()

    return callback($coreAPI as CoreAPI)
}
