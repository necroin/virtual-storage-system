window.request_url = "%s"
window.storage_url = null

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
    url = window.request_url
    if (window.storage_url != null) {
        url = window.storage_url
    }
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

function set_storage(url) {
    window.storage_url = url
    open(document.getElementById("filesystem-address-line").value)
}

function back() {
    let splitted_path = document.getElementById("filesystem-address-line").value.split("/");
    let path = splitted_path.slice(0, splitted_path.length - 1).join("/");
    if (path == "") {
        path = "/"
    }
    open(path)
}

function set_focus_item(table, item) {
    table.focus_item = item
}

function open_create_options() {
    document.getElementById("create-options").classList.toggle("show");
}

function create(type) {
    var req = new XMLHttpRequest();
    req.open("POST", window.request_url + "/insert/" + type, false);
    req.send(document.getElementById("filesystem-address-line").value)
    open(document.getElementById("filesystem-address-line").value)
}

function remove() {
    var req = new XMLHttpRequest();
    req.open("POST", window.request_url + "/delete", false);
    let focus_item = document.getElementById("filesystem-explorer-table").focus_item
    if (focus_item != null) {
        let focus_item_name = focus_item.attributes.name.value
        let path = [document.getElementById("filesystem-address-line").value, focus_item_name].join("/")
        if (document.getElementById("filesystem-address-line").value == "/") {
            path = "/" + focus_item
        }
        req.send(path)
        open(document.getElementById("filesystem-address-line").value)
    }
}

function copy() {
    var req = new XMLHttpRequest();
    let focus_item = document.getElementById("filesystem-explorer-table").focus_item
    if (focus_item != null) {
        let focus_item_name = focus_item.attributes.name.value
        let focus_item_type = focus_item.attributes.custom_type.value
        let path = [document.getElementById("filesystem-address-line").value, focus_item_name].join("/")
        if (document.getElementById("filesystem-address-line").value == "/") {
            path = "/" + focus_item_name
        }

        req.open("POST", window.request_url + "/copy/" + focus_item_type, false);
        req.send(path)
        document.getElementById("status-bar-text").innerHTML = req.responseText
    }
}

function paste() {
    var req = new XMLHttpRequest();
    req.open("POST", window.request_url + "/paste", false);
    req.send(document.getElementById("filesystem-address-line").value)
    open(document.getElementById("filesystem-address-line").value)
    document.getElementById("status-bar-text").innerHTML = req.responseText
}

