async function showProfilePage() {
    $("login-page").classList.add("hidden");
    $("profile-page").classList.remove("hidden");
    $("loading").classList.remove("hidden");

    try {
        const [userData, xpData, progressData] = await Promise.all([
            gql(USER_QUERY),
            gql(XP_QUERY),
            gql(PROGRESS_QUERY),
        ]);

        const user = userData.user[0];
        const transactions = xpData.transaction;
        const progresses = progressData.progress;

        renderUser(user);
        renderXP(transactions);
        renderProjects(progresses);
        drawXpOverTime(transactions);
        drawXpByProject(transactions);
        drawAuditDonut(user);
        $("loading").classList.add("hidden");
    } catch (err) {
        // If the token is old or broken, go back to login
        console.error(err);
        logout();
        showLoginError("Session expired. Please sign in again.");
    }
}

function renderUser(user) {
    const attrs = user.attrs || {};
    const fullName = [attrs.firstName, attrs.lastName].filter(Boolean).join(" ");
    $("user-name").textContent = fullName || user.login;
    $("user-sub").textContent = `@${user.login} · id ${user.id}`;
    $("prompt-login").textContent = user.login;

    $("audit-ratio").textContent = (user.auditRatio || 0).toFixed(1);
    $("audit-detail").textContent =
        `Done ${formatXP(user.totalUp)} · Received ${formatXP(user.totalDown)}`;
}

function renderXP(transactions) {
    const total = transactions.reduce((sum, t) => sum + t.amount, 0);
    $("total-xp").textContent = formatXP(total);
}

function renderProjects(progresses) {
    // One project can have several attempts (retries),
    // so keep only the BEST grade for each project name.
    const best = {};
    progresses.forEach((p) => {
        const name = p.object.name;
        const current = best[name] ?? -1;
        const grade = p.grade ?? -1;
        if (!(name in best) || grade > current) best[name] = p.grade;
    });

    const grades = Object.values(best);
    const passed = grades.filter((g) => g !== null && g >= 1).length;
    const failed = grades.filter((g) => g !== null && g < 1).length;
    const inProgress = grades.filter((g) => g === null).length;

    $("projects-count").textContent = passed;
    $("projects-detail").textContent =
        `${passed} passed · ${failed} failed · ${inProgress} in progress`;
}