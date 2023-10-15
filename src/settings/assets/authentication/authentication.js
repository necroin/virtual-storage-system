window.request_url = '%s'

function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, "https://" + url, false);
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
    let response = request("POST", window.request_url+"/auth/token", data)
    open(response)
}

function open(token) {
    if (token != "") {
        window.location =  "https://"+window.request_url+"/"+token+"/router"
    }
}