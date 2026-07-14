# About the Project — GraphQL Profile

## What is this project?

A personal profile website that shows my own Reboot01 school data:
my XP, my audit ratio, my projects, and graphs of my progress.

The data is NOT stored in my project. It lives on the school's server,
and my website **asks for it** using the **GraphQL API**.

## The big picture (how it works)

```
[Login page] → send username+password → [Reboot01 signin API]
                                              │
                                              ▼
                                     returns a JWT (a token)
                                              │
                                              ▼
[Profile page] → send queries + JWT → [Reboot01 GraphQL API]
                                              │
                                              ▼
                                    returns MY data as JSON
                                              │
                                              ▼
                       JavaScript fills the cards + draws SVG graphs
```

## The 4 key concepts

### 1. Authentication (login)

- I send my username (or email) + password to the signin endpoint
- They are sent using **Basic auth**: the text `"username:password"`
  encoded in **base64** (a simple text encoding, NOT encryption)
- If the credentials are correct, the server returns a **JWT**

### 2. JWT (JSON Web Token)

- A long text string that works like a **key card**: it proves who I am
- I save it in the browser (`localStorage`) so the user stays logged in
- Every GraphQL request sends it in a header:
  `Authorization: Bearer <token>` — this is called **Bearer auth**
- The server reads the JWT and only returns data belonging to ME
  (this is **authorization**: what am I allowed to see)
- **Logout** = simply delete the token from localStorage

> Authentication = proving who you are (login).
> Authorization = what you are allowed to access (only your own data).

### 3. GraphQL

- A query language for APIs. Instead of fixed endpoints (REST),
  there is ONE endpoint, and I describe exactly what data I want
- Example: `{ user { id login } }` returns only id and login — nothing more
- The project requires all 3 query types, and I use them all:

| Type | Meaning | Where in my code |
|---|---|---|
| Normal | just ask for fields | `USER_QUERY` (user info) |
| With arguments | filter/sort with `where`, `order_by` | `XP_QUERY` (only xp type, only my module events, sorted by date) |
| Nested | ask for a related table inside another | `object { name type }` inside transaction and progress |

### 4. SVG graphs

- SVG = drawing shapes (lines, circles, rectangles, text) with code
- I build 3 graphs with pure SVG — no chart library:
  1. **XP over time** — a cumulative line: each XP transaction adds a point
  2. **XP by project** — horizontal bars for the top 10 projects
  3. **Audit ratio donut** — a ring comparing audits given vs received

## The data I display (and where it comes from)

| Card / Graph | Data source |
|---|---|
| Name, login, id | `user` table (`attrs` has the full name) |
| Total XP | sum of `transaction` rows with type `"xp"` from the module + JS piscine events |
| Audit ratio | `user.auditRatio`, `user.totalUp` (given), `user.totalDown` (received) |
| Projects passed | `progress` rows for projects, counting each project once by its BEST grade (retries don't count as fails) |

> Important detail: XP is filtered by event path
> `/bahrain/bh-module` + `/bahrain/bh-module/piscine-js`.
> I verified with aggregate queries that this matches the platform's
> total exactly (958 kB). Piscine Go XP is not counted by the platform.

## The files

```
index.html      the page skeleton: login section + profile section
style.css       terminal-style dark design (teal accent, mono font)
js/config.js    API URLs + small helpers
js/queries.js   the gql() function + the 3 GraphQL queries
js/auth.js      login / logout (JWT handling)
js/profile.js   fetches the data and fills the cards
js/graphs.js    draws the 3 SVG graphs
js/main.js      connects buttons + starts the app
```

## Hosting

- The site is **static** (only HTML/CSS/JS, no server code of mine)
- Hosted on **GitHub Pages**: the files are stored on GitHub's servers
  and served to any visitor at
  https://manar-13.github.io/graphql/
- Anyone can log in there with their own Reboot01 account and see
  their own profile (the JWT decides whose data is returned)

---
