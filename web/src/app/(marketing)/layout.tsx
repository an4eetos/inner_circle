import Link from "next/link";

export default function MarketingLayout({ children }: LayoutProps<"/">) {
  return (
    <>
      <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-4 py-5 sm:px-6">
        <Link href="/" className="flex items-center gap-2 text-lg font-semibold tracking-tight">
          <span aria-hidden className="inline-block h-3 w-3 rounded-full bg-mint" />
          Inner Circle
        </Link>
        <span className="rounded-full border border-border px-3 py-1 text-xs text-muted">
          Concept · first experiment soon
        </span>
      </header>
      <main className="flex-1">{children}</main>
      <footer className="border-t border-border">
        <div className="mx-auto flex w-full max-w-6xl flex-col gap-3 px-4 py-8 text-sm text-muted sm:flex-row sm:items-center sm:justify-between sm:px-6">
          <p>Friends suggest, doctors decide. Inner Circle does not give medical advice.</p>
          <nav className="flex gap-5">
            <Link href="/privacy" className="hover:text-text">
              Privacy
            </Link>
            <Link href="/terms" className="hover:text-text">
              Terms
            </Link>
          </nav>
        </div>
      </footer>
    </>
  );
}
