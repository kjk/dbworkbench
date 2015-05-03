/* jshint -W097,-W117 */
'use strict';

// Loosely inspired by flux ideas.
// One part of the code can trigger an action by calling a function in this
// module. Other parts of the code can provide callbacks to be called when
// action is triggered.

// index in subscribers array for a given action
var tableSelectedIdx = 0;
var viewSelectedIdx = 1;
var executeQueryIdx = 2;
var explainQueryIdx = 3;

// must be in same order as *Idx above
var actionNames = [
  "tableSelected",
  "viewSelected",
  "executeQuery",
  "explainQuery"
];

// index is one of the above constants.
// value at a given index is [[cbFunc, cbId], ...]
var actionCallbacks = [];

// current global callback id to hand out in on()
// we don't bother recycling them after off()
var currCid = 0;

function getActionName(idx) {
  return actionNames[idx] + " (" + idx + ")";
}

function broadcast(actionIdx) {
  var callbacks = actionCallbacks[actionIdx];
  if (!callbacks || callbacks.length === 0) {
    console.log("action.broadcast: no callback for action", getActionName(actionIdx));
    return;
  }

  var args = Array.prototype.slice.call(arguments, 1);
  callbacks.map(function(cbInfo) {
    var cb = cbInfo[0];
    console.log("broadcastAction: calling callback for action", getActionName(actionIdx), "with", args.length, "args");
    if (args.length > 0) {
      cb.apply(null, args);
    } else {
      cb();
    }
  });
}

// subscribe to be notified about an action.
// returns an id that can be used to unsubscribe with off()
function on(action, cb) {
  currCid++;
  var callbacks = actionCallbacks[action];
  var el = [cb, currCid];
  if (!callbacks) {
    actionCallbacks[action] = [el];
  } else {
    callbacks.push(el);
  }
  return currCid;
}

function off(actionIdx, cbId) {
  var callbacks = actionCallbacks[actionIdx];
  if (callbacks && callbacks.length > 0) {
    var n = callbacks.length;
    for (var i = 0; i < n; i++) {
      if (callbacks[i][1] === cbId) {
        callbacks.splice(i, 1);
        return
      }
    }
  }
  console.log("action.off: didn't find callback id", cbId, "for action", getActionName(actionIdx));
}

/* actions */

function tableSelected(name) {
  broadcast(tableSelectedIdx, name);
}

function onTableSelected(cb) {
  return on(tableSelectedIdx, cb);
}

function offTableSelected(cbId) {
  off(tableSelectedIdx, cbId);
}

function viewSelected(view) {
  broadcast(viewSelectedIdx, view);
}

function onViewSelected(cb) {
  return on(viewSelectedIdx, cb);
}

function offViewSelected(cbId) {
  off(viewSelectedIdx, cbId);
}

function executeQuery(query) {
  broadcast(executeQueryIdx, query);
}

function onExecuteQuery(cb) {
  return on(executeQueryIdx, cb);
}

function offExecuteQuery(cbId) {
  off(executeQueryIdx, cbId);
}

function explainQuery(query) {
  broadcast(explainQueryIdx, query);
}

function onExplainQuery(cb) {
  return on(explainQueryIdx, cb);
}

function offExplainQuery(cbId) {
  off(explainQueryIdx, cbId);
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
