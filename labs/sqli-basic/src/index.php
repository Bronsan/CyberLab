<?php
/**
 * CyberLab — Basic SQL Injection Challenge
 * Vulnerable login page with intentional SQL injection flaw
 */

$flag = trim(file_get_contents('/flag.txt'));

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'] ?? '';
    $password = $_POST['password'] ?? '';

    // VULNERABILITY: SQL injection via direct string concatenation
    $mysqli = new mysqli('mysql', getenv('MYSQL_USER') ?: 'root', getenv('MYSQL_PASSWORD') ?: '', 'cyberlab');

    if (!$mysqli->connect_error) {
        $query = "SELECT * FROM users WHERE username = '$username' AND password = '$password'";
        $result = $mysqli->query($query);

        if ($result && $result->num_rows > 0) {
            // Login successful — show flag
            $message = "<div class='success'>🎉 Login successful! Flag: <code>$flag</code></div>";
        } else {
            $message = "<div class='error'>❌ Invalid credentials</div>";
        }
    } else {
        // Fallback for standalone mode
        if ($username === "admin' OR '1'='1" && strlen($username) > 10) {
            $message = "<div class='success'>🎉 Login successful! Flag: <code>$flag</code></div>";
        } else {
            $message = "<div class='error'>❌ Invalid credentials</div>";
        }
    }
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CyberLab — SQL Injection Challenge</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: monospace; background: #0a0a0a; color: #00ff41; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
        .container { background: #111; border: 1px solid #1a1a1a; border-radius: 12px; padding: 40px; width: 400px; }
        h1 { font-size: 1.2em; margin-bottom: 20px; color: #fff; }
        input { width: 100%; padding: 10px; margin-bottom: 12px; background: #0a0a0a; border: 1px solid #1a1a1a; border-radius: 6px; color: #00ff41; font-family: monospace; }
        button { width: 100%; padding: 10px; background: #00ff41; color: #000; border: none; border-radius: 6px; cursor: pointer; font-weight: bold; }
        .success { background: #00ff4122; color: #00ff41; padding: 10px; border-radius: 6px; margin-top: 10px; }
        .error { background: #ff000022; color: #ff3333; padding: 10px; border-radius: 6px; margin-top: 10px; }
        code { background: #1a1a1a; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔐 Admin Login</h1>
        <form method="POST">
            <input type="text" name="username" placeholder="Username" />
            <input type="password" name="password" placeholder="Password" />
            <button type="submit">Login</button>
        </form>
        <?php if (isset($message)) echo $message; ?>
    </div>
</body>
</html>
