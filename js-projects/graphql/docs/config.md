# Notes — js/config.js

The smallest file. It holds the **API addresses** and two **helpers** that
every other file uses. It loads FIRST (see the `<script>` order in
`index.html`), so everything here is available to the rest of the code.

## The URLs

```javascript
const DOMAIN = "learn.reboot01.com";
const SIGNIN_URL = `https://${DOMAIN}/api/auth/signin`;
const GRAPHQL_URL = `https://${DOMAIN}/api/graphql-engine/v1/graphql`;
```

Three constants (`const` = a value that never changes):

- `DOMAIN` — the school's domain, written once
- `SIGNIN_URL` — where we send username+password to get the JWT
  (used in `auth.js`)
- `GRAPHQL_URL` — where we send all GraphQL queries
  (used in `queries.js`)

The backticks `` ` `` make a **template string**: `${DOMAIN}` gets replaced
by the value of the variable. So `SIGNIN_URL` becomes
`https://learn.reboot01.com/api/auth/signin`.

> Why put the domain in ONE variable? If the school ever changes its domain,
> I fix one line and the whole project follows. This is the rule:
> **don't repeat yourself**.

> These two URLs come directly from the project subject: the signin endpoint
> and the GraphQL endpoint, with `((DOMAIN))` replaced by our school's
> real domain.

## $ — the tiny selector helper

```javascript
// Shortcut: $("id") instead of document.getElementById("id")
const $ = (id) => document.getElementById(id);
```

An **arrow function** stored in a constant named `$`. It takes an element's
id and returns that element from the page.

Without it: `document.getElementById("total-xp")`
With it: `$("total-xp")`

Same result, much shorter — and we grab elements dozens of times in this
project.

> `$` is just a normal variable name in JavaScript — nothing magical.
> (The famous jQuery library used the same idea, but this is OUR one-line
> version, no library involved.)

## formatXP(amount)

```javascript
// Turn 154000 into "154 kB", like the platform shows XP
function formatXP(amount) {
    if (amount >= 1_000_000) return (amount / 1_000_000).toFixed(2) + " MB";
    if (amount >= 1_000) return Math.round(amount / 1_000) + " kB";
    return amount + " B";
}
```

XP amounts come from the API as plain numbers like `957963` — hard to read.
This function turns them into the same style the platform uses:

| Input | Rule that fires | Output |
|---|---|---|
| `1400000` | ≥ 1,000,000 → divide by a million, keep 2 decimals | `"1.40 MB"` |
| `957963` | ≥ 1,000 → divide by a thousand, round | `"958 kB"` |
| `250` | small number | `"250 B"` |

Line by line:

- `1_000_000` — the underscores are only for human eyes; JavaScript reads it
  as `1000000`
- `.toFixed(2)` — keep exactly 2 digits after the decimal point (`1.4` →
  `"1.40"`)
- `Math.round(...)` — round to the nearest whole number (`957.963` → `958`)
- The `if`s are checked top to bottom: big numbers exit at MB, medium at kB,
  and anything else falls through to the last line (B)

> Why B / kB / MB? The platform treats XP like **bytes** (1000 B = 1 kB,
> 1000 kB = 1 MB), so I follow the same convention — that's why my
> "958 kB" matches the platform's "958.0K" exactly.

---
