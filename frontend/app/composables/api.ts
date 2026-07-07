import { CoreAPI } from "~/core/api";
import { FetchError } from 'ofetch'

type Callback<T, R = T> = (api: CoreAPI) => Promise<T | R>

export type APIResult<T, R = T> = {
    data: Ref<T | R>,
    error: Ref<FetchError | null>
    pending: Ref<boolean>
    refresh: () => Promise<void>
}

export const useAPI = <T, R = T>(callback: Callback<T, R>, options?: { default?: () => T, finally?: () => void }): APIResult<T, R> => {
    const { $coreAPI } = useNuxtApp()

    const pending = ref<boolean>(false)
    const error = ref<FetchError | null>(null)
    const data = ref<T | R | null>(null) as Ref<T | R>

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
            console.error('API error', err)

            if (options?.default) {
                data.value = options.default()
            } else {
                data.value = null as unknown as T | R
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
