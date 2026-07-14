async function login() {
    const identifier = $("identifier").value.trim();
    const password = $("password").value;
    $("login-error").classList.add("hidden");

    if (!identifier || !password) {
        showLoginError("Please fill in both fields.");
        return;
    }

    try {
        // Basic auth = "username:password" encoded in base64
        const res = await fetch(SIGNIN_URL, {
            method: "POST",
            headers: { Authorization: "Basic " + btoa(`${identifier}:${password}`) },
        });

        if (!res.ok) {
            showLoginError("Invalid username/email or password. Please try again.");
            return;
        }

        const token = await res.json();
        localStorage.setItem("jwt", token);

        $("identifier").value = "";
        $("password").value = "";
        showProfilePage();
    } catch (err) {
        showLoginError("Network error. Please check your connection.");
    }
}

function showLoginError(message) {
    const errorBox = $("login-error");
    errorBox.textContent = message;
    errorBox.classList.remove("hidden");
}

function logout() {
    localStorage.removeItem("jwt");
    $("profile-page").classList.add("hidden");
    $("login-page").classList.remove("hidden");
}