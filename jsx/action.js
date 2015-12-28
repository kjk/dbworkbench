/* jshint -W097,-W117 */
'use strict';

/* reusable part */

// Loosely inspired by flux ideas.
// One part of the code can trigger an action by calling a function in this
// module. Other parts of the code can provide callbacks to be called when
// action is triggered.

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
    console.log("broadcastAction: calling callback for action", getActionName(actionIdx), "with", args, "args");
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
        return;
      }
    }
  }
  console.log("action.off: didn't find callback id", cbId, "for action", getActionName(actionIdx));
}

/* actions specific to an app */

// index in actionCallbacks array for a given action
var tableSelectedIdx = 0;
var viewSelectedIdx = 1;
var executeQueryIdx = 2;
var explainQueryIdx = 3;
var disconnectDatabaseIdx = 4;
var alertBarIdx = 5;
var spinnerIdx = 6;
var resetPaginationIdx = 7;
var selectedCellPositionIdx = 8;
var editedCellsIdx = 9;

// must be in same order as *Idx above
var actionNames = [
  "tableSelected",
  "viewSelected",
  "executeQuery",
  "explainQuery",
  "disconnectDatabase",
  "alertBar",
  "spinner",
  "resetPagination",
  "selectedCellPosition",
  "editedCells",
];

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

function disconnectDatabase(query) {
  broadcast(disconnectDatabaseIdx, query);
}

function onDisconnectDatabase(cb) {
  return on(disconnectDatabaseIdx, cb);
}

function offDisconnectDatabase(cbId) {
  off(disconnectDatabaseIdx, cbId);
}

function alertBar(message) {
  broadcast(alertBarIdx, message);
}

function onAlertBar(cb) {
  return on(alertBarIdx, cb);
}

function offAlertBar(cbId) {
  off(alertBarIdx, cbId);
}

function spinner(toggle) {
  broadcast(spinnerIdx, toggle);
}

function onSpinner(cb) {
  return on(spinnerIdx, cb);
}

function offSpinner(cbId) {
  off(spinnerIdx, cbId);
}

function resetPagination(toggle) {
  broadcast(resetPaginationIdx, toggle);
}

function onResetPagination(cb) {
  return on(resetPaginationIdx, cb);
}

function offResetPagination(cbId) {
  off(resetPaginationIdx, cbId);
}

function selectedCellPosition(newPosition) {
  broadcast(selectedCellPositionIdx, newPosition);
}

function onSelectedCellPosition(cb) {
  return on(selectedCellPositionIdx, cb);
}

function offSelectedCellPosition(cbId) {
  off(selectedCellPositionIdx, cbId);
}

function editedCells(newCells) {
  broadcast(editedCellsIdx, newCells);
}

function onEditedCells(cb) {
  return on(editedCellsIdx, cb);
}

function offEditedCells(cbId) {
  off(editedCellsIdx, cbId);
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
  offExplainQuery: offExplainQuery,

  disconnectDatabase: disconnectDatabase,
  onDisconnectDatabase: onDisconnectDatabase,
  offDisconnectDatabase: offDisconnectDatabase,

  alertBar: alertBar,
  onAlertBar: onAlertBar,
  offAlertBar: offAlertBar,

  spinner: spinner,
  onSpinner: onSpinner,
  offSpinner: offSpinner,

  resetPagination: resetPagination,
  onResetPagination: onResetPagination,
  offResetPagination: offResetPagination,

  selectedCellPosition: selectedCellPosition,
  onSelectedCellPosition: onSelectedCellPosition,
  offSelectedCellPosition: offSelectedCellPosition,

  editedCells: editedCells,
  onEditedCells: onEditedCells,
  offEditedCells: offEditedCells,
};
