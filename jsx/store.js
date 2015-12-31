/* reusable part */

/*
Store is a collection global variables and ability to get
notified when the value changes.

Apis:

cid = store.on(key, cb);
cid = store.onMap(key, subkey, cb);

store.off(cid);

store.set(key, value);
store.setMap(key, subkey, cb);

store.del(key);
store.delMap(key, subkey);

del() and delMap() are useful for freeing memory for a given value.

Callbacks always take one argument: new value of the variable.
*/

/*
Maps keys to their values. For non-map values, key is a known string.
For map values, key is key + subkey.
Value is [val, [cb1, cbId1], [cb2, cbId2], ...] i.e. current value followed by zero or more callbacks
*/
let store = {};

// current global callback id to hand out in on()
// we don't bother recycling them after off()
let currCid = 0;

function getFullKey(key, subkey) {
  if (subkey) {
    return key + "-" + subkey;
  }
  return key;
}

function broadcast(key, val, subkey) {
  const fullKey = getFullKey(key, subkey);
  const valAndCbs = store[fullKey];
  const n = valAndCbs.length;
  if (n < 2) {
    console.log(`store.broadcast: no callbacks for key '${fullKey}'`);
    return;
  }

  console.log(`store.broadcast: key: '${fullKey}', val: '${val}', nObservers: ${n-1}`);
  for (let i = 1; i < n; i++) {
    const cb = valAndCbs[i][0];
    cb(val);
  }
}

export function onMap(key, subkey, cb) {
  currCid++;
  const fullKey = getFullKey(key, subkey);
  const cbInfo = [cb, currCid];
  let valAndCbs = store[fullKey];
  if (!valAndCbs) {
    const defVal = defValues[key];
    store[fullKey] = [defVal, cbInfo];
  } else {
    valAndCbs.push(cbInfo);
  }
  return currCid;
}

export function on(key, cb) {
  return onMap(key, null, cb);
}

export function offMap(key, subkey, cbId) {
  const fullKey = getFullKey(key, subkey);
  const valAndCbs = store[fullKey];
  const n = valAndCbs.length;
  for (let i = 1; i < n; i++) {
    const cbId2 = valAndCbs[i][1];
    if (cbId == cbId2) {
      valAndCbs.splice(i, 1);
      return;
    }
  }
  console.log(`store.off: didn't find callback '${cbId}' for '{fullKey}'`);
}

export function off(key, cbId) {
  offMap(key, null, cbId);
}

export function getMap(key, subkey) {
  let fullKey = getFullKey(key, subkey);
  const valAndCbs = store[fullKey];
  if (!valAndCbs) {
    return defValues[key];
  }
  return valAndCbs[0];
}

export function get(key) {
  return getMap(key);
}

/*
shouldBroadcast: some values are synthetic values i.e. the value we broadcast
                 is not the same as raw value
*/
function set2(key, newVal, subkey, shouldBroadcast) {
  let fullKey = getFullKey(key, subkey);
  let valAndCbs = store[fullKey];
  let prevVal;
  if (!valAndCbs) {
    store[fullKey] = [newVal];
    prevVal = defValues[key]; 
  } else {
    prevVal = valAndCbs[0];
    valAndCbs[0] = newVal;
  }

  // optimization: don't notify if those are exactly the same objects
  // TODO: also do a by-value compare
  if (prevVal === newVal) {
    return;
  }
  if (shouldBroadcast) {
    broadcast(key, newVal, subkey);
  }
}

export function setMap(key, newVal, subkey) {
  set2(key, newVal, subkey, true);
}

export function set(key, newVal) {
  set2(key, newVal, null, true);
}

// TODO: not sure if should broadcast deletions or not
export function delMap(key, subkey) {
  const fullKey = getFullKey(key, subkey);
  delete store[fullKey];
}

// TODO: not sure if should broadcast deletions or not
export function del(key) {
  delMap(key);
}

/* things specific to an app */

const queryCmd = "query";
const spinnerCmd = "spinner";

var defValues = {
  "query": null,
  "spinner": 0
};

export function onQuery(queryId, cb) { 
  return onMap(queryCmd, queryId, cb);
}

export function offQuery(queryId, cbId) {
  offMap(queryCmd, queryId, cbId);
}

export function spinnerIsVisible() {
  return get(spinnerCmd) > 0;
}

export function spinnerShow() {
  const newSpinnerState = get(spinnerCmd) + 1;
  set2(spinnerCmd, newSpinnerState, null, false);
  if (newSpinnerState == 1) {
    // we transitioned from 'not visible' to 'visible' state 
    broadcast(spinnerCmd, true);
  }
  //console.log(`spinnerShow: ${newSpinnerState}`);
}

export function spinnerHide() {
  const newSpinnerState = get(spinnerCmd) - 1;
  set2(spinnerCmd, newSpinnerState, null, false);
  if (newSpinnerState == 0) {
    // we transitioned from 'visible' to 'not visible' state
    broadcast(spinnerCmd, false);
  }
  if (newSpinnerState < 0) {
    throw new Error(`negative spinnerState (${newSpinnerState}))`);
  }
  //console.log(`spinnerHide: ${newSpinnerState}`);
}

export function onSpinner(cb) {
  return on(spinnerCmd, cb);
}

export function offSpinner(cbId) {
  off(spinnerCmd, cbId);
}
