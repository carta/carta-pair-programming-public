# React + TypeScript + Vite

## Package manager: pnpm

This frontend uses [pnpm](https://pnpm.io). The version is pinned in `package.json` (`packageManager`). Running the app with Docker (`docker compose up --build` from the repo root) needs nothing else. To work on the frontend outside Docker:

1. Use Node 22.13 or later, and run `corepack enable` once per Node install. On Node 25 or later, Corepack is no longer bundled: run `npm install -g corepack` first.
2. If your checkout was installed with npm before, delete `frontend/node_modules` once. pnpm installs over an npm-built `node_modules` without warning and keeps npm's leftover copies.
3. From `frontend/`:

   ```bash
   pnpm install --frozen-lockfile    # install exactly what pnpm-lock.yaml records
   pnpm run dev                      # Vite dev server on :5173
   pnpm run build                    # tsc -b && vite build
   pnpm run lint                     # eslint
   pnpm add <package>                # add a dependency and update pnpm-lock.yaml
   pnpm exec shadcn add <component>  # add a shadcn/ui component
   ```

What differs from npm:

- pnpm reads its settings from `pnpm-workspace.yaml`. It ignores settings in `package.json` and in `.npmrc`, except registry and auth lines.
- A dependency's install script runs only if `allowBuilds` in `pnpm-workspace.yaml` sets it to `true`. `pnpm install` fails with `ERR_PNPM_IGNORED_BUILDS` on any dependency that has an install script but no entry. Add it as `true` if it needs the script, or `false` if it doesn't.
- pnpm won't pick a version published less than 24 hours ago (`minimumReleaseAge`). A range or a tag such as `latest` resolves to the newest version older than that. An exact newer version is still added, but pnpm then writes `minimumReleaseAgeExclude` entries into `pnpm-workspace.yaml`. Commit them, or `pnpm install --frozen-lockfile` (and the Docker build) fails with `ERR_PNPM_MINIMUM_RELEASE_AGE_VIOLATION`.
- Pass script arguments without `--`, for example `pnpm run dev --port 3000`. pnpm passes a literal `--` through, and Vite ignores every flag after it.
- After a dependency change, rebuild the Docker frontend's `node_modules` volume (see the root README).

---

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Oxc](https://oxc.rs)
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/)

## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...

      // Remove tseslint.configs.recommended and replace with this
      tseslint.configs.recommendedTypeChecked,
      // Alternatively, use this for stricter rules
      tseslint.configs.strictTypeChecked,
      // Optionally, add this for stylistic rules
      tseslint.configs.stylisticTypeChecked,

      // Other configs...
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
