/* reusable part */

/*
Store is a collection global variables and ability to get
notified when the value changes.

Apis:

cid = store.on(key, cb, owner);
cid = store.onMap(key, subkey, cb, owner);

store.off(cidOrOwner);

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
Value is [val, [cb1, cbId1, cbOwner1], [cb2, cbId2, cbOwner2], ...]
i.e. current value followed by zero or more callbacks
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

  if (watchingBroadcast[key]) {
    console.log(`store.broadcast: key: '${fullKey}', val: '${val}', nObservers: ${n-1}`);
  }
  for (let i = 1; i < n; i++) {
    const cb = valAndCbs[i][0];
    cb(val);
  }
}

export function onMap(key, subkey, cb, owner) {
  currCid++;
  const fullKey = getFullKey(key, subkey);
  const cbInfo = [cb, currCid, owner];
  let valAndCbs = store[fullKey];
  if (!valAndCbs) {
    const defVal = defValues[key];
    store[fullKey] = [defVal, cbInfo];
  } else {
    valAndCbs.push(cbInfo);
  }
  return currCid;
}

export function on(key, cb, owner) {
  return onMap(key, null, cb, owner);
}

export function offFullKey(fullKey, cbIdOrOwner) {
  const valAndCbs = store[fullKey];
  if (!valAndCbs) {
    throw new Error("offFullKey for: ", fullKey, " valAndCbs is: ", valAndCbs);
  }
  const n = valAndCbs.length;
  for (let i = 1; i < n; i++) {
    const cbId = valAndCbs[i][1];
    const cbOwner = valAndCbs[i][2];
    if (cbIdOrOwner === cbId || cbIdOrOwner === cbOwner) {
      valAndCbs.splice(i, 1);
      return;
    }
  }
  //console.log(`store.off: didn't find callback '${cbId}' for '{fullKey}'`);
}

export function offMap(key, subkey, cbIdOrOwner) {
  const fullKey = getFullKey(key, subkey);
  offFullKey(fullKey, cbIdOrOwner);
}

export function off(key, cbId) {
  offFullKey(key, cbId);
}

export function offAllForOwner(owner) {
  for (let key in store) {
    offFullKey(key, owner);
  }
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

const spinnerKey = "spinner";
const sidebarDxKey = "sidebarDx";
const queryEditDyKey = "queryEditDy";

// for debugging: keys that we're watching i.e.
// we'll log broadcasting new value
var watchingBroadcast = {
  "queryEditDy": false,
};

var defValues = {
  "spinner": 0,
  "sidebarDx": 250,
  "queryEditDy": 200,
};

export function spinnerIsVisible() {
  return get(spinnerKey) > 0;
}

export function spinnerShow() {
  const newVal = get(spinnerKey) + 1;
  set2(spinnerKey, newVal, null, false);
  if (newVal == 1) {
    // we transitioned from 'not visible' to 'visible' state
    broadcast(spinnerKey, true);
  }
  //console.log(`spinnerShow: ${newVal}`);
}

export function spinnerHide() {
  const newVal = get(spinnerKey) - 1;
  set2(spinnerKey, newVal, null, false);
  if (newVal == 0) {
    // we transitioned from 'visible' to 'not visible' state
    broadcast(spinnerKey, false);
  }
  if (newVal < 0) {
    throw new Error(`negative spinnerState (${newVal}))`);
  }
  //console.log(`spinnerHide: ${newVal}`);
}

export function onSpinner(cb, owner) {
  return on(spinnerKey, cb, owner);
}

export function offSpinner(cbId) {
  off(spinnerKey, cbId);
}

export function onSidebarDx(cb, owner) {
  return on(sidebarDxKey, cb, owner);
}

export function offSidebarDx(cbId) {
  return off(sidebarDxKey, cbId);
}

export function getSidebarDx() {
  return get(sidebarDxKey);
}

export function setSidebarDx(newVal) {
  set(sidebarDxKey, newVal);
}

export function onQueryEditDy(cb, owner) {
  return on(queryEditDyKey, cb, owner);
}

export function offQueryEditDy(cbId) {
  return off(queryEditDyKey, cbId);
}

export function getQueryEditDy() {
  return get(queryEditDyKey);
}

export function setQueryEditDy(newVal) {
  set(queryEditDyKey, newVal);
}
