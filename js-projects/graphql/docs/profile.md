# Notes — js/profile.js

The "director" file: it fetches all the data (using the queries), then sends
each piece to the right place — the cards and the graphs.

## showProfilePage()

```javascript
async function showProfilePage() {
    $("login-page").classList.add("hidden");
    $("profile-page").classList.remove("hidden");
    $("loading").classList.remove("hidden");
```

Switch the screens: hide the login section, show the profile section, and
show the "loading data_" message (the data is not here yet — fetching takes
a moment).

```javascript
    try {
        const [userData, xpData, progressData] = await Promise.all([
            gql(USER_QUERY),
            gql(XP_QUERY),
            gql(PROGRESS_QUERY),
        ]);
```

Run the THREE queries **at the same time**:

- Each `gql(...)` call starts a request to the GraphQL API
- `Promise.all([...])` means: "start all of them together, and wait until
  ALL are finished"
- Without it (three separate `await` lines), the queries would run one
  after another — three waits instead of one. With it, the total wait is
  only as long as the slowest query
- `const [a, b, c] = ...` unpacks the three answers, in the same order we
  sent them

> Audit question: *"Why Promise.all?"* — "The three queries don't depend on
> each other, so I run them in parallel. The page loads faster."

```javascript
        const user = userData.user[0];
        const transactions = xpData.transaction;
        const progresses = progressData.progress;
```

Unwrap the answers. GraphQL always returns tables as **arrays**, even when
there is only one row — the `user` query returns a list with exactly one
user (me), so `user[0]` takes that single item.

```javascript
        renderUser(user);
        renderXP(transactions);
        renderProjects(progresses);
        drawXpOverTime(transactions);
        drawXpByProject(transactions);
        drawAuditDonut(user);
        $("loading").classList.add("hidden");
```

Distribute the data: three functions fill the cards, three draw the graphs.
When everything is on screen, hide the "loading" message.

Notice how the same data is reused: `transactions` feeds one card AND two
graphs; `user` feeds two cards AND the donut. Fetch once, use many times.

```javascript
    } catch (err) {
        // If the token is old or broken, go back to login
        console.error(err);
        logout();
        showLoginError("Session expired. Please sign in again.");
    }
}
```

If anything fails — most commonly an **expired JWT** — we land here:
log the real error to the console (for developers), log the user out
(deletes the bad token), and show a clear message on the login page.

> This makes the app self-healing: a dead token never leaves the user stuck
> on a broken profile page.

## renderUser(user)

```javascript
    const attrs = user.attrs || {};
    const fullName = [attrs.firstName, attrs.lastName].filter(Boolean).join(" ");
    $("user-name").textContent = fullName || user.login;
```

`attrs` is a JSON column in the user table holding extra info (first name,
last name, ...). Building the full name, step by step:

1. `[attrs.firstName, attrs.lastName]` — a list of the two names, but any
   of them might be missing (`undefined`)
2. `.filter(Boolean)` — keep only the values that exist (throws away
   `undefined`, `null`, and empty strings)
3. `.join(" ")` — glue what's left with a space: `"Manar Mohamed"`

If BOTH are missing, `fullName` is an empty string → `|| user.login` falls
back to showing the username instead. The page never shows an empty title.

```javascript
    $("user-sub").textContent = `@${user.login} · id ${user.id}`;
    $("prompt-login").textContent = user.login;
```

Fill the small line under the name (`@manmohamed · id 5916`) and put the
login into the terminal prompt at the top
(`manmohamed@reboot01:~$ profile`).

```javascript
    $("audit-ratio").textContent = (user.auditRatio || 0).toFixed(1);
    $("audit-detail").textContent =
        `Done ${formatXP(user.totalUp)} · Received ${formatXP(user.totalDown)}`;
```

The audit card: `auditRatio` comes ready from the API; `.toFixed(1)` rounds
it to one decimal (`1.2666…` → `"1.3"` — same as the platform shows).
Below it, the two totals formatted as MB/kB.

> `textContent` (not `innerHTML`) is used everywhere on purpose: it treats
> the value as plain TEXT, never as HTML code. Safer — nothing coming from
> the API can inject markup into my page.

## renderXP(transactions)

```javascript
    const total = transactions.reduce((sum, t) => sum + t.amount, 0);
    $("total-xp").textContent = formatXP(total);
```

`reduce` boils a whole array down to one value:

- Start `sum` at `0` (the last argument)
- For every transaction `t`, do `sum + t.amount`
- The final `sum` = total XP

Then format it (`957963` → `"958 kB"`) and put it in the big card.

> Same job as a `for` loop adding into a variable — just written in one
> line. If `reduce` feels strange, translate it in your head to:
> "start at 0, add every amount".

## renderProjects(progresses)

The smartest logic in the file. Problem: the `progress` table has one row
per ATTEMPT, and a project can be tried several times. Example from my real
data: groupie-tracker-filters appears 3 times (grades 0.96, 0, 1.2). Counting
rows directly would say "2 failed" for a project I eventually PASSED.

Solution: keep only the **best grade per project**, then count projects.

```javascript
    const best = {};
    progresses.forEach((p) => {
        const name = p.object.name;
        const current = best[name] ?? -1;
        const grade = p.grade ?? -1;
        if (!(name in best) || grade > current) best[name] = p.grade;
    });
```

Walk over every attempt and remember the best one per name:

- `best` is an object used as a map: `{ "forum": 1.2, "graphql": null, ... }`
- `?? -1` — the `??` operator means "if the left side is null/undefined, use
  the right side". We treat "no grade yet" as `-1` ONLY for comparing, so
  any real grade beats it
- `!(name in best)` — first time seeing this project? store it
- `grade > current` — found a better attempt? replace it

After the loop, each project appears exactly ONCE with its best grade.

> Why `??` and not `||`? Because `0 || -1` gives `-1` (zero counts as
> false!) but `0 ?? -1` gives `0`. A grade of 0 is a real grade — we must
> not lose it. Nice detail to mention in the audit.

```javascript
    const grades = Object.values(best);
    const passed = grades.filter((g) => g !== null && g >= 1).length;
    const failed = grades.filter((g) => g !== null && g < 1).length;
    const inProgress = grades.filter((g) => g === null).length;
```

`Object.values` takes just the grades (ignores the names). Then three
filters count them:

- **passed** — has a grade and it's ≥ 1 (on the platform, grade ≥ 1 = pass)
- **failed** — has a grade and it's < 1 (best attempt still failing)
- **in progress** — grade is `null` = the project has no grade yet
  (fun fact: the "1 in progress" on my page is the graphql project itself)

```javascript
    $("projects-count").textContent = passed;
    $("projects-detail").textContent =
        `${passed} passed · ${failed} failed · ${inProgress} in progress`;
```

Put the numbers in the card: the big number is passed projects, the small
line shows the full breakdown.

> Ready audit answer: *"Is this number accurate?"* — "Yes. The progress
> table has one row per attempt; I verified my raw rows in the console.
> I group attempts by project name, keep the best grade, and count each
> project once: 19 passed, and graphql itself is the one in progress."

---
