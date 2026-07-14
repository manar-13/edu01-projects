// Sends any query to the GraphQL endpoint using the saved JWT
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

    const data = await res.json();
    if (data.errors) throw new Error(data.errors[0].message);
    return data.data;
}

// Query 1 — NORMAL query: basic user info
const USER_QUERY = `{
  user {
    id
    login
    attrs
    auditRatio
    totalUp
    totalDown
  }
}`;

// Query 2 — query WITH ARGUMENTS + NESTED: all XP transactions
// from the Bahrain module and the JS piscine (matches the platform total)
const XP_QUERY = `{
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
}`;

// Query 3 — NESTED query with arguments: project results (pass/fail)
const PROGRESS_QUERY = `{
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
}`;