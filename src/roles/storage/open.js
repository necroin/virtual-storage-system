function open(path) {
    var req = new XMLHttpRequest();
    req.open("GET", "%s", false);
    req.send(path);
    document.body.innerHTML = req.responseText 
}