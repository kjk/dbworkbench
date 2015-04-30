/* jshint -W097,-W117 */
'use strict';

// Loosely inspired by flux ideas.
// An action is a function in this module.
// One can subscribe to get notified when action happened.

// array of callbacks.
var subscribers = [];

// index in subscribers array for a given action
var tableSelectedIdx = 0;
var viewSelectedIdx = 1;
var executeQueryIdx = 2;
var explainQueryIdx = 3;

// TODO: multiple subscribers
function broadcast(action) {
  var cb = subscribers[action];
  if (cb) {
    var args = Array.prototype.slice.call(arguments, 1);
    console.log("broadcastAction: calling callback ", cb, " for action ", action, " with ", args.length, " args");
    if (args.length > 0) {
      cb.apply(null, args);
    } else {
      cb();
    }
  } else {
    console.log("broadcastAction: no callback for action ", action);
  }
}

// TODO: multiple subscribers
// TODO: should return callback id that can be used with unsubscribeFromAction
function subscribe(action, cb) {
  var currentCb = subscribers[action];
  if (currentCb) {
    console.log("subscribeToAction: already has a callback for action ", action, " will over-write");
  }
  subscribers[action] = cb;
}

function unsubscribe(action, cb) {
  var currentCb = subscribers[action];
  if (currentCb === cb) {
    subscribers[action] = null;
  }
}

/* actions */

function tableSelected(name) {
  broadcast(tableSelectedIdx, name);
}

function onTableSelected(cb) {
  subscribe(tableSelectedIdx, cb);
}

function offTableSelected(cb) {
  unsubscribe(tableSelectedIdx, cb);
}

function viewSelected(view) {
  broadcast(viewSelectedIdx, view);
}

function onViewSelected(cb) {
  subscribe(viewSelectedIdx, cb);
}

function offViewSelected(cb) {
  unsubscribe(viewSelectedIdx, cb);
}

function executeQuery(query) {
  broadcast(executeQueryIdx, query);
}

function onExecuteQuery(cb) {
  subscribe(executeQueryIdx, cb);
}

function offExecuteQuery(cb) {
  unsubscribe(executeQueryIdx, cb);
}

function explainQuery(query) {
  broadcast(explainQueryIdx, query);
}

function onExplainQuery(cb) {
  subscribe(explainQueryIdx, cb);
}

function offExplainQuery(cb) {
  unsubscribe(explainQueryIdx, cb);
}

module.exports = {
  tableSelected: tableSelected,
  onTableSelected: onTableSelected,
  offTableSelected: offTableSelected,

  viewSelected: viewSelected,
  onViewSelected: onViewSelected,
  offViewSelected: offViewSelected,

  executeQuery: executeQuery,
  onExecuteQuery: onExecuteQuery,
  offExecuteQuery: offExecuteQuery,

  explainQuery: explainQuery,
  onExplainQuery: onExplainQuery,
  offExplainQuery: offExplainQuery
};
