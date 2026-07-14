# Audit Preparation — Questions & Answers

Likely questions from the audit mentor, with short answers in simple
English. Practice saying the answers OUT LOUD in your own words — don't
memorize word by word.

---

## 1. Login & Authentication

**Q: How does your login work?**
A: I send a POST request to the signin endpoint with the username (or
email) and password joined as `"username:password"`, encoded in base64,
inside a `Authorization: Basic ...` header. If the credentials are correct,
the server returns a JWT, and I save it in localStorage.

**Q: What is a JWT?**
A: A JSON Web Token — a signed token that proves who I am. It has 3 parts
separated by dots: header, payload (which contains my user id), and
signature. I send it with every GraphQL request, and the server uses it to
return only MY data.

**Q: What is the difference between authentication and authorization?**
A: Authentication is proving who you are (the login). Authorization is what
you are allowed to access — the GraphQL API uses my JWT to give me access
only to my own data.

**Q: What is the difference between Basic auth and Bearer auth?**
A: Basic auth sends username:password in base64 — I use it ONCE, at login.
Bearer auth sends the token — I use it with EVERY GraphQL request:
`Authorization: Bearer <token>`.

**Q: Is base64 encryption? Is it secure?**
A: No — base64 is just encoding, anyone can decode it back. It is safe here
only because the connection is HTTPS, which encrypts everything in transit.

**Q: What happens if I enter wrong credentials?**
A: The server responds with an error status, so `res.ok` is false, and I
show a clear message: "Invalid username/email or password." I also handle a
different case: if the network itself fails, the `catch` block shows a
"Network error" message instead.

**Q: How does logout work?**
A: I remove the JWT from localStorage and switch back to the login page.
Without the token, no GraphQL request can be made on my behalf.

**Q: Where do you store the token, and why?**
A: In localStorage, so the user stays logged in after a page refresh. On
page load I check: if a token exists, I go straight to the profile.

**Q: What if the saved token is expired?**
A: The GraphQL request fails, my `gql()` function throws an error, the
`catch` in `showProfilePage` logs the user out and shows "Session expired.
Please sign in again."

---

## 2. GraphQL

**Q: What is GraphQL?**
A: A query language for APIs. Instead of many fixed endpoints like REST,
there is ONE endpoint, and the client describes exactly which fields it
wants. The server returns only those fields — nothing more.

**Q: Show me your normal query.**
A: `USER_QUERY` in js/queries.js — it asks the `user` table for id, login,
attrs, auditRatio, totalUp and totalDown. No filters, no nesting — just
fields.

**Q: Show me your query with arguments.**
A: `XP_QUERY` — it filters the `transaction` table with `where`
(type equals "xp", and the event path is IN a list of two paths) and sorts
with `order_by: createdAt asc`. `PROGRESS_QUERY` also uses `where`.

**Q: Show me your nested query.**
A: Inside `XP_QUERY` and `PROGRESS_QUERY` I ask for `object { name type }`
— the object is a related table, and GraphQL returns its fields inside the
same answer. In `PROGRESS_QUERY` I even FILTER by a nested field:
`object: { type: { _eq: "project" } }`.

**Q: What do `_eq` and `_in` mean?**
A: `_eq` = equals one value. `_in` = the value is any one of a list. I use
`_in` to accept XP from two events at once.

**Q: Why is the GraphQL request a POST and not a GET?**
A: The query itself travels in the request BODY as JSON
(`{ "query": "..." }`), and bodies belong to POST requests.

**Q: How do you know the data you show is correct?**
A: I verified it against the platform. For example, my XP: I ran aggregate
queries in the console comparing filters — module only = 621,300, piscine
JS = 336,663, together = 957,963 ≈ 958 kB, exactly the platform's number.
All XP including Piscine Go would be 1,250,171, but the platform doesn't
count Piscine Go, so neither do I.

**Q: Why does your code do `user[0]`?**
A: GraphQL returns tables as arrays even for one row. Because of the JWT,
the user table returns exactly one user — me — so I take index 0.

---

## 3. The data / the cards

**Q: Where does the audit ratio come from?**
A: Directly from the user table: `auditRatio`, plus `totalUp` (XP from
audits I did) and `totalDown` (XP from audits done on my work). I round it
to 1 decimal with `.toFixed(1)` — same as the platform.

**Q: How do you calculate Total XP?**
A: I sum the `amount` of all my XP transactions using `reduce`, then format
the number as kB/MB, treating XP like bytes — the same convention as the
platform.

**Q: How do you count passed projects? Is it accurate?**
A: The progress table has one row per ATTEMPT, and a project can be retried.
In my data, groupie-tracker-filters appears 3 times (0.96, 0, 1.2). So I
group attempts by project name and keep only the BEST grade, then count:
grade ≥ 1 = passed, < 1 = failed, null = in progress. Result: 19 passed,
0 failed, 1 in progress — and the one in progress is this graphql project
itself.

---

## 4. SVG graphs

**Q: Why SVG?**
A: The subject requires SVG. It's also sharp at any zoom level, scales
responsively with the viewBox, and every shape can hold its own tooltip.

**Q: Did you use a chart library?**
A: No — every line, bar and circle is created by my own code with
`document.createElementNS`.

**Q: How does the XP-over-time graph work?**
A: I turn the transactions into cumulative points (a running total), then
two scaling functions convert data into pixels: x maps time to horizontal
position, y maps XP to vertical position. The line is one SVG path built
from "M" (move) and "L" (line to) commands. The shaded area is the same
path closed down to the bottom with "Z".

**Q: Why is the y formula "flipped" (H minus something)?**
A: In SVG, y = 0 is the TOP and y grows downward. So a bigger XP value must
produce a SMALLER y to appear higher on the screen.

**Q: How does the donut chart work?**
A: Two circles with a thick stroke and no fill. The back circle is the
"received" color. The front circle uses `stroke-dasharray`: one visible
dash whose length equals given's share of the circumference — that draws an
arc without arc math. I rotate it -90° so it starts at the top.

**Q: Are the graphs interactive / dynamic?**
A: Yes. Hovering any dot, bar, or ring shows a tooltip with exact values —
done with SVG `<title>` elements. And everything is computed from the
logged-in user's data, so any student sees their own graphs.

---

## 5. Hosting

**Q: Where is it hosted and how?**
A: On GitHub Pages: https://manar-13.github.io/graphql/ — the site is
static (only HTML/CSS/JS), GitHub's servers serve my files to any visitor.

**Q: If I log in on your hosted site, whose data do I see?**
A: Your own. The hosting only serves my FILES; the JavaScript runs in your
browser and fetches data from the Reboot01 API with YOUR JWT.

**Q: Why don't you need a backend?**
A: Because Reboot01's API is the backend. My site is only a frontend that
authenticates and queries it.

---

## 6. Code / JavaScript details

**Q: Why did you split the code into 6 files?**
A: One job per file: config, queries, auth, profile rendering, graphs, and
startup. It's easier to read, easier to find things, and it was also the
feedback I received — the load order in index.html matters: config first,
main last.

**Q: What does Promise.all do in your code?**
A: My three queries don't depend on each other, so `Promise.all` runs them
in parallel and waits for all of them. The page loads faster than running
them one by one.

**Q: What's the difference between `login` and `login()` in
addEventListener?**
A: `login` (no parentheses) passes the function itself, to be called later
when the click happens. `login()` would RUN it immediately during page load
— a bug.

**Q: Why `??` instead of `||` in renderProjects?**
A: Because a grade of 0 is a real grade. `0 || -1` gives -1 (zero counts as
false), but `0 ?? -1` keeps the 0. `??` only replaces null/undefined.

**Q: Why textContent and not innerHTML?**
A: textContent treats the value as plain text, never as HTML — so nothing
coming from the API can inject markup into my page. Safer.

---

## 7. Honest answers if you don't know something

It's okay to not know everything. Good honest answers:

- "I'm not sure about that exact detail — but let me show you where it is
  in the code and we can look together."
- "I didn't implement that because the subject didn't require it, but I
  understand the idea: ..."

Never invent an answer. Auditors respect honesty + willingness to look at
the code far more than fake confidence.

---
