function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function auth() {
    let username = document.getElementById("auth-name").value
    let password = document.getElementById("auth-password").value
    let data = JSON.stringify(
        {
            "username": username,
            "password": password,
        }
    );
    let response = request("POST", window.location.href+"/token", data)
    open(response)
}

function open(token) {
    if (token != "") {
        window.location.pathname = "/"+token+"/router/explorer"
    }
}