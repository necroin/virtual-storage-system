function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function GetTopology(url) {
    let response = request("GET", "https://"+url+"/router/topology")
    let data = JSON.parse(response)

    let storagesList = document.getElementById("storages")
    for (index in data.storages) {
        let storage = data.storages[index]

        let item = document.createElement("div")
        item.className = "list-item"
        let nameItemElement = document.createElement("span")
        nameItemElement.innerText = storage.hostname

        let urlItemElement = document.createElement("span")
        urlItemElement.innerText = storage.url

        item.appendChild(nameItemElement)
        item.appendChild(urlItemElement)

        storagesList.appendChild(item)
    }

    let runnersList = document.getElementById("runners")
    for (index in data.runners) {
        let runner = data.runners[index]

        let item = document.createElement("div")
        item.className = "list-item"
        let nameItemElement = document.createElement("span")
        nameItemElement.innerText = runner.hostname

        let urlItemElement = document.createElement("span")
        urlItemElement.innerText = runner.url

        item.appendChild(nameItemElement)
        item.appendChild(urlItemElement)

        runnersList.appendChild(item)
    }
}