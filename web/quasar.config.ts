import { defineConfig } from "#q-app";

export default defineConfig(() => {
  return {
    boot: [],

    css: ["app.scss"],

    extras: ["roboto-font", "material-icons"],

    build: {
      target: {},

      typescript: {
        strict: true,
        vueShim: true
      },

      filenameBasedRouting: {
        exclude: ["**/components/**"]
      },

      vueRouterMode: "hash"
    },

    devServer: {
      open: true,
      proxy: {
        "/api": {
          target: "http://localhost:8080",
          changeOrigin: true,
          rewrite: (path: string) => path.replace(/^\/api/, "")
        }
      }
    },

    framework: {
      config: {},

      lang: "pt-BR",

      plugins: []
    },

    animations: [],

    ssr: {
      prodPort: 3000,
      middlewares: ["render"]
    },

    ssg: {},

    pwa: {
      workboxMode: "GenerateSW"
    },

    cordova: {},

    capacitor: {
      hideSplashscreen: true
    },

    electron: {
      preloadScripts: ["electron-preload"],

      inspectPort: 5858,

      bundler: "packager",

      packager: {},

      builder: {
        appId: "project-manager-clicksign"
      }
    },

    bex: {
      extraScripts: []
    }
  };
});
