<h3>hostmark</h3>
<h4>Create your account</h4>

<form>
    <label for="username">Username</label>
    <input id="username" name="username" type="text" placeholder="Username" autocomplete="username" required />

    <label for="password">Password</label>
    <input id="password" name="password" type="password" placeholder="Password (at least 15 characters)" minlength="15" maxlength="64" required />

    <label for="confirm">Confirm Password</label>
    <input id="confirm" name="confirm" type="password" placeholder="Password" required />

    <button type="submit">Sign up</button>
    <small>Already have an account? <a href="/login">Sign in</a></small>
</form>

<style>
    h3 {
        text-align: center;
    }
</style>
