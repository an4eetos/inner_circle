import "server-only";

import createClient from "openapi-fetch";

import type { paths } from "@/gen/api/schema";

/**
 * Typed client for the Go API, generated from api/openapi.yaml.
 * Server-only: the base URL and (later) service credentials never reach the browser.
 */
export const api = createClient<paths>({
  baseUrl: process.env.API_BASE_URL ?? "http://localhost:18080",
});

export type { components } from "@/gen/api/schema";
