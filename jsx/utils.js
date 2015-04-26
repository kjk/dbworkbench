function runOnLoad(f) {
  if (window.addEventListener) {
    window.addEventListener('DOMContentLoaded', f);
  } else {
    window.attachEvent('onload', f);
  }
}

function apiCall(method, path, params, cb) {
  $.ajax({
    url: "/api" + path,
    method: method,
    cache: false,
    data: params,
    success: function(data) {
      cb(data);
    },
    error: function(xhr, status, data) {
      cb(jQuery.parseJSON(xhr.responseText));
    }
  });
}

function getTables(cb)                 { apiCall("get", "/tables", {}, cb); }
function getTableRows(table, opts, cb) { apiCall("get", "/tables/" + table + "/rows", opts, cb); }
function getTableStructure(table, cb)  { apiCall("get", "/tables/" + table, {}, cb); }
function getTableIndexes(table, cb)    { apiCall("get", "/tables/" + table + "/indexes", {}, cb); }
function getHistory(cb)                { apiCall("get", "/history", {}, cb); }
function getBookmarks(cb)              { apiCall("get", "/bookmarks", {}, cb); }

function executeQuery(query, cb) {
  apiCall("post", "/query", { query: query }, cb);
}

function explainQuery(query, cb) {
  apiCall("post", "/explain", { query: query }, cb);
}


window.runOnLoad = runOnLoad;
window.apiCall = apiCall;
window.getTables = getTables;
window.getTableRows = getTableRows;
window.getTableStructure = getTableStructure;
window.getTableIndexes = getTableIndexes;
window.getHistory = getHistory;
window.getBookmarks = getBookmarks;
