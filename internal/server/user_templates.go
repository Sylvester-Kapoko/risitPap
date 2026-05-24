package server

const changePasswordTemplate = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Change Password</title>
<style>body{font-family:sans-serif;max-width:400px;margin:100px auto;padding:20px;}
input{display:block;width:100%;padding:8px;margin:10px 0;}
button{padding:10px 20px;}</style></head><body>
<h1>Change Password</h1>
<form method="POST">
    <input type="password" name="old_password" placeholder="Old Password" required>
    <input type="password" name="new_password" placeholder="New Password" required>
    <button type="submit">Change</button>
</form>
{{.Error}}</body></html>` // #nosec G101

const userListTemplate = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>User Management</title>
<style>body{font-family:sans-serif;max-width:600px;margin:auto;padding:20px;}
table{width:100%;border-collapse:collapse;}td,th{padding:8px;border-bottom:1px solid #ddd;}
</style></head><body>
<h1>User Management</h1>
<table>
<tr><th>Username</th><th>Action</th></tr>
{{range .Users}}<tr><td>{{.Username}}</td>
<td>{{if $.IsSupervisor}}<a href="/admin/users/delete?username={{.Username}}" onclick="return confirm('Delete user?')">Delete</a>{{else}}—{{end}}</td></tr>{{end}}
</table>
{{if .IsSupervisor}}
<h2>Add User</h2>
<form method="POST" action="/admin/users/add">
    <input name="username" placeholder="Username" required>
    <input name="password" type="password" placeholder="Password" required>
    <button type="submit">Add</button>
</form>
{{end}}
</body></html>`
