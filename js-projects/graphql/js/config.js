const DOMAIN = "learn.reboot01.com";
const SIGNIN_URL = `https://${DOMAIN}/api/auth/signin`;
const GRAPHQL_URL = `https://${DOMAIN}/api/graphql-engine/v1/graphql`;

// Shortcut: $("id") instead of document.getElementById("id")
const $ = (id) => document.getElementById(id);

// Turn 154000 into "154 kB", like the platform shows XP
function formatXP(amount) {
    if (amount >= 1_000_000) return (amount / 1_000_000).toFixed(2) + " MB";
    if (amount >= 1_000) return Math.round(amount / 1_000) + " kB";
    return amount + " B";
}