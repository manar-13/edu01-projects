# Notes — js/auth.js

This file handles **login** and **logout** (the JWT key card).

## login()

```javascript
async function login() {
    const identifier = $("identifier").value.trim();
    const password = $("password").value;
    $("login-error").classList.add("hidden");
```

Read what the user typed in the two input boxes. `.trim()` removes accidental
spaces around the username (but NOT from the password — spaces can be part of
a real password). Then hide any old error message before we start fresh.

> `async` before the function means: this function will **wait** for slow
> things (like talking to a server) using `await`.

```javascript
    if (!identifier || !password) {
        showLoginError("Please fill in both fields.");
        return;
    }
```

If either box is empty, show an error and `return` = stop here, don't contact
the server at all.

```javascript
    try {
        // Basic auth = "username:password" encoded in base64
        const res = await fetch(SIGNIN_URL, {
            method: "POST",
            headers: { Authorization: "Basic " + btoa(`${identifier}:${password}`) },
        });
```

The heart of the login:

- `fetch` sends an HTTP request to the signin endpoint, and `await` pauses
  until the server answers
- `method: "POST"` — we are sending data, not just reading a page
- `btoa(...)` converts the text `"myusername:mypassword"` into **base64**
  (example: `bWFuYXI6MTIzNA==`)
- The header `Authorization: Basic <base64>` is the standard format for
  **Basic authentication** — exactly what the project subject requires

> Base64 is **encoding, not encryption** — anyone can decode it back. It is
> safe here only because the connection is HTTPS (encrypted).
> Good audit answer if asked "is base64 secure?"

```javascript
        if (!res.ok) {
            showLoginError("Invalid username/email or password. Please try again.");
            return;
        }
```

`res.ok` is true when the server answers with a success code (200).
If credentials are wrong, the server answers with an error code (like 401),
so `res.ok` is false → show the error message.
This is the audit item *"Is an appropriate error shown?"*

```javascript
        const token = await res.json();
        localStorage.setItem("jwt", token);
```

If login succeeded, the server's answer body **is the JWT**. We read it with
`res.json()`, then save it in `localStorage` — a small storage box inside the
browser that survives page refreshes. That's why you stay logged in after
refreshing.

```javascript
        $("identifier").value = "";
        $("password").value = "";
        showProfilePage();
```

Clear the two input boxes (so the password doesn't stay visible on the login
form), then switch to the profile page.

```javascript
    } catch (err) {
        showLoginError("Network error. Please check your connection.");
    }
}
```

`try...catch`: if anything inside `try` **crashes** (for example: no internet,
server unreachable), the code jumps here instead of breaking the page.

> Note the difference: wrong password = server answered "no" (`!res.ok`).
> Network error = server never answered at all (`catch`).
> Two different failures, two different messages.

## showLoginError(message)

```javascript
function showLoginError(message) {
    const errorBox = $("login-error");
    errorBox.textContent = message;
    errorBox.classList.remove("hidden");
}
```

A small helper: put the message text inside the error element, then remove
the `hidden` class so it becomes visible. We use it from 4 different places,
so making it a function avoids repeating code.

## logout()

```javascript
function logout() {
    localStorage.removeItem("jwt");
    $("profile-page").classList.add("hidden");
    $("login-page").classList.remove("hidden");
}
```

Logout is beautifully simple: **delete the JWT** from localStorage (the key
card is thrown away — no more data access), hide the profile, show the login
page again.

> Audit question: *"How does logout work?"* → "I remove the JWT from
> localStorage. Without the token, no GraphQL request can be made on my
> behalf."
---