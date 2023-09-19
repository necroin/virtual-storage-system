function open(path) {
    var req = new XMLHttpRequest();
    req.open("POST", "%s", false);
    req.send(path);
    if (path == "") {
        document.body.innerHTML = req.responseText
    } else {
        document.getElementById("filesystem-explorer-table-body").innerHTML = req.responseText
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

function open_create_options() {
    document.getElementById("create-options").classList.toggle("show");
}

window.onclick = function (event) {
    if (!event.target.matches('.bar-text-button')) {
        var dropdowns = document.getElementsByClassName("dropdown-content");
        var i;
        for (i = 0; i < dropdowns.length; i++) {
            var openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
}