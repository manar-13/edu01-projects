# Notes — js/queries.js

The GraphQL heart of the project: ONE function that talks to the API, and
the THREE queries. This is the file the audit form checks directly:
*"Does the project have at least the mandatory queries (nested, normal and
using arguments)?"* — answer: yes, all three, in this file.

## gql(query) — the GraphQL client

```javascript
async function gql(query) {
    const token = localStorage.getItem("jwt");
    const res = await fetch(GRAPHQL_URL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + token,
        },
        body: JSON.stringify({ query }),
    });
```

One reusable function sends ANY query. Piece by piece:

- Get the saved JWT from localStorage
- `method: "POST"` — GraphQL requests are always POST: we SEND the query
  text in the body (unlike REST where reading usually uses GET)
- `"Content-Type": "application/json"` — tells the server "the body is
  JSON"
- `Authorization: "Bearer " + token` — **Bearer authentication**: the JWT
  goes with every request, and the server uses it to know WHO is asking →
  it only returns MY rows. That's why the query `{ user { id } }` returns
  one user, not the whole school
- `JSON.stringify({ query })` — GraphQL expects the body to be JSON shaped
  like `{ "query": "..." }`. The shorthand `{ query }` is the same as
  `{ query: query }`

> Login used **Basic** auth (username:password, once). Everything after
> uses **Bearer** auth (the token, every request). Two different schemes
> for two different jobs — great audit answer.

```javascript
    const data = await res.json();
    if (data.errors) throw new Error(data.errors[0].message);
    return data.data;
}
```

GraphQL answers always look like `{ "data": {...} }` or
`{ "errors": [...] }` (or both):

- If there are errors, `throw` — this jumps straight to the `catch` in
  `showProfilePage`, which logs the user out (this is how an expired token
  is handled)
- Otherwise return only the useful part: `data.data`

> Why is this function so valuable? **Reusability.** All three queries —
> and any test query in the browser console — go through this one function.
> During development I even used it in the console to compare XP filters.

## Query 1 — USER_QUERY (the NORMAL query)

```graphql
{
  user {
    id
    login
    attrs
    auditRatio
    totalUp
    totalDown
  }
}
```

The simplest query type: name the table (`user`), list the fields you want.
No filters, no nesting. The API returns ONLY these fields — that's the core
idea of GraphQL: *ask exactly for what you need*.

The fields:

| Field | What it is | Used for |
|---|---|---|
| `id`, `login` | my user id and username | identification header |
| `attrs` | JSON with extra info (first/last name) | showing my full name |
| `auditRatio` | ready-computed ratio | the audit card + donut center |
| `totalUp` | XP from audits I did for others | "Given" in the donut |
| `totalDown` | XP from audits done on my work | "Received" in the donut |

> Thanks to the JWT, `user` returns only the authenticated user — an array
> with exactly one item. That's why the code does `user[0]`.

## Query 2 — XP_QUERY (ARGUMENTS + NESTED together)

```graphql
{
  transaction(
    where: {
      type: { _eq: "xp" }
      event: { path: { _in: ["/bahrain/bh-module", "/bahrain/bh-module/piscine-js"] } }
    }
    order_by: { createdAt: asc }
  ) {
    amount
    createdAt
    path
    object {
      name
      type
    }
  }
}
```

The transaction table holds many row types (xp, up, down, level...). The
**arguments** in parentheses filter and sort:

- `where:` — only rows matching ALL these conditions:
  - `type: { _eq: "xp" }` — `_eq` = equals. Only XP rows
  - `event: { path: { _in: [...] } }` — `_in` = "is any value in this
    list". Only XP from two events: the Bahrain module AND the JS piscine
- `order_by: { createdAt: asc }` — sort oldest → newest (`asc` =
  ascending). The XP-over-time graph needs them in time order

The **nested** part: `object { name type }` — each transaction points to an
object (the project/exercise it came from). Instead of a separate request
per transaction, GraphQL joins the related table inside the same answer.
That's how the bar chart knows the project NAMES.

> **Why these two event paths?** I verified with aggregate queries in the
> console: module alone = 621,300 XP, piscine JS = 336,663. Together =
> 957,963 ≈ 958 kB — exactly the platform's number. All XP ever (including
> Piscine Go) is 1,250,171 — the platform does NOT count Piscine Go, so
> neither do I. I can demonstrate this live in the audit.

## Query 3 — PROGRESS_QUERY (NESTED with arguments)

```graphql
{
  progress(
    where: {
      object: { type: { _eq: "project" } }
      event: { path: { _eq: "/bahrain/bh-module" } }
    }
  ) {
    grade
    object {
      name
    }
  }
}
```

The progress table has one row per attempt at anything (exercises,
checkpoints, projects...). The filter keeps only what I need for the
projects card:

- `object: { type: { _eq: "project" } }` — this is filtering by a field of
  the NESTED table: "only rows whose related object is a project" (not an
  exercise). Filtering through a relationship — a nice GraphQL power
- `event: { path: { _eq: "/bahrain/bh-module" } }` — only the main module

And I ask back for `grade` (≥ 1 pass, < 1 fail, `null` = not graded yet)
plus the nested `object { name }` so retries of the same project can be
grouped by name in `renderProjects`.

## The audit checklist for this file

| Requirement | Where | Status |
|---|---|---|
| Normal query | `USER_QUERY` | ✅ |
| Query using arguments | `XP_QUERY` (`where`, `order_by`), `PROGRESS_QUERY` (`where`) | ✅ |
| Nested query | `object { ... }` inside both `XP_QUERY` and `PROGRESS_QUERY` | ✅ |

The subject also says the types can be used *"together or separately"* —
queries 2 and 3 use arguments AND nesting together.

## Mini GraphQL dictionary (for quick review)

| Word | Meaning |
|---|---|
| `where` | filter: which rows to return |
| `_eq` | equals one value |
| `_in` | equals ANY value in a list |
| `order_by` | sort the rows |
| `asc` / `desc` | ascending / descending |
| nested block `x { ... }` | bring fields from a RELATED table in the same answer |
| Bearer auth | sending the JWT with each request to prove identity |

---
