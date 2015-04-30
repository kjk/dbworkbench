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

function getTables(cb) {
  apiCall("get", "/tables", {}, cb);
}

function getTableRows(table, opts, cb) {
  apiCall("get", "/tables/" + table + "/rows", opts, cb);
}

function getTableStructure(table, cb) {
  apiCall("get", "/tables/" + table, {}, cb);
}

function getTableIndexes(table, cb) {
  apiCall("get", "/tables/" + table + "/indexes", {}, cb);
}

function getHistory(cb) {
  apiCall("get", "/history", {}, function(data) {
    var rows = [];
    for (var i in data) {
      rows.unshift([parseInt(i) + 1, data[i].query, data[i].timestamp]);
    }
    cb({ columns: ["id", "query", "timestamp"], rows: rows });
  });
}

function getBookmarks(cb) {
  apiCall("get", "/bookmarks", {}, cb);
}

function getActivity(cb) {
  apiCall("get", "/activity", {}, cb);
}

function executeQuery(query, cb) {
  apiCall("post", "/query", {
    query: query
  }, cb);
}

function explainQuery(query, cb) {
  apiCall("post", "/explain", {
    query: query
  }, cb);
}

function getConnectionInfo(cb) {
  apiCall("get", "/connection", {}, function(data) {
    var rows = [];
    for (var key in data) {
      rows.push([key, data[key]]);
    }

    cb({
      columns: ["attribute", "value"],
      rows: rows
    });
  });
}

module.exports = {
  call: apiCall,
  getTables: getTables,
  getTableRows: getTableRows,
  getTableStructure: getTableStructure,
  getTableIndexes: getTableIndexes,
  getHistory: getHistory,
  getBookmarks: getBookmarks,
  getActivity: getActivity,
  executeQuery: executeQuery,
  explainQuery: explainQuery,
  getConnectionInfo: getConnectionInfo
};
