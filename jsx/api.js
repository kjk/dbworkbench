import $ from 'jquery';
import action from './action.js';

function apiCall(method, path, params, cb) {
  action.spinnerShow();

  $.ajax({
    url: "/api" + path,
    method: method,
    cache: false,
    data: params,
    async: true,
    success: function(data) {
      action.spinnerHide();
      if (cb) {
        cb(data);
      }
    },
    error: function(xhr, status, data) {
      action.spinnerHide();
      if (xhr.status == "0") {
        // Backend is down
        action.alertBar("Something is wrong. Please restart the application.");
      } else {
        // API call failed
      }
      if (cb) {
        cb($.parseJSON(xhr.responseText));
      }
    }
  });
}

export function connect(type, url, urlSafe, cb) {
  var opts = {
    type: type,
    url: url,
    urlSafe: urlSafe
  };
  apiCall("post", "/connect", opts, cb);
}

export function disconnect(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("post", "/disconnect", opts, cb);
}

export function getTables(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables", opts, cb);
}

export function getTableStructure(connId, table, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables/" + table, opts, cb);
}

export function getTableIndexes(connId, table, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables/" + table + "/indexes", opts, cb);
}

export function getTableInfo(connId, table, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/tables/" + table + "/info", opts, cb);
}

export function getHistory(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/history", opts, function(data) {
    var rows = [];
    for (var i in data) {
      rows.unshift([parseInt(i) + 1, data[i].query, data[i].timestamp]);
    }
    cb({ columns: ["id", "query", "timestamp"], rows: rows });
  });
}

export function queryAsync(connId, query, cb) {
  apiCall("post", "/queryasync", {
    conn_id: connId,
    query: query
  }, cb);
}

export function queryAsyncStatus(connId, queryId, cb) {
  apiCall("post", "/queryasyncstatus", {
    conn_id: connId,
    query_id: queryId
  }, cb);
}

export function queryAsyncData(connId, queryId, start, count, cb) {
  apiCall("post", "/queryasyncdata", {
    conn_id: connId,
    query_id: queryId,
    start: start,
    count: count
  }, cb);
}

export function getBookmarks(cb) {
  apiCall("get", "/getbookmarks", {}, cb);
}

export function addBookmark(bookmark, cb) {
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

export function removeBookmark(id, cb) {
  var opts = { id: id };
  apiCall("post", "/removebookmark", opts, cb);
}

export function getActivity(connId, cb) {
  var opts = { conn_id : connId };
  apiCall("get", "/activity", opts, cb);
}

export function executeQuery(connId, query, cb) {
  apiCall("post", "/query", {
    conn_id : connId,
    query: query
  }, cb);
}

export function explainQuery(connId, query, cb) {
  apiCall("post", "/explain", {
    conn_id: connId,
    query: query
  }, cb);
}

export function getConnectionInfo(connId, cb) {
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
