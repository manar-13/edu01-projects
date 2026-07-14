# GraphQL Profile

A personal profile page for the Reboot01 (01-edu) platform, built with plain
HTML, CSS, and JavaScript. It authenticates against the platform's signin
endpoint, then uses the **GraphQL API** to fetch and display the logged-in
user's data, including SVG statistic graphs.

## Features

- **Login** with username **or** email + password (Basic auth → JWT), with a
  clear error message for invalid credentials, and a **logout** button
- **Overview cards**: total XP, audit ratio (done vs received), and projects
  passed / failed / in progress
- **Statistics section** with three graphs, all hand-built with **SVG**
  (no chart libraries):
  1. XP over time — cumulative line + area chart
  2. XP by project — top 10 horizontal bar chart
  3. Audit ratio — donut chart (given vs received)
- Hover any point, bar, or ring to see exact values (native SVG tooltips)

## GraphQL usage

All three required query types are used (see `js/queries.js`):

| Type | Where |
|---|---|
| Normal | `user { id login ... }` |
| With arguments | `transaction(where: ..., order_by: ...)` |
| Nested | `object { name type }` inside `transaction` and `progress` |

## Project structure

```
.
├── index.html          # page structure: login + profile sections
├── style.css           # terminal-inspired dark theme
└── js/
    ├── config.js       # endpoints + small helpers
    ├── queries.js      # gql() client + the three GraphQL queries
    ├── auth.js         # login / logout (JWT handling)
    ├── profile.js      # fetches data and fills the info cards
    ├── graphs.js       # the three SVG graphs
    └── main.js         # event listeners + startup
```

## How to run locally

No build step and no dependencies are needed — only a small static file
server (opening `index.html` directly from the file system may block API
requests in some browsers).

**Option 1 — Python (pre-installed on macOS/Linux):**

```bash
# from the project root folder:
python3 -m http.server 8080
```

Then open <http://localhost:8080> in your browser.

**Option 2 — Node.js:**

```bash
npx serve .
```

Then open the URL it prints (usually <http://localhost:3000>).

**Option 3 — VS Code:** install the *Live Server* extension, right-click
`index.html` → *Open with Live Server*.

Log in with your Reboot01 credentials (username or email + password).

## Live demo

The profile is hosted on GitHub Pages:

**<https://manar-13.github.io/graphql/>**

Log in with your own Reboot01 credentials to see your profile.