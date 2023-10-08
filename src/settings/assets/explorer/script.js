window.request_url = "%s"

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

function open(path) {
    var req = new XMLHttpRequest();
    req.open("POST", window.request_url, false);
    req.send(path);
    if (path == "") {
        document.body.innerHTML = req.responseText
    } else {
        document.getElementById("filesystem-explorer-table-body").innerHTML = req.responseText
        document.getElementById("filesystem-address-line").value = path
    }
    document.getElementById("filesystem-explorer-table").focus_item = null
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

function create(type) {
    var req = new XMLHttpRequest();
    req.open("POST", window.request_url + "/insert/"+type, false);
    req.send(document.getElementById("filesystem-address-line").value)
    open(document.getElementById("filesystem-address-line").value)
}

function remove() {
    var req = new XMLHttpRequest();
    req.open("POST", window.request_url + "/delete", false);
    let focus_item = document.getElementById("filesystem-explorer-table").focus_item
    if (focus_item != null) {
        focus_item = focus_item.attributes.name.value
        let path = [document.getElementById("filesystem-address-line").value, focus_item].join("/")
        if (document.getElementById("filesystem-address-line").value == "/") {
            path = "/" + focus_item
        }
        req.send(path)
        open(document.getElementById("filesystem-address-line").value)
    }
}

function copy(type) {

}

function paste(type) {

}