import React from 'react';
import { Table } from './table.jsx';
import { Tr } from './tr.jsx';
import { Td } from './td.jsx';
import { Tfoot } from './tfoot.jsx';
import { Thead } from './thead.jsx';
import { Sort } from './sort.jsx';

React.Children.children = function(children) {
  return React.Children.map(children, function(x) {
      return x;
    }) || [];
};

// Array.prototype.find polyfill - see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/find
if (!Array.prototype.find) {
  Object.defineProperty(Array.prototype, 'find', {
    enumerable: false,
    configurable: true,
    writable: true,
    value: function(predicate) {
      if (this === null) {
        throw new TypeError('Array.prototype.find called on null or undefined');
      }
      if (typeof predicate !== 'function') {
        throw new TypeError('predicate must be a function');
      }
      var list = Object(this);
      var length = list.length >>> 0;
      var thisArg = arguments[1];
      var value;
      for (var i = 0; i < length; i++) {
        if (i in list) {
          value = list[i];
          if (predicate.call(thisArg, value, i, list)) {
            return value;
          }
        }
      }
      return undefined;
    }
  });
}

const Reactable = {
  Table,
  Tr,
  Td,
  Tfoot,
  Thead,
  Sort,
};

export default Reactable;

if (typeof (window) !== 'undefined') {
  window.Reactable = Reactable;
}
