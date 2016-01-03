import * as action from './action.js';
import * as store from './store.js';
import 'whatwg-fetch';

function mapObject(o, cb) {
  let res = [];
  for (let k in o) {
    const v = o[k];
    res.push(cb(k, v));
  }
  return res;
}

function formDataFromObject(params) {
  let form = new FormData();
  for (let k in params) {
    const v = params[k];
    form.append(k, v);
  }
  return form;
}

function urlFragmentFromObject(params) {
  let s = '';
  for (let k in params) {
    const v = encodeURIComponent(params[k]);
    k = encodeURIComponent(k);
    if (s == '') {
      s = '?';
    } else {
      s += '&';
    }
    s += `${k}=${v}`;
  }
  return s;
}

function json_ok(json, cb, passJsonErrorToCallback) {
  if (json.error && passJsonErrorToCallback) {
    action.alertBox(`${json.error}`);
    return;
  }
  if (cb) {
    cb(json);
  }
}

function json_failed(error) {
  const msg = error.message;
  action.alertBox(`${msg}`);
}

function fetch_ok(resp, cb, method, url, passJsonErrorToCallback) {
  store.spinnerHide();
  if (!resp.ok) {
    action.alertBox(`Something is wrong. Please restart the application. ${method.toUpperCase()} ${url} Status: ${resp.status} "${resp.statusText}"`);
    return null;
  }
  resp.json().then(
    (json) => json_ok(json, cb, passJsonErrorToCallback),
    (error) => json_failed
  );
}

function fetch_failed(error) {
  store.spinnerHide();
  const msg = error.message;
  action.alertBox(`${msg}`);
}

// passJsonErrorToCallback : if true, a callback will receive json response
//     even if it has error field (used for /api/connect). Otherwise we show
//     the error globally and don't call the callback (most api calls)
function apiCall(method, url, params, cb, passJsonErrorToCallback) {
  const opts = {
    method: method,
    cache: 'no-cache'
  };
  if (method == 'post') {
    opts.body = formDataFromObject(params);
  } else {
    url += urlFragmentFromObject(params);
  }

  store.spinnerShow();
  fetch(url, opts).then(
    (resp) => fetch_ok(resp, cb, method, url, passJsonErrorToCallback),
    (error) => fetch_failed(error)
  );
}

export function connect(type, url, urlSafe, cb) {
  const params = {
    type: type,
    url: url,
    urlSafe: urlSafe
  };
  apiCall('post', '/api/connect', params, cb, true);
}

export function disconnect(connId, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('post', '/api/disconnect', params, cb);
}

export function getTables(connId, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/tables', params, cb);
}

export function getTableStructure(connId, table, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/tables/' + table, params, cb);
}

export function getTableIndexes(connId, table, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/tables/' + table + '/indexes', params, cb);
}

export function getTableInfo(connId, table, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/tables/' + table + '/info', params, cb);
}

export function getHistory(connId, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/history', params, function(data) {
    let rows = [];
    for (let i in data) {
      rows.unshift([parseInt(i) + 1, data[i].query, data[i].timestamp]);
    }
    cb({
      columns: ['id', 'query', 'timestamp'],
      rows: rows
    });
  });
}

export function queryAsync(connId, query, cb) {
  const params = {
    conn_id: connId,
    query: query
  };
  apiCall('post', '/api/queryasync', params, cb);
}

export function queryAsyncStatus(connId, queryId, cb) {
  const params = {
    conn_id: connId,
    query_id: queryId
  };
  apiCall('post', '/api/queryasyncstatus', params, cb);
}

export function queryAsyncData(connId, queryId, start, count, cb) {
  const params = {
    conn_id: connId,
    query_id: queryId,
    start: start,
    count: count
  };
  apiCall('post', '/api/queryasyncdata', params, cb);
}

export function getBookmarks(cb) {
  apiCall('get', '/api/getbookmarks', {}, cb);
}

export function addBookmark(bookmark, cb) {
  const params = {
    id: bookmark['id'],
    nick: bookmark['nick'],
    type: bookmark['type'],
    database: bookmark['database'],
    host: bookmark['host'],
    port: bookmark['port'],
    user: bookmark['user'],
    password: bookmark['password'],
  };
  apiCall('post', '/api/addbookmark', params, cb);
}

export function removeBookmark(id, cb) {
  const params = {
    id: id
  };
  apiCall('post', '/api/removebookmark', params, cb);
}

export function getActivity(connId, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/activity', params, cb);
}

export function executeQuery(connId, query, cb) {
  const params = {
    conn_id: connId,
    query: query
  };
  apiCall('post', '/api/query', params, cb);
}

export function explainQuery(connId, query, cb) {
  const params = {
    conn_id: connId,
    query: query
  };
  apiCall('post', '/api/explain', params, cb);
}

export function getConnectionInfo(connId, cb) {
  const params = {
    conn_id: connId
  };
  apiCall('get', '/api/connection', params, function(data) {
    //const rows = mapObject(data, (k, v) => [k, v]);
    let rows = [];
    for (let key in data) {
      rows.push([key, data[key]]);
    }
    const res = {
      columns: ['attribute', 'value'],
      rows: rows
    };
    cb(res);
  });
}

export function launchBrowserWithURL(url) {
  const params = {
    url: url
  };
  apiCall('get', '/api/launchbrowser', params);
}
