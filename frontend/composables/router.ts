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

    return (inp: string[]) => {
        if (equal(selected.value, inp)) return;

        selected.value = inp

        router.push({ name: route.name as string, query: { ...route.query, [key]: inp }, params: route.params })
    }
}

const equal = <T>(a: T[], b: T[]): boolean => {
    if (a.length !== b.length) return false

    return a.every(item => b.includes(item))
}