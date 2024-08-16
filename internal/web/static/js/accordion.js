// @ts-nocheck

function toggle(el) {
  if (!el.classList.contains("accordion")) {
    el = el.parentElement
  }
  el.classList.toggle("active")
  const panel = el.nextElementSibling
  if (panel.style.display === "flex") {
    panel.style.display = "none"
    return false
  } else {
    panel.style.display = "flex"
    return true
  }
}
