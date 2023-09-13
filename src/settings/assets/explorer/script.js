function open(path) {
    var req = new XMLHttpRequest();
    req.open("POST", "%s", false);
    req.send(path);
    if (path == "") {
        document.body.innerHTML = req.responseText
    } else {
        document.getElementById("filesystem-explorer-table-body").innerHTML= req.responseText
        document.getElementById("filesystem-address-line").value = path
    }
}

function back() {
    let splitted_path = document.getElementById("filesystem-address-line").value.split("/");
    let path = splitted_path.slice(0, splitted_path.length - 1).join("/");
    if (path == "") {
        path = "/"
    }
    open(path)
}