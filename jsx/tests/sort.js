var assert = require('assert');

function multiply(x, y) {
  return x*y;
};

describe("multiply", function () {
  it("returns the correct multiplied value", function () {
    var res = multiply(2, 4);
    assert.equal(res, 8);
  });
});


