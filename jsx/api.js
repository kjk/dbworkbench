var action = require('./action.js');

function apiCall2(showSpinner, method, path, params, cb) {
  if (showSpinner) {
    action.spinner(true);    
  }

  $.ajax({
    url: "/api" + path,
    method: method,
    cache: false,
    data: params,
    async: true,
    success: function(data) {
      if (showSpinner) {
        action.spinner(false);      
      }
      if (cb) {
        cb(data);
      }
    },
    error: function(xhr, status, data) {
      if (showSpinner) {
        action.spinner(false);      
      }
      if (xhr.status == "0") {
        // Backend is down
        action.alertBar("Something is wrong. Please restart");
      } else {
        // API call failed
      }
      if (cb) {
        cb(jQuery.parseJSON(xhr.responseText));
      }
    }
  });
}

function apiCall(method, path, params, cb) {
  apiCall2(true, method, path, params, cb);
}

function apiCallNoSpinner(method, path, params, cb) {
  apiCall2(false, method, path, params, cb);
}

function connect(type, url, cb) {
  var opts = {
    type: type,
    url: url
  };
  apiCallNoSpinner("post", "/connect", opts, cb);
}

function disconnect(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("post", "/disconnect", opts, cb);
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

function getTableInfo(connId, table, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables/" + table + "/info", opts, cb);
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
  apiCall("get", "/getbookmarks", {}, cb);
}

function addBookmark(bookmark, cb) {
  var opts = {
    id: bookmark["id"],
    nick: bookmark["nick"],
    type: bookmark["type"],
    database: bookmark["database"],
    host: bookmark["host"],
    port: bookmark["port"],
    user: bookmark["user"],
    password: bookmark["password"],
  };
  apiCall("post", "/addbookmark", opts, cb);
}

function removeBookmark(id, cb) {
  var opts = { id: id };
  apiCall("post", "/removebookmark", opts, cb);
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
  connect: connect,
  disconnect: disconnect,
  getTables: getTables,
  getTableRows: getTableRows,
  getTableStructure: getTableStructure,
  getTableIndexes: getTableIndexes,
  getTableInfo: getTableInfo,
  getHistory: getHistory,
  getBookmarks: getBookmarks,
  addBookmark: addBookmark,
  removeBookmark: removeBookmark,
  getActivity: getActivity,
  executeQuery: executeQuery,
  explainQuery: explainQuery,
  getConnectionInfo: getConnectionInfo
};
