import Image from "next/image";

const ideas = [
  {
    title: "3D body map",
    body: "A see-through body with markers on the organs you are working on: gallbladder, pancreas, heart, brain, gut, a knee in rehab. Colour shows the track: Health, Biohacking or Longevity.",
  },
  {
    title: "A program per organ",
    body: "Tap a marker to open its short-term program (this week's checklist) and long-term plan (phases over months, with a goal).",
  },
  {
    title: "Your inner circle",
    body: "Friends you accept can follow a program, cheer milestones, nudge you when you go quiet, and suggest recipes, swaps or questions for your doctor.",
  },
  {
    title: "Privacy per organ",
    body: "Each organ has its own audience and its own level of detail. Your knee rehab can be public while your medication stays with your closest people.",
  },
];

const screens = [
  {
    src: "/screens/01-body-map.png",
    title: "Body map",
    caption:
      "Height, weight and BMI next to the body, a marker per organ, a summary card for the selected program.",
  },
  {
    src: "/screens/03-gallbladder-program.png",
    title: "Gallbladder care",
    caption:
      "Nightly check-in with a streak, today's plate, weekly checklist, 6-month plan, notes from friends.",
  },
  {
    src: "/screens/04-quit-sugary-drinks.png",
    title: "Quit sugary drinks",
    caption:
      "Craving SOS with swaps, cravings chart, step-down plan, a Coke-free bet with a friend.",
  },
  {
    src: "/screens/05-circle-feed.png",
    title: "Circle feed",
    caption:
      "Friends' progress rings, milestones to cheer, nudges for friends who went quiet, questions to answer.",
  },
  {
    src: "/screens/06-sharing-controls.png",
    title: "Who sees this organ",
    caption: "Pick the audience and what they see: progress, meds and meals, results, comments.",
  },
  {
    src: "/screens/02-body-map-pancreas.png",
    title: "Every marker opens a program",
    caption: "Pancreas, heart, brain, gut: each with its own plan and its own circle.",
  },
];

const mechanics = [
  [
    "Nightly check-in",
    "One tap for a daily medication or habit. If you miss it, the circle can get a ping, but only when you have turned that on.",
  ],
  ["Nudge", "A friend reminds you, but only on the programs where you asked for nudges."],
  ["Cheer", "Light support on milestones and streaks."],
  [
    "Bet / challenge",
    "Two friends, one goal, a stake: first to 30 days Coke-free, loser buys a gym month.",
  ],
  [
    "Recommendation",
    "A friend suggests a recipe, a swap or a question for your doctor. Labelled as a suggestion, never as medical advice.",
  ],
  ["Copy program", "Start a friend's protocol as your own."],
];

const principles = [
  ["Programs over metrics", "The daily reason to open the app is the check-in, not the 3D model."],
  [
    "Treatment first, optimization second",
    "People on a real course (medication, diet, rehab) have higher stakes and a daily habit.",
  ],
  [
    "Privacy by design",
    "Health data is sensitive by default. Sharing is opt-in, per organ and per layer.",
  ],
  ["Friends suggest, doctors decide", "Social recommendations never replace medical advice."],
];

export default function LandingPage() {
  return (
    <>
      <section className="mx-auto grid w-full max-w-6xl items-center gap-10 px-4 pt-10 pb-16 sm:px-6 md:grid-cols-[1fr_1.1fr] md:pt-16">
        <div>
          <p className="mb-4 text-sm font-medium tracking-widest text-mint uppercase">
            Health · Biohacking · Longevity
          </p>
          <h1 className="text-4xl leading-tight font-semibold tracking-tight sm:text-5xl">
            See inside.
            <br />
            Bring your inner circle.
          </h1>
          <p className="mt-6 max-w-xl text-lg text-muted">
            Most health apps show you numbers. Inner Circle shows you{" "}
            <span className="text-text">your body</span> and what you are doing about it: a program
            for every organ you are working on, and the friends you choose keeping you on track.
          </p>
          <div className="mt-8 flex flex-wrap gap-3 text-sm">
            <Pill color="bg-amber">Health</Pill>
            <Pill color="bg-periwinkle">Biohacking</Pill>
            <Pill color="bg-mint">Longevity</Pill>
          </div>
        </div>
        <Image
          src="/screens/00-overview.png"
          alt="Inner Circle concept screens: body map, organ program, circle feed and sharing controls"
          width={1600}
          height={1000}
          priority
          className="h-auto w-full rounded-2xl border border-border"
        />
      </section>

      <section className="mx-auto w-full max-w-6xl px-4 py-12 sm:px-6">
        <h2 className="text-2xl font-semibold tracking-tight">The idea</h2>
        <div className="mt-6 grid gap-4 sm:grid-cols-2">
          {ideas.map((i) => (
            <article key={i.title} className="rounded-2xl border border-border bg-surface p-6">
              <h3 className="font-semibold">{i.title}</h3>
              <p className="mt-2 text-muted">{i.body}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="mx-auto w-full max-w-6xl px-4 py-12 sm:px-6">
        <h2 className="text-2xl font-semibold tracking-tight">Screens</h2>
        <div className="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {screens.map((s) => (
            <figure key={s.src} className="rounded-2xl border border-border bg-surface p-4">
              <div className="relative aspect-[9/16] overflow-hidden rounded-xl bg-surface-2">
                <Image
                  src={s.src}
                  alt={s.title}
                  fill
                  sizes="(min-width: 1024px) 33vw, (min-width: 640px) 50vw, 100vw"
                  className="object-cover object-top"
                />
              </div>
              <figcaption className="mt-4">
                <p className="font-semibold">{s.title}</p>
                <p className="mt-1 text-sm text-muted">{s.caption}</p>
              </figcaption>
            </figure>
          ))}
        </div>
      </section>

      <section className="mx-auto w-full max-w-6xl px-4 py-12 sm:px-6">
        <h2 className="text-2xl font-semibold tracking-tight">Core mechanics</h2>
        <dl className="mt-6 divide-y divide-border rounded-2xl border border-border bg-surface">
          {mechanics.map(([name, body]) => (
            <div key={name} className="grid gap-1 p-5 sm:grid-cols-[14rem_1fr] sm:gap-6">
              <dt className="font-semibold">{name}</dt>
              <dd className="text-muted">{body}</dd>
            </div>
          ))}
        </dl>
      </section>

      <section className="mx-auto w-full max-w-6xl px-4 pt-12 pb-20 sm:px-6">
        <h2 className="text-2xl font-semibold tracking-tight">Principles</h2>
        <div className="mt-6 grid gap-4 sm:grid-cols-2">
          {principles.map(([title, body]) => (
            <article key={title} className="rounded-2xl border border-border p-6">
              <h3 className="font-semibold text-mint">{title}</h3>
              <p className="mt-2 text-muted">{body}</p>
            </article>
          ))}
        </div>
      </section>
    </>
  );
}

function Pill({ color, children }: { color: string; children: React.ReactNode }) {
  return (
    <span className="inline-flex items-center gap-2 rounded-full border border-border px-4 py-2">
      <span aria-hidden className={`h-2.5 w-2.5 rounded-full ${color}`} />
      {children}
    </span>
  );
}
