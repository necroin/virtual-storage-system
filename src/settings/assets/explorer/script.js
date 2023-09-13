function open(path) {
    var req = new XMLHttpRequest();
    req.open("POST", "%s", false);
    req.send(path);
    document.body.innerHTML = req.responseText
}