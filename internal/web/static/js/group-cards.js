// @ts-nocheck

document.addEventListener("htmx:afterRequest", function (e) {
  if (!e.target.classList.contains("card")) return

  const data = JSON.parse(e.detail.xhr.responseText)
  if (data.new_state) {
    e.target.classList.add("owned")
  } else {
    e.target.classList.remove("owned")
  }
})

const urlParams = new URLSearchParams(window.location.search)

for (const el of document.querySelectorAll("select")) {
  if (urlParams.has(el.name)) {
    el.value = urlParams.get(el.name)
  }
  el.addEventListener("change", (_) => {
    for (const e of document.getElementsByClassName("accordion")) {
      e.classList.remove("active")
    }
    for (const e of document.getElementsByClassName("panel")) {
      e.style = "display: none;"
    }
  })
}
