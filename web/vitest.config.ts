import vue from "@vitejs/plugin-vue";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: [
      {
        find: "@",
        replacement: fileURLToPath(new URL("./src", import.meta.url))
      },
      { find: /^quasar$/, replacement: "quasar/dist/quasar.client.js" }
    ]
  },
  test: {
    environment: "happy-dom",
    include: ["src/**/*.test.ts"],
    css: false
  }
});
