import type { Metadata } from "next";

export const metadata: Metadata = { title: "Terms of use" };

// Placeholder until the lawyer-reviewed terms for the launch country are ready (Phase 5).
export default function TermsPage() {
  return (
    <article className="mx-auto w-full max-w-3xl px-4 py-12 sm:px-6">
      <h1 className="text-3xl font-semibold tracking-tight">Terms of use</h1>
      <p className="mt-6 text-muted">
        Inner Circle is in the concept stage. Terms of use will be published here before the first
        participants join.
      </p>
      <p className="mt-4 text-muted">
        Inner Circle is not a medical service and does not give medical advice. Suggestions from
        friends are labelled as suggestions; always check changes to your treatment with your
        doctor.
      </p>
    </article>
  );
}
