function runOnLoad(f) {
  if (window.addEventListener) {
    window.addEventListener('DOMContentLoaded', f);
  } else {
    window.attachEvent('onload', f);
  }
}

window.runOnLoad = runOnLoad;
