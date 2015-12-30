/* reusable part */

/*
Store is a collection global variables and ability to get
notified when the value changes.

Apis:
cid = store.on(variableIdx, cb);
store.off(cid);

store.set(variableIdx, value);
store.delete(variableIdx);

cid = store.onMap(variableIdx, key, cb);
store.setMap(variableIdx, key, cb);

Callbacks always take one argument: new value of the variable.
*/

/*
Maps keys to their values. For non-map values, key is a string keyNames[keyIdx]
For map values, key is keyNames[keyIdx] + subkey.
Value is [val, [cb1, cbId1], [cb2, cbId2], ...] i.e. current value followed by zero or more callbacks
*/
let store = {};

// current global callback id to hand out in on()
// we don't bother recycling them after off()
let currCid = 0;

function keyNamePretty(keyIdx) {
  return `${keyNames[keyIdx]} (${keyIdx})`;
}

function keyToStr(keyIdx, subkey) {
  if (subkey) {
    return keyNames[keyIdx] + "-" + subkey;
  }
  return keyNames[keyIdx];
}

function broadcast(keyIdx, val, subkey) {
  const keyStr = keyToStr(keyIdx, subkey);
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

function on2(keyIdx, cb, subkey) {
  currCid++;
  const keyStr = keyToStr(keyIdx, subkey);
  const cbInfo = [cb, currCid];
  let newValAndCbs = store[keyStr];
  if (!newValAndCbs) {
    const defVal = defValues[keyIdx];
    newValAndCbs = [defVal, cbInfo];
  } else {
    newValAndCbs.push(cbInfo);
  }
  store[keyStr] = newValAndCbs;
  return currCid;
}

export function on(keyIdx, cb) {
  return on2(keyIdx, cb);
}

export function onMap(keyIdx, subkey, cb) {
  return on2(keyIdx, cb, subkey);
}

function off2(keyIdx, cbId, subkey) {
  const keyStr = keyToStr(keyIdx, subkey);
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

export function off(keyIdx, cbId) {
  off2(keyIdx, cbId);
}

export function offMap(keyIdx, subkey, cbId) {
  off2(keyIdx, cbId, subkey); 
}

function get2(keyIdx, subkey) {
  let keyStr = keyToStr(keyIdx, subkey);
  const valAndCbs = store[keyStr];
  if (!valAndCbs) {
    const defVal = defValues[keyIdx]; 
    store[keyStr] = [defVal];
    return defVal;
  }
  return valAndCbs[0];
}

export function get(keyIdx) {
  return get2(keyIdx);
}

export function getMap(keyIdx, subkey) {
  return get2(keyIdx, subkey);
}

/*
shouldBroadcast: some values are synthetic values i.e. the value we broadcast
                 is not the same as raw value
*/
function set2(keyIdx, newVal, subkey, shouldBroadcast) {
  let keyStr = keyToStr(keyIdx, subkey);
  const valAndCbs = store[keyStr];
  if (!valAndCbs) {
    const defVal = defValues[keyIdx]; 
    store[keyStr] = [defVal];
    if (shouldBroadcast) {
      broadcast(keyIdx, defVal, subkey);
    }
  }
  const prevVal = valAndCbs[0];
  // optimization: don't notify if those are exactly the same objects
  // TODO: also do a by-value compare
  if (prevVal === newVal) {
    return;
  }
  valAndCbs[0] = newVal;
  if (shouldBroadcast) {
    broadcast(keyIdx, newVal, subkey);
  }
}

export function set(keyIdx, newVal) {
  set2(keyIdx, newVal, null, true);
}

export function setMap(keyIdx, newVal, subkey) {
  set2(keyIdx, newVal, subkey, true);
}

/* things specific to an app */

// index into keyNames
const queryIdx = 0;
const spinnerIdx = 1;

var keyNames = [
  "query",
  "spinner",
];

var defValues = [
  null,   // query
  0,      //spinner
];

export function onQuery(queryId, cb) { 
  return onMap(queryIdx, queryId, cb);
}

export function offQuery(queryId, cbId) {
  offMap(queryIdx, queryId, cbId);
}

export function spinnerIsVisible() {
  return get(spinnerIdx) > 0;
}

export function spinnerShow() {
  const newSpinnerState = get(spinnerIdx) + 1;
  set2(spinnerIdx, newSpinnerState, null, false);
  if (newSpinnerState == 1) {
    // we transitioned from 'not visible' to 'visible' state 
    broadcast(spinnerIdx, true);
  }
  //console.log(`spinnerShow: ${newSpinnerState}`);
}

export function spinnerHide() {
  const newSpinnerState = get(spinnerIdx) - 1;
  set2(spinnerIdx, newSpinnerState, null, false);
  if (newSpinnerState == 0) {
    // we transitioned from 'visible' to 'not visible' state
    broadcast(spinnerIdx, false);
  }
  if (newSpinnerState < 0) {
    throw new Error(`negative spinnerState (${newSpinnerState}))`);
  }
  //console.log(`spinnerHide: ${newSpinnerState}`);
}

export function onSpinner(cb) {
  return on(spinnerIdx, cb);
}

export function offSpinner(cbId) {
  off(spinnerIdx, cbId);
}

