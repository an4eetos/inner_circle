import type { Metadata } from "next";

export const metadata: Metadata = { title: "Privacy policy" };

// Placeholder until the lawyer-reviewed policy for the launch country is ready (Phase 5).
export default function PrivacyPage() {
  return (
    <article className="mx-auto w-full max-w-3xl px-4 py-12 sm:px-6">
      <h1 className="text-3xl font-semibold tracking-tight">Privacy policy</h1>
      <p className="mt-6 text-muted">
        Inner Circle is in the concept stage and does not collect any personal or health data yet.
        The full privacy policy will be published here before the first participants join.
      </p>
      <p className="mt-4 text-muted">
        Our commitments: health data is private by default, sharing is opt-in per organ and per
        layer, health data is never sold or used for advertising, and you can export or delete
        everything at any time.
      </p>
    </article>
  );
}
