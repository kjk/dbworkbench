/* reusable part */

// Loosely inspired by flux ideas.
// One part of the code can trigger an action by calling a function in this
// module. Other parts of the code can provide callbacks to be called when
// action is triggered.

// index is one of the above constants.
// value at a given index is [[cbFunc, cbId], ...]
let actionCallbacks = [];

// current global callback id to hand out in on()
// we don't bother recycling them after off()
let currCid = 0;

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
const tableSelectedIdx = 0;
const viewSelectedIdx = 1;
const executeQueryIdx = 2;
const explainQueryIdx = 3;
const disconnectDatabaseIdx = 4;
const alertBarIdx = 5;
const spinnerIdx = 6;
const resetPaginationIdx = 7;
const selectedCellPositionIdx = 8;
const editedCellsIdx = 9;

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

export function tableSelected(name) {
  broadcast(tableSelectedIdx, name);
}

export function onTableSelected(cb) {
  return on(tableSelectedIdx, cb);
}

export function offTableSelected(cbId) {
  off(tableSelectedIdx, cbId);
}

export function viewSelected(view) {
  broadcast(viewSelectedIdx, view);
}

export function onViewSelected(cb) {
  return on(viewSelectedIdx, cb);
}

export function offViewSelected(cbId) {
  off(viewSelectedIdx, cbId);
}

export function executeQuery(query) {
  broadcast(executeQueryIdx, query);
}

export function onExecuteQuery(cb) {
  return on(executeQueryIdx, cb);
}

export function offExecuteQuery(cbId) {
  off(executeQueryIdx, cbId);
}

export function explainQuery(query) {
  broadcast(explainQueryIdx, query);
}

export function onExplainQuery(cb) {
  return on(explainQueryIdx, cb);
}

export function offExplainQuery(cbId) {
  off(explainQueryIdx, cbId);
}

export function disconnectDatabase(query) {
  broadcast(disconnectDatabaseIdx, query);
}

export function onDisconnectDatabase(cb) {
  return on(disconnectDatabaseIdx, cb);
}

export function offDisconnectDatabase(cbId) {
  off(disconnectDatabaseIdx, cbId);
}

export function alertBar(message) {
  broadcast(alertBarIdx, message);
}

export function onAlertBar(cb) {
  return on(alertBarIdx, cb);
}

export function offAlertBar(cbId) {
  off(alertBarIdx, cbId);
}

// there's only one logical spinner whose state is
// reflected by spinner component. For simplicity we
// want to nest calls to spinnerShow()/spinnerHide()
// and we only notify subscribers on state transitions
let spinnerState = 0;

export function spinnerIsVisible() {
  return spinnerState > 0;
}

export function spinnerShow() {
  spinnerState += 1;
  if (spinnerState == 1) {
    // we transitioned from 'not visible' to 'visible' state 
    broadcast(spinnerIdx, true);
  }
}

export function spinnerHide() {
  spinnerState -= 1;
  if (spinnerState == 0) {
    // we transitioned from 'visible' to 'not visible' state
    broadcast(spinnerIdx, false);
  }
  if (spinnerState < 0) {
    throw new Error(`negative spinnerState (${spinnerState}))`);
  }
}

export function onSpinner(cb) {
  return on(spinnerIdx, cb);
}

export function offSpinner(cbId) {
  off(spinnerIdx, cbId);
}

export function resetPagination(toggle) {
  broadcast(resetPaginationIdx, toggle);
}

export function onResetPagination(cb) {
  return on(resetPaginationIdx, cb);
}

export function offResetPagination(cbId) {
  off(resetPaginationIdx, cbId);
}

export function selectedCellPosition(newPosition) {
  broadcast(selectedCellPositionIdx, newPosition);
}

export function onSelectedCellPosition(cb) {
  return on(selectedCellPositionIdx, cb);
}

export function offSelectedCellPosition(cbId) {
  off(selectedCellPositionIdx, cbId);
}

export function editedCells(newCells) {
  broadcast(editedCellsIdx, newCells);
}

export function onEditedCells(cb) {
  return on(editedCellsIdx, cb);
}

export function offEditedCells(cbId) {
  off(editedCellsIdx, cbId);
}
