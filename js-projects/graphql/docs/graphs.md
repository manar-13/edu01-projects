# Notes — js/graphs.js

The biggest file: it draws the **3 SVG graphs**. No chart library — every
line, bar, and circle is created by our own code.

## One idea to understand first: how SVG coordinates work

- An SVG is a drawing area. Position `(0, 0)` is the **TOP-LEFT** corner
- `x` grows to the **right** (normal)
- `y` grows **DOWNWARD** (opposite of math class!)

> This is why the y-scaling formulas below start with `H - ...`:
> a BIG XP value must become a SMALL y number (near the top).

All three graphs use `viewBox: "0 0 W H"` — this means "my drawing space is
W wide and H tall". The real size on screen stretches with CSS
(`width: 100%`), but positions inside stay proportional. That's why the
graphs are responsive.

## svgEl(tag, attrs) — the factory helper

```javascript
function svgEl(tag, attrs) {
    const el = document.createElementNS("http://www.w3.org/2000/svg", tag);
    for (const key in attrs) el.setAttribute(key, attrs[key]);
    return el;
}
```

Creates any SVG element in one call.

- SVG elements must be created with `createElementNS` (NS = namespace) —
  the browser needs to know this is an SVG `<circle>`, not a normal HTML tag
- The `for...in` loop copies every property from the `attrs` object onto the
  element as an attribute
- Example: `svgEl("circle", { cx: 10, cy: 20, r: 5 })` builds
  `<circle cx="10" cy="20" r="5"/>`

> Without this helper we would repeat 3 lines for every single shape.
> With it: one line per shape.

## Graph 1 — drawXpOverTime(transactions)

Draws the cumulative XP line: every transaction adds a point, and the line
only goes UP.

```javascript
    const holder = $("graph-xp-time");
    holder.innerHTML = "";
    if (transactions.length === 0) {
        holder.textContent = "No XP data yet.";
        return;
    }
```

Find the box in the HTML where the graph lives, and empty it (so re-drawing
never stacks two graphs). If there is no data, show a friendly message and
stop — never draw an empty graph.

```javascript
    const W = 640, H = 300;
    const pad = { top: 20, right: 20, bottom: 40, left: 60 };
```

The drawing space: 640×300. `pad` = empty margins inside the drawing, kept
free for the labels (left needs 60px for "958 kB", bottom needs 40px for
dates).

```javascript
    let sum = 0;
    const points = transactions.map((t) => {
        sum += t.amount;
        return { time: new Date(t.createdAt).getTime(), xp: sum };
    });
```

Turn transactions into **cumulative** points. `sum` keeps growing: if the
amounts are 100, 50, 75 → the points hold 100, 150, 225. Each point =
`{ time, xp so far }`. `getTime()` converts a date into a plain number
(milliseconds) so we can do math with it.

> "Cumulative" is the key word: the graph shows my TOTAL so far at each
> moment, not each separate gain. That's why it looks like climbing stairs.

```javascript
    const minT = points[0].time;
    const maxT = points[points.length - 1].time;
    const maxXP = sum;
```

The data's borders: first date, last date, and final total. (The
transactions are already sorted by date — the query used
`order_by: { createdAt: asc }` — so first item = oldest.)

```javascript
    const x = (t) => pad.left + ((t - minT) / (maxT - minT || 1)) * (W - pad.left - pad.right);
    const y = (v) => H - pad.bottom - (v / maxXP) * (H - pad.top - pad.bottom);
```

**The most important two lines in the file.** They translate data values
into pixel positions:

- `x(t)`: "how far along the time range is this moment, as a fraction 0→1"
  (`(t - minT) / (maxT - minT)`), then stretch that fraction across the
  drawable width, and shift right past the left padding
- `y(v)`: same idea for XP, but flipped with `H - pad.bottom - ...` because
  SVG y grows downward: XP = 0 sits at the bottom, XP = max sits at the top
- `|| 1` protects against dividing by zero if all transactions happen at the
  same instant

> Audit answer for "how did you position the points?": "I map each value
> to a fraction between 0 and 1, then multiply by the drawable area.
> This is called linear scaling."

```javascript
    [0, 0.5, 1].forEach((f) => {
        const gy = y(maxXP * f);
        svg.appendChild(svgEl("line", { ... }));
        const label = svgEl("text", { ... "text-anchor": "end" });
        label.textContent = formatXP(maxXP * f);
        ...
    });
```

Draw 3 horizontal grid lines at 0%, 50% and 100% of the max XP, each with a
label on its left ("0 B", "479 kB", "958 kB"). `text-anchor: "end"` means
the text's END sticks to the given x — so labels align nicely to the right.

The next block adds the first date (bottom-left) and last date
(bottom-right) with `toLocaleDateString()` — formats the date in the
visitor's local style.

```javascript
    const linePath = points.map((p, i) => `${i === 0 ? "M" : "L"}${x(p.time)},${y(p.xp)}`).join(" ");
```

SVG paths are drawn with a mini-language:

- `M x,y` = **M**ove the pen to a point (no drawing)
- `L x,y` = draw a **L**ine to a point

So the first point gets `M` (place the pen) and every other point gets `L`
(draw to here). Result: `"M60,260 L95,240 L130,180 ..."` — the whole line in
one string.

```javascript
    const areaPath = linePath +
        ` L${x(maxT)},${H - pad.bottom} L${x(minT)},${H - pad.bottom} Z`;
```

The shaded area under the line = the SAME path, then: line down to the
bottom-right corner, line across to the bottom-left corner, and `Z` =
close the shape. A closed shape can be filled with color.

```javascript
    svg.appendChild(svgEl("path", { d: areaPath, class: "area" }));
    svg.appendChild(svgEl("path", { d: linePath, class: "line" }));
```

Add the area FIRST, then the line on top (later elements draw on top of
earlier ones — like layers). The colors come from CSS classes.

```javascript
    points.forEach((p, i) => {
        const dot = svgEl("circle", { cx: x(p.time), cy: y(p.xp), r: 3.5, class: "dot" });
        const tip = svgEl("title", {});
        tip.textContent = `${transactions[i].object?.name || "XP"} · ...`;
        dot.appendChild(tip);
        svg.appendChild(dot);
    });
```

A small circle on every point. Inside each circle we put a `<title>`
element — in SVG, a `<title>` child becomes a **native browser tooltip**:
hover the dot and the browser shows the text. Free interactivity, zero
JavaScript events!

- `transactions[i].object?.name` — the `?.` is **optional chaining**: if
  `object` is missing (null), don't crash, just give `undefined`
- `|| "XP"` — and in that case, show "XP" as a fallback name

Finally `holder.appendChild(svg)` puts the finished drawing into the page.

## Graph 2 — drawXpByProject(transactions)

Horizontal bars for the top 10 projects by XP.

```javascript
    const byProject = {};
    transactions.forEach((t) => {
        const name = t.object?.name || t.path.split("/").pop();
        byProject[name] = (byProject[name] || 0) + t.amount;
    });
```

**Group** the XP by project name using a plain object as a counter:

- Get the name from the nested `object`, or (fallback) take the last piece
  of the path: `"/bahrain/bh-module/forum"` → split by `/` → take the last
  part → `"forum"`
- `(byProject[name] || 0) + t.amount` — "the total so far (or 0 if this is
  the first time we see this name) plus this amount"
- Result example: `{ "make-your-game": 147000, "forum": 76000, ... }`

```javascript
    const top = Object.entries(byProject)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 10);
```

Three steps chained:

1. `Object.entries` turns the object into pairs:
   `[["make-your-game", 147000], ["forum", 76000], ...]`
2. `.sort((a, b) => b[1] - a[1])` sorts by the amount (`[1]` = second item
   of the pair), biggest first
3. `.slice(0, 10)` keeps only the first 10

```javascript
    const H = pad.top + pad.bottom + top.length * rowH;
    const maxXP = top[0][1];
```

The height is calculated from the number of bars (34px per row) — if a
student has fewer than 10 projects, the graph shrinks. `maxXP` = the biggest
project (first after sorting); its bar will be full-width and all others are
sized relative to it.

```javascript
    top.forEach(([name, amount], i) => {
        const yPos = pad.top + i * rowH;
        const barW = (amount / maxXP) * (W - pad.left - pad.right);
```

For each project (the `[name, amount]` syntax unpacks the pair):

- `yPos` — each row sits 34px lower than the previous one
- `barW` — the bar's width as a fraction of the biggest bar (same linear
  scaling idea as Graph 1)

Then three shapes per row:

1. **The name** on the left (`text-anchor: "end"` aligns it to the right,
   next to the bar). Long names get cut at 17 letters + `…`
2. **The bar**: a `rect` starting at `pad.left`, with width `barW` and
   rounded corners (`rx: 5`), plus a `<title>` tooltip with the exact value
3. **The value** text just after the bar's end (`pad.left + barW + 8`)

## Graph 3 — drawAuditDonut(user)

A ring showing audits **given** vs **received**, with the ratio in the
middle. This one uses a clever trick instead of complicated arc math.

```javascript
    const given = user.totalUp || 0;
    const received = user.totalDown || 0;
    const total = given + received;
```

The data comes from the user table: `totalUp` = XP from audits I did for
others, `totalDown` = XP from audits done on my projects. (`|| 0` = if the
value is missing, use 0 so the math never breaks.)

```javascript
    const r = 90;
    const stroke = 30;
    const circumference = 2 * Math.PI * r;
    const givenPart = (given / total) * circumference;
```

- The donut is just a **circle with a very thick border** (stroke = 30px)
  and no fill
- `circumference` = the full length around the circle (2πr — school math!)
- `givenPart` = the share of that length that belongs to "given".
  Example: given is 56% of total → givenPart = 56% of the circumference

```javascript
    const ringReceived = svgEl("circle", {
        cx, cy, r, fill: "none",
        stroke: "#33415E", "stroke-width": stroke,
    });
```

First layer: a full circle in the dim blue color = "received". It acts as
the background of the ring.

```javascript
    const ringGiven = svgEl("circle", {
        cx, cy, r, fill: "none",
        stroke: "#56E1CE", "stroke-width": stroke,
        "stroke-dasharray": `${givenPart} ${circumference - givenPart}`,
        transform: `rotate(-90 ${cx} ${cy})`,
        "stroke-linecap": "round",
    });
```

Second layer on top: the SAME circle in teal, but only partially visible.
**The trick is `stroke-dasharray`** — it makes a dashed border:
"draw X pixels, skip Y pixels". We set it to
`givenPart, (circumference - givenPart)` = ONE visible dash exactly as long
as the given share, then invisible for the rest. Result: an arc, without any
arc math.

- `rotate(-90 ...)` — a circle's border starts at the RIGHT (3 o'clock);
  rotating -90° makes our arc start at the TOP (12 o'clock), which looks
  natural
- `stroke-linecap: "round"` — rounded ends for the arc

Both rings have `<title>` tooltips with the exact values.

```javascript
    const ratioText = svgEl("text", {
        x: cx, y: cy + 8, "text-anchor": "middle", class: "donut-center",
    });
    ratioText.textContent = (user.auditRatio || 0).toFixed(1);
```

The big number in the hole of the donut: `auditRatio` comes ready-made from
the API, rounded to 1 decimal (`1.2666...` → `"1.3"`) — the same number the
platform shows. `text-anchor: "middle"` centers the text on the given x.

The last block draws the **legend**: for each of the two colors, a small
rounded square + a text label, positioned on the right side of the donut.

## Quick audit answers for this file

- *"Why SVG and not canvas or a library?"* — The subject requires SVG.
  SVG is also sharp at any size, responsive via viewBox, and each shape can
  carry its own tooltip.
- *"How does the donut show the ratio?"* — Two circles: a full one behind
  (received) and a partial arc on top (given) made with `stroke-dasharray`,
  sized as given's fraction of the circumference.
- *"Are the graphs dynamic?"* — Yes: everything is calculated from the
  logged-in user's data (scales, bar count, graph height, arc size). Any
  student who logs in sees their own graphs.

---
