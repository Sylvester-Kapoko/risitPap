package server

const loginTemplate = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>Login</title>
<style>body{font-family:sans-serif;max-width:400px;margin:100px auto;padding:20px;}
input{display:block;width:100%;padding:8px;margin:10px 0;}
button{padding:10px 20px;}</style></head><body>
<h1>RisitPap Login</h1>
<form method="POST">
  <input name="username" placeholder="Username" required>
  <input name="password" type="password" placeholder="Password" required>
  <button type="submit">Log in</button>
</form>
{{.Error}}</body></html>`