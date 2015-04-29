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

// TODO: unsubscribe

// TODO: multiple subscribers
function broadcastAction(action) {
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
function subscribeToAction(action, cb) {
  var currentCb = subscribers[action];
  if (currentCb) {
    console.log("subscribeToAction: already has a callback for action ", action, " will over-write");
  }
  subscribers[action] = cb;
}

function unsubscribeFromAction(action, cb) {
  var currentCb = subscribers[action];
  if (currentCb === cb) {
    subscribers[action] = null;
  }
}

/* actions */

function tableSelected(name) {
  broadcastAction(tableSelectedIdx, name);
}
function onTableSelected(cb) {
  subscribeToAction(tableSelectedIdx, cb);
}
function viewSelected(view) {
  broadcastAction(viewSelectedIdx, view);
}
function onViewSelected(cb) {
  subscribeToAction(viewSelectedIdx, cb);
}

exports.tableSelected = tableSelected;
exports.onTableSelected = onTableSelected;
exports.viewSelected = viewSelected;
exports.onViewSelected = onViewSelected;
