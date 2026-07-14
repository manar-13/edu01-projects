# Notes — js/main.js

The smallest file, but it's the **starting point**: it connects the buttons
to the functions and decides which page to show first. It loads LAST in
`index.html`, because it uses functions defined in all the other files.

## Connecting the buttons

```javascript
$("login-btn").addEventListener("click", login);
$("logout-btn").addEventListener("click", logout);
```

`addEventListener("click", login)` means: "when this element is clicked,
run the `login` function".

- Note we write `login` — **without parentheses**. We are handing the
  function itself to the browser to call LATER, when the click happens
- Writing `login()` (with parentheses) would be a bug: it would run the
  function immediately, once, while the page loads

> Classic audit / interview question: *"What's the difference between
> `login` and `login()` here?"* — `login` passes the function as a value;
> `login()` calls it right now and passes its RESULT.

## Enter key = sign in

```javascript
// Pressing Enter in the password field also signs in
$("password").addEventListener("keydown", (e) => {
    if (e.key === "Enter") login();
});
```

A small quality-of-life feature: listen to every key press (`keydown`)
inside the password box.

- The browser gives us an event object `e` with details about what happened
- `e.key` is the name of the pressed key: `"a"`, `"Shift"`, `"Enter"`...
- If it's `"Enter"` → call `login()` — here WITH parentheses, because now
  we do want to run it (we are already inside the "later" function)

> This is a UI good practice: users expect Enter to submit a login form.
> Small details like this help the bonus item "Does the UI respect good
> practices?"

## Deciding the first page

```javascript
// If a token is already saved, skip the login page
if (localStorage.getItem("jwt")) {
    showProfilePage();
}
```

This runs ONCE, when the page loads:

- `localStorage.getItem("jwt")` returns the saved token, or `null` if
  there isn't one
- In an `if`, a non-empty string counts as true and `null` counts as false
- So: token exists → the user logged in before → go straight to the
  profile. No token → do nothing, and the login page (visible by default
  in the HTML) stays on screen

> This is why refreshing the page keeps you logged in: the JWT survives
> in localStorage, and this line finds it.

> And if the saved token is EXPIRED? `showProfilePage()` will try to fetch
> data, the API will answer with an error, and the `catch` inside
> `showProfilePage` calls `logout()` and returns the user to the login
> page with a "Session expired" message. So even the bad case is handled.

## The whole app flow in one picture

```
page loads
   │
   ├── token in localStorage? ──── yes ──→ showProfilePage()
   │                                             │
   no                                            ▼
   │                                    queries + cards + graphs
   ▼
login page waits
   │
   ├── user clicks "Sign in" ──→ login()
   ├── user presses Enter ─────→ login()
   │                               │ success
   │                               ▼
   │                        showProfilePage()
   │
   └── (on profile) clicks "log out" ──→ logout() ──→ back to login page
```

---
