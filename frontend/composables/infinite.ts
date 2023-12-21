import type { Ref, UnwrapRef } from "vue";
import type { UnwrapRefSimple } from "@vue/reactivity";

export const useInfinite = <T>(list: Ref<T[] | null>) => {
  const loaded = ref<T[]>([])
  const index = ref(1)

  watch(list, (l) => {
    if (!list.value?.length) return

    loaded.value = (l || []).slice(0, 1) as UnwrapRefSimple<T>[]
  }, { immediate: true })

  const load = ({ done }: any) => {
    const sum = (list.value || []).length
    if (!sum) { return done('ok') }

    if (sum === 1) { return done('empty') }

    const last = index.value
    const next = index.value + 2 > sum ? sum :  index.value + 2
    loaded.value = [...loaded.value, ...(list.value || []).slice(last, next)] as UnwrapRefSimple<T>[]

    index.value = next
    if (next === sum) {
      done('empty')
    } else {
      done('ok')
    }
  }


  return { load, loaded }
}

export const useInfiniteScroll = <T>(list: T[]) => {
  const loaded = ref<T[]>([])
  const index = ref(1)

  watch(list, (l) => {
    if (!list.length) return

    loaded.value = (l || []).slice(0, 1) as UnwrapRefSimple<T>[]
  }, { immediate: true })

  const load = ({ done }: any) => {
    const sum = (list || []).length
    if (!sum) { return done('ok') }

    if (sum === 1) { return done('empty') }

    const last = index.value
    const next = index.value + 2 > sum ? sum :  index.value + 2
    loaded.value = [...loaded.value, ...(list || []).slice(last, next)] as UnwrapRefSimple<T>[]

    index.value = next
    if (next === sum) {
      done('empty')
    } else {
      done('ok')
    }
  }


  return { load, loaded }
}
