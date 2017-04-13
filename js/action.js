/* reusable part */

// Loosely inspired by flux ideas.
// One part of the code can trigger an action by calling a function in this
// module. Other parts of the code can provide callbacks to be called when
// action is triggered.

// index is one of the above constants.
// value at a given index is [[cbFunc, cbId, owner], ...]
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
    console.log(
      "action.broadcast: no callback for action",
      actionName(actionIdx)
    );
    return;
  }

  const args = Array.prototype.slice.call(arguments, 1);
  for (let cbInfo of callbacks) {
    const cb = cbInfo[0];
    console.log(
      "action.broadcast: action: ",
      actionName(actionIdx),
      "args: ",
      args
    );
    if (args.length > 0) {
      cb.apply(null, args);
    } else {
      cb();
    }
  }
}

// subscribe to be notified about an action.
// returns an id that can be used to unsubscribe with off()
function on(actionIdx, cb, owner) {
  currCid++;
  const callbacks = actionCallbacks[actionIdx];
  const el = [cb, currCid, owner];
  if (!callbacks) {
    actionCallbacks[actionIdx] = [el];
  } else {
    callbacks.push(el);
  }
  return currCid;
}

function off(actionIdx, cbIdOrOwner) {
  if (actionIdx == 10) {
    console.log("off: clearFilterIdx");
  }
  const callbacks = actionCallbacks[actionIdx] || [];
  const n = callbacks.length;
  for (let i = 0; i < n; i++) {
    const cbInfo = callbacks[i];
    if (cbInfo[1] === cbIdOrOwner || cbInfo[2] === cbIdOrOwner) {
      callbacks.splice(i, 1);
      return 1 + off(actionIdx, cbIdOrOwner);
    }
  }
  return 0;
  //console.log(`action.off: didn't find callback '${cbIdOrOwner}' for '${actionName(actionIdx)}'`);
}

export function offAllForOwner(owner) {
  let n = 0;
  for (let i = 0; i <= lastIdx; i++) {
    n += off(i, owner, 0);
  }
  if (n == 0) {
    throw Error("didn't find any callbacks for ", owner);
  }
}

/* actions specific to an app */

// index in actionCallbacks array for a given action
const tableSelectedIdx = 0;
const viewSelectedIdx = 1;
const executeQueryIdx = 2;
const explainQueryIdx = 3;
const disconnectDatabaseIdx = 4;
const alertBoxIdx = 5;
const resetPaginationIdx = 6;
const selectedCellPositionIdx = 7;
const editedCellsIdx = 8;
const filterChangedIdx = 9;
const clearFilterIdx = 10;
var lastIdx = 10;

// must be in same order as *Idx above
var actionNames = [
  "tableSelected",
  "viewSelected",
  "executeQuery",
  "explainQuery",
  "disconnectDatabase",
  "alertBox",
  "resetPagination",
  "selectedCellPosition",
  "editedCells",
  "filterChanged",
  "clearFilter",
];

export function tableSelected(name) {
  broadcast(tableSelectedIdx, name);
}

export function onTableSelected(cb, owner) {
  return on(tableSelectedIdx, cb, owner);
}

export function viewSelected(view) {
  broadcast(viewSelectedIdx, view);
}

export function onViewSelected(cb, owner) {
  return on(viewSelectedIdx, cb, owner);
}

export function executeQuery(query) {
  broadcast(executeQueryIdx, query);
}

export function onExecuteQuery(cb, owner) {
  return on(executeQueryIdx, cb, owner);
}

export function explainQuery(query) {
  broadcast(explainQueryIdx, query);
}

export function onExplainQuery(cb, owner) {
  return on(explainQueryIdx, cb, owner);
}

export function disconnectDatabase(query) {
  broadcast(disconnectDatabaseIdx, query);
}

export function onDisconnectDatabase(cb, owner) {
  return on(disconnectDatabaseIdx, cb, owner);
}

export function alertBox(message) {
  broadcast(alertBoxIdx, message);
}

export function onAlertBox(cb, owner) {
  return on(alertBoxIdx, cb, owner);
}

export function resetPagination(toggle) {
  broadcast(resetPaginationIdx, toggle);
}

export function onResetPagination(cb, owner) {
  return on(resetPaginationIdx, cb, owner);
}

export function selectedCellPosition(newPosition) {
  broadcast(selectedCellPositionIdx, newPosition);
}

export function onSelectedCellPosition(cb, owner) {
  return on(selectedCellPositionIdx, cb, owner);
}

export function editedCells(newCells) {
  broadcast(editedCellsIdx, newCells);
}

export function onEditedCells(cb, owner) {
  return on(editedCellsIdx, cb, owner);
}

// Maybe: could be in store
export function filterChanged(s) {
  broadcast(filterChangedIdx, s);
}

export function onFilterChanged(cb, owner) {
  return on(filterChangedIdx, cb, owner);
}

export function clearFilter() {
  broadcast(clearFilterIdx);
}

export function onClearFilter(cb, owner) {
  return on(clearFilterIdx, cb, owner);
}
