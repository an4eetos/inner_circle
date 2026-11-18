# Inner Circle web

Next.js (App Router, TypeScript, Tailwind). In Phase 0 it serves the marketing landing page;
later phases add invite links (`/i/{code}`), the doctor portal and the admin area.
See `docs/architecture/` for the plan.

```bash
npm ci
npm run dev        # http://localhost:3000
npm run gen:api    # regenerate src/gen/api/schema.d.ts from ../api/openapi.yaml
npm run lint && npm run typecheck && npm run build
```

API access goes through `src/lib/api.ts` (typed by the generated schema) and runs on the
server only; `API_BASE_URL` is read from the environment (see `.env.example`).
