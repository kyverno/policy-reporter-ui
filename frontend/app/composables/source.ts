import { SourceContext } from "~/provider/dashboard";

export const useSourceContext = (source: Ref<string | undefined>) => {
    provide(SourceContext, source)
}

export const injectSourceContext = () => inject(SourceContext, ref(undefined))