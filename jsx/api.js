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

function getTables(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables", opts, cb);
}

function getTableRows(connId, table, opts, cb) {
  opts.conn_id = connId;
  apiCall("get", "/tables/" + table + "/rows", opts, cb);
}

function getTableStructure(connId, table, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables/" + table, opts, cb);
}

function getTableIndexes(connId, table, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables/" + table + "/indexes", opts, cb);
}

function getHistory(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/history", opts, function(data) {
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

function getActivity(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/activity", opts, cb);
}

function executeQuery(connId, query, cb) {
  apiCall("post", "/query", {
    conn_id : connId,
    query: query
  }, cb);
}

function explainQuery(connId, query, cb) {
  apiCall("post", "/explain", {
    conn_id: connId,
    query: query
  }, cb);
}

function getConnectionInfo(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/connection", opts, function(data) {
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
