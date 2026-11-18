# Inner Circle mobile (Flutter: iOS + Android)

```bash
make mobile-run FLAVOR=dev        # from the repo root; needs `make dev` running for the API
cd mobile && flutter test && flutter analyze
make gen-mobile                   # l10n + drift codegen
make gen-dart                     # API client in packages/api_client (from api/openapi.yaml)
```

- **Config** comes from `config/<flavor>.json` via `--dart-define-from-file` (`APP_ENV`, `API_BASE_URL`).
- **Android flavors** `dev` / `staging` / `prod` install side by side (`com.innercircle.app.dev`, `.stg`, none).
  Dev allows plain HTTP to the local API; `make mobile-run` runs `adb reverse tcp:18080 tcp:18080`.
- **iOS** uses one bundle ID (`com.innercircle.app`) for now; per-flavor schemes and bundle IDs need Xcode
  (see Phase 0 notes). Select the environment with the config file.
- **Structure**: `lib/app` (config, theme, router, providers, drift DB), `lib/features/<feature>/{data,presentation}`,
  `lib/l10n` (ARB strings; generated code in `lib/l10n/gen`).
- Generated code (`*.g.dart`, `lib/l10n/gen`, `packages/api_client/lib`) is committed; CI fails if it is stale.
