const isTouch = () => 'ontouchstart' in window || window.DocumentTouch && document instanceof window.DocumentTouch || navigator.maxTouchPoints > 0 || window.navigator.msMaxTouchPoints > 0

function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function Init(url, pid) {
    let canvas = document.getElementById("canvas")
    canvas.__custom__ = {}

    canvas.onmousedown = (event) => {
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "leftDown", coords: { x: event.offsetX, y: event.offsetY } }))
    }
    canvas.onmouseup = (event) => {
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "leftUp", coords: { x: event.offsetX, y: event.offsetY } }))
    }
    canvas.onmousemove = (event) => {
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "move", coords: { x: event.offsetX, y: event.offsetY } }))
    }
    canvas.onwheel = (event) => {
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "scroll", coords: { x: event.offsetX, y: event.offsetY }, scroll_delta: { x: -event.deltaX, y: -event.deltaY } }))
    }

    canvas.onclick = (event) => { }
    canvas.ondblclick = (event) => { }

    canvas.ontouchstart = (event) => {
        canvas.contentEditable = true
        canvas.__custom__.touchContext = null
        if (event.touches.length == 1) {
            let offsetX = event.touches[0].pageX
            let offsetY = event.touches[0].pageY
            canvas.__custom__.touchContext = { x: offsetX, y: offsetY }
        }
        if (event.touches.length == 2) {
            canvas.contentEditable = false
        }
    }
    canvas.ontouchmove = (event) => {
        if (event.touches.length == 1 && canvas.__custom__.touchContext != null) {
            event.preventDefault()

            let offsetX = event.touches[0].pageX
            let offsetY = event.touches[0].pageY

            let deltaX = (canvas.__custom__.touchContext.x - offsetX) / 2
            let deltaY = (canvas.__custom__.touchContext.y - offsetY) / 2

            request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "scroll", coords: { x: offsetX, y: offsetY }, scroll_delta: { x: deltaX, y: deltaY } }))
        }
    }
    canvas.ontouchend = (event) => { }

    canvas.onkeydown = (event) => {
        event.preventDefault()
        let keycode = event.keyCode
        request("POST", "https://" + url + "/runner/keyboard/" + pid, JSON.stringify({ type: "keydown", keycode: keycode }))
    }
    canvas.onkeyup = (event) => {
        event.preventDefault()
        let keycode = event.keyCode
        request("POST", "https://" + url + "/runner/keyboard/" + pid, JSON.stringify({ type: "keyup", keycode: keycode }))
    }

    setTimeout(window.LaunchStream, 0, url, pid)
}

function LaunchStream(url, pid) {
    let canvas = document.getElementById("canvas")
    let ctx = canvas.getContext('2d')
    newImage = new Image();
    newImage.onload = function () {
        canvas.width = newImage.width
        canvas.height = newImage.height
        // let scale = window.screen.width / newImage.width
        // canvas.__custom__.scale = scale
        // ctx.scale(scale, scale)
        ctx.drawImage(newImage, 0, 0, newImage.width, newImage.height);
        setTimeout(window.LaunchStream, 0, url, pid)
    }
    newImage.src = "https://" + url + "/runner/image/" + pid + "?time=" + new Date().getTime();
}

