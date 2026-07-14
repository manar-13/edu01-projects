$("login-btn").addEventListener("click", login);
$("logout-btn").addEventListener("click", logout);

// Pressing Enter in the password field also signs in
$("password").addEventListener("keydown", (e) => {
    if (e.key === "Enter") login();
});

// If a token is already saved, skip the login page
if (localStorage.getItem("jwt")) {
    showProfilePage();
}