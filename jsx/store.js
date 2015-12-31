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

// is used to mark deleted values
const deletedValue = {};

function keyToStr(key, subkey) {
  if (subkey) {
    return key + "-" + subkey;
  }
  return key;
}

function broadcast(key, val, subkey) {
  const keyStr = keyToStr(key, subkey);
  const valAndCbs = store[keyStr];
  const n = valAndCbs.length;
  if (n < 2) {
    console.log(`store.broadcast: no callbacks for key '${keyStr}'`);
    return;
  }

  console.log(`store.broadcast: key: '${keyStr}', val: '${val}', nObservers: ${n-1}`);
  for (let i = 1; i < n; i++) {
    const cb = valAndCbs[i][0];
    cb(val);
  }
}

function on2(key, cb, subkey) {
  currCid++;
  const keyStr = keyToStr(key, subkey);
  const cbInfo = [cb, currCid];
  let valAndCbs = store[keyStr];
  if (!valAndCbs) {
    const defVal = defValues[key];
    store[keyStr] = [defVal, cbInfo];
  } else {
    if (deletedValue === valAndCbs[0]) {
      throw new Error(`trying to subscribe to delete value with key ${keyStr}`);
    }
    valAndCbs.push(cbInfo);
  }
  return currCid;
}

export function on(key, cb) {
  return on2(key, cb);
}

export function onMap(key, subkey, cb) {
  return on2(key, cb, subkey);
}

function off2(key, cbId, subkey) {
  const keyStr = keyToStr(key, subkey);
  const valAndCbs = store[keyStr];
  const n = valAndCbs.length;
  for (let i = 1; i < n; i++) {
    const cbId2 = valAndCbs[i][1];
    if (cbId == cbId2) {
      valAndCbs.splice(i, 1);
      return;
    }
  }
  console.log(`store.off: didn't find callback '${cbId}' for '{keyStr}'`);
}

export function off(key, cbId) {
  off2(key, cbId);
}

export function offMap(key, subkey, cbId) {
  off2(key, cbId, subkey); 
}

function get2(key, subkey) {
  let keyStr = keyToStr(key, subkey);
  const valAndCbs = store[keyStr];
  if (!valAndCbs) {
    const defVal = defValues[key]; 
    store[keyStr] = [defVal];
    return defVal;
  }
  if (deletedValue === valAndCbs[0]) {
    throw new Error(`trying to get delete value with key ${keyStr}`);
  }
  return valAndCbs[0];
}

export function get(key) {
  return get2(key);
}

export function getMap(key, subkey) {
  return get2(key, subkey);
}

/*
shouldBroadcast: some values are synthetic values i.e. the value we broadcast
                 is not the same as raw value
*/
function set2(key, newVal, subkey, shouldBroadcast) {
  if (deletedValue === newVal) {
    shouldBroadcast = false;
  }
  let keyStr = keyToStr(key, subkey);
  const valAndCbs = store[keyStr];
  if (!valAndCbs) {
    const defVal = defValues[key]; 
    store[keyStr] = [defVal];
    if (shouldBroadcast) {
      broadcast(key, defVal, subkey);
    }
    return;
  }
  const prevVal = valAndCbs[0];
  if (deletedValue === prevVal) {
    throw new Error(`trying to set delete value with key ${keyStr}`);
  }

  // optimization: don't notify if those are exactly the same objects
  // TODO: also do a by-value compare
  if (prevVal === newVal) {
    return;
  }
  valAndCbs[0] = newVal;
  if (shouldBroadcast) {
    broadcast(key, newVal, subkey);
  }
}

export function set(key, newVal) {
  set2(key, newVal, null, true);
}

export function setMap(key, newVal, subkey) {
  set2(key, newVal, subkey, true);
}

// TODO: not sure if should broadcast deletions or not
export function del(key) {
  set2(key, deletedValue, null, false);
}

export function delMap(key, subkey) {
  set2(key, deletedValue, subkey, false);
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
