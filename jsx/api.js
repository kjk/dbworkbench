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
function getActivity(cb)               { apiCall("get", "/activity", {}, cb); }

function executeQuery(query, cb) {
  apiCall("post", "/query", { query: query }, cb);
}

function explainQuery(query, cb) {
  apiCall("post", "/explain", { query: query }, cb);
}

module.exports.call = apiCall;
module.exports.getTables = getTables;
module.exports.getTableRows = getTableRows;
module.exports.getTableStructure = getTableStructure;
module.exports.getTableIndexes = getTableIndexes;
module.exports.getHistory = getHistory;
module.exports.getBookmarks = getBookmarks;
