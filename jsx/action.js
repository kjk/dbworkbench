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

function actionName(idx) {
  return `${actionNames[idx]} (${idx})`;
}

function broadcast(actionIdx) {
  const callbacks = actionCallbacks[actionIdx];
  if (!callbacks || callbacks.length == 0) {
    console.log("action.broadcast: no callback for action", actionName(actionIdx));
    return;
  }

  const args = Array.prototype.slice.call(arguments, 1);
  for(let cbInfo of callbacks) {
    const cb = cbInfo[0];
    console.log("action.broadcast: action: ", actionName(actionIdx), "args: ", args);
    if (args.length > 0) {
      cb.apply(null, args);
    } else {
      cb();
    }
  }
}

// subscribe to be notified about an action.
// returns an id that can be used to unsubscribe with off()
function on(action, cb) {
  currCid++;
  const callbacks = actionCallbacks[action];
  let el = [cb, currCid];
  if (!callbacks) {
    actionCallbacks[action] = [el];
  } else {
    callbacks.push(el);
  }
  return currCid;
}

function off(actionIdx, cbId) {
  const callbacks = actionCallbacks[actionIdx];
  if (callbacks && callbacks.length > 0) {
    const n = callbacks.length;
    for (let i = 0; i < n; i++) {
      if (callbacks[i][1] === cbId) {
        callbacks.splice(i, 1);
        return;
      }
    }
  }
  console.log(`action.off: didn't find callback '${cbId}' for '${actionName(actionIdx)}'`);
}

/* actions specific to an app */

// index in actionCallbacks array for a given action
const tableSelectedIdx = 0;
const viewSelectedIdx = 1;
const executeQueryIdx = 2;
const explainQueryIdx = 3;
const disconnectDatabaseIdx = 4;
const alertBarIdx = 5;
const resetPaginationIdx = 6;
const selectedCellPositionIdx = 7;
const editedCellsIdx = 8;

// must be in same order as *Idx above
var actionNames = [
  "tableSelected",
  "viewSelected",
  "executeQuery",
  "explainQuery",
  "disconnectDatabase",
  "alertBar",
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
