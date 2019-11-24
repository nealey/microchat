// jshint asi:true

function μchatInit(initEvent) {
  let updateInterval
  let lastMsg
  
  function toast(message, timeout=5000) {
    let p = document.createElement("p")
  
    p.innerText = message
    document.querySelector("#chattoast").appendChild(p)
    setTimeout(e => { p.remove() }, timeout)
  }

  function input(event) {
    event.preventDefault()
    
    let form = event.target
    let inp = form.elements.text
    let body = new FormData(form)
    fetch("say", {
      method: "POST",
      body: body,
    })
    .then(resp => {
      if (resp.ok) {
        // Yay it was okay, reset input
        inp.value = ""
      } else {
        toast("ERROR: DOES NOT COMPUTE")
        console.log(resp)
      }
    })
    .catch(err => {
      toast("ERROR: DOES NOT COMPUTE")
      console.log(err)
    })
    
    window.localStorage["who"] = form.elements.who.value
  }
  
  function updateLog(log) {
    let lastLogMsg = log[log.length - 1]
    if (
      lastMsg &&
      (lastLogMsg.When == lastMsg.When) &&
      (lastLogMsg.Who == lastMsg.Who) &&
      (lastLogMsg.Text == lastMsg.Text)
    ) {
      return
    }
    lastMsg = lastLogMsg
    
    let chatlog = document.querySelector("#chatlog")
    while (chatlog.firstChild) {
      chatlog.firstChild.remove()
    }
    for (let ll of log) {
      let line = chatlog.appendChild(document.createElement("div"))
      line.classList.add("line")
      
      let when = line.appendChild(document.createElement("span"))
      when.classList.add("when")
      when.innerText = (new Date(ll.When * 1000)).toISOString()
      
      let who = line.appendChild(document.createElement("span"))
      who.classList.add("who")
      who.innerText = ll.Who
      
      let text = line.appendChild(document.createElement("span"))
      text.classList.add("text")
      text.innerText = ll.Text
    }
    chatlog.lastChild.scrollIntoView()
    
    // trigger a fetch of chat log so the user gets some feedback
    update()
  }
  
  function update(event) {
    fetch("read")
    .then(resp => {
      if (resp.ok) {
        resp.json()
        .then(updateLog)
      }
    })
    .catch(err => {
      toast("Server error: " + err)
    })
  }

  for (let f of document.forms) {
    let who = window.localStorage["who"]
    f.addEventListener("submit", input)
    if (who) {
      f.elements.who.value = who
    }
  }
  updateInterval = setInterval(update, 3000)
  update()
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", μchatInit)
} else {
  μchatInit()
}
