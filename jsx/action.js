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
    console.log('action.broadcast: no callback for action', actionName(actionIdx));
    return;
  }

  const args = Array.prototype.slice.call(arguments, 1);
  for (let cbInfo of callbacks) {
    const cb = cbInfo[0];
    console.log('action.broadcast: action: ', actionName(actionIdx), 'args: ', args);
    if (args.length > 0) {
      cb.apply(null, args);
    } else {
      cb();
    }
  }
}

// subscribe to be notified about an action.
// returns an id that can be used to unsubscribe with off()
function on(action, cb, owner) {
  currCid++;
  const callbacks = actionCallbacks[action];
  let el = [cb, currCid, owner];
  if (!callbacks) {
    actionCallbacks[action] = [el];
  } else {
    callbacks.push(el);
  }
  return currCid;
}

function off(actionIdx, cbIdOrOwner) {
  const callbacks = actionCallbacks[actionIdx];
  if (callbacks && callbacks.length > 0) {
    const n = callbacks.length;
    for (let i = 0; i < n; i++) {
      if (callbacks[i][1] === cbIdOrOwner || callbacks[i][2] === cbIdOrOwner) {
        callbacks.splice(i, 1);
        return;
      }
    }
  }
  //console.log(`action.off: didn't find callback '${cbIdOrOwner}' for '${actionName(actionIdx)}'`);
}

export function offAllForOwner(owner) {
  for (let i = 0; i < lastIdx; i++) {
    off(i, owner);
  }
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
var lastIdx = 8;

// must be in same order as *Idx above
var actionNames = [
  'tableSelected',
  'viewSelected',
  'executeQuery',
  'explainQuery',
  'disconnectDatabase',
  'alertBar',
  'resetPagination',
  'selectedCellPosition',
  'editedCells',
];

export function tableSelected(name) {
  broadcast(tableSelectedIdx, name);
}

export function onTableSelected(cb, owner) {
  return on(tableSelectedIdx, cb, owner);
}

export function offTableSelected(cbIdOrOwner) {
  off(tableSelectedIdx, cbIdOrOwner);
}

export function viewSelected(view) {
  broadcast(viewSelectedIdx, view);
}

export function onViewSelected(cb, owner) {
  return on(viewSelectedIdx, cb, owner);
}

export function offViewSelected(cbIdOrOwner) {
  off(viewSelectedIdx, cbIdOrOwner);
}

export function executeQuery(query) {
  broadcast(executeQueryIdx, query);
}

export function onExecuteQuery(cb, owner) {
  return on(executeQueryIdx, cb, owner);
}

export function offExecuteQuery(cbIdOrOwner) {
  off(executeQueryIdx, cbIdOrOwner);
}

export function explainQuery(query) {
  broadcast(explainQueryIdx, query);
}

export function onExplainQuery(cb, owner) {
  return on(explainQueryIdx, cb, owner);
}

export function offExplainQuery(cbIdOrOwner) {
  off(explainQueryIdx, cbIdOrOwner);
}

export function disconnectDatabase(query) {
  broadcast(disconnectDatabaseIdx, query);
}

export function onDisconnectDatabase(cb, owner) {
  return on(disconnectDatabaseIdx, cb, owner);
}

export function offDisconnectDatabase(cbIdOrOwner) {
  off(disconnectDatabaseIdx, cbIdOrOwner);
}

export function alertBar(message) {
  broadcast(alertBarIdx, message);
}

export function onAlertBar(cb, owner) {
  return on(alertBarIdx, cb, owner);
}

export function offAlertBar(cbIdOrOwner) {
  off(alertBarIdx, cbIdOrOwner);
}

export function resetPagination(toggle) {
  broadcast(resetPaginationIdx, toggle);
}

export function onResetPagination(cb, owner) {
  return on(resetPaginationIdx, cb, owner);
}

export function offResetPagination(cbIdOrOwner) {
  off(resetPaginationIdx, cbIdOrOwner);
}

export function selectedCellPosition(newPosition) {
  broadcast(selectedCellPositionIdx, newPosition);
}

export function onSelectedCellPosition(cb, owner) {
  return on(selectedCellPositionIdx, cb, owner);
}

export function offSelectedCellPosition(cbIdOrOwner) {
  off(selectedCellPositionIdx, cbIdOrOwner);
}

export function editedCells(newCells) {
  broadcast(editedCellsIdx, newCells);
}

export function onEditedCells(cb, owner) {
  return on(editedCellsIdx, cb, owner);
}

export function offEditedCells(cbIdOrOwner) {
  off(editedCellsIdx, cbIdOrOwner);
}
