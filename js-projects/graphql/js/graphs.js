// Small helper to create SVG elements with attributes
function svgEl(tag, attrs) {
    const el = document.createElementNS("http://www.w3.org/2000/svg", tag);
    for (const key in attrs) el.setAttribute(key, attrs[key]);
    return el;
}

// Graph 1: XP over time (cumulative line + area)
function drawXpOverTime(transactions) {
    const holder = $("graph-xp-time");
    holder.innerHTML = "";
    if (transactions.length === 0) {
        holder.textContent = "No XP data yet.";
        return;
    }

    const W = 640, H = 300;
    const pad = { top: 20, right: 20, bottom: 40, left: 60 };

    // Build cumulative points: [time, total xp so far]
    let sum = 0;
    const points = transactions.map((t) => {
        sum += t.amount;
        return { time: new Date(t.createdAt).getTime(), xp: sum };
    });

    const minT = points[0].time;
    const maxT = points[points.length - 1].time;
    const maxXP = sum;

    // Scale data values into pixel positions
    const x = (t) => pad.left + ((t - minT) / (maxT - minT || 1)) * (W - pad.left - pad.right);
    const y = (v) => H - pad.bottom - (v / maxXP) * (H - pad.top - pad.bottom);

    const svg = svgEl("svg", { viewBox: `0 0 ${W} ${H}`, class: "chart" });

    // Horizontal grid lines + labels (0%, 50%, 100% of max XP)
    [0, 0.5, 1].forEach((f) => {
        const gy = y(maxXP * f);
        svg.appendChild(svgEl("line", { x1: pad.left, y1: gy, x2: W - pad.right, y2: gy, class: "grid" }));
        const label = svgEl("text", { x: pad.left - 8, y: gy + 4, class: "axis-label", "text-anchor": "end" });
        label.textContent = formatXP(maxXP * f);
        svg.appendChild(label);
    });

    // Date labels (start and end)
    const startLabel = svgEl("text", { x: pad.left, y: H - 12, class: "axis-label" });
    startLabel.textContent = new Date(minT).toLocaleDateString();
    const endLabel = svgEl("text", { x: W - pad.right, y: H - 12, class: "axis-label", "text-anchor": "end" });
    endLabel.textContent = new Date(maxT).toLocaleDateString();
    svg.appendChild(startLabel);
    svg.appendChild(endLabel);

    // The line path: "M x,y L x,y L x,y ..."
    const linePath = points.map((p, i) => `${i === 0 ? "M" : "L"}${x(p.time)},${y(p.xp)}`).join(" ");

    // The area under the line (same path, closed down to the bottom)
    const areaPath = linePath +
        ` L${x(maxT)},${H - pad.bottom} L${x(minT)},${H - pad.bottom} Z`;

    svg.appendChild(svgEl("path", { d: areaPath, class: "area" }));
    svg.appendChild(svgEl("path", { d: linePath, class: "line" }));

    // A small circle on each point, with a native tooltip (<title>)
    points.forEach((p, i) => {
        const dot = svgEl("circle", { cx: x(p.time), cy: y(p.xp), r: 3.5, class: "dot" });
        const tip = svgEl("title", {});
        tip.textContent =
            `${transactions[i].object?.name || "XP"} · +${formatXP(transactions[i].amount)}\n` +
            `Total: ${formatXP(p.xp)} · ${new Date(p.time).toLocaleDateString()}`;
        dot.appendChild(tip);
        svg.appendChild(dot);
    });

    holder.appendChild(svg);
}

// Graph 2: XP by project (horizontal bars, top 10)
function drawXpByProject(transactions) {
    const holder = $("graph-xp-project");
    holder.innerHTML = "";
    if (transactions.length === 0) {
        holder.textContent = "No XP data yet.";
        return;
    }

    // Group XP by project name
    const byProject = {};
    transactions.forEach((t) => {
        const name = t.object?.name || t.path.split("/").pop();
        byProject[name] = (byProject[name] || 0) + t.amount;
    });

    // Sort biggest first, keep top 10
    const top = Object.entries(byProject)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 10);

    const W = 640;
    const rowH = 34;
    const pad = { top: 10, right: 70, bottom: 10, left: 150 };
    const H = pad.top + pad.bottom + top.length * rowH;
    const maxXP = top[0][1];

    const svg = svgEl("svg", { viewBox: `0 0 ${W} ${H}`, class: "chart" });

    top.forEach(([name, amount], i) => {
        const yPos = pad.top + i * rowH;
        const barW = (amount / maxXP) * (W - pad.left - pad.right);

        // Project name on the left
        const label = svgEl("text", {
            x: pad.left - 10, y: yPos + rowH / 2 + 4,
            class: "axis-label", "text-anchor": "end",
        });
        label.textContent = name.length > 18 ? name.slice(0, 17) + "…" : name;
        svg.appendChild(label);

        // The bar itself, with a native tooltip
        const bar = svgEl("rect", {
            x: pad.left, y: yPos + 6,
            width: barW, height: rowH - 12,
            rx: 5, class: "bar",
        });
        const tip = svgEl("title", {});
        tip.textContent = `${name}: ${formatXP(amount)}`;
        bar.appendChild(tip);
        svg.appendChild(bar);

        // XP value at the end of the bar
        const value = svgEl("text", {
            x: pad.left + barW + 8, y: yPos + rowH / 2 + 4,
            class: "bar-value",
        });
        value.textContent = formatXP(amount);
        svg.appendChild(value);
    });

    holder.appendChild(svg);
}

// Graph 3: Audit ratio donut (XP given in audits vs received)
function drawAuditDonut(user) {
    const holder = $("graph-audit");
    holder.innerHTML = "";

    const given = user.totalUp || 0;
    const received = user.totalDown || 0;
    const total = given + received;
    if (total === 0) {
        holder.textContent = "No audit data yet.";
        return;
    }

    const W = 640, H = 260;
    const cx = 160, cy = H / 2;
    const r = 90;
    const stroke = 30;
    const circumference = 2 * Math.PI * r;

    // How much of the ring belongs to "given"
    const givenPart = (given / total) * circumference;

    const svg = svgEl("svg", { viewBox: `0 0 ${W} ${H}`, class: "chart" });

    // Full ring = "received" color (background circle)
    const ringReceived = svgEl("circle", {
        cx, cy, r, fill: "none",
        stroke: "#33415E", "stroke-width": stroke,
    });
    const tipR = svgEl("title", {});
    tipR.textContent = `Received: ${formatXP(received)}`;
    ringReceived.appendChild(tipR);
    svg.appendChild(ringReceived);

    // On top: the "given" part, drawn with stroke-dasharray
    // (dasharray = "visible length, invisible length")
    const ringGiven = svgEl("circle", {
        cx, cy, r, fill: "none",
        stroke: "#56E1CE", "stroke-width": stroke,
        "stroke-dasharray": `${givenPart} ${circumference - givenPart}`,
        transform: `rotate(-90 ${cx} ${cy})`,
        "stroke-linecap": "round",
    });
    const tipG = svgEl("title", {});
    tipG.textContent = `Given: ${formatXP(given)}`;
    ringGiven.appendChild(tipG);
    svg.appendChild(ringGiven);

    // The ratio number in the middle of the donut
    const ratioText = svgEl("text", {
        x: cx, y: cy + 8, "text-anchor": "middle", class: "donut-center",
    });
    ratioText.textContent = (user.auditRatio || 0).toFixed(1);
    svg.appendChild(ratioText);

    // Legend on the right side
    const legend = [
        { color: "#56E1CE", label: `Given (audits you did): ${formatXP(given)}` },
        { color: "#33415E", label: `Received (audits on you): ${formatXP(received)}` },
    ];
    legend.forEach((item, i) => {
        const ly = cy - 15 + i * 30;
        svg.appendChild(svgEl("rect", { x: 320, y: ly - 10, width: 14, height: 14, rx: 4, fill: item.color }));
        const text = svgEl("text", { x: 344, y: ly + 2, class: "legend-label" });
        text.textContent = item.label;
        svg.appendChild(text);
    });

    holder.appendChild(svg);
}