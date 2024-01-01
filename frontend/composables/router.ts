import debounce from 'lodash.debounce'

export const useDebounce = (wait: number = 600) => debounce((emit: () => void) => { emit() }, wait)

export const defineRouteQuery = (key: string, selected: Ref<string[]>) => {
    const router = useRouter()
    const route = useRoute()

    let values: string[] = []

    if (Array.isArray(route.query[key])) {
        values = (route.query[key] as string[]).filter(c => !!c)
    }

    if (typeof route.query[key] === 'string') {
        values = [route.query[key] as string]
    }

    if (values.length) {
        selected.value = values
    }

    const debounced = useDebounce()

    return (inp: string[]) => {
        selected.value = inp

        debounced(() => {
            router.push({ name: route.name as string, query: { ...route.query, [key]: inp }, params: route.params })
        })
    }
}
