window.request_url = "%s"
window.storage_url = null

function remove_dropdown(event, event_id, dropdown_id) {
    if (!event.target.matches('#'+event_id)){
        document.getElementById(dropdown_id).classList.remove('show');
    }
}

window.onclick = function (event) {
    remove_dropdown(event, "bar-create-button", "create-options")
    remove_dropdown(event, "bar-options-button", "options")
}

function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, "https://" + url, false);
    req.send(data);
    return req.responseText
}

function get_request_url() {
    var url = window.request_url
    if (window.storage_url != null) {
        url = window.storage_url
    }
    return url
}

function open(path) {
    let response = request("POST", get_request_url(), path);
    if (path == "") {
        document.body.innerHTML = response
    } else {
        document.getElementById("filesystem-explorer-table-body").innerHTML = response
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

function open_create_dialog() {
    document.getElementById("create-dialog").style.display = "flex";
    document.getElementById("create-dialog-overlay").style.display = "block";
}

function close_create_dialog() {
    document.getElementById("create-dialog").style.display = "none";
    document.getElementById("create-dialog-overlay").style.display = "none";
}

function open_options() {
    document.getElementById("options").classList.toggle("show");
}

function update_status_bar(raw_data) {
    let data = JSON.parse(raw_data)
    document.getElementById("status-bar-progress").innerHTML = data.status
    document.getElementById("status-bar-text").innerHTML = data.text
}

function create(type) {
    let data = JSON.stringify(
        {
            "type": type,
            "path": document.getElementById("filesystem-address-line").value,
            "name": document.getElementById("create-dialog-name").value
        }
    );
    let response = request("POST", get_request_url() + "/insert/" + type, data);
    open(document.getElementById("filesystem-address-line").value)
    close_create_dialog()
    update_status_bar(response)
}

function remove() {
    let focus_item = document.getElementById("filesystem-explorer-table").focus_item
    if (focus_item != null) {
        let focus_item_name = focus_item.attributes.name.value
        let path = [document.getElementById("filesystem-address-line").value, focus_item_name].join("/")
        if (document.getElementById("filesystem-address-line").value == "/") {
            path = "/" + focus_item_name
        }
        let response = request("POST", get_request_url() + "/delete", path);
        open(document.getElementById("filesystem-address-line").value)
        update_status_bar(response)
    }
}

function copy() {
    let focus_item = document.getElementById("filesystem-explorer-table").focus_item
    if (focus_item != null) {
        let focus_item_name = focus_item.attributes.name.value
        let focus_item_type = focus_item.attributes.custom_type.value
        let path = [document.getElementById("filesystem-address-line").value, focus_item_name].join("/")
        if (document.getElementById("filesystem-address-line").value == "/") {
            path = "/" + focus_item_name
        }
        let response = request("POST", get_request_url() + "/copy/" + focus_item_type, path);
        update_status_bar(response)
    }
}

function paste() {
    let response = request("POST", get_request_url() + "/paste", document.getElementById("filesystem-address-line").value);
    open(document.getElementById("filesystem-address-line").value)
    update_status_bar(response)
}

