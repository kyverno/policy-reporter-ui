import { defineNuxtModule, addImportsDir, addPlugin } from "@nuxt/kit";
import { resolve, join } from "pathe";
import type { Nuxt } from "@nuxt/schema";

export default defineNuxtModule({
  setup(options: any, nuxt: Nuxt) {
    nuxt.hook("components:dirs", (dirs) => {
      dirs.push({
        path: join(__dirname, "components"),
      });
    });

    addImportsDir(resolve(__dirname, "./composables"));
    addPlugin(resolve(__dirname, "./plugins/apis.ts"));

    // nuxt.hook("pages:extend", (pages) => {
    //   pages.push({
    //     name: "dashboard",
    //     path: "/",
    //     file: resolve(__dirname, "./pages/index.vue"),
    //   });
    // });
  },
});
