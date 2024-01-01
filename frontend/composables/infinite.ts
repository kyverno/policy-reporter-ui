import type { Ref, UnwrapRef } from "vue";
import type { UnwrapRefSimple } from "@vue/reactivity";

export const useInfinite = <T>(list: Ref<T[] | null>, defaultLoadings = 3) => {
  const loaded = ref<T[]>([])
  const index = ref(0)

  watch(list, (newValue, oldValue) => {
    const l = newValue || []
    const length = (l || []).length
    const oldLength = (oldValue || []).length

    if (!length || (length === loaded.value.length && loaded.value.every(i => (l.includes(i as T))))) return;

    if (length - oldLength === 1) {
      loaded.value = l.slice(0, length) as UnwrapRefSimple<T>[]
      index.value = length

      return
    }

    if (oldLength > 0 && oldLength < length) {
      loaded.value = l.slice(0, oldLength + 1) as UnwrapRefSimple<T>[]
      index.value = oldLength + 1

      return
    }

    if (oldLength > length) {
      loaded.value = l.slice(0, length) as UnwrapRefSimple<T>[]
      index.value = length

      return
    }

    let loadCounter = defaultLoadings
    if (length < defaultLoadings) {
      loadCounter = length
    }

    loaded.value = l.slice(0, loadCounter) as UnwrapRefSimple<T>[]
    index.value = loadCounter
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
